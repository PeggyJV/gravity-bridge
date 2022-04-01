package v1

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	ctx.Logger().Info("Gravity v1 to v2: Beginning store migration")

	store := ctx.KVStore(storeKey)

	migrateCosmosOriginatedERC20ToDenom(store)

	ctx.Logger().Info("Gravty v1 to v2: Store migration complete")

	return nil
}

func migrateCosmosOriginatedERC20ToDenom(store storetypes.KVStore) error {
	prefixStore := prefix.NewStore(store, []byte{ERC20ToDenomKey})
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
