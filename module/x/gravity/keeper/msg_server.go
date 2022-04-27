package keeper

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the gov MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) SetDelegateKeys(c context.Context, msg *types.MsgDelegateKeys) (*types.MsgDelegateKeysResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	orchAddr, err := sdk.AccAddressFromBech32(msg.OrchestratorAddress)
	if err != nil {
		return nil, err
	}

	ethAddr := common.HexToAddress(msg.EVMAddress)

	// ensure that the validator exists
	if k.Keeper.StakingKeeper.Validator(ctx, valAddr) == nil {
		return nil, sdkerrors.Wrap(stakingtypes.ErrNoValidatorFound, valAddr.String())
	}

	// check if the EVM address is currently not used
	validators := k.getValidatorsByEVMAddress(ctx, ethAddr)
	if len(validators) > 0 {
		return nil, sdkerrors.Wrapf(types.ErrDelegateKeys, "ethereum address %s in use", ethAddr)
	}

	// check if the orchestrator address is currently not used
	ethAddrs := k.getEVMAddressesByOrchestrator(ctx, orchAddr)
	if len(ethAddrs) > 0 {
		return nil, sdkerrors.Wrapf(types.ErrDelegateKeys, "orchestrator address %s in use", orchAddr)
	}

	valAccAddr := sdk.AccAddress(valAddr)
	valAccSeq, err := k.accountKeeper.GetSequence(ctx, valAccAddr)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDelegateKeys, "failed to get sequence for validator account %s", valAccAddr)
	}

	var nonce uint64
	if valAccSeq > 0 {
		nonce = valAccSeq - 1
	}

	signMsgBz := k.cdc.MustMarshal(&types.DelegateKeysSignMsg{
		ValidatorAddress: valAddr.String(),
		// We decrement since we process the message after the ante-handler which
		// increments the nonce.
		Nonce: nonce,
	})

	hash := crypto.Keccak256Hash(signMsgBz).Bytes()

	if err = types.ValidateEVMSignature(hash, msg.EthSignature, ethAddr); err != nil {
		return nil, sdkerrors.Wrapf(
			types.ErrDelegateKeys,
			"failed to validate delegate keys signature for EVM address %X; %s ;%d",
			ethAddr, err, nonce,
		)
	}

	k.SetOrchestratorValidatorAddress(ctx, valAddr, orchAddr)
	k.setValidatorEVMAddress(ctx, valAddr, ethAddr)
	k.setEVMOrchestratorAddress(ctx, ethAddr, orchAddr)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeySetOrchestratorAddr, orchAddr.String()),
			sdk.NewAttribute(types.AttributeKeySetEVMAddr, ethAddr.Hex()),
			sdk.NewAttribute(types.AttributeKeyValidatorAddr, valAddr.String()),
		),
	)

	return &types.MsgDelegateKeysResponse{}, nil

}

// SubmitEVMTxConfirmation handles MsgSubmitEVMTxConfirmation
func (k msgServer) SubmitEVMTxConfirmation(c context.Context, msg *types.MsgSubmitEVMTxConfirmation) (*types.MsgSubmitEVMTxConfirmationResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	confirmation, err := types.UnpackConfirmation(msg.Confirmation)
	if err != nil {
		return nil, err
	}

	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	chainID := types.ChainIDOrDefault(msg.ChainId)
	if !k.chainIDsContains(ctx, chainID) {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedEVM, fmt.Sprintf("unsupport chain ID: %d", chainID))
	}

	otx := k.GetOutgoingTx(ctx, chainID, confirmation.GetStoreIndex(chainID))
	if otx == nil {
		k.Logger(ctx).Error(
			"no outgoing tx",
			"store index", fmt.Sprintf("%x", confirmation.GetStoreIndex(chainID)),
		)
		return nil, sdkerrors.Wrap(types.ErrInvalid, "couldn't find outgoing tx")
	}

	gravityID := k.getGravityID(ctx)
	checkpoint := otx.GetCheckpoint([]byte(gravityID))

	ethAddress := k.GetValidatorEVMAddress(ctx, val)
	if ethAddress != confirmation.GetSigner() {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "eth address does not match signer eth address")
	}

	if err = types.ValidateEVMSignature(checkpoint, confirmation.GetSignature(), ethAddress); err != nil {
		k.Logger(ctx).Error("error validating signature",
			"eth addr", ethAddress.String(),
			"gravityID", gravityID,
			"checkpoint", hex.EncodeToString(checkpoint),
			"type url", msg.Confirmation.TypeUrl,
			"signature", hex.EncodeToString(confirmation.GetSignature()),
			"error", err)
		return nil, sdkerrors.Wrap(types.ErrInvalid, fmt.Sprintf(
			"signature verification failed ethAddress %s gravityID %s checkpoint %s typeURL %s signature %s err %s",
			ethAddress.Hex(),
			gravityID,
			hex.EncodeToString(checkpoint),
			msg.Confirmation.TypeUrl,
			hex.EncodeToString(confirmation.GetSignature()),
			err,
		))
	}
	// TODO: should validators be able to overwrite their signatures?
	if k.getEVMSignature(ctx, chainID, confirmation.GetStoreIndex(chainID), val) != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "signature duplicate")
	}

	key := k.SetEVMSignature(ctx, chainID, confirmation, val)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyEVMSignatureKey, string(key)),
		),
	)
	return &types.MsgSubmitEVMTxConfirmationResponse{}, nil
}

