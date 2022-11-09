package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// GravityDenomPrefix indicates the prefix for all assests minted by this module
	GravityDenomPrefix = ModuleName

	// GravityDenomSeparator is the separator for gravity denoms
	GravityDenomSeparator = "/"
)

// EVMAddrLessThan migrates the Ethereum address less than function
func EVMAddrLessThan(e, o string) bool {
	return bytes.Compare([]byte(e)[:], []byte(o)[:]) == -1
}

// ValidateEthereumAddress validates the ethereum address strings
// func ValidateEthereumAddress(a string) error {
// 	if a == "" {
// 		return fmt.Errorf("empty")
// 	}
// 	if !regexp.MustCompile("^0x[0-9a-fA-F]{40}$").MatchString(a) {
// 		return fmt.Errorf("address(%s) doesn't pass regex", a)
// 	}
// 	if len(a) != EthereumContractAddressLen {
// 		return fmt.Errorf("address(%s) of the wrong length exp(%d) actual(%d)", a, len(a), EthereumContractAddressLen)
// 	}
// 	return nil
// }

/////////////////////////
//     ERC20Token      //
/////////////////////////

// NewERC20Token returns a new instance of an ERC20
func NewERC20Token(chainID uint32, amount uint64, contract common.Address) ERC20Token {
	return NewSDKIntERC20Token(chainID, sdk.NewIntFromUint64(amount), contract)
}

func NewSDKIntERC20Token(chainID uint32, amount sdk.Int, contract common.Address) ERC20Token {
	return ERC20Token{
		Amount:   amount,
		Contract: contract.Hex(),
		ChainId:  chainID,
	}
}

func GravityDenom(chainID uint32, contract common.Address) string {
	return strings.Join([]string{GravityDenomPrefix, strconv.Itoa(int(chainID)), contract.Hex()}, GravityDenomSeparator)
}

// GravityCoin returns the gravity representation of the ERC20
func (e ERC20Token) GravityCoin() sdk.Coin {
	return sdk.Coin{Amount: e.Amount, Denom: GravityDenom(e.ChainId, common.HexToAddress(e.Contract))}
}

func GravityDenomToERC20(denom string) (uint32, string, error) {
	if !strings.HasPrefix(denom, GravityDenomPrefix+GravityDenomSeparator) {
		return 0, "", fmt.Errorf("denom prefix(%s) not equal to expected(%s)", denom, GravityDenomPrefix+GravityDenomSeparator)
	}

	gDenomComponents := strings.Split(denom, GravityDenomSeparator)
	if len(gDenomComponents) != 3 {
		return 0, "", fmt.Errorf("unable to properly split denom into components: %s", denom)
	}

	chainIDStr := gDenomComponents[1]
	chainID, err := strconv.ParseUint(chainIDStr, 0, 32)
	if err != nil {
		return 0, "", fmt.Errorf("error parsing chain ID: %s. error: %e", chainIDStr, err)
	}

	contract := gDenomComponents[2]
	switch {
	case !common.IsHexAddress(contract):
		return 0, "", fmt.Errorf("error validating evm contract address")
	default:
		return uint32(chainID), common.HexToAddress(contract).Hex(), nil
	}
}

func NormalizeCoinDenom(coin *sdk.Coin) {
	coin.Denom = NormalizeDenom(coin.Denom)
}

func NormalizeDenom(denom string) string {
	if chainID, contract, err := GravityDenomToERC20(denom); err == nil {
		return GravityDenom(chainID, common.HexToAddress(contract))
	}

	return denom
}

func NewSendToEVMTx(chainID uint32, id uint64, tokenContract common.Address, sender sdk.AccAddress, recipient common.Address, amount, feeAmount uint64) *SendToEVM {
	return &SendToEVM{
		Id:           id,
		Erc20Fee:     NewERC20Token(chainID, feeAmount, tokenContract),
		Sender:       sender.String(),
		EVMRecipient: recipient.Hex(),
		Erc20Token:   NewERC20Token(chainID, amount, tokenContract),
		ChainId:      chainID,
	}
}
