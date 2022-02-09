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

	// QuerierRoute to be used for querierer msgs
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

// BigEndianToUint32 returns a uint32 from big endian encoded bytes. If encoding
// is empty, zero is returned.
func BigEndianToUint32(bz []byte) uint32 {
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint32(bz)
}

// MakeOrchestratorValidatorAddressKey returns the following key format
// prefix
// [0xe8][0001][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeOrchestratorValidatorAddressKey(orc sdk.AccAddress) []byte {
	return bytes.Join([][]byte{{OrchestratorToValidatorAddressKey}, orc.Bytes()}, []byte{})
}

// MakeValidatorEVMAddressKeyForValidator returns the following key format
// prefix              cosmos-validator
// [0x0][0001][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeValidatorEVMAddressKeyForValidator(validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{ValidatorToEVMAddressKey}, validator.Bytes()}, []byte{})
}

// MakeEVMOrchestratorAddressKey returns the following key format
// prefix              cosmos-validator
// [0x0][0001][cosmosvaloper1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEVMOrchestratorAddressKey(eth common.Address) []byte {
	return bytes.Join([][]byte{{EVMToOrchestratorAddressKey}, eth.Bytes()}, []byte{})
}

/////////////////////////
// Etheruem Signatures //
/////////////////////////

// MakeEVMSignatureKeyForValidator returns the following key format
// prefix   nonce                    validator-address
// [0x0][0 0 0 0 0 0 0 1][cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn]
func MakeEVMSignatureKeyForValidator(chainID uint32, storeIndex []byte, validator sdk.ValAddress) []byte {
	return bytes.Join([][]byte{{EVMSignatureKey}, Uint32ToBigEndian(chainID), storeIndex, validator.Bytes()}, []byte{})
}

func EVMSignatureKeyStoreIndexPrefix(chainID uint32, storeIndex []byte) []byte {
	return bytes.Join([][]byte{{EVMSignatureKey}, Uint32ToBigEndian(chainID), storeIndex}, []byte{})
}

/////////////////////////////////
// Etheruem Event Vote Records //
/////////////////////////////////

func MakeEVMEventVoteRecordKey(chainID uint32, eventNonce uint64, claimHash []byte) []byte {
	return bytes.Join([][]byte{{EVMEventVoteRecordKey}, Uint32ToBigEndian(chainID), sdk.Uint64ToBigEndian(eventNonce), claimHash}, []byte{})
}

func EVMEventVoteRecordPrefix(chainID uint32) []byte {
	return bytes.Join([][]byte{{EVMEventVoteRecordKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

//////////////////
// Outgoing Txs //
//////////////////

// MakeOutgoingTxKey returns the store index passed with a prefix
func MakeOutgoingTxKey(chainID uint32, storeIndex []byte) []byte {
	return bytes.Join([][]byte{{OutgoingTxKey}, Uint32ToBigEndian(chainID), storeIndex}, []byte{})
}

func OutgoingTxKeyPrefixWithPrefixByte(chainID uint32, prefix byte) []byte {
	return bytes.Join([][]byte{{OutgoingTxKey}, Uint32ToBigEndian(chainID), {prefix}}, []byte{})
}

func OutgoingTxKeyPrefix(chainID uint32) []byte {
	return bytes.Join([][]byte{{OutgoingTxKey}, Uint32ToBigEndian(chainID)}, []byte{})
}

//////////////////////
// Send To Etheruem //
//////////////////////

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

func MakeSignerSetTxKey(chainID uint32, nonce uint64) []byte {
	return bytes.Join([][]byte{{SignerSetTxPrefixByte}, Uint32ToBigEndian(chainID), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeBatchTxKey(chainID uint32, addr common.Address, nonce uint64) []byte {
	return bytes.Join([][]byte{{BatchTxPrefixByte}, Uint32ToBigEndian(chainID), addr.Bytes(), sdk.Uint64ToBigEndian(nonce)}, []byte{})
}

func MakeContractCallTxKey(chainID uint32, invalscope []byte, invalnonce uint64) []byte {
	return bytes.Join([][]byte{{ContractCallTxPrefixByte}, Uint32ToBigEndian(chainID), invalscope, sdk.Uint64ToBigEndian(invalnonce)}, []byte{})
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

func MakeLastUnBondingBlockHeightKey(chainID uint32) []byte {
	return bytes.Join([][]byte{{LastUnBondingBlockHeightKey}, Uint32ToBigEndian(chainID)}, []byte{})
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
