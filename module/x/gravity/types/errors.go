package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalid                     = sdkerrors.Register(ModuleName, 103, "invalid")
	ErrSupplyOverflow              = sdkerrors.Register(ModuleName, 104, "malicious ERC20 with invalid supply sent over bridge")
	ErrDelegateKeys                = sdkerrors.Register(ModuleName, 105, "failed to delegate keys")
	ErrEmptyEVMSig                 = sdkerrors.Register(ModuleName, 106, "empty EVM signature")
	ErrInvalidERC20Event           = sdkerrors.Register(ModuleName, 107, "invalid ERC20 deployed event")
	ErrInvalidEVMProposalRecipient = sdkerrors.Register(ModuleName, 108, "invalid community pool EVM spend proposal recipient")
	ErrInvalidEVMProposalAmount    = sdkerrors.Register(ModuleName, 109, "invalid community pool EVM spend proposal amount")
	ErrInvalidEVMProposalBridgeFee = sdkerrors.Register(ModuleName, 110, "invalid community pool EVM spend proposal bridge fee")
	ErrEVMProposalDenomMismatch    = sdkerrors.Register(ModuleName, 111, "community pool EVM spend proposal amount and bridge fee denom mismatch")
	ErrUnsupportedEVM              = sdkerrors.Register(ModuleName, 112, "invalid chain ID supplied")
	ErrDuplicateGravityID          = sdkerrors.Register(ModuleName, 113, "gravity ID matches one on another evm chain")
	ErrDuplicateChainID            = sdkerrors.Register(ModuleName, 114, "chain ID matches one on another evm chain")
)
