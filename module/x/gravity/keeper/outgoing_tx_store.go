package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type OutgoingTxStore struct {
	gravityStoreKey sdk.StoreKey
	cdc             codec.BinaryMarshaler
}

func (s OutgoingTxStore) getStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(s.gravityStoreKey), []byte{types.OutgoingTxKey})
}

func (s OutgoingTxStore) Set(ctx sdk.Context, outgoing types.OutgoingTx) {
	any, err := types.PackOutgoingTx(outgoing)
	if err != nil {
		panic(err)
	}
	k := outgoing.GetStoreIndex()
	v := s.cdc.MustMarshalBinaryBare(any)
	s.getStore(ctx).Set(k, v)
}

func (s OutgoingTxStore) Get(ctx sdk.Context, storeIndex []byte) types.OutgoingTx {
	v := s.getStore(ctx).Get(storeIndex)
	var out types.OutgoingTx
	if err := s.cdc.UnmarshalInterface(v, &out); err != nil {
		panic(err)
	}
	return out
}

func (s OutgoingTxStore) LatestSignerSetTx(ctx sdk.Context) (*types.SignerSetTx, error) {
	/*
		var otx types.OutgoingTx

		store := prefix.NewStore(ctx.KVStore(k.storeKey), append([]byte{types.OutgoingTxKey}, types.SignerSetTxPrefixByte))
		iter := store.ReverseIterator(nil, nil)
		defer iter.Close()

		var any cdctypes.Any
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &any)

		if err := k.cdc.UnpackAny(&any, &otx); err != nil {
			return nil, err
		}

		ss, ok := otx.(*types.SignerSetTx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "couldn't cast to signer set for latest")
		}
	*/
	return nil, fmt.Errorf("TODO(levi) not implemented")
}
