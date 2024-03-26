package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Hooks Create new gravity hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, _ sdk.ConsAddress, valAddress sdk.ValAddress) error {

	// When Validator starts Unbonding, Persist the block height in the store if their power is greater
	// than 1% of the total power.
	// Later in endblocker, check if this persisted block height is the current one and create a signer set tx if it is.
	// The reason for creating signer set txs in endblock is to create only one valset request per block,
	// if multiple validators starts unbonding at same block.

	lastUnbondingBlockHeight := h.k.GetLastUnbondingBlockHeight(ctx)
	if lastUnbondingBlockHeight == uint64(ctx.BlockHeight()) {
		return nil
	}

	latestSignerSet := h.k.GetLatestSignerSetTx(ctx)
	ethAddress := h.k.GetValidatorEthereumAddress(ctx, valAddress).Hex()
	power := uint64(0)
	totalPower := uint64(0)
	for _, s := range latestSignerSet.Signers {
		if s.EthereumAddress == ethAddress {
			power = s.Power
			break
		}

		totalPower += s.Power
	}

	if totalPower == 0 {
		return nil
	}

	proportion := float64(power) / float64(totalPower)
	if proportion > 0.01 {
		h.k.setLastUnbondingBlockHeight(ctx, uint64(ctx.BlockHeight()))
	}

	return nil
}

func (h Hooks) BeforeDelegationCreated(_ sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error { return nil }
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress) error       { return nil }
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterUnbondingInitiated(ctx sdk.Context, _ uint64) error {
	return nil
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
