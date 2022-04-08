package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/migrations/v1/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info("Gravity v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	migrateCosmosOriginatedERC20ToDenom(store)
	migrateContractCallTxTimeouts(store, cdc)

	ctx.Logger().Info("Gravty v1 to v2: Store migration complete")

	return nil
}

func migrateCosmosOriginatedERC20ToDenom(store storetypes.KVStore) error {
	prefixStore := prefix.NewStore(store, []byte{types.ERC20ToDenomKey})
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		oldKey := iter.Key()
		newKey := common.HexToAddress(string(oldKey)).Bytes()

		prefixStore.Delete(oldKey)
		prefixStore.Set(newKey, iter.Value())
	}

	return nil
}

func migrateContractCallTxTimeouts(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	lastObservedEthereumBlockHeight := getLastObservedEthereumBlockHeight(store, cdc).EthereumHeight

	prefixStore := prefix.NewStore(store, types.MakeOutgoingTxKey([]byte{types.ContractCallTxPrefixByte}))
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var any cdctypes.Any
		cdc.MustUnmarshal(iter.Value(), &any)
		var otx types.OutgoingTx
		if err := cdc.UnpackAny(&any, &otx); err != nil {
			panic(err)
		}

		cctx, _ := otx.(*types.ContractCallTx)
		// adding 800 blocks to the last observed height to clean all old stuff out in a few hours
		cctx.Timeout = lastObservedEthereumBlockHeight + uint64(800)
		prefixStore.Delete(iter.Key())

		newAny, err := types.PackOutgoingTx(otx)
		if err != nil {
			panic(err)
		}
		prefixStore.Set(iter.Key(), cdc.MustMarshal(newAny))
	}
	return nil
}

func getLastObservedEthereumBlockHeight(store storetypes.KVStore, cdc codec.BinaryCodec) types.LatestEthereumBlockHeight {
	bytes := store.Get([]byte{types.LastEthereumBlockHeightKey})

	if len(bytes) == 0 {
		return types.LatestEthereumBlockHeight{
			CosmosHeight:   0,
			EthereumHeight: 0,
		}
	}
	height := types.LatestEthereumBlockHeight{}
	cdc.MustUnmarshal(bytes, &height)

	return height
}
