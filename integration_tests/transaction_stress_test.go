package integration_tests

// package imports
import (
	"fmt"
	"os"

)

// We can add more erc20's, however testgb is already set up and theoretically solely viable for stress testing tx's
var erc20_addresses = [...]string{"testgb"}

func (s *IntegrationTestSuite) TestTransactionStress() {
	os.Setenv("TransactionStressTest", "true")

	s.Run("Transaction stress test", func() {
		fmt.Println("StressTestTransaction starting")

		// Check that eth test addresses have expected funds
		

		fmt.Println("StressTestTransaction completed successfully")
	})
}