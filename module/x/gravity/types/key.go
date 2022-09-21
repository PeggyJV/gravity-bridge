package types

import (
	"bytes"
	"encoding/binary"

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

	// QuerierRoute to be used for query msgs
	QuerierRoute = ModuleName
)

const (
	_ = byte(iota)
	// Key Delegation
	ValidatorToEVMAddressKey
	OrchestratorToValidatorAddressKey
	EVMToOrchestratorAddressKey

	// Core types
	EVMSignatureKey
	EVMEventVoteRecordKey
	OutgoingTxKey
	SendToEVMKey

	// Latest nonce indexes
	LastEventNonceByValidatorKey
	LastObservedEventNonceKey
	LatestSignerSetTxNonceKey
	LastSlashedOutgoingTxBlockKey
	LastOutgoingBatchNonceKey

	// LastSendToEVMIDKey indexes the lastTxPoolID
	LastSendToEVMIDKey

	// LastEVMBlockHeightKey indexes the latest EVM block height
	LastEVMBlockHeightKey

	// DenomToERC20Key prefixes the index of Cosmos originated asset denoms to ERC20s
	DenomToERC20Key

	// ERC20ToDenomKey prefixes the index of Cosmos originated assets ERC20s to denoms
	ERC20ToDenomKey

	// LastUnBondingBlockHeightKey indexes the last validator unbonding block height
	LastUnBondingBlockHeightKey

	// LastObservedSignerSetKey index the last observed signer set
	LastObservedSignerSetKey

	// EVMHeightVoteKey indexes the latest heights observed by each validator
	EVMHeightVoteKey
)

////////////////////
// Key Delegation //
////////////////////

// Uint32ToBigEndian - marshals uint64 to a bigendian byte slice so it can be sorted
func Uint32ToBigEndian(i uint32) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint32(b, i)
	return b
}

// MakeOrchestratorValidatorAddressKey returns the following key format
// [0x2][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeOrchestratorValidatorAddressKey(orc sdk.AccAddress) []byte {
	return bytes.Join([][]byte{{OrchestratorToValidatorAddressKey}, orc.Bytes()}, []byte{})
}

// MakeValidatorEVMAddressKey returns the following key format
// prefix chain  cosmos-validator
// [0x1][0001][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeValidatorEVMAddressKey(validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{ValidatorToEVMAddressKey}, validator.Bytes()}, []byte{})
}

// MakeEVMOrchestratorAddressKey returns the following key format
// prefix chain  cosmos-validator
// [0x3][0001][0xc783df8a850f42e7F7e57013759C285caa701eB6]
func MakeEVMOrchestratorAddressKey(eth common.Address) []byte {
	return bytes.Join([][]byte{{EVMToOrchestratorAddressKey}, eth.Bytes()}, []byte{})
}

/////////////////////////
// Ethereum Signatures //
/////////////////////////

// MakeEVMSignatureKeyForValidator returns the following key format
// prefix chain  nonce                    validator-address
// [0x4][0001][0 0 0 0 0 0 0 1][][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEVMSignatureKeyForValidator(storeIndex []byte, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{EVMSignatureKeyStoreIndexPrefix(storeIndex), validator.Bytes()}, []byte{})
}

func EVMSignatureKeyStoreIndexPrefix(storeIndex []byte) []byte {
	return bytes.Join([][]byte{EVMSignatureKeyPrefix(), storeIndex}, []byte{})
}

func EVMSignatureKeyPrefix() []byte {
	return []byte{EVMSignatureKey}
}

/////////////////////////////////
// Ethereum Event Vote Records //
/////////////////////////////////

func MakeEVMEventVoteRecordKey(chainID uint32, eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{EVMEventVoteRecordPrefix(chainID), sdk.Uint64ToBigEndian(eventNonce), claimHash}, []byte{})
}

