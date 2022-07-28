package integration_tests

// package imports
import (
	"context"
	"fmt"
	"crypto/ecdsa"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
)

// We can add more erc20's, however testgb is already set up and theoretically solely viable for stress testing tx's
var erc20s = [...]string{"testgb"}

func (s *IntegrationTestSuite) TestTransactionStress() {
	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		// Check that eth test addresses have expected funds
		for _, acct := range stress_test_eth_addresses {
			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(acct.address), testERC20contract)
			s.Require().NoError(err, "error getting balance")
			s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())	
		}

		// TODO: need to approve spend first?

		// Send many tx's through to cosmos
		for _, acct := range stress_test_eth_addresses {
			s.T().Logf("sending to cosmos..")

			ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
			s.Require().NoError(err)
		
			privateKey, err := crypto.HexToECDSA(acct.privateKey[2:])
			s.Require().NoError(err)
		
			publicKey := privateKey.Public()
			publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
		
			fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
			nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
			s.Require().NoError(err)

			value := big.NewInt(0)
			gasLimit := uint64(1000000)
			gasPrice, err := ethClient.SuggestGasPrice(context.Background())
			s.Require().NoError(err)

			// Send to existing validator
			tx := types.NewTransaction(nonce, gravityContract, value, gasLimit, gasPrice, PackSendToCosmos(testERC20contract, s.chain.validators[1].keyInfo.GetAddress(), sdk.NewInt(200)))
		
			chainID, err := ethClient.NetworkID(context.Background())
			s.Require().NoError(err)

			signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
			s.Require().NoError(err)

			err = ethClient.SendTransaction(context.Background(), signedTx)
			s.Require().NoError(err)
		}

		fmt.Println("StressTestTransaction completed successfully")
	})
}