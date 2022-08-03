package keeper

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// createSendToEVM
// - checks a counterpart denominator exists for the given voucher type
// - burns the voucher for transfer amount and fees
// - persists an OutgoingTx
// - adds the TX to the `available` TX pool via a second index
func (k Keeper) createSendToEVM(ctx sdk.Context, chainID uint32, sender sdk.AccAddress, counterpartReceiver string, amount sdk.Coin, fee sdk.Coin) (uint64, error) {
	totalAmount := amount.Add(fee)
	totalInVouchers := sdk.Coins{totalAmount}

	// If the coin is a gravity voucher, burn the coins. If not, check if there is a deployed ERC20 contract representing it.
	// If there is, lock the coins.

	isCosmosOriginated, tokenContract, err := k.DenomToERC20Lookup(ctx, chainID, totalAmount.Denom)
	if err != nil {
		return 0, err
	}

	if senderModule, ok := k.SenderModuleAccounts[sender.String()]; ok {
		if err := k.BankKeeper.SendCoinsFromModuleToModule(ctx, senderModule, types.ModuleName, totalInVouchers); err != nil {
			return 0, err
		}
	} else {
		if err := k.BankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, totalInVouchers); err != nil {
			return 0, err
		}
	}

	// If it is not a cosmos-originated asset we burn
	if !isCosmosOriginated {
		if err := k.BankKeeper.BurnCoins(ctx, types.ModuleName, totalInVouchers); err != nil {
			panic(err)
		}
	}

	// get next tx id from keeper
	nextID := k.incrementLastSendToEVMIDKey(ctx, chainID)

	// construct outgoing tx, as part of this process we represent
	// the token as an ERC20 token since it is preparing to go to ETH
	// rather than the denom that is the input to this function.

	// set the outgoing tx in the pool index
	k.setUnbatchedSendToEVM(ctx, &types.SendToEVM{
		Id:           nextID,
		Sender:       sender.String(),
		EVMRecipient: counterpartReceiver,
		Erc20Token:   types.NewSDKIntERC20Token(amount.Amount, tokenContract),
		Erc20Fee:     types.NewSDKIntERC20Token(fee.Amount, tokenContract),
		ChainId:      chainID,
	})

	return nextID, nil
}

// cancelSendToEVM
// - checks that the provided tx actually exists
// - deletes the unbatched tx from the pool
// - issues the tokens back to the sender
func (k Keeper) cancelSendToEVM(ctx sdk.Context, chainID uint32, id uint64, s string) error {
	sender, _ := sdk.AccAddressFromBech32(s)

	var send *types.SendToEVM
	for _, ste := range k.getUnbatchedSendToEVMs(ctx, chainID) {
		if ste.Id == id {
			send = ste
		}
	}
	if send == nil {
		// NOTE: this case will also be hit if the transaction is in a batch
		return sdkerrors.Wrap(types.ErrInvalid, "id not found in send to ethereum pool")
	}

	if sender.String() != send.Sender {
		return fmt.Errorf("can't cancel a message you didn't send")
	}

	isCosmosOriginated, denom := k.ERC20ToDenomLookup(ctx, chainID, common.HexToAddress(send.Erc20Token.Contract))
	amountToRefund := send.Erc20Token.Amount.Add(send.Erc20Fee.Amount)
	coinsToRefund := sdk.NewCoins(sdk.NewCoin(denom, amountToRefund))

	// If it is not cosmos-originated the coins are minted
	if !isCosmosOriginated {
		if err := k.BankKeeper.MintCoins(ctx, types.ModuleName, coinsToRefund); err != nil {
			return sdkerrors.Wrapf(err, "mint vouchers coins: %s", coinsToRefund)
		}
	}

	if err := k.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coinsToRefund); err != nil {
		return sdkerrors.Wrap(err, "sending coins from module account")
	}

	k.deleteUnbatchedSendToEVM(ctx, chainID, send.Id, send.Erc20Fee)
	return nil
}

func (k Keeper) setUnbatchedSendToEVM(ctx sdk.Context, ste *types.SendToEVM) {
	ctx.KVStore(k.StoreKey).Set(types.MakeSendToEVMKeyForEvent(ste.ChainId, ste.Id, ste.Erc20Fee), k.Cdc.MustMarshal(ste))
}

func (k Keeper) deleteUnbatchedSendToEVM(ctx sdk.Context, chainID uint32, id uint64, fee types.ERC20Token) {
	ctx.KVStore(k.StoreKey).Delete(types.MakeSendToEVMKeyForEvent(chainID, id, fee))
}

func (k Keeper) iterateUnbatchedSendToEVMsByContract(ctx sdk.Context, chainID uint32, contract common.Address, cb func(*types.SendToEVM) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.StoreKey), types.MakeSendToEVMKeyForContract(chainID, contract)).ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ste types.SendToEVM
		k.Cdc.MustUnmarshal(iter.Value(), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (k Keeper) IterateUnbatchedSendToEVMs(ctx sdk.Context, chainID uint32, cb func(*types.SendToEVM) bool) {
	iter := prefix.NewStore(ctx.KVStore(k.StoreKey), types.MakeSendToEVMKey(chainID)).ReverseIterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ste types.SendToEVM
		k.Cdc.MustUnmarshal(iter.Value(), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (k Keeper) getUnbatchedSendToEVMs(ctx sdk.Context, chainID uint32) []*types.SendToEVM {
	var out []*types.SendToEVM
	k.IterateUnbatchedSendToEVMs(ctx, chainID, func(ste *types.SendToEVM) bool {
		out = append(out, ste)
		return false
	})
	return out
}

func (k Keeper) incrementLastSendToEVMIDKey(ctx sdk.Context, chainID uint32) uint64 {
	store := ctx.KVStore(k.StoreKey)
	bz := store.Get(types.MakeLastSendToEVMIDKey(chainID))
	var id uint64 = 0
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	newId := id + 1
	bz = sdk.Uint64ToBigEndian(newId)
	store.Set(types.MakeLastSendToEVMIDKey(chainID), bz)
	return newId
}
