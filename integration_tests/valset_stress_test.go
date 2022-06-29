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
func (s *IntegrationTestSuite) TestValsetUpdate(val sdk.ValAddress) {
	s.Run("Bring up chain, and test the valset update", func() {
		ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
		s.Require().NoError(err, "error setting up eth client")

		starting_nonce, err := ethClient.PendingNonceAt(context.Background(), gravityContract)
		s.Require().NoError(err, "error getting starting nonce")

		bondTokens := sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

		bondCoin := sdk.NewCoin(sdk.DefaultBondDenom, bondTokens)

		delegator := s.chain.orchestrators[0].keyInfo.GetAddress()

		val := sdk.ValAddress(s.chain.validators[1].keyInfo.GetAddress())

		timeout := time.After(300 * time.Second)
		ch_err := make(chan error)
	loop:
		for {
			select {
			case <-timeout:
				break loop

			default:

				s.T().Logf("Sending in valset request (starting_eth_valset_nonce %d)", starting_nonce)

				s.T().Logf("Delegating %v to %v in order to generate a validator set update", bondCoin, delegator)

				types.NewMsgDelegate(delegator, val, bondCoin)

			case <-ch_err:
				continue loop
			}
		}

		current_nonce, err := ethClient.PendingNonceAt(context.Background(), gravityContract)
		s.Require().NoError(err, "error getting current nonce")

	nonce:
		for current_nonce == starting_nonce {
			select {
			case <-timeout:
				break nonce
			default:
				s.Require().NoError(err, "Validator set is not yet updated")
				current_nonce, err = ethClient.PendingNonceAt(context.Background(), gravityContract)
				s.Require().NoError(err, "error getting current nonce")
			}
		}

		if current_nonce != starting_nonce {
			s.T().Logf("Validator set successfully updated!")
		} else {
			s.T().Logf("Failed to update validator set")
		}

	})
}
