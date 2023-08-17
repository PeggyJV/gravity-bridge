package v2

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func MigrateParamStore(ctx sdktypes.Context, subspace paramstypes.Subspace) error {
	// Don't want to overwrite values if they were set in an upgrade handler
	if !subspace.Has(ctx, types.ParamStoreConfirmedOutgoingTxWindow) {
		subspace.Set(ctx, types.ParamStoreConfirmedOutgoingTxWindow, types.DefaultParams().ConfirmedOutgoingTxWindow)
	}
	if !subspace.Has(ctx, types.ParamStoreEventVoteWindow) {
		subspace.Set(ctx, types.ParamStoreEventVoteWindow, types.DefaultParams().EthereumEventVoteWindow)
	}

	return nil
}
