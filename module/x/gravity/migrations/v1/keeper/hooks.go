package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/migrations/v1/types"
)

type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks Create new gravity hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)                    {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)          {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) {}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
}

var _ types.GravityHooks = Keeper{}

func (k Keeper) AfterContractCallExecutedEvent(ctx sdk.Context, event types.ContractCallExecutedEvent) {
	if k.hooks != nil {
		k.hooks.AfterContractCallExecutedEvent(ctx, event)
	}
}

func (k Keeper) AfterERC20DeployedEvent(ctx sdk.Context, event types.ERC20DeployedEvent) {
	if k.hooks != nil {
		k.hooks.AfterERC20DeployedEvent(ctx, event)
	}
}

func (k Keeper) AfterSignerSetExecutedEvent(ctx sdk.Context, event types.SignerSetTxExecutedEvent) {
	if k.hooks != nil {
		k.hooks.AfterSignerSetExecutedEvent(ctx, event)
	}
}

func (k Keeper) AfterBatchExecutedEvent(ctx sdk.Context, event types.BatchExecutedEvent) {
	if k.hooks != nil {
		k.hooks.AfterBatchExecutedEvent(ctx, event)
	}
}

func (k Keeper) AfterSendToCosmosEvent(ctx sdk.Context, event types.SendToCosmosEvent) {
	if k.hooks != nil {
		k.hooks.AfterSendToCosmosEvent(ctx, event)
	}
}

func (k *Keeper) SetHooks(sh types.GravityHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set gravity hooks twice")
	}

	k.hooks = sh

	return k
}
