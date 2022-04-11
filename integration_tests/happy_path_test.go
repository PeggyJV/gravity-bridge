package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

func (s *IntegrationTestSuite) TestHappyPath() {
	s.Run("Bring up chain, and test the happy path", func() {
		s.T().Logf("approving Gravity to spend ERC 20")
		err := s.approveERC20()
		s.Require().NoError(err, "error approving spending balance for the gravity contract on behalf of the first validator")

		allowance, err := s.getERC20AllowanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), gravityContract)
		s.Require().NoError(err, "error getting allowance of gravity contract spending on behalf of first validator")
		s.Require().Equal(UInt256Max(), allowance.BigInt(), "spending allowance not set correctly, got: %s", allowance.String())

		balance, err := s.getEthBalanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address))
		s.Require().NoError(err, "error getting first validator balance")
		s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())

		for _, val := range s.chain.validators {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].keyInfo.GetAddress().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.keyInfo.GetAddress().String(), res.Balances)
		}

		// send from val 0 on eth to val 1 on cosmos
		s.T().Logf("sending to cosmos")
		err = s.sendToCosmos(s.chain.validators[1].keyInfo.GetAddress(), sdk.NewInt(200))
		s.Require().NoError(err, "error sending test denom to cosmos")

		for _, val := range s.chain.validators {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].keyInfo.GetAddress().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.keyInfo.GetAddress().String(), res.Balances)
		}

		var gravityDenom string
		s.Require().Eventuallyf(func() bool {
			val := s.chain.validators[0]
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

			bankQueryClient := banktypes.NewQueryClient(clientCtx)
			res, err := bankQueryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].keyInfo.GetAddress().String(),
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
				if coin.Denom == gravityDenom && coin.Amount.Equal(sdk.NewInt(200)) {
					return true
				}
			}

			s.T().Logf("balance not found, received %v", res.Balances)

			return false
		}, 105*time.Second, 10*time.Second, "balance never found on cosmos")

		s.T().Logf("sending to ethereum")

		sendToEthereumMsg := types.NewMsgSendToEthereum(
			s.chain.validators[1].keyInfo.GetAddress(),
			s.chain.validators[1].ethereumKey.address,
			sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(100)},
			sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
		)

		s.Require().Eventuallyf(func() bool {
			val := s.chain.validators[1]
			keyring, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.keyInfo.GetAddress())
			s.Require().NoError(err)

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
			return true
		}, 105*time.Second, 10*time.Second, "unable to send to ethereum")

	})
}
