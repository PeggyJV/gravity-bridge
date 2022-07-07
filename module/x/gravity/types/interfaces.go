package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gogo/protobuf/proto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// EVMTxConfirmation represents one validtors signature for a given
// outgoing EVM transaction
type EVMTxConfirmation interface {
	proto.Message

	GetSigner() common.Address
	GetSignature() []byte
	GetStoreIndex(chainID uint32) []byte
	Validate() error
}

// EVMEvent represents a event from the gravity contract
// on the counterparty EVM chain
type EVMEvent interface {
	proto.Message

	GetEventNonce() uint64
	GetEVMHeight() uint64
	Hash() tmbytes.HexBytes
	Validate() error
	ChainID() uint32
}

type OutgoingTx interface {
	// NOTE: currently the function signatures here don't match, figure out how to do this proprly
	// maybe add an interface arg here and typecheck in each implementation?

	// The only one that will be problematic is BatchTx, which needs to pull all the constituent
	// transactions before calculating the checkpoint

	GetCheckpoint([]byte) []byte
	GetStoreIndex() []byte
	GetCosmosHeight() uint64
}