// SubmitEVMEvent handles MsgSubmitEVMEvent
func (k msgServer) SubmitEVMEvent(c context.Context, msg *types.MsgSubmitEVMEvent) (*types.MsgSubmitEVMEventResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	event, err := types.UnpackEvent(msg.Event)
	if err != nil {
		return nil, err
	}

	// return an error if the validator isn't in the active set
	val, err := k.getSignerValidator(ctx, msg.Signer)
	if err != nil {
		return nil, err
	}

	chainID := types.ChainIDOrDefault(msg.ChainId)
	if !k.chainIDsContains(ctx, chainID) {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedEVM, fmt.Sprintf("unsupport chain ID: %d", chainID))
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
			sdk.NewAttribute(types.AttributeKeyEVMEventVoteRecordID, string(types.MakeEVMEventVoteRecordKey(chainID, event.GetEventNonce(), event.Hash()))),
		),
	)

	return &types.MsgSubmitEVMEventResponse{}, nil
}

// SendToEVM handles MsgSendToEVM
func (k msgServer) SendToEVM(c context.Context, msg *types.MsgSendToEVM) (*types.MsgSendToEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	chainID := types.ChainIDOrDefault(msg.GetChainId())
	if !k.chainIDsContains(ctx, chainID) {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedEVM, fmt.Sprintf("unsupport chain ID: %d", chainID))
	}

	txID, err := k.createSendToEVM(ctx, chainID, sender, msg.EVMRecipient, msg.Amount, msg.BridgeFee)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents([]sdk.Event{
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawalReceived,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(chainID))),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, strconv.Itoa(int(txID))),
			sdk.NewAttribute(types.AttributeKeyNonce, fmt.Sprint(txID)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(txID)),
		),
	})

	return &types.MsgSendToEVMResponse{Id: txID}, nil
}

// RequestBatchTx handles MsgRequestBatchTx
func (k msgServer) RequestBatchTx(c context.Context, msg *types.MsgRequestBatchTx) (*types.MsgRequestBatchTxResponse, error) {
	// TODO: limit this to only orchestrators and validators?
	ctx := sdk.UnwrapSDKContext(c)

	chainID := types.ChainIDOrDefault(msg.GetChainId())
	if !k.chainIDsContains(ctx, chainID) {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedEVM, fmt.Sprintf("unsupport chain ID: %d", chainID))
	}

	// Check if the denom is a gravity coin, if not, check if there is a deployed ERC20 representing it.
	// If not, error out
	_, tokenContract, err := k.DenomToERC20Lookup(ctx, chainID, msg.Denom)
	if err != nil {
		return nil, err
	}

	batchID := k.BuildBatchTx(ctx, chainID, tokenContract, BatchTxSize)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyContract, tokenContract.Hex()),
			sdk.NewAttribute(types.AttributeKeyBatchNonce, fmt.Sprint(batchID.BatchNonce)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(chainID))),
		),
	)

	return &types.MsgRequestBatchTxResponse{}, nil
}

func (k msgServer) CancelSendToEVM(c context.Context, msg *types.MsgCancelSendToEVM) (*types.MsgCancelSendToEVMResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	chainID := types.ChainIDOrDefault(msg.GetChainId())
	if !k.chainIDsContains(ctx, chainID) {
		return nil, sdkerrors.Wrap(types.ErrUnsupportedEVM, fmt.Sprintf("unsupport chain ID: %d", chainID))
	}

	err := k.Keeper.cancelSendToEVM(ctx, chainID, msg.Id, msg.Sender)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents([]sdk.Event{
		sdk.NewEvent(
			types.EventTypeBridgeWithdrawCanceled,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyContract, k.getBridgeContractAddress(ctx)),
			sdk.NewAttribute(types.AttributeKeyBridgeChainID, strconv.Itoa(int(chainID))),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyOutgoingTXID, fmt.Sprint(msg.Id)),
		),
	})

	return &types.MsgCancelSendToEVMResponse{}, nil
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
