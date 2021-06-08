package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetDelegateKeys(c context.Context, msg *types.MsgDelegateKeys) (*types.MsgDelegateKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// these errors are checked in validate basic
	val, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	orch, _ := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	eth := common.HexToAddress(msg.EthereumAddress)

	// ensure that the validator exists
	if k.Keeper.StakingKeeper.Validator(ctx, val) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, val.String())
	}

	// set the three indexes
	k.SetOrchestratorValidatorAddress(ctx, val, orch)
	k.setValidatorEthereumAddress(ctx, val, eth)
	k.setEthereumOrchestratorAddress(ctx, eth, orch)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeySetOrchestratorAddr, orch.String()),
			sdk.NewAttribute(types.AttributeKeySetEthereumAddr, eth.Hex()),
			sdk.NewAttribute(types.AttributeKeyValidatorAddr, val.String()),
		),
	)

	return &types.MsgDelegateKeysResponse{}, nil

}

// SubmitEthereumTxConfirmation handles MsgSubmitEthereumTxConfirmation
func (k msgServer) SubmitEthereumTxConfirmation(c context.Context, msg *types.MsgSubmitEthereumTxConfirmation) (*types.MsgSubmitEthereumTxConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	confirmation, err := types.UnpackConfirmation(msg.Confirmation)
	if err != nil {
		return nil, err
	}

	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	otx := k.GetOutgoingTx(ctx, confirmation.GetStoreIndex())
	if otx == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find outgoing tx")
	}

	gravityID := k.getGravityID(ctx)
	checkpoint := otx.GetCheckpoint([]byte(gravityID))

	ethAddress := k.GetValidatorEthereumAddress(ctx, val)
	if ethAddress != confirmation.GetSigner() {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "eth address does not match signer eth address")
	}

	if err = types.ValidateEthereumSignature(checkpoint, confirmation.GetSignature(), ethAddress); err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf(
			"signature verification failed expected sig by %s with gravity-id %s with checkpoint %s found %s %s",
			ethAddress.Hex(),
			gravityID,
			hex.EncodeToString(checkpoint),
			msg.Confirmation.TypeUrl,
			hex.EncodeToString(msg.Confirmation.Value),
		))
	}

	// TODO: should validators be able to overwrite their signatures?
	if k.getEthereumSignature(ctx, confirmation.GetStoreIndex(), val) != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature duplicate")
	}

	key := k.SetEthereumSignature(ctx, confirmation, val)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyEthereumSignatureKey, string(key)),
		),
	)

	return &types.MsgSubmitEthereumTxConfirmationResponse{}, nil
}

// func (k Keeper) ValidateEthereumSignature

// SubmitEthereumEvent handles MsgSubmitEthereumEvent
func (k msgServer) SubmitEthereumEvent(c context.Context, msg *types.MsgSubmitEthereumEvent) (*types.MsgSubmitEthereumEventResponse, error) {
	log.Println(":==: msgServer.SubmitEthereumEvent")

	ctx := sdk.UnwrapSDKContext(c)

	event, err := types.UnpackEvent(msg.Event)
	if err != nil {
		return nil, err
	}
	log.Println(":==: msgServer.SubmitEthereumEvent: event.GetEventNonce()", event.GetEventNonce())

	// return an error if the validator isn't in the active set
	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	// Add the claim to the store
	_, err = k.recordEventVote(ctx, event, val)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "create event vote record")
	}

	// Emit the handle message event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, fmt.Sprintf("%T", event)),
			// TODO: maybe return something better here? is this the right string representation?
			sdk.NewAttribute(types.AttributeKeyEthereumEventVoteRecordID, string(types.MakeEthereumEventVoteRecordKey(event.GetEventNonce(), event.Hash()))),
		),
	)

	return &types.MsgSubmitEthereumEventResponse{}, nil
}

// SendToEthereum handles MsgSendToEthereum
func (k msgServer) SendToEthereum(c context.Context, msg *types.MsgSendToEthereum) (*types.MsgSendToEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	txID, err := k.createSendToEthereum(ctx, sender, msg.EthereumRecipient, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents([]sdk.Event{
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawalReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, strconv.Itoa(int(txID))),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(txID)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	})

	return &types.MsgSendToEthereumResponse{}, nil
}

// RequestBatchTx handles MsgRequestBatchTx
func (k msgServer) RequestBatchTx(c context.Context, msg *types.MsgRequestBatchTx) (*types.MsgRequestBatchTxResponse, error) {
	// TODO: limit this to only orchestrators and validators?
	ctx := sdk.UnwrapSDKContext(c)

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID := k.BuildBatchTx(ctx, tokenContract, BatchTxSize)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyContract, tokenContract.Hex()),
			sdk.NewAttribute(types.AttributeKeyBatchNonce, fmt.Sprint(batchID.BatchNonce)),
		),
	)

	return &types.MsgRequestBatchTxResponse{}, nil
}

func (k msgServer) CancelSendToEthereum(c context.Context, msg *types.MsgCancelSendToEthereum) (*types.MsgCancelSendToEthereumResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	err := k.Keeper.cancelSendToEthereum(ctx, msg.Id, msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents([]sdk.Event{
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(k.getBridgeChainID(ctx)))),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(msg.Id)),
		),
	})

	return &types.MsgCancelSendToEthereumResponse{}, nil
}

// getSignerValidator takes an sdk.AccAddress that represents either a validator or orchestrator address and returns
// the assoicated validator address
func (k Keeper) getSignerValidator(ctx sdk.Context, signerString string) (sdk.ValAddress, error) {
	signer, err := sdk.AccAddressFromBech32(signerString)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signer address")
	}
	var validatorI stakingtypes.ValidatorI
	if validator := k.GetOrchestratorValidatorAddress(ctx, signer); validator == nil {
		validatorI = k.StakingKeeper.Validator(ctx, sdk.ValAddress(signer))
	} else {
		validatorI = k.StakingKeeper.Validator(ctx, validator)
	}

	if validatorI == nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "not orchestrator or validator")
	} else if !validatorI.IsBonded() {
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf("validator is not bonded: %s", validatorI.GetOperator()))
	}

	return validatorI.GetOperator(), nil
}
