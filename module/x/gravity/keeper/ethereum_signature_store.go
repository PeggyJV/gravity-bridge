package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type EthereumSignatureStore struct {
	gravityStoreKey sdk.StoreKey
	cdc             codec.BinaryMarshaler
}

func (s EthereumSignatureStore) getStore(ctx sdk.Context, checkpoint []byte) prefix.Store {
	return prefix.NewStore(ctx.KVStore(s.gravityStoreKey), append([]byte{types.EthereumSignatureKey}, checkpoint...))
}

// func (s EthereumSignatureStore) makeKey(checkpoint []byte, ethereumSigner string) []byte {
// 	return append(checkpoint, []byte(ethereumSigner)...)
// }

func (s EthereumSignatureStore) Set(ctx sdk.Context, conf types.EthereumTxConfirmation) {
	s.getStore(ctx, conf.Checkpoint).Set([]byte(conf.EthereumSigner), conf.EthereumSignature)
}

func (s EthereumSignatureStore) Get(ctx sdk.Context, checkpoint []byte, ethereumSigner string) []byte {
	return s.getStore(ctx, checkpoint).Get([]byte(ethereumSigner))
}

func (s EthereumSignatureStore) Delete(ctx sdk.Context, checkpoint []byte, ethereumSigner string) {
	s.getStore(ctx, checkpoint).Delete([]byte(ethereumSigner))
}

func (s EthereumSignatureStore) Iterate(ctx sdk.Context, checkpoint []byte, cb func(ethereumSigner string, signature []byte) bool) {
	iter := s.getStore(ctx, checkpoint).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		// cb returns true to stop early
		if cb(string(iter.Key()), iter.Value()) {
			return
		}
	}
}
