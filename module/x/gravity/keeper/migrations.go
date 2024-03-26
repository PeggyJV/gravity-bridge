//go:build exclude

package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v1 "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/migrations/v1"
	v2 "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/migrations/v2"
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
	return v1.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}

// Migrate2to3 migrates from consensus version 2 to 3.
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v2.MigrateParamStore(ctx, m.keeper.paramSpace)
}

// Migrate3to4 migrates from consensus version 3 to 4.
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	ctx.Logger().Info("v3 to v4 migration is a no-op")

	return nil
}
