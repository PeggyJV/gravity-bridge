package types

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	_ OutgoingTx = &SignerSetTx{}
	_ OutgoingTx = &BatchTx{}
	_ OutgoingTx = &ContractCallTx{}
)

const (
	_ = iota
	SignerSetTxPrefixByte
	BatchTxPrefixByte
	ContractCallTxPrefixByte
)

///////////////////
// GetStoreIndex //
///////////////////

// TODO: do we need a prefix byte for the different types?
func (sstx *SignerSetTx) GetStoreIndex() []byte {
	return MakeSignerSetTxKey(sstx.Nonce)
}

func (btx *BatchTx) GetStoreIndex() []byte {
	return MakeBatchTxKey(gethcommon.HexToAddress(btx.TokenContract), btx.BatchNonce)
}

func (cctx *ContractCallTx) GetStoreIndex() []byte {
	return MakeContractCallTxKey(cctx.InvalidationScope, cctx.InvalidationNonce)
}

///////////////////
// GetCheckpoint //
///////////////////

func (sstx *SignerSetTx) GetCosmosHeight() uint64 {
	return sstx.Height
}

func (btx *BatchTx) GetCosmosHeight() uint64 {
	return btx.Height
}

func (cctx *ContractCallTx) GetCosmosHeight() uint64 {
	return cctx.Height
}

///////////////////
// GetCheckpoint //
///////////////////

// GetCheckpoint returns the checkpoint
func (u SignerSetTx) GetCheckpoint(gravityID []byte) []byte {
	return []byte{}
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (b BatchTx) GetCheckpoint(gravityID []byte) []byte {
	return []byte{}
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (c ContractCallTx) GetCheckpoint(gravityID []byte) []byte {
	return []byte{}
}
