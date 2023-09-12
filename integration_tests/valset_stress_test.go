package integration_tests

// package imports
import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	gravity "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

// call stress_test for a range of nonce values

// / Write test_valset_update test to get latest nonce value
func (s *IntegrationTestSuite) TestValsetStressUpdate() {
	s.Run("Bring up chain, and test the valset update", func() {
		orchKey := s.chain.orchestrators[1]
		keyring := orchKey.keyring

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", keyring, "orch", orchKey.address())
		s.Require().NoError(err)

		startingNonce, err := s.getLastValsetNonce(gravityContract)
		s.Require().NoError(err, "error getting starting nonce")

		bondTokens := sdk.TokensFromConsensusPower(50000, sdk.DefaultPowerReduction)
		bondCoin := sdk.NewCoin("testgb", bondTokens)

		delegator := s.chain.orchestrators[1].address()
		val := sdk.ValAddress(s.chain.validators[3].address())

		// Delegate about 5% of the total staking power.
		s.Require().Eventuallyf(func() bool {
			s.T().Logf("Sending in valset request (starting_eth_valset_nonce %d)", startingNonce)

			s.T().Logf("Delegating %v to %v in order to generate a validator set update", bondCoin, delegator)

			delegate := types.NewMsgDelegate(delegator, val, bondCoin)
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

		// Verify that delegation went through.
		s.T().Logf("verifying delegation")
		s.Require().Eventuallyf(func() bool {

			s.Require().NoError(err, "error querying delegator bonded validators")
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Delegation(context.Background(), &types.QueryDelegationRequest{DelegatorAddr: delegator.String(), ValidatorAddr: val.String()})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("Here's the delegation response: %s", res.DelegationResponse)
			return true
		}, 20*time.Second, 1*time.Second, "Delegation wasn't successful")

		// Query signer set, to make sure valset was updated.
		s.T().Logf("verifying signerset")
		s.Require().Eventuallyf(func() bool {

			s.Require().NoError(err, "error querying signerset")
			queryClient := gravity.NewQueryClient(clientCtx)
			res, err := queryClient.LatestSignerSetTx(context.Background(), &gravity.LatestSignerSetTxRequest{})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("Here's the last signerset response: %s", res.SignerSet)
			return true
		}, 20*time.Second, 1*time.Second, "Signerset can't be retrieved")

		// Grab current nonce.
		currentNonce, err := s.getLastValsetNonce(gravityContract)
		s.Require().NoError(err, "error getting current nonce")

		// Run a loop until you get a nonce higher than the initial nonce
		s.Require().Eventuallyf(func() bool {
			for currentNonce == startingNonce {
				currentNonce, err = s.getLastValsetNonce(gravityContract)
				if currentNonce != startingNonce {
					return true
				}
			}
			return true
		}, 5*time.Minute, 10*time.Second, "Validator set is not yet updated")

		if currentNonce != startingNonce {
			s.T().Logf("Validator set successfully updated! nonce: %d", currentNonce)
		} else {
			s.T().Logf("Failed to update validator set")
		}

	})
}

// loop and test valset update from a range of nonces
func (s *IntegrationTestSuite) TestValsetUpdate() {
	s.Run("Bring up chain, and test the valset update", func() {
		for nonce := 0; nonce <= 10; nonce++ {
			s.TestValsetStressUpdate()
		}
	})

}
