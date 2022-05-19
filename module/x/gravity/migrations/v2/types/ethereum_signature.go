package types

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ EthereumTxConfirmation = &SignerSetTxConfirmation{}
	_ EthereumTxConfirmation = &ContractCallTxConfirmation{}
	_ EthereumTxConfirmation = &BatchTxConfirmation{}
)

///////////////
// GetSigner //
///////////////

func (u *SignerSetTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *ContractCallTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

func (u *BatchTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EthereumSigner)
}

///////////////////
// GetStoreIndex //
///////////////////

func (sstx *SignerSetTxConfirmation) GetStoreIndex() []byte {
	return MakeSignerSetTxKey(sstx.SignerSetNonce)
}

func (btx *BatchTxConfirmation) GetStoreIndex() []byte {
	return MakeBatchTxKey(common.HexToAddress(btx.TokenContract), btx.BatchNonce)
}

func (cctx *ContractCallTxConfirmation) GetStoreIndex() []byte {
	return MakeContractCallTxKey(cctx.InvalidationScope, cctx.InvalidationNonce)
}

//////////////
// Validate //
//////////////

func (u *SignerSetTxConfirmation) Validate() error {
	return nil
}

func (u *ContractCallTxConfirmation) Validate() error {
	return nil

}

func (u *BatchTxConfirmation) Validate() error {
	return nil

}
