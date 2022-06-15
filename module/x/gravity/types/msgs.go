package types

import (
	"fmt"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ sdk.Msg = &MsgDelegateKeys{}
	_ sdk.Msg = &MsgSendToEVM{}
	_ sdk.Msg = &MsgCancelSendToEVM{}
	_ sdk.Msg = &MsgRequestBatchTx{}
	_ sdk.Msg = &MsgSubmitEVMEvent{}
	_ sdk.Msg = &MsgSubmitEVMTxConfirmation{}
	_ sdk.Msg = &MsgEVMHeightVote{}

	_ cdctypes.UnpackInterfacesMessage = &MsgSubmitEVMEvent{}
	_ cdctypes.UnpackInterfacesMessage = &MsgSubmitEVMTxConfirmation{}
	_ cdctypes.UnpackInterfacesMessage = &EVMEventVoteRecord{}
)

// NewMsgDelegateKeys returns a reference to a new MsgDelegateKeys.
func NewMsgDelegateKeys(val sdk.ValAddress, orchAddr sdk.AccAddress, ethAddr string, ethSig []byte) *MsgDelegateKeys {
	return &MsgDelegateKeys{
		ValidatorAddress:    val.String(),
		OrchestratorAddress: orchAddr.String(),
		EVMAddress:          ethAddr,
		EVMSignature:        ethSig,
	}
}

// Route should return the name of the module
func (msg *MsgDelegateKeys) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgDelegateKeys) Type() string { return "delegate_keys" }

// ValidateBasic performs stateless checks
func (msg *MsgDelegateKeys) ValidateBasic() (err error) {
	if _, err = sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.ValidatorAddress)
	}
	if _, err = sdk.AccAddressFromBech32(msg.OrchestratorAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.OrchestratorAddress)
	}
	if !common.IsHexAddress(msg.EVMAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "EVM address")
	}
	if len(msg.EVMSignature) == 0 {
		return ErrEmptyEVMSig
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg *MsgDelegateKeys) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg *MsgDelegateKeys) GetSigners() []sdk.AccAddress {
	acc, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sdk.AccAddress(acc)}
}

// Route should return the name of the module
func (msg *MsgSubmitEVMEvent) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEVMEvent) Type() string { return "submit_EVM_event" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEVMEvent) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event, err := UnpackEvent(msg.Event)
	if err != nil {
		return err
	}

	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}
	return event.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEVMEvent) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEVMEvent) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgSubmitEVMEvent) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var event EVMEvent
	return unpacker.UnpackAny(msg.Event, &event)
}

// Route should return the name of the module
func (msg *MsgSubmitEVMTxConfirmation) Route() string { return RouterKey }

// Type should return the action
func (msg *MsgSubmitEVMTxConfirmation) Type() string { return "submit_EVM_signature" }

// ValidateBasic performs stateless checks
func (msg *MsgSubmitEVMTxConfirmation) ValidateBasic() (err error) {
	if _, err = sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}

	event, err := UnpackConfirmation(msg.Confirmation)

	if err != nil {
		return err
	}

	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}

	return event.Validate()
}

// GetSignBytes encodes the message for signing
func (msg *MsgSubmitEVMTxConfirmation) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg *MsgSubmitEVMTxConfirmation) GetSigners() []sdk.AccAddress {
	// TODO: figure out how to convert between AccAddress and ValAddress properly
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{acc}
}

func (msg *MsgSubmitEVMTxConfirmation) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var sig EVMTxConfirmation
	return unpacker.UnpackAny(msg.Confirmation, &sig)
}

// NewMsgSendToEVM returns a new MsgSendToEVM
func NewMsgSendToEVM(chainID uint32, sender sdk.AccAddress, destAddress string, send sdk.Coin, bridgeFee sdk.Coin) *MsgSendToEVM {
	return &MsgSendToEVM{
		Sender:       sender.String(),
		EVMRecipient: destAddress,
		Amount:       send,
		BridgeFee:    bridgeFee,
		ChainId:      chainID,
	}
}

// Route should return the name of the module
func (msg MsgSendToEVM) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSendToEVM) Type() string { return "send_to_eth" }

// ValidateBasic runs stateless checks on the message
// Checks if the Eth address is valid
func (msg MsgSendToEVM) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}

	// fee and send must be of the same denom
	// this check is VERY IMPORTANT
	if msg.Amount.Denom != msg.BridgeFee.Denom {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins,
			fmt.Sprintf("fee and amount must be the same type %s != %s", msg.Amount.Denom, msg.BridgeFee.Denom))
	}

	if !msg.Amount.IsValid() || msg.Amount.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "amount")
	}
	if !msg.BridgeFee.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "fee")
	}
	if !common.IsHexAddress(msg.EVMRecipient) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "EVM address")
	}
	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSendToEVM) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgSendToEVM) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgRequestBatchTx returns a new msgRequestBatch
func NewMsgRequestBatchTx(chainID uint32, denom string, signer sdk.AccAddress) *MsgRequestBatchTx {
	return &MsgRequestBatchTx{
		Denom:   denom,
		Signer:  signer.String(),
		ChainId: chainID,
	}
}

// Route should return the name of the module
func (msg MsgRequestBatchTx) Route() string { return RouterKey }

// Type should return the action
func (msg MsgRequestBatchTx) Type() string { return "request_batch" }

// ValidateBasic performs stateless checks
func (msg MsgRequestBatchTx) ValidateBasic() error {
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return sdkerrors.Wrap(err, "denom is invalid")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgRequestBatchTx) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgRequestBatchTx) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgCancelSendToEVM returns a new MsgCancelSendToEVM
func NewMsgCancelSendToEVM(chainID uint32, id uint64, orchestrator sdk.AccAddress) *MsgCancelSendToEVM {
	return &MsgCancelSendToEVM{
		Id:      id,
		Sender:  orchestrator.String(),
		ChainId: chainID,
	}
}

// Route should return the name of the module
func (msg MsgCancelSendToEVM) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCancelSendToEVM) Type() string { return "cancel_send_to_EVM" }

// ValidateBasic performs stateless checks
func (msg MsgCancelSendToEVM) ValidateBasic() error {
	if msg.Id == 0 {
		return sdkerrors.Wrap(ErrInvalid, "Id cannot be 0")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender)
	}
	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCancelSendToEVM) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgCancelSendToEVM) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}

// NewMsgEVMHeightVote returns a new MsgEVMHeightVote
func NewMsgEVMHeightVote(chainID uint32, ethereumHeight uint64, signer sdk.AccAddress) *MsgEVMHeightVote {
	return &MsgEVMHeightVote{
		EvmHeight: ethereumHeight,
		Signer:    signer.String(),
		ChainId:   chainID,
	}
}

// Route should return the name of the module
func (msg MsgEVMHeightVote) Route() string { return RouterKey }

// Type should return the action
func (msg MsgEVMHeightVote) Type() string { return "ethereum_height_vote" }

// ValidateBasic performs stateless checks
func (msg MsgEVMHeightVote) ValidateBasic() error {
	if msg.EvmHeight == 0 {
		return sdkerrors.Wrap(ErrInvalid, "ethereum height cannot be 0")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	if msg.ChainId == 0 {
		return sdkerrors.Wrap(ErrUnsupportedEVM, "chain id cannot be zero")
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgEVMHeightVote) GetSignBytes() []byte {
	panic(fmt.Errorf("deprecated"))
}

// GetSigners defines whose signature is required
func (msg MsgEVMHeightVote) GetSigners() []sdk.AccAddress {
	acc, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{acc}
}
