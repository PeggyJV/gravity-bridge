package integration_tests

// package imports
import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

// We have 4 validators running so this totals to 100 tx's
const transactions_per_validator int64 = 25
const cosmos_sent_amt int64 = 100
const eth_sent_amt int64 = 1

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

		sendAmt := sdk.NewInt(cosmos_sent_amt)
		// Send many tx's through to cosmos
		for i, validator := range s.chain.validators {
			s.T().Logf("sending %d tx's to cosmos for validator %d ..", transactions_per_validator, i+1)

			for j := 0; j < int(transactions_per_validator); j++ {
				s.Require().NoError(s.SendEthTransaction(&validator.ethereumKey, gravityContract, PackSendToCosmos(testERC20contract, s.chain.validators[len(s.chain.validators)-1-i].address(), sendAmt)))
			}

			s.T().Logf("%d Tx sent.", transactions_per_validator)
		}

		var gravityDenom string
		for i, validator := range s.chain.validators {
			s.Require().Eventuallyf(func() bool {
				s.T().Logf("Checking validator %d", i+1)

				kb, err := validator.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", validator.address())
				s.Require().NoError(err)

				bankQueryClient := banktypes.NewQueryClient(clientCtx)
				res, err := bankQueryClient.AllBalances(context.Background(),
					&banktypes.QueryAllBalancesRequest{
						Address: validator.address().String(),
					})
				if err != nil {
					s.T().Logf("error: %s", err)
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
					if coin.Denom == gravityDenom && coin.Amount.Equal(sdk.NewInt(cosmos_sent_amt*transactions_per_validator)) {
						s.T().Logf("Expected funds recieved for validator %d, balance: %v", i+1, coin)
						return true
					}
				}

				s.T().Logf("balance not found, received %v", res.Balances)

				return false
			}, 105*time.Second, 10*time.Second, "balance never found on cosmos")
		}
		fmt.Println("Ethereum -> Cosmos stress test completed.")

		for i, validator := range s.chain.validators {
			s.Require().Eventuallyf(func() bool {
				s.T().Logf("sending %d tx's to ethereum for validator %d ..", transactions_per_validator, i+1)

				sendToEthereumMsg := types.NewMsgSendToEthereum(
					validator.address(),
					s.chain.validators[len(s.chain.validators)-1-i].ethereumKey.address,
					sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(eth_sent_amt)},
					sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
				)

				keyring, err := validator.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", validator.address())
				s.Require().NoError(err)

				for j := 0; j < int(transactions_per_validator); j++ {
					response, err := s.chain.sendMsgs(*clientCtx, sendToEthereumMsg)
					if err != nil {
						s.T().Logf("error: %s", err)
						return false
					}
					if response.Code != 0 {
						if response.Code != 32 {
							s.T().Log(response)
						}
						return false
					}
				}

				s.T().Logf("%d Tx sent.", transactions_per_validator)
				return true
			}, 105*time.Second, 10*time.Second, "unable to send to ethereum")
		}

		for i, validator := range s.chain.validators {
			s.Require().Eventuallyf(func() bool {
				s.T().Logf("Checking validator %d", i+1)

				balance, err := s.getEthTokenBalanceOf(common.HexToAddress(validator.ethereumKey.address), testERC20contract)
				s.Require().NoError(err, "error getting destination balance")

				if balance.LT(sdk.NewInt(10000 - (cosmos_sent_amt * transactions_per_validator) + (eth_sent_amt * transactions_per_validator))) {
					s.T().Logf("funds not received yet, dest balance: %s", balance.String())
					return false
				}

				s.T().Logf("Funds recieved for validator %d, current balance: %v", i+1, balance.String())
				return true
			}, time.Second*180, time.Second*10, "balance never found")
		}
		fmt.Println("Cosmos -> Ethereum stress test completed.")

		fmt.Println("StressTestTransaction completed successfully")
	})
}
