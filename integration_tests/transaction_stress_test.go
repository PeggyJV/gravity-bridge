package integration_tests

// package imports
import (
	"fmt"
	"time"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
)

const NUM_USERS int = 100
// We can add more erc20's, however testgb is already set up and theoretically solely viable for stress testing tx's
var erc20_addresses = [...]string{"testgb"}

func (s *IntegrationTestSuite) TestTransactionStress() {
	os.Setenv("TransactionStressTest", "true")

	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		// Create 100 eth addresses to send funds from 
		eth_addresses := make([]*ethereumKey, NUM_USERS)

		for i := 0; i < 100; i++ {
			mnemonic, err := createMnemonic()
			s.Require().NoError(err)

			eth_addresses[i], err = ethereumKeyFromMnemonic(mnemonic)
			s.Require().NoError(err)

			// Send ERC20(s) to each address
			for _, denom := range erc20_addresses {
				s.T().Logf("sending to ethereum")
				sendToEthereumMsg := types.NewMsgSendToEthereum(
					s.chain.validators[1].keyInfo.GetAddress(),
					s.chain.validators[1].ethereumKey.address,
					sdk.Coin{Denom: denom, Amount: sdk.NewInt(1000)},
					sdk.Coin{Denom: denom, Amount: sdk.NewInt(1)},
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
			}
		}

		fmt.Println("StressTestTransaction completed successfully")
	})
}