package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("v3 upgrade: entering handler")

		// todo: is this needed if this is the second upgrade?
		fromVM := make(map[string]uint64)
		for moduleName, module := range mm.Modules {
			fromVM[moduleName] = module.ConsensusVersion()
		}

		// Overwrite the gravity module's version back to 1 so the migration will run to v2
		fromVM[gravitytypes.ModuleName] = 2

		ctx.Logger().Info("v3 upgrade: running migrations and exiting handler")
		return mm.RunMigrations(ctx, configurator, fromVM)
	}
}
