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
	SendToEthereumFeeIdIndexKey

	//// Single value keys
	// STORE REFACTOR: figure out which of these are directly related to a table and add their functionality into that table.
	// For the ones that are genuinely standalone values, store them in a separate SingleValueStore which is more like a true KV store
	// Latest nonce indexes
	LastEventNonceByValidatorKey
	LastObservedEventNonceKey
	LatestSignerSetTxNonceKey
	LastSlashedOutgoingTxBlockKey
	LastSlashedSignerSetTxNonceKey
	LastOutgoingBatchNonceKey
	// LastUnBondingBlockHeightKey indexes the last validator unbonding block height
	LastUnBondingBlockHeightKey
	LastObservedSignerSetKey
	// LastEthereumBlockHeightKey indexes the latest Ethereum block height
	LastEthereumBlockHeightKey
	///// End single value keys

	// LastSendToEthereumIDKey indexes the lastTxPoolID
	// This should be integrated into the SendToEthereumStore and completely private
	LastSendToEthereumIDKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey
)

////////////////////
// Key Delegation //
////////////////////
// STORE REFACTOR: nice and easy (hopefully). We want to unify these into one table which will require
// modification of any calling code.

// MakeOrchestratorValidatorAddressKey returns the following key format
// prefix
// [0xe8][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeOrchestratorValidatorAddressKey(orc sdk.AccAddress) []byte {
	return append([]byte{OrchestratorValidatorAddressKey}, orc.Bytes()...)
}

// MakeValidatorEthereumAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeValidatorEthereumAddressKey(validator sdk.ValAddress) []byte {
	return append([]byte{ValidatorEthereumAddressKey}, validator.Bytes()...)
}

// MakeEthereumOrchestratorAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEthereumOrchestratorAddressKey(eth common.Address) []byte {
	return append([]byte{EthereumOrchestratorAddressKey}, eth.Bytes()...)
}

/////////////////////////
// Etheruem Signatures //
/////////////////////////

// MakeEthereumSignatureKey returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEthereumSignatureKey(storeIndex []byte, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{EthereumSignatureKey}, storeIndex, validator.Bytes()}, []byte{})
}

//////////////////
// Outgoing Txs //
//////////////////

// MakeOutgoingTxKey returns the store index passed with a prefix
// func MakeOutgoingTxKey(storeIndex []byte) []byte {
// 	return append([]byte{OutgoingTxKey}, storeIndex...)
// }

//////////////////////
// Send To Etheruem //
//////////////////////

// MakeSendToEthereumKey returns the following key format
//  fee_amount        id
// [1000000000][0 0 0 0 0 0 0 1]
// func MakeSendToEthereumKey(id uint64, fee ERC20Token) []byte {
// 	amount := make([]byte, 32)
// 	return append(fee.Amount.BigInt().FillBytes(amount), sdk.Uint64ToBigEndian(id)...)
// }

// MakeLastEventNonceByValidatorKey indexes lateset event nonce by validator
// MakeLastEventNonceByValidatorKey returns the following key format
// prefix              cosmos-validator
// [0x0][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeLastEventNonceByValidatorKey(validator sdk.ValAddress) []byte {
	return append([]byte{LastEventNonceByValidatorKey}, validator.Bytes()...)
}

// STORE REFACTOR: Change from two byte indexes to a table with
// two indexes to rows that have two fields
///////
func MakeDenomToERC20Key(denom string) []byte {
	return append([]byte{DenomToERC20Key}, []byte(denom)...)
}

func MakeERC20ToDenomKey(erc20 string) []byte {
	return append([]byte{ERC20ToDenomKey}, []byte(erc20)...)
}

//////

func MakeSignerSetTxKey(nonce uint64) []byte {
	return append([]byte{SignerSetTxPrefixByte}, sdk.Uint64ToBigEndian(nonce)...)
}

func MakeBatchTxKey(addr common.Address, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxPrefixByte}, addr.Bytes(), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeContractCallTxKey(invalscope []byte, invalnonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxPrefixByte}, invalscope, sdk.Uint64ToBigEndian(invalnonce)}, []byte{})
}
