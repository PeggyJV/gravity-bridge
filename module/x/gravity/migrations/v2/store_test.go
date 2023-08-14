package v2_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	v2 "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestStoreMigration(t *testing.T) {
	gravityKey := storetypes.NewKVStoreKey(v2.ModuleName)
	tGravityKey := storetypes.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(gravityKey, tGravityKey)
	aminoCodec := codec.NewLegacyAmino()
	paramstore := paramtypes.NewSubspace(nil, aminoCodec, gravityKey, tGravityKey, v2.ModuleName)

	// Check no params
	require.False(t, paramstore.Has(ctx, types.ParamStoreConfirmedOutgoingTxWindow))
	require.False(t, paramstore.Has(ctx, types.ParamStoreEventVoteWindow))

	// Run migrations.
	err := v2.MigrateParamStore(ctx, paramstore)
	require.NoError(t, err)

	// Make sure the new params are set.
	require.True(t, paramstore.Has(ctx, types.ParamStoreConfirmedOutgoingTxWindow))
	require.True(t, paramstore.Has(ctx, types.ParamStoreEventVoteWindow))
}
