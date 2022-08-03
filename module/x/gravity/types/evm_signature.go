package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

var (
	_ EVMTxConfirmation = &SignerSetTxConfirmation{}
	_ EVMTxConfirmation = &ContractCallTxConfirmation{}
	_ EVMTxConfirmation = &BatchTxConfirmation{}
)

///////////////
// GetSigner //
///////////////

func (u *SignerSetTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EVMSigner)
}

func (u *ContractCallTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EVMSigner)
}

func (u *BatchTxConfirmation) GetSigner() common.Address {
	return common.HexToAddress(u.EVMSigner)
}

///////////////////
// GetStoreIndex //
///////////////////

func (sstx *SignerSetTxConfirmation) GetStoreIndex() []byte {
	return MakeSignerSetTxStoreIndex(sstx.ChainId, sstx.SignerSetNonce)
}

func (btx *BatchTxConfirmation) GetStoreIndex() []byte {
	return MakeBatchTxStoreIndex(btx.ChainId, common.HexToAddress(btx.TokenContract), btx.BatchNonce)
}

func (cctx *ContractCallTxConfirmation) GetStoreIndex() []byte {
	return MakeContractCallTxStoreIndex(cctx.ChainId, cctx.InvalidationScope, cctx.InvalidationNonce)
}

//////////////
// Validate //
//////////////

func (u *SignerSetTxConfirmation) Validate() error {
	if u.SignerSetNonce == 0 {
		return fmt.Errorf("nonce must be set")
	}
	if !common.IsHexAddress(u.EVMSigner) {
		return sdkerrors.Wrap(ErrInvalid, "EVM signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	if u.ChainId == 0 {
		return fmt.Errorf("chain id can not be zero")
	}
	return nil
}

func (u *ContractCallTxConfirmation) Validate() error {
	if u.InvalidationNonce == 0 {
		return fmt.Errorf("invalidation nonce must be set")
	}
	if u.InvalidationScope == nil {
		return fmt.Errorf("invalidation scope must be set")
	}
	if !common.IsHexAddress(u.EVMSigner) {
		return sdkerrors.Wrap(ErrInvalid, "EVM signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	if u.ChainId == 0 {
		return fmt.Errorf("chain id can not be zero")
	}
	return nil
}

func (u *BatchTxConfirmation) Validate() error {
	if u.BatchNonce == 0 {
		return fmt.Errorf("nonce must be set")
	}
	if !common.IsHexAddress(u.TokenContract) {
		return fmt.Errorf("token contract address must be valid EVM address")
	}
	if !common.IsHexAddress(u.EVMSigner) {
		return sdkerrors.Wrap(ErrInvalid, "EVM signer must be address")
	}
	if u.Signature == nil {
		return fmt.Errorf("signature must be set")
	}
	if u.ChainId == 0 {
		return fmt.Errorf("chain id can not be zero")
	}
	return nil
}
