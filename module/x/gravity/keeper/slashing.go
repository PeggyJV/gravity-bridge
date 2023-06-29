package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetValidatorInfo returns the consensus key address, signing info, and whether or not the validator exists, for the purposes of slashing/jailing
func (k Keeper) GetValidatorSlashingCriteria(ctx sdk.Context, validator stakingtypes.Validator) (signingStartHeight int64, exists bool) {
	consensusKeyAddress, err := validator.GetConsAddr()
	if err != nil {
		panic(fmt.Sprintf("failed to get consensus address: %s", err))
	}
	signingInfo, exists := k.SlashingKeeper.GetValidatorSigningInfo(ctx, consensusKeyAddress)
	signingStartHeight = signingInfo.StartHeight

	return
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
