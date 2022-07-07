package types

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "gravity"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querierer msgs
	QuerierRoute = ModuleName
)

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

	// EthereumHeightVoteKey indexes the latest heights observed by each validator
	EthereumHeightVoteKey
)

//////////////////
// Outgoing Txs //
//////////////////

// MakeOutgoingTxKey returns the store index passed with a prefix
func MakeOutgoingTxKey(storeIndex []byte) []byte {
	return append([]byte{OutgoingTxKey}, storeIndex...)
}

//////////////////////
// Send To Etheruem //
//////////////////////

func MakeERC20ToDenomKey(erc20 common.Address) []byte {
	return append([]byte{ERC20ToDenomKey}, erc20.Bytes()...)
}

func MakeSignerSetTxKey(nonce uint64) []byte {
	return append([]byte{SignerSetTxPrefixByte}, sdk.Uint64ToBigEndian(nonce)...)
}

func MakeBatchTxKey(addr common.Address, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxPrefixByte}, addr.Bytes(), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeContractCallTxKey(invalscope []byte, invalnonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxPrefixByte}, invalscope, sdk.Uint64ToBigEndian(invalnonce)}, []byte{})
}
