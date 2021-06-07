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

func (s EthereumVoteRecordStore) Get(ctx sdk.Context, eventHash []byte) *types.EthereumEventVoteRecord {
	if bz := s.getStore(ctx).Get(eventHash); bz != nil {
		var r types.EthereumEventVoteRecord
		s.cdc.MustUnmarshalBinaryBare(bz, &r)
		return &r
	}
	return nil
}

func (s EthereumVoteRecordStore) Set(ctx sdk.Context, eventHash []byte, eventVoteRecord *types.EthereumEventVoteRecord) {
	v := s.cdc.MustMarshalBinaryBare(eventVoteRecord)
	s.getStore(ctx).Set(eventHash, v)
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

// GetEthereumEventVoteRecordMapping returns a mapping of eventnonce -> attestations at that nonce
func (s EthereumVoteRecordStore) GetEventNonceMapping(ctx sdk.Context) (out map[uint64][]*types.EthereumEventVoteRecord) {
	out = make(map[uint64][]*types.EthereumEventVoteRecord)
	s.IterateAll(ctx, func(_ []byte, eventVoteRecord *types.EthereumEventVoteRecord) bool {
		event, err := types.UnpackEvent(eventVoteRecord.Event)
		if err != nil {
			panic(err)
		}
		if val, ok := out[event.GetEventNonce()]; !ok {
			out[event.GetEventNonce()] = []*types.EthereumEventVoteRecord{eventVoteRecord}
		} else {
			out[event.GetEventNonce()] = append(val, eventVoteRecord)
		}
		return false
	})
	return
}
