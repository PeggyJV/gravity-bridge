package types

type ChainSpecific interface {
	ChainIDOrDefault() uint32
}

const EthereumChainID = 1

func ChainIDOrDefault(chainID uint32) uint32 {
	if chainID == 0 {
		return EthereumChainID
	}
	return chainID
}
