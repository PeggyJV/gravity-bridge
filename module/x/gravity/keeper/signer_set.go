package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func (k Keeper) signerSetTxExecuted(ctx sdk.Context, event *types.SignerSetTxExecutedEvent) {
	k.setLastObservedSignerSetTx(ctx, types.SignerSetTx{
		Nonce:   event.SignerSetTxNonce,
		Signers: event.Members,
	})
	otx := k.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(event.SignerSetTxNonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean up signer set txs",
			"signer_set_nonce", event.SignerSetTxNonce)
		return
	}
	signerSetTx, _ := otx.(*types.SignerSetTx)
	k.IterateOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(key []byte, otx types.OutgoingTx) bool {
		sstx, _ := otx.(*types.SignerSetTx)
		if sstx.Nonce < signerSetTx.Nonce {
			k.DeleteEthereumSignatures(ctx, sstx.GetStoreIndex())
			k.DeleteOutgoingTx(ctx, sstx.GetStoreIndex())
		}
		return false
	})

	k.SetCompletedOutgoingTx(ctx, signerSetTx)
	k.DeleteOutgoingTx(ctx, signerSetTx.GetStoreIndex())
}
