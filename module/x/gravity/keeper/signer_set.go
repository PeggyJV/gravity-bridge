package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// TODO here we should check the contents of the validator set against
// the store, if they differ we should take some action to indicate to the
// user that bridge highjacking has occurred
func (k Keeper) signerSetExecuted(ctx sdk.Context, nonce uint64) {
	otx := k.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(nonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean signer sets",
			"signer set nonce", nonce)
		return
	}

	completedSignerSetTx, _ := otx.(*types.SignerSetTx)
	k.setLastObservedSignerSetTx(ctx, *completedSignerSetTx)
	k.CompleteOutgoingTx(ctx, completedSignerSetTx)
}

func (k Keeper) GetUnsignedSignerSetTxs(ctx sdk.Context, val sdk.ValAddress) []*types.SignerSetTx {
	var unconfirmed []*types.SignerSetTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, cotx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, cotx.GetStoreIndex(), val)
		if len(sig) == 0 {
			signerSet, ok := cotx.(*types.SignerSetTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for completed tx %s", cotx))
			}
			unconfirmed = append(unconfirmed, signerSet)
		}
		return false
	})
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, otx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, otx.GetStoreIndex(), val)
		if len(sig) == 0 {
			signerSet, ok := otx.(*types.SignerSetTx)
			if !ok {
				panic(sdkerrors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
			}
			unconfirmed = append(unconfirmed, signerSet)
		}
		return false
	})

	return unconfirmed
}
