package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		// If the iterated contract call's nonce is lower than the one that was just executed, delete it
		sstx, _ := otx.(*types.SignerSetTx)
		if sstx.Nonce < completedSignerSetTx.Nonce {
			k.DeleteEthereumSignatures(ctx, sstx.GetStoreIndex())
			k.DeleteOutgoingTx(ctx, sstx.GetStoreIndex())
		}
		return false
	})

	k.setLastObservedSignerSetTx(ctx, *completedSignerSetTx)
	k.Complete(ctx, completedSignerSetTx)
}
