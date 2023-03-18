package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v1"
	v2 "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// Migrate1to2 migrates from consensus version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v1.MigrateStore(ctx, m.keeper.StoreKey, m.keeper.Cdc)
}

// Migrate2to3 migrates from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	newKeeper := v2.NewKeeper{
		StakingKeeper:          m.keeper.StakingKeeper,
		StoreKey:               m.keeper.StoreKey,
		ParamSpace:             m.keeper.ParamSpace,
		Cdc:                    m.keeper.Cdc,
		AccountKeeper:          m.keeper.AccountKeeper,
		BankKeeper:             m.keeper.BankKeeper,
		SlashingKeeper:         m.keeper.SlashingKeeper,
		DistributionKeeper:     m.keeper.DistributionKeeper,
		PowerReduction:         m.keeper.PowerReduction,
		ReceiverModuleAccounts: m.keeper.ReceiverModuleAccounts,
		SenderModuleAccounts:   m.keeper.SenderModuleAccounts,
	}

	return v2.MigrateStore(ctx, &newKeeper)
}
