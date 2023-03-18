package keeper

import (
	"encoding/binary"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// BatchTxSize TODO: should we make this a parameter or a call arg?
const BatchTxSize = 100

// BuildBatchTx starts the following process chain:
//   - find bridged denominator for given voucher type
//   - determine if an unexecuted batch is already waiting for this token type, if so confirm the new batch would
//     have a higher total fees. If not exit withtout creating a batch
//   - select available transactions from the outgoing transaction pool sorted by fee desc
//   - persist an outgoing batch object with an incrementing ID = nonce
//   - emit an event
func (k Keeper) BuildBatchTx(ctx sdk.Context, chainID uint32, contractAddress common.Address, maxElements int) *types.BatchTx {
	// if there is a more profitable batch for this token type do not create a new batch
	if lastBatch := k.getLastOutgoingBatchByTokenType(ctx, chainID, contractAddress); lastBatch != nil {
		if lastBatch.GetFees().GTE(k.getBatchFeesByTokenType(ctx, chainID, contractAddress, maxElements)) {
			return nil
		}
	}

	var selectedStes []*types.SendToEVM
	k.iterateUnbatchedSendToEVMsByContract(ctx, chainID, contractAddress, func(ste *types.SendToEVM) bool {
		selectedStes = append(selectedStes, ste)
		k.deleteUnbatchedSendToEVM(ctx, chainID, ste.Id, ste.Erc20Fee)
		return len(selectedStes) == maxElements
	})

	// do not create batches that would contain no transactions, even if they are requested
	if len(selectedStes) == 0 {
		return nil
	}

	batch := &types.BatchTx{
		BatchNonce:    k.incrementLastOutgoingBatchNonce(ctx, chainID),
		Timeout:       k.getTimeoutHeight(ctx, chainID),
		Transactions:  selectedStes,
		TokenContract: contractAddress.Hex(),
		Height:        uint64(ctx.BlockHeight()),
		ChainId:       chainID,
	}
	k.SetOutgoingTx(ctx, batch)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeOutgoingBatch,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyChainID, fmt.Sprint(chainID)),
		sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(batch.BatchNonce)),
		sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(batch.BatchNonce)),
	))

	return batch
}

// batchTxExecuted is run when the Cosmos chain detects that a batch has been executed on EVM
// It deletes all the transactions in the batch, then cancels all earlier batches
func (k Keeper) batchTxExecuted(ctx sdk.Context, chainID uint32, tokenContract common.Address, nonce uint64) {
	otx := k.GetOutgoingTx(ctx, types.MakeBatchTxStoreIndex(chainID, tokenContract, nonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean batches",
			"token contract", tokenContract.Hex(),
			"nonce", nonce)
		return
	}
	batchTx, _ := otx.(*types.BatchTx)
	k.IterateOutgoingTxsByType(ctx, chainID, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		// If the iterated batches nonce is lower than the one that was just executed, cancel it
		btx, _ := otx.(*types.BatchTx)
		if (btx.BatchNonce < batchTx.BatchNonce) && (btx.TokenContract == batchTx.TokenContract) {
			k.CancelBatchTx(ctx, chainID, btx)
		}
		return false
	})
	k.DeleteOutgoingTx(ctx, batchTx)
}

// getBatchFeesByTokenType gets the fees the next batch of a given token type would
// have if created. This info is both presented to relayers for the purpose of determining
// when to request batches and also used by the batch creation process to decide not to create
// a new batch
func (k Keeper) getBatchFeesByTokenType(ctx sdk.Context, chainID uint32, tokenContractAddr common.Address, maxElements int) sdk.Int {
	feeAmount := sdk.ZeroInt()
	i := 0
	k.iterateUnbatchedSendToEVMsByContract(ctx, chainID, tokenContractAddr, func(tx *types.SendToEVM) bool {
		feeAmount = feeAmount.Add(tx.Erc20Fee.Amount)
		i++
		return i == maxElements
	})

	return feeAmount
}

// CancelBatchTx releases all TX in the batch and deletes the batch
func (k Keeper) CancelBatchTx(ctx sdk.Context, chainID uint32, batch *types.BatchTx) {
	// free transactions from batch and reindex them
	for _, tx := range batch.Transactions {
		k.setUnbatchedSendToEVM(ctx, tx)
	}

	// Delete batch since it is finished
	k.DeleteOutgoingTx(ctx, batch)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeOutgoingBatchCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyChainID, fmt.Sprint(chainID)),
			sdk.NewAttribute(types.AttributeKeyOutgoingBatchID, fmt.Sprint(batch.BatchNonce)),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(batch.BatchNonce)),
		),
	)
}

// getLastOutgoingBatchByTokenType gets the latest outgoing tx batch by token type
func (k Keeper) getLastOutgoingBatchByTokenType(ctx sdk.Context, chainID uint32, token common.Address) *types.BatchTx {
	var lastBatch *types.BatchTx = nil
	lastNonce := uint64(0)
	k.IterateOutgoingTxsByType(ctx, chainID, types.BatchTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		btx, _ := otx.(*types.BatchTx)
		if common.HexToAddress(btx.TokenContract) == token && btx.BatchNonce > lastNonce {
			lastBatch = btx
			lastNonce = btx.BatchNonce
		}
		return false
	})
	return lastBatch
}

// SetLastSlashedOutgoingTxBlockHeight sets the latest slashed Batch block height
func (k Keeper) SetLastSlashedOutgoingTxBlockHeight(ctx sdk.Context, chainID uint32, blockHeight uint64) {
	key := types.MakeLastSlashedOutgoingTxBlockKey(chainID)
	ctx.KVStore(k.StoreKey).Set(key, sdk.Uint64ToBigEndian(blockHeight))
}

// GetLastSlashedOutgoingTxBlockHeight returns the latest slashed Batch block
func (k Keeper) GetLastSlashedOutgoingTxBlockHeight(ctx sdk.Context, chainID uint32) uint64 {
	key := types.MakeLastSlashedOutgoingTxBlockKey(chainID)
	if bz := ctx.KVStore(k.StoreKey).Get(key); bz == nil {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}

func (k Keeper) GetUnSlashedOutgoingTxs(ctx sdk.Context, chainID uint32, maxHeight uint64) (out []types.OutgoingTx) {
	lastSlashed := k.GetLastSlashedOutgoingTxBlockHeight(ctx, chainID)
	k.iterateOutgoingTxs(ctx, chainID, func(key []byte, otx types.OutgoingTx) bool {
		if (otx.GetCosmosHeight() < maxHeight) && (otx.GetCosmosHeight() > lastSlashed) {
			out = append(out, otx)
		}
		return false
	})
	return
}

func (k Keeper) incrementLastOutgoingBatchNonce(ctx sdk.Context, chainID uint32) uint64 {
	store := ctx.KVStore(k.StoreKey)
	key := types.MakeLastOutgoingBatchNonceKey(chainID)
	bz := store.Get(key)
	var id uint64 = 0
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	newId := id + 1
	bz = sdk.Uint64ToBigEndian(newId)
	store.Set(key, bz)
	return newId
}
