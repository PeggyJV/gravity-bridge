package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrInvalid                          = errors.Register(ModuleName, 3, "invalid")
	ErrSupplyOverflow                   = errors.Register(ModuleName, 4, "malicious ERC20 with invalid supply sent over bridge")
	ErrDelegateKeys                     = errors.Register(ModuleName, 5, "failed to delegate keys")
	ErrEmptyEthSig                      = errors.Register(ModuleName, 6, "empty Ethereum signature")
	ErrInvalidERC20Event                = errors.Register(ModuleName, 7, "invalid ERC20 deployed event")
	ErrInvalidEthereumProposalRecipient = errors.Register(ModuleName, 8, "invalid community pool Ethereum spend proposal recipient")
	ErrInvalidEthereumProposalAmount    = errors.Register(ModuleName, 9, "invalid community pool Ethereum spend proposal amount")
	ErrInvalidEthereumProposalBridgeFee = errors.Register(ModuleName, 10, "invalid community pool Ethereum spend proposal bridge fee")
	ErrEthereumProposalDenomMismatch    = errors.Register(ModuleName, 11, "community pool Ethereum spend proposal amount and bridge fee denom mismatch")
	ErrInvalidValidatorAddress          = errors.Register(ModuleName, 12, "invalid validator address")
	ErrInvalidOrchestratorAddress       = errors.Register(ModuleName, 13, "invalid orchestrator address")
)
