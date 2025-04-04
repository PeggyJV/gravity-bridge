package keeper

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

// TODO: should we make this a parameter or a a call arg?
const BatchTxSize = 100

// CreateBatchTx starts the following process chain:
//   - find bridged denominator for given voucher type
//   - select available transactions from the unbatched SendToEthereums sorted by fee desc, keeping track of the total fees.
//   - if the total fees overflow, stop the iteration and use the SendToEthereums that we have.
//   - if the existing batch is more profitable than the new batch would be, do not create a new batch
//   - persist an OutgoingTx (BatchTx) object with an incrementing ID = nonce
//   - emit an event
func (k Keeper) CreateBatchTx(ctx sdk.Context, contractAddress common.Address, maxElements int) *types.BatchTx {
	var selectedStes []*types.SendToEthereum
	k.iterateUnbatchedSendToEthereumsByContract(ctx, contractAddress, func(ste *types.SendToEthereum) bool {
		selectedStes = append(selectedStes, ste)
		return len(selectedStes) == maxElements
	})

	// do not create batches that would contain no transactions, even if they are requested
	if len(selectedStes) == 0 {
		return nil
	}

	// Sort selected transactions by fee in descending order so we get the most profitable batch possible
	sort.SliceStable(selectedStes, func(i, j int) bool {
		return selectedStes[i].Erc20Fee.Amount.GT(selectedStes[j].Erc20Fee.Amount)
	})

	var finalStes []*types.SendToEthereum
	totalFees := sdk.ZeroInt()
	for _, ste := range selectedStes {
		// Prevent integer overflows by stopping the iteration if the total fees overflow
		newTotal, err := totalFees.SafeAdd(ste.Erc20Fee.Amount)
		if err != nil {
			break
		}
		totalFees = newTotal
		finalStes = append(finalStes, ste)
	}

	if len(finalStes) == 0 {
		return nil
	}

	// if there is a more profitable batch for this token type do not create a new batch
	if lastBatch := k.getLastOutgoingBatchByTokenType(ctx, contractAddress); lastBatch != nil {
		lastBatchFees, err := lastBatch.GetFees()
		if err != nil {
			k.Logger(ctx).Error("Overflowed attempting to get fees for last batch.", "error", err)
			return nil
		}
		if lastBatchFees.GTE(totalFees) {
			return nil
		}
	}

	// Delete the newly batched SendToEthereums
	for _, ste := range finalStes {
		k.deleteUnbatchedSendToEthereum(ctx, ste.Id, ste.Erc20Fee)
	}

	batch := &types.BatchTx{
		BatchNonce:    k.incrementLastOutgoingBatchNonce(ctx),
		Timeout:       k.getTimeoutHeight(ctx),
		Transactions:  finalStes,
		TokenContract: contractAddress.Hex(),
		Height:        uint64(ctx.BlockHeight()),
	}
	k.SetOutgoingTx(ctx, batch)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
		sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(batch.BatchNonce)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(batch.BatchNonce)),
	))

	return batch
}

// batchTxExecuted is run when the Cosmos chain detects that a batch has been executed on Ethereum
// It deletes all the transactions in the batch, then cancels all earlier batches
func (k Keeper) batchTxExecuted(ctx sdk.Context, tokenContract common.Address, nonce uint64) {
	otx := k.GetOutgoingTx(ctx, types.MakeBatchTxKey(tokenContract, nonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean batches",
			"token contract", tokenContract.Hex(),
			"nonce", nonce)
		return
	}
	batchTx, _ := otx.(*types.BatchTx)
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		btx, ok := otx.(*types.BatchTx)
		if !ok {
			panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
		}

		if (btx.BatchNonce < batchTx.BatchNonce) && (btx.TokenContract == batchTx.TokenContract) {
			k.DeleteEthereumSignatures(ctx, btx.GetStoreIndex())
			k.CancelBatchTx(ctx, btx)
		}
		return false
	})

	k.CompleteOutgoingTx(ctx, batchTx)
}

// CancelBatchTx releases all TX in the batch and deletes the batch
func (k Keeper) CancelBatchTx(ctx sdk.Context, batch *types.BatchTx) {
	// free transactions from batch and reindex them
	for _, tx := range batch.Transactions {
		k.setUnbatchedSendToEthereum(ctx, tx)
	}

	// Delete batch since it is finished
	k.DeleteOutgoingTx(ctx, batch.GetStoreIndex())

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatchCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
			sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(batch.BatchNonce)),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(batch.BatchNonce)),
		),
	)
}

// getLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) getLastOutgoingBatchByTokenType(ctx sdk.Context, token common.Address) *types.BatchTx {
	var lastBatch *types.BatchTx = nil
	lastNonce := uint64(0)
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)
		if common.HexToAddress(btx.TokenContract) == token && btx.BatchNonce > lastNonce {
			lastBatch = btx
			lastNonce = btx.BatchNonce
		}
		return false
	})
	return lastBatch
}

// GetUnsignedBatchTxs returns all batches for which the specified validator has not submitted confirmations in ascending nonce order
func (k Keeper) GetUnsignedBatchTxs(ctx sdk.Context, val sdk.ValAddress) []*types.BatchTx {
	var unconfirmed []*types.BatchTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, cotx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, cotx.GetStoreIndex(), val)
		if len(sig) == 0 {
			batch, ok := cotx.(*types.BatchTx)
			if !ok {
				panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for completed tx %s", cotx))
			}
			unconfirmed = append(unconfirmed, batch)
		}
		return false
	})
	k.IterateOutgoingTxsByType(ctx, types.BatchTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 {
			batch, ok := otx.(*types.BatchTx)
			if !ok {
				panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to batch tx for %s", otx))
			}
			unconfirmed = append(unconfirmed, batch)
		}
		return false
	})

	return orderBatchesByNonceAscending(unconfirmed)
}

func (k Keeper) incrementLastOutgoingBatchNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LastOutgoingBatchNonceKey})
	var id uint64 = 0
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	newId := id + 1
	bz = sdk.Uint64ToBigEndian(newId)
	store.Set([]byte{types.LastOutgoingBatchNonceKey}, bz)
	return newId
}

// orderBatchesByNonceAscending orders the batches by their BatchNonce in ascending order
func orderBatchesByNonceAscending(batches []*types.BatchTx) []*types.BatchTx {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].BatchNonce < batches[j].BatchNonce
	})

	return batches
}
