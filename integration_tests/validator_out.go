package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	gravity "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

// Validator out tests a validator that is not running the mandatory Ethereum node. This validator will be slashed and the bridge will remain functioning.

// Start the chain with validators
func (s *IntegrationTestSuite) TestValidatorOut() {
	s.Run("Bring up chain, and test the valset update", func() {
		// remove fourth validator's orchestrator
		s.dockerPool.RemoveContainerByName("orchestrator3")

		validator := s.chain.validators[3]

		valAddr := validator.keyInfo.GetAddress()

		firstValidator := sdk.ValAddress(valAddr).String()

		ethereumAddr := s.chain.validators[3].ethereumKey.address

		var gravityDenom string

		orchKey := s.chain.orchestrators[1]
		keyring := orchKey.keyring

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", keyring, "orch", orchKey.keyInfo.GetAddress())
		s.Require().NoError(err)

		// Check jail status of validators
		s.Require().Eventuallyf(func() bool {
			newQ := types.NewQueryClient(clientCtx)
			res, err := newQ.Validator(context.Background(), &types.QueryValidatorRequest{ValidatorAddr: firstValidator})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("validator response: %s", res.GetValidator())
			return true
		}, 20*time.Second, 1*time.Second, "can't find slashing info")

		// Query ERC20 Denom
		s.Require().Eventuallyf(func() bool {
			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(validator.ethereumKey.address), testERC20contract)
			s.Require().NoError(err, "error getting balance")
			s.T().Logf("balance: %s", balance)

			gbQueryClient := gravity.NewQueryClient(clientCtx)
			denomRes, err := gbQueryClient.ERC20ToDenom(context.Background(),
				&gravity.ERC20ToDenomRequest{
					Erc20: testERC20contract.String(),
				})
			if err != nil {
				s.T().Logf("error querying ERC20 denom %s, %e", testERC20contract.String(), err)
				return false
			}
			s.Require().False(denomRes.CosmosOriginated, "ERC20-originated denom marked as cosmos originated")
			gravityDenom = denomRes.Denom

			return true

		}, 5*time.Minute, 1*time.Second, "can't query ERC2O Denom")

		// Send Cosmos -> Ethereum transaction
		s.Require().Eventuallyf(func() bool {
			s.T().Logf("sending transaction between Cosmos and Ethereum account")
			Transaction := gravity.NewMsgSendToEthereum(
				valAddr,
				ethereumAddr,
				sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
				sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
			)

			keyRing, err := validator.keyring()
			s.Require().NoError(err)

			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyRing, "val", valAddr)
			s.Require().NoError(err)

			response, err := s.chain.sendMsgs(*clientCtx, Transaction)
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

		}, 5*time.Minute, 1*time.Second, "can't send funds from Cosmos to Eth account successfully")

		// Create Transaction batch
		s.Require().Eventuallyf(func() bool {
			batchTx := gravity.NewMsgRequestBatchTx(gravityDenom, valAddr)

			keyRing, err := validator.keyring()
			s.Require().NoError(err)

			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyRing, "val", valAddr)
			s.Require().NoError(err)

			response, err := s.chain.sendMsgs(*clientCtx, batchTx)
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
		}, 5*time.Minute, 1*time.Second, "can't create TX batch successfully")

		// Check jail status of validators
		s.Require().Eventuallyf(func() bool {
			newQ := types.NewQueryClient(clientCtx)
			res, err := newQ.Validator(context.Background(), &types.QueryValidatorRequest{ValidatorAddr: firstValidator})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("validator response: %s", res.GetValidator())
			return true
		}, 20*time.Second, 1*time.Second, "can't find slashing info")

	})
}
