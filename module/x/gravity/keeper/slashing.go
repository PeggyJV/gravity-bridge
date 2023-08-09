package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type ValidatorSlashingInfo struct {
	Validator   stakingtypes.Validator
	Exists      bool
	SigningInfo slashingtypes.ValidatorSigningInfo
	ConsAddress sdk.ConsAddress
}

// GetUnbondingValidatorSlashingInfos returns the information needed for slashing for each unbonding validator
func (k Keeper) GetUnbondingValidatorSlashingInfos(ctx sdk.Context) ([]stakingtypes.Validator, []ValidatorSlashingInfo) {
	blockTime := ctx.BlockTime().Add(k.StakingKeeper.GetParams(ctx).UnbondingTime)
	blockHeight := ctx.BlockHeight()

	var unbondingValInfos []ValidatorSlashingInfo
	var unbondingValidators []stakingtypes.Validator
	unbondingValIterator := k.StakingKeeper.ValidatorQueueIterator(ctx, blockTime, blockHeight)
	defer unbondingValIterator.Close()
	for ; unbondingValIterator.Valid(); unbondingValIterator.Next() {
		unbondingValidatorsAddr := stakingtypes.ValAddresses{}
		k.cdc.MustUnmarshal(unbondingValIterator.Value(), &unbondingValidatorsAddr)
		for _, valAddr := range unbondingValidatorsAddr.Addresses {
			addr, err := sdk.ValAddressFromBech32(valAddr)
			if err != nil {
				panic(fmt.Sprintf("failed to bech32 decode validator address: %s", err))
			}

			validator, _ := k.StakingKeeper.GetValidator(ctx, addr)
			unbondingValidators = append(unbondingValidators, validator)
			unbondingValInfos = append(unbondingValInfos, k.GetValidatorSlashingInfo(ctx, validator))
		}
	}

	return unbondingValidators, unbondingValInfos
}

// GetBondedValidatorSlashingInfos returns the information needed for slashing for each bonded validator
func (k Keeper) GetBondedValidatorSlashingInfos(ctx sdk.Context) ([]stakingtypes.Validator, []ValidatorSlashingInfo) {
	var bondedValInfos []ValidatorSlashingInfo
	bondedValidators := k.StakingKeeper.GetBondedValidatorsByPower(ctx)
	for _, validator := range bondedValidators {
		bondedValInfos = append(bondedValInfos, k.GetValidatorSlashingInfo(ctx, validator))
	}

	return bondedValidators, bondedValInfos
}

// GetValidatorInfo returns the consensus key address, signing info, and whether or not the validator exists, for the purposes of slashing/jailing
func (k Keeper) GetValidatorSlashingInfo(ctx sdk.Context, validator stakingtypes.Validator) ValidatorSlashingInfo {
	consensusKeyAddress, err := validator.GetConsAddr()
	if err != nil {
		panic(fmt.Sprintf("failed to get consensus address: %s", err))
	}
	signingInfo, exists := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consensusKeyAddress)

	return ValidatorSlashingInfo{validator, exists, signingInfo, consensusKeyAddress}
}

// SlashAndJail slashes the validator and sets the validator to jailed if they are not already jailed
func (k Keeper) SlashAndJail(ctx sdk.Context, validator stakingtypes.Validator, reason string) {
	// Retrieve the validator afresh in case it has been jailed since the first retrieval
	validator, _ = k.StakingKeeper.GetValidator(ctx, validator.GetOperator())
	if validator.IsJailed() {
		return
	}

	consensusKeyAddress, err := validator.GetConsAddr()
	if err != nil {
		panic(fmt.Sprintf("failed to get consensus address: %s", err))
	}

	params := k.GetParams(ctx)
	power := validator.ConsensusPower(k.PowerReduction)

	k.StakingKeeper.Slash(
		ctx,
		consensusKeyAddress,
		ctx.BlockHeight(),
		power,
		// TODO: Differentiate between otx types for slashing fraction in future slashing rework
		params.SlashFractionBatch,
	)
	k.StakingKeeper.Jail(ctx, consensusKeyAddress)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			slashingtypes.EventTypeSlash,
			sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consensusKeyAddress.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consensusKeyAddress.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyReason, reason),
			sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", power)),
		),
	)
}
