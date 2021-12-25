package integration_tests

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"strings"
)

type EthereumConfig struct {
	ChainID             uint `json:"chainId"`
	HomesteadBlock      uint `json:"homesteadBlock"`
	EIP150Block         uint `json:"eip150Block"`
	EIP155Block         uint `json:"eip155Block"`
	EIP158Block         uint `json:"eip158Block"`
	ByzantiumBlock      uint `json:"byzantiumBlock"`
	ConstantinopleBlock uint `json:"constantinopleBlock"`
	PetersburgBlock     uint `json:"petersburgBlock"`
	IstanbulBlock       uint `json:"istanbulBlock"`
	BerlinBlock         uint `json:"berlinBlock"`
}

type Allocation struct {
	Balance string `json:"balance"`
}

type EthereumGenesis struct {
	Difficulty string                `json:"difficulty"`
	GasLimit   string                `json:"gasLimit"`
	Config     EthereumConfig        `json:"config"`
	Alloc      map[string]Allocation `json:"alloc"`
}

const approveERC20ABIJSON = `
[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]
`

func packCall(abiString, method string, args []interface{}) []byte {
	encodedCall, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		panic(sdkerrors.Wrap(err, "bad ABI definition in code"))
	}
	abiEncodedCall, err := encodedCall.Pack(method, args...)
	if err != nil {
		panic(sdkerrors.Wrap(err, "error packing calling"))
	}
	return abiEncodedCall
}

func PackDeployERC20(denom string, name string, symbol string, decimals uint8) []byte {
	return packCall(gravitytypes.DeployERC20ABIJSON, "deployERC20", []interface{}{
		denom,
		name,
		symbol,
		decimals,
	})
}

func PackSendToCosmos(tokenContract common.Address, destination sdk.AccAddress, amount sdk.Int) []byte {
	destinationBytes, _ := byteArrayToFixByteArray(destination.Bytes())
	return packCall(gravitytypes.SendToCosmosABIJSON, "sendToCosmos", []interface{}{
		tokenContract,
		destinationBytes,
		amount.BigInt(),
	})
}

func PackApproveERC20(spender common.Address) []byte {
	u256max := sdk.RelativePow(sdk.NewUint(2), sdk.NewUint(255), sdk.OneUint())
	return packCall(approveERC20ABIJSON, "approve", []interface{}{
		spender,
		u256max.BigInt(),
	})
}

func byteArrayToFixByteArray(b []byte) (out [32]byte, err error) {
	if len(b) > 32 {
		return out, fmt.Errorf("array too long")
	}
	copy(out[:], b)
	return out, nil
}
