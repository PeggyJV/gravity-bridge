package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

// Validator out tests a validator that is not running the mandatory Ethereum node. This validator will be slashed and the bridge will remain functioning.

// Start the chain with validators
func (s *IntegrationTestSuite) TestValidatorOut() {
	s.Run("Bring up chain, and test the valset update", func() {
		s.dockerPool.RemoveContainerByName("orchestrator3")

		s.T().Logf("approving Gravity to spend ERC 20")
		err := s.approveERC20()
		s.Require().NoError(err, "error approving spending balance for the gravity contract on behalf of the first validator")

		allowance, err := s.getERC20AllowanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), gravityContract)
		s.Require().NoError(err, "error getting allowance of gravity contract spending on behalf of first validator")
		s.Require().Equal(UInt256Max(), allowance.BigInt(), "spending allowance not set correctly, got: %s", allowance.String())

		balance, err := s.getEthTokenBalanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), testERC20contract)
		s.Require().NoError(err, "error getting first validator balance")
		s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())

		// send from val 0 on eth to val 1 on cosmos
		s.T().Logf("sending to cosmos")
		err = s.sendToCosmos(s.chain.validators[1].keyInfo.GetAddress(), sdk.NewInt(200))
		s.Require().NoError(err, "error sending test denom to cosmos")

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
		}, 20*time.Minute, 10*time.Second, "unable to send to ethereum")

		// Create Transaction batch
		s.Require().Eventuallyf(func() bool {
			batchTx := types.NewMsgRequestBatchTx(gravityDenom, s.chain.validators[2].keyInfo.GetAddress())

			keyRing, err := s.chain.validators[2].keyring()
			s.Require().NoError(err)

			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyRing, "val", s.chain.validators[2].keyInfo.GetAddress())
			s.Require().NoError(err)

			response, err := s.chain.sendMsgs(*clientCtx, batchTx)
			s.T().Logf("batch response: %s", response)
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

			s.Require().NoError(err, "error querying delegator bonded validators")

			return true
		}, 30*time.Minute, 1*time.Second, "can't create TX batch successfully")

		orchKey := s.chain.orchestrators[1]
		keyring := orchKey.keyring

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", keyring, "orch", orchKey.keyInfo.GetAddress())
		s.Require().NoError(err)

		startingNonce, err := s.getLastValsetNonce(gravityContract)
		s.Require().NoError(err, "error getting starting nonce")

		bondCoin := sdk.NewCoin("testgb", sdk.NewIntFromUint64(9))

		delegator := s.chain.orchestrators[1].keyInfo.GetAddress()
		val := sdk.ValAddress(s.chain.validators[3].keyInfo.GetAddress())

		// Delegate about 5% of the total staking power.
		s.Require().Eventuallyf(func() bool {
			s.T().Logf("Sending in valset request (starting_eth_valset_nonce %d)", startingNonce)

			s.T().Logf("Delegating %v to %v in order to generate a validator set update", bondCoin, delegator)

			delegate := stakingtypes.NewMsgDelegate(delegator, val, bondCoin)
			response, err := s.chain.sendMsgs(*clientCtx, delegate)
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
		}, 5*time.Minute, 10*time.Second, "Delegate to validator failed will retry")

		keyRing, err := s.chain.validators[3].keyring()
		s.Require().NoError(err)

		clientCtx, err = s.chain.clientContext("tcp://localhost:26657", &keyRing, "val", s.chain.validators[3].keyInfo.GetAddress())
		s.Require().NoError(err)
		queryClient := types.NewQueryClient(clientCtx)

		respon, err := queryClient.UnsignedBatchTxs(context.Background(), &types.UnsignedBatchTxsRequest{Address: s.chain.orchestrators[0].keyInfo.GetAddress().String()})
		s.T().Logf("Unsigned batch tx for val 1: %s", respon)

		errorred, err := queryClient.UnsignedBatchTxs(context.Background(), &types.UnsignedBatchTxsRequest{Address: s.chain.orchestrators[1].keyInfo.GetAddress().String()})
		s.T().Logf("Unsigned batch tx for val 2: %s", errorred)

		errorre, err := queryClient.UnsignedBatchTxs(context.Background(), &types.UnsignedBatchTxsRequest{Address: s.chain.orchestrators[2].keyInfo.GetAddress().String()})
		s.T().Logf("Unsigned batch tx for val 3: %s", errorre)

		errorr, err := queryClient.UnsignedBatchTxs(context.Background(), &types.UnsignedBatchTxsRequest{Address: s.chain.orchestrators[3].keyInfo.GetAddress().String()})
		s.T().Logf("Unsigned batch tx for val 4: %s", errorr)

		// Check jail status of validators
		s.Require().Eventuallyf(func() bool {
			orchKey := s.chain.validators[3]
			keyring, err := orchKey.keyring()
			s.Require().NoError(err)

			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", s.chain.validators[3].keyInfo.GetAddress())
			s.Require().NoError(err)
			newQ := stakingtypes.NewQueryClient(clientCtx)
			res, err := newQ.Validator(context.Background(), &stakingtypes.QueryValidatorRequest{ValidatorAddr: sdk.ValAddress(s.chain.validators[3].keyInfo.GetAddress()).String()})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.Require().False(res.GetValidator().IsJailed())
			return true
		}, 20*time.Minute, 1*time.Second, "can't find slashing info")

	})
}
