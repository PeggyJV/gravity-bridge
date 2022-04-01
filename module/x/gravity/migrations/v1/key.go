package v1

import "github.com/ethereum/go-ethereum/common"

const (
	// StoreKey to be used when creating the KVStore
	StoreKey = "gravity"
)

// just replicating the current keyset in full rather than trying to hand count and
// manually assign the iota values
const (
	_ = byte(iota)
	// Key Delegation
	ValidatorEthereumAddressKey
	OrchestratorValidatorAddressKey
	EthereumOrchestratorAddressKey

	// Core types
	EthereumSignatureKey
	EthereumEventVoteRecordKey
	OutgoingTxKey
	SendToEthereumKey

	// Latest nonce indexes
	LastEventNonceByValidatorKey
	LastObservedEventNonceKey
	LatestSignerSetTxNonceKey
	LastSlashedOutgoingTxBlockKey
	LastSlashedSignerSetTxNonceKey
	LastOutgoingBatchNonceKey

	// LastSendToEthereumIDKey indexes the lastTxPoolID
	LastSendToEthereumIDKey

	// LastEthereumBlockHeightKey indexes the latest Ethereum block height
	LastEthereumBlockHeightKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey

	// LastUnBondingBlockHeightKey indexes the last validator unbonding block height
	LastUnBondingBlockHeightKey

	LastObservedSignerSetKey
)

func MakeNewERC20ToDenomKey(erc20 common.Address) []byte {
	return append([]byte{ERC20ToDenomKey}, erc20.Bytes()...)
}
