package integration_tests

// package imports
import (
	"context"
	"time"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

// We have 4 validators running so this totals to 100 tx's
const transactions_per_validator int64 = 25
const sent_amt int64 = 100

func (s *IntegrationTestSuite) TestTransactionStress() {
	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		// Approve spend & verify funds
		for _, validator := range s.chain.validators {			
			err := s.SendEthTransaction(&validator.ethereumKey, testERC20contract, PackApproveERC20(gravityContract))
			s.Require().NoError(err, "error approving spend")

			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(validator.ethereumKey.address), testERC20contract)
			s.Require().NoError(err, "error getting balance")
			s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())	
		}

		sendAmt := sdk.NewInt(sent_amt)
		// Send many tx's through to cosmos
		for i := 0; i < len(s.chain.validators); i++ {
			s.T().Logf("sending %d tx's to cosmos for validator %d ..", transactions_per_validator, i+1)

			for j := 0; j < int(transactions_per_validator); j++ {
				s.Require().NoError(s.SendEthTransaction(&s.chain.validators[i].ethereumKey, gravityContract, PackSendToCosmos(testERC20contract, s.chain.validators[len(s.chain.validators)-1-i].keyInfo.GetAddress(), sendAmt)))
			}
			
			s.T().Logf("%d Tx sent.", transactions_per_validator)
		}

		var gravityDenom string
		for i := 0; i < len(s.chain.validators); i++ {
			s.Require().Eventuallyf(func() bool {
				s.T().Logf("Checking validator %d", i+1)

				validator := s.chain.validators[i]
				kb, err := validator.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", validator.keyInfo.GetAddress())
				s.Require().NoError(err)

				bankQueryClient := banktypes.NewQueryClient(clientCtx)
				res, err := bankQueryClient.AllBalances(context.Background(),
					&banktypes.QueryAllBalancesRequest{
						Address: validator.keyInfo.GetAddress().String(),
					})
				if err != nil {
					return false
				}

				gbQueryClient := types.NewQueryClient(clientCtx)
				denomRes, err := gbQueryClient.ERC20ToDenom(context.Background(),
					&types.ERC20ToDenomRequest{
						Erc20: testERC20contract.String(),
					})
				if err != nil {
					s.T().Logf("error querying ERC20 denom %s, %e", testERC20contract.String(), err)
					return false
				}
				s.Require().False(denomRes.CosmosOriginated, "ERC20-originated denom marked as cosmos originated")
				gravityDenom = denomRes.Denom

				for _, coin := range res.Balances {
					if coin.Denom == gravityDenom && coin.Amount.Equal(sdk.NewInt(sent_amt * transactions_per_validator)) {
						s.T().Logf("Expected funds recieved for validator %d, balance: %v", i + 1, coin)
						return true
					}
				}

				s.T().Logf("balance not found, received %v", res.Balances)

				return false
			}, 105*time.Second, 10*time.Second, "balance never found on cosmos")
		}
		fmt.Println("Ethereum -> Cosmos stress test completed.")

	
		fmt.Println("StressTestTransaction completed successfully")
	})
}