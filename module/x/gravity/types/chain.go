package types

type ChainSpecific interface {
	ChainID() uint32
}

const EthereumChainID = 1

var (
	_ ChainSpecific = &MsgSendToEthereum{}
	_ ChainSpecific = &MsgCancelSendToEthereum{}
	_ ChainSpecific = &MsgRequestBatchTx{}
	_ ChainSpecific = &MsgSubmitEthereumEvent{}
	_ ChainSpecific = &MsgSubmitEthereumTxConfirmation{}

	_ ChainSpecific = &SignerSetTx{}
	_ ChainSpecific = &BatchTx{}
	_ ChainSpecific = &ContractCallTx{}
)

// messages

func (msg *MsgSendToEthereum) ChainID() uint32 {
	chainID := msg.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (msg *MsgCancelSendToEthereum) ChainID() uint32 {
	chainID := msg.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (msg *MsgRequestBatchTx) ChainID() uint32 {
	chainID := msg.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (msg *MsgSubmitEthereumEvent) ChainID() uint32 {
	chainID := msg.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (msg *MsgSubmitEthereumTxConfirmation) ChainID() uint32 {
	chainID := msg.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

// outgoing txs

func (otx *SignerSetTx) ChainID() uint32 {
	chainID := otx.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (otx *BatchTx) ChainID() uint32 {
	chainID := otx.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}

func (otx *ContractCallTx) ChainID() uint32 {
	chainID := otx.GetChainId()
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}
