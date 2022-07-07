package keeper

import (
	"encoding/binary"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2/types"
)

// GetLastSlashedOutgoingTxBlockHeight returns the latest slashed Batch block
func (k Keeper) GetLastSlashedOutgoingTxBlockHeight(ctx sdk.Context) uint64 {
	if bz := ctx.KVStore(k.storeKey).Get([]byte{types.LastSlashedOutgoingTxBlockKey}); bz == nil {
		return 0
	} else {
		return binary.BigEndian.Uint64(bz)
	}
}
