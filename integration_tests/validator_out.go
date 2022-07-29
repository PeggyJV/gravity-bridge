package integration_tests

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	gravity "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

// Validator out tests a validator that is not running the mandatory Ethereum node. This validator will be slashed and the bridge will remain functioning.

// Start the chain with validators
func (s *IntegrationTestSuite) TestValidatorOut() {
	s.Run("Bring up chain, and test the valset update", func() {
		// remove fourth validator's orchestrator
		s.dockerPool.RemoveContainerByName("orchestrator3")
		val := s.chain.validators[1]

		firstValidator := sdk.ValAddress(s.chain.validators[3].keyInfo.GetAddress()).String()

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

		bondTokens := sdk.TokensFromConsensusPower(900000, sdk.DefaultPowerReduction)
		bondCoin := sdk.NewCoin("testgb", bondTokens)

		delegator := s.chain.orchestrators[1].keyInfo.GetAddress()

		// Delegate about 5% of the total staking power.
		s.Require().Eventuallyf(func() bool {
			s.T().Logf("Delegating %v to %v in order to generate a validator set update", bondCoin, delegator)

			delegate := types.NewMsgDelegate(delegator, sdk.ValAddress(val.keyInfo.GetAddress()), bondCoin)
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

		// Query or send batchtx.
		s.Require().Eventuallyf(func() bool {
			queryClient := gravity.BatchTx{
				BatchNonce: 1,
				Transactions: []*gravity.SendToEthereum{
					gravity.NewSendToEthereumTx(2, myTokenContractAddr, mySender, myReceiver, 101, 3),
					gravity.NewSendToEthereumTx(3, myTokenContractAddr, mySender, myReceiver, 102, 2),
				},
				TokenContract: myTokenContractAddr.Hex(),
				Height:        1234567,
			}

			mm := gravity.NewMsgRequestBatchTx()
			res, err := queryClient.BatchTx(context.Background(), &gravity.BatchTxRequest{})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("validator response: %s", res)
			return true
		}, 20*time.Second, 1*time.Second, "can't batch tx")

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