func EVMEventVoteRecordPrefix(chainID uint32) []byte {
	return bytes.Join([][]byte{{EVMEventVoteRecordKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

//////////////////
// Outgoing Txs //
//////////////////

// MakeOutgoingTxKey returns the store index passed with a prefix
func MakeOutgoingTxKey(storeIndex []byte) []byte {
	return bytes.Join([][]byte{{OutgoingTxKey}, storeIndex}, []byte{})
}

func OutgoingTxKeyPrefixWithPrefixByte(chainID uint32, prefix byte) []byte {
	return bytes.Join([][]byte{OutgoingTxKeyPrefix(), {prefix}, Uint32ToBigEndian(chainID)}, []byte{})
}

func OutgoingTxKeyPrefix() []byte {
	return []byte{OutgoingTxKey}
}

/////////////////
// Send To EVM //
/////////////////

// MakeSendToEVMKey returns the following key format
// prefix chain           eth-contract-address            fee_amount        id
// [0x7][0001][0xc783df8a850f42e7F7e57013759C285caa701eB6][1000000000][0 0 0 0 0 0 0 1]
func MakeSendToEVMKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{SendToEVMKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeSendToEVMKeyForContract(chainID uint32, contract common.Address) []byte {
	return bytes.Join([][]byte{MakeSendToEVMKey(chainID), contract.Bytes()}, []byte{})
}

func MakeSendToEVMKeyForEvent(chainID uint32, id uint64, fee ERC20Token) []byte {
	amount := make([]byte, 32)
	return bytes.Join([][]byte{MakeSendToEVMKeyForContract(chainID, common.HexToAddress(fee.Contract)), fee.Amount.BigInt().FillBytes(amount), sdk.Uint64ToBigEndian(id)}, []byte{})
}

// MakeLastEventNonceByValidatorKey indexes lateset event nonce by validator
// MakeLastEventNonceByValidatorKey returns the following key format
// prefix chain              cosmos-validator
// [0x8][0001][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeLastEventNonceByValidatorKey(chainID uint32, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{LastEventNonceByValidatorKey}, Uint32ToBigEndian(chainID), validator.Bytes()}, []byte{})
}

func MakeDenomToERC20Key(chainID uint32, denom string) []byte {
	return bytes.Join([][]byte{{DenomToERC20Key}, Uint32ToBigEndian(chainID), []byte(denom)}, []byte{})
}

func MakeERC20ToDenomKey(chainID uint32, erc20 string) []byte {
	return bytes.Join([][]byte{{ERC20ToDenomKey}, Uint32ToBigEndian(chainID), []byte(erc20)}, []byte{})
}

func MakeERC20ToDenomKeyPrefix(chainID uint32) []byte {
	return bytes.Join([][]byte{{ERC20ToDenomKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeSignerSetTxStoreIndex(chainID uint32, nonce uint64) []byte {
	return bytes.Join([][]byte{{SignerSetTxPrefixByte}, Uint32ToBigEndian(chainID), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeBatchTxStoreIndex(chainID uint32, addr common.Address, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxPrefixByte}, Uint32ToBigEndian(chainID), addr.Bytes(), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeContractCallTxStoreIndex(chainID uint32, invalidationScope []byte, invalidationNonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxPrefixByte}, Uint32ToBigEndian(chainID), invalidationScope, sdk.Uint64ToBigEndian(invalidationNonce)}, []byte{})
}

func MakeLatestSignerSetTxNonceKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LatestSignerSetTxNonceKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastSlashedOutgoingTxBlockKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastSlashedOutgoingTxBlockKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastOutgoingBatchNonceKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastOutgoingBatchNonceKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastObservedEventNonceKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastObservedEventNonceKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastSendToEVMIDKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastSendToEVMIDKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastEVMBlockHeightKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastEVMBlockHeightKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeLastUnBondingBlockHeightKey() []byte {
	return []byte{LastUnBondingBlockHeightKey}
}

func MakeLastObservedSignerSetKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastObservedSignerSetKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

// Prefix conversions

func ValidatorToEVMAddressKeyPrefix() []byte {
	return []byte{ValidatorToEVMAddressKey}
}

func EVMToOrchestratorAddressKeyPrefix() []byte {
	return []byte{EVMToOrchestratorAddressKey}
}

func MakeEVMHeightVoteKeyPrefix(chainID uint32) []byte {
	return bytes.Join([][]byte{{EVMHeightVoteKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

func MakeEVMHeightVoteKey(chainID uint32, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{MakeEVMHeightVoteKeyPrefix(chainID), validator.Bytes()}, []byte{})
}
