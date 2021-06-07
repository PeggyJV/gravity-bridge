package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common"
)

type SendToEthereumStore struct {
	gravityStoreKey sdk.StoreKey
	cdc             codec.BinaryMarshaler
}

func (s SendToEthereumStore) getStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(s.gravityStoreKey), []byte{types.SendToEthereumKey})
}

func (s SendToEthereumStore) getFeeIdIndexStore(ctx sdk.Context, tokenContract string) prefix.Store {
	return prefix.NewStore(ctx.KVStore(s.gravityStoreKey), append([]byte{types.SendToEthereumFeeIdIndexKey}, []byte(tokenContract)...))
}

func (s SendToEthereumStore) makeFeeIdIndexKey(fee types.ERC20Token, id uint64) []byte {
	amount := make([]byte, 32)
	return append(fee.Amount.BigInt().FillBytes(amount), sdk.Uint64ToBigEndian(id)...)
}

func (s SendToEthereumStore) Set(ctx sdk.Context, ste *types.SendToEthereum) {
	index := s.makeFeeIdIndexKey(ste.Erc20Fee, ste.Id)
	s.getFeeIdIndexStore(ctx, ste.Erc20Fee.Contract).Set(index, sdk.Uint64ToBigEndian(ste.Id))
	s.getStore(ctx).Set(sdk.Uint64ToBigEndian(ste.Id), s.cdc.MustMarshalBinaryBare(ste))
}

func (s SendToEthereumStore) Get(ctx sdk.Context, id uint64) *types.SendToEthereum {
	if bz := s.getStore(ctx).Get(sdk.Uint64ToBigEndian(id)); bz != nil {
		var ste types.SendToEthereum
		s.cdc.MustUnmarshalBinaryBare(bz, &ste)
		return &ste
	}
	return nil
}

func (s SendToEthereumStore) Delete(ctx sdk.Context, id uint64) {
	if row := s.Get(ctx, id); row != nil {
		s.getFeeIdIndexStore(ctx, row.Erc20Fee.Contract).Delete(s.makeFeeIdIndexKey(row.Erc20Fee, id))
		s.getStore(ctx).Delete(sdk.Uint64ToBigEndian(id))
	}
}

func (s SendToEthereumStore) IterateAll(ctx sdk.Context, cb func(*types.SendToEthereum) bool) {
	iter := s.getStore(ctx).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ste types.SendToEthereum
		s.cdc.MustUnmarshalBinaryBare(iter.Value(), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (s SendToEthereumStore) GetAll(ctx sdk.Context) []*types.SendToEthereum {
	var out []*types.SendToEthereum
	s.IterateAll(ctx, func(ste *types.SendToEthereum) bool {
		out = append(out, ste)
		return false
	})
	return out
}

func (s SendToEthereumStore) IterateOrderedByFeeAndId(ctx sdk.Context, contract common.Address, cb func(*types.SendToEthereum) bool) {
	store := s.getStore(ctx)
	indexIter := s.getFeeIdIndexStore(ctx, contract.Hex()).ReverseIterator(nil, nil)
	defer indexIter.Close()
	for ; indexIter.Valid(); indexIter.Next() {
		var ste types.SendToEthereum
		s.cdc.MustUnmarshalBinaryBare(store.Get(indexIter.Value()), &ste)
		if cb(&ste) {
			break
		}
	}
}

func (s SendToEthereumStore) PaginateBySender(ctx sdk.Context, senderAddress string, pageReq *query.PageRequest) ([]*types.SendToEthereum, *query.PageResponse, error) {
	var sendToEthereums []*types.SendToEthereum
	pageRes, err := query.FilteredPaginate(s.getStore(ctx), pageReq, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var ste types.SendToEthereum
		s.cdc.MustUnmarshalBinaryBare(value, &ste)
		if ste.Sender == senderAddress {
			sendToEthereums = append(sendToEthereums, &ste)
			return true, nil
		}
		return false, nil
	})
	return sendToEthereums, pageRes, err
}
