package keeper

import (
	"sort"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

// TODO here we should check the contents of the validator set against
// the store, if they differ we should take some action to indicate to the
// user that bridge highjacking has occurred
func (k Keeper) SignerSetExecuted(ctx sdk.Context, nonce uint64) {
	otx := k.GetOutgoingTx(ctx, types.MakeSignerSetTxKey(nonce))
	if otx == nil {
		k.Logger(ctx).Error("Failed to clean signer sets",
			"signer set nonce", nonce)
		return
	}

	sstx, ok := otx.(*types.SignerSetTx)
	if !ok {
		panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for outgoing tx %s", otx))
	}

	k.setLastObservedSignerSetTx(ctx, *sstx)

	// We don't use CompleteOutgoingTx here so that BeginBlocker can handle
	// pruning of old signer set txs
	k.SetCompletedOutgoingTx(ctx, sstx)
}

func (k Keeper) GetUnsignedSignerSetTxs(ctx sdk.Context, val sdk.ValAddress) []*types.SignerSetTx {
	var unconfirmed []*types.SignerSetTx
	k.IterateCompletedOutgoingTxsByType(ctx, types.SignerSetTxPrefixByte, func(_ []byte, cotx types.OutgoingTx) bool {
		sig := k.getEthereumSignature(ctx, cotx.GetStoreIndex(), val)
		if len(sig) == 0 {
			signerSet, ok := cotx.(*types.SignerSetTx)
			if !ok {
				panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for completed tx %s", cotx))
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
				panic(errors.Wrapf(types.ErrInvalid, "couldn't cast to signer set for %s", otx))
			}
			unconfirmed = append(unconfirmed, signerSet)
		}
		return false
	})

	return orderSignerSetsByNonceAscending(unconfirmed)
}

// orderSignerSetsByNonceAscending orders the batches by their BatchNonce in ascending order
func orderSignerSetsByNonceAscending(txs []*types.SignerSetTx) []*types.SignerSetTx {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Nonce < txs[j].Nonce
	})

	return txs
}
