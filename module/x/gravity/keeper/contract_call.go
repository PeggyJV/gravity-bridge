package keeper

import (
	"bytes"
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func (k Keeper) GetUnconfirmedContractCallTxs(ctx sdk.Context, val sdk.ValAddress) []*types.ContractCallTx {
	var unconfirmed []*types.ContractCallTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, cotx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, cotx.GetStoreIndex(), val)
		if len(sig) == 0 {
			call, ok := cotx.(*types.ContractCallTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for completed tx %s", cotx))
			}
			unconfirmed = append(unconfirmed, call)
		}
		return false
	})
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 {
			call, ok := otx.(*types.ContractCallTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to contract call for %s", otx))
			}
			unconfirmed = append(unconfirmed, call)
		}
		return false
	})

	return unconfirmed
}

func (k Keeper) contractCallExecuted(ctx sdk.Context, invalidationScope []byte, invalidationNonce uint64) {
	otx := k.GetOutgoingTx(ctx, types.MakeContractCallTxKey(invalidationScope, invalidationNonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean contract calls",
			"invalidation scope", hex.EncodeToString(invalidationScope),
			"invalidation nonce", invalidationNonce)
		return
	}

	completedCallTx, _ := otx.(*types.ContractCallTx)
	k.IterateOutgoingTxsByType(ctx, types.ContractCallTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		// If the iterated contract call's nonce is lower than the one that was just executed, delete it
		cctx, _ := otx.(*types.ContractCallTx)
		if (cctx.InvalidationNonce < completedCallTx.InvalidationNonce) &&
			bytes.Equal(cctx.InvalidationScope, completedCallTx.InvalidationScope) {
			k.DeleteEthereumSignatures(ctx, cctx.GetStoreIndex())
			k.DeleteOutgoingTx(ctx, cctx.GetStoreIndex())
		}
		return false
	})

	k.Complete(ctx, completedCallTx)
}
