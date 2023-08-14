package v2

import (
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func MigrateParamStore(ctx sdktypes.Context, subspace paramstypes.Subspace) error {
	if subspace.HasKeyTable() {
		subspace.Set(ctx, types.ParamStoreConfirmedOutgoingTxWindow, types.DefaultParams().ConfirmedOutgoingTxWindow)
		subspace.Set(ctx, types.ParamStoreEventVoteWindow, types.DefaultParams().EthereumEventVoteWindow)
	} else {
		subspace.WithKeyTable(types.ParamKeyTable())
		subspace.Set(ctx, types.ParamStoreConfirmedOutgoingTxWindow, types.DefaultParams().ConfirmedOutgoingTxWindow)
		subspace.Set(ctx, types.ParamStoreEventVoteWindow, types.DefaultParams().EthereumEventVoteWindow)
	}

	return nil
}
