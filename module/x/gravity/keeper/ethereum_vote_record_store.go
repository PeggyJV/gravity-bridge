package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type EthereumVoteRecordStore struct {
	gravityStoreKey sdk.StoreKey
	cdc             codec.BinaryMarshaler
}

func (s EthereumVoteRecordStore) getStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(s.gravityStoreKey), []byte{types.EthereumEventVoteRecordKey})
}

func (s EthereumVoteRecordStore) makeKey(eventNonce uint64, eventHash []byte) []byte {
	return append(sdk.Uint64ToBigEndian(eventNonce), eventHash...)
}

func (s EthereumVoteRecordStore) Get(ctx sdk.Context, eventNonce uint64, eventHash []byte) *types.EthereumEventVoteRecord {
	k := s.makeKey(eventNonce, eventHash)
	if bz := s.getStore(ctx).Get(k); bz != nil {
		var r types.EthereumEventVoteRecord
		s.cdc.MustUnmarshalBinaryBare(bz, &r)
		return &r
	}
	return nil
}

func (s EthereumVoteRecordStore) Set(ctx sdk.Context, eventNonce uint64, eventHash []byte, eventVoteRecord *types.EthereumEventVoteRecord) {
	k := s.makeKey(eventNonce, eventHash)
	v := s.cdc.MustMarshalBinaryBare(eventVoteRecord)
	s.getStore(ctx).Set(k, v)
}

func (s EthereumVoteRecordStore) IterateAll(ctx sdk.Context, cb func([]byte, *types.EthereumEventVoteRecord) bool) {
	iter := s.getStore(ctx).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		att := &types.EthereumEventVoteRecord{}
		s.cdc.MustUnmarshalBinaryBare(iter.Value(), att)
		// cb returns true to stop early
		if cb(iter.Key(), att) {
			return
		}
	}
}
