package integration_tests

// package imports
import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// We can add more erc20's, however testgb is already set up and theoretically solely viable for stress testing tx's
var erc20s = [...]string{"testgb"}

func (s *IntegrationTestSuite) TestTransactionStress() {
	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		// Approve spend & check that eth test addresses have expected funds
		for _, acct := range stress_test_eth_addresses {
			err := s.SendEthTransaction(acct, testERC20contract, PackApproveERC20(gravityContract))
			s.Require().NoError(err, "error approving spend")

			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(acct.address), testERC20contract)
			s.Require().NoError(err, "error getting balance")
			s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())	
		}

		sendAmt := sdk.NewInt(200)
		// Send many tx's through to cosmos
		for _, acct := range stress_test_eth_addresses {
			s.T().Logf("sending to cosmos..")
			// Send to existing validator
			s.Require().NoError(s.SendEthTransaction(acct, gravityContract, PackSendToCosmos(testERC20contract, s.chain.validators[1].keyInfo.GetAddress(), sendAmt)))
			s.T().Logf("Tx sent.")
		}

		fmt.Println("StressTestTransaction completed successfully")
	})
}