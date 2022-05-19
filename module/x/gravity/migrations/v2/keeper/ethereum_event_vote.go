package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/types"
)

// GetLastObservedEventNonce returns the latest observed event nonce
func (k Keeper) GetLastObservedEventNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte{types.LastObservedEventNonceKey})

	if len(bytes) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(bytes)
}

// GetLastObservedEthereumBlockHeight height gets the block height to of the last observed attestation from
// the store
func (k Keeper) GetLastObservedEthereumBlockHeight(ctx sdk.Context) types.LatestEthereumBlockHeight {
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get([]byte{types.LastEthereumBlockHeightKey})

	if len(bytes) == 0 {
		return types.LatestEthereumBlockHeight{
			CosmosHeight:   0,
			EthereumHeight: 0,
		}
	}
	height := types.LatestEthereumBlockHeight{}
	k.Cdc.MustUnmarshal(bytes, &height)
	return height
}
