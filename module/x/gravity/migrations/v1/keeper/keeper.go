//go:build exclude

package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/migrations/v1/types"
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	StakingKeeper  types.StakingKeeper
	storeKey       storetypes.StoreKey
	paramSpace     paramtypes.Subspace
	cdc            codec.Codec
	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	SlashingKeeper types.SlashingKeeper
	PowerReduction sdk.Int
	hooks          types.GravityHooks
}

// NewKeeper returns a new instance of the gravity keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accKeeper types.AccountKeeper,
	stakingKeeper types.StakingKeeper,
	bankKeeper types.BankKeeper,
	slashingKeeper types.SlashingKeeper,
	powerReduction sdk.Int,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	k := Keeper{
		cdc:            cdc,
		paramSpace:     paramSpace,
		storeKey:       storeKey,
		accountKeeper:  accKeeper,
		StakingKeeper:  stakingKeeper,
		bankKeeper:     bankKeeper,
		SlashingKeeper: slashingKeeper,
		PowerReduction: powerReduction,
	}

	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, ps types.Params) {
	k.paramSpace.SetParamSet(ctx, &ps)
}
