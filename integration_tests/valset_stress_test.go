package integration_tests

// package imports
import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// call stress_test for a range of nonce values

/// Write test_valset_update test to get latest nonce value
func (s *IntegrationTestSuite) TestValsetUpdate() {
	s.Run("Bring up chain, and test the valset update", func() {
		ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
		s.Require().NoError(err, "error setting up eth client")

		validator := s.chain.orchestrators[1]
		keyring := validator.keyring

		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", keyring, "orch", validator.keyInfo.GetAddress())
		s.Require().NoError(err)

		startingNonce, err := ethClient.NonceAt(context.Background(), gravityContract, nil)
		s.Require().NoError(err, "error getting starting nonce")

		bondTokens := sdk.TokensFromConsensusPower(500000, sdk.DefaultPowerReduction)

		bondCoin := sdk.NewCoin("testgb", bondTokens)

		delegator := s.chain.orchestrators[1].keyInfo.GetAddress()

		val := sdk.ValAddress(s.chain.validators[1].keyInfo.GetAddress())

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
		}, 300*time.Second, 10*time.Second, "Delegate to validator failed will retry")

		s.T().Logf("verifying delegation")
		s.Require().Eventuallyf(func() bool {

			s.Require().NoError(err, "error querying delegator bonded validators")
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Delegation(context.Background(), &types.QueryDelegationRequest{DelegatorAddr: delegator.String(), ValidatorAddr: val.String()})
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			s.T().Logf("Here's the delegation response: %s", res.DelegationResponse.Delegation)
			return true
		}, 200*time.Second, 1*time.Second, "Delegation wasn't successful")

		currentNonce, err := ethClient.NonceAt(context.Background(), gravityContract, nil)
		s.Require().NoError(err, "error getting current nonce")

		s.Require().Eventuallyf(func() bool {
			for currentNonce == startingNonce {
				currentNonce, err = ethClient.NonceAt(context.Background(), gravityContract, nil)
				if currentNonce != startingNonce {
					return true
				}
			}
			return true
		}, 300*time.Second, 10*time.Second, "Validator set is not yet updated")

		currentNonce, err = ethClient.NonceAt(context.Background(), gravityContract, nil)

		if currentNonce != startingNonce {
			s.T().Log(currentNonce)
			s.T().Logf("Validator set successfully updated!")
		} else {
			s.T().Logf("Failed to update validator set")
		}

	})
}
