package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                     = sdkerrors.Register(ModuleName, 3, "invalid")
	ErrSupplyOverflow              = sdkerrors.Register(ModuleName, 4, "malicious ERC20 with invalid supply sent over bridge")
	ErrDelegateKeys                = sdkerrors.Register(ModuleName, 5, "failed to delegate keys")
	ErrEmptyEVMSig                 = sdkerrors.Register(ModuleName, 6, "empty EVM signature")
	ErrInvalidERC20Event           = sdkerrors.Register(ModuleName, 7, "invalid ERC20 deployed event")
	ErrInvalidEVMProposalRecipient = sdkerrors.Register(ModuleName, 8, "invalid community pool EVM spend proposal recipient")
	ErrInvalidEVMProposalAmount    = sdkerrors.Register(ModuleName, 9, "invalid community pool EVM spend proposal amount")
	ErrInvalidEVMProposalBridgeFee = sdkerrors.Register(ModuleName, 10, "invalid community pool EVM spend proposal bridge fee")
	ErrEVMProposalDenomMismatch    = sdkerrors.Register(ModuleName, 11, "community pool EVM spend proposal amount and bridge fee denom mismatch")
	ErrUnsupportedEVM              = sdkerrors.Register(ModuleName, 12, "invalid chain ID supplied")
	ErrDuplicateGravityID          = sdkerrors.Register(ModuleName, 13, "gravity ID matches one on another evm chain")
)
