package integration_tests

// package imports
import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// We have 4 validators running so this totals to 100 tx's
const transactionsPerValidator int64 = 25
const cosmosSentAmt int64 = 100
const ethSentAmt int64 = 1

func (s *IntegrationTestSuite) TestTransactionStress() {
	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		ethereum := *s.evmResources[0]
		gravityContract := gravityContracts[0]
		testERC20contract := testERC20contracts[0]

		// Approve spend & verify funds
		for _, validator := range s.chain.validators {
			err := sendEthTransaction(ethereum, &validator.ethereumKey, testERC20contract, PackApproveERC20(gravityContract))
			s.Require().NoError(err, "error approving spend")

			balance, err := s.getEthTokenBalanceOf(ethereum, common.HexToAddress(validator.ethereumKey.address), testERC20contract)
			s.Require().NoError(err, "error getting balance")
			s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())
		}

		sendAmt := sdk.NewInt(cosmosSentAmt)
		// Send many tx's through to cosmos
		for i, validator := range s.chain.validators {
			s.T().Logf("sending %d tx's to cosmos for validator %d ..", transactionsPerValidator, i+1)

			for j := 0; j < int(transactionsPerValidator); j++ {
				s.Require().NoError(sendEthTransaction(ethereum, &validator.ethereumKey, gravityContract, PackSendToCosmos(testERC20contract, s.chain.validators[len(s.chain.validators)-1-i].keyInfo.GetAddress(), sendAmt)))
			}

			s.T().Logf("%d Tx sent.", transactionsPerValidator)
		}

		var gravityDenom string
		queryingVal := s.chain.validators[2]
		s.Require().Eventuallyf(func() bool {
			kb, err := queryingVal.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", queryingVal.keyInfo.GetAddress())
			s.Require().NoError(err)

			gbQueryClient := types.NewQueryClient(clientCtx)
			denomRes, err := gbQueryClient.ERC20ToDenom(context.Background(),
				&types.ERC20ToDenomRequest{
					Erc20:   testERC20contract.String(),
					ChainId: types.EthereumChainID,
				})
			s.Require().NoError(err)
			s.Require().False(denomRes.CosmosOriginated, "ERC20-originated denom marked as cosmos originated")
			gravityDenom = denomRes.Denom

			bankQueryClient := banktypes.NewQueryClient(clientCtx)
			balanceMap := make(map[int]sdk.Coins)

			for i, validator := range s.chain.validators {
				res, err := bankQueryClient.AllBalances(context.Background(),
					&banktypes.QueryAllBalancesRequest{
						Address: validator.keyInfo.GetAddress().String(),
					})
				s.Require().NoError(err)
				balanceMap[i] = res.Balances

				balance, err := s.getEthTokenBalanceOf(ethereum, common.HexToAddress(validator.ethereumKey.address), testERC20contract)
				s.Require().NoError(err)
				s.T().Logf("test erc20 balance for validator %d: %v", i, balance)
			}

			met := true
			for i, _ := range s.chain.validators {
				for _, coin := range balanceMap[i] {
					if coin.Denom == gravityDenom {
						if coin.Amount.Equal(sdk.NewInt(cosmosSentAmt * transactionsPerValidator)) {
							s.T().Logf("correct funds recieved for validator %d, balance: %v", i, coin)
						} else {
							s.T().Logf("incorrect funds received for validator %d, got %d, expected %d", i, coin.Amount.Int64(),
								cosmosSentAmt*transactionsPerValidator)
							met = false
						}
					}
				}
			}

			return met
		}, 2*time.Minute, 10*time.Second, "balance never found on cosmos")
		fmt.Println("Ethereum -> Cosmos stress test completed.")

		for i, validator := range s.chain.validators {
			s.T().Logf("sending %d tx's to ethereum for validator %d ..", transactionsPerValidator, i+1)

			sendToEthereumMsg := types.NewMsgSendToEVM(
				types.EthereumChainID,
				validator.keyInfo.GetAddress(),
				s.chain.validators[len(s.chain.validators)-1-i].ethereumKey.address,
				sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(ethSentAmt)},
				sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
			)

			keyring, err := validator.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", validator.keyInfo.GetAddress())
			s.Require().NoError(err)

			for j := 0; j < int(transactionsPerValidator); j++ {
				response, err := s.chain.sendMsgs(*clientCtx, sendToEthereumMsg)
				s.Require().NoError(err)
				if response.Code != 0 {
					s.T().Logf("response non-zero: %s", response.RawLog)
					s.Require().Zero(response.Code)
				}
			}

			s.T().Logf("%d Tx sent.", transactionsPerValidator)
		}

		for i, validator := range s.chain.validators {
			s.Require().Eventuallyf(func() bool {
				s.T().Logf("Checking validator %d", i+1)

				balance, err := s.getEthTokenBalanceOf(ethereum, common.HexToAddress(validator.ethereumKey.address), testERC20contract)
				s.Require().NoError(err, "error getting destination balance")

				if balance.LT(sdk.NewInt(10000 - (cosmosSentAmt * transactionsPerValidator) + (ethSentAmt * transactionsPerValidator))) {
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
