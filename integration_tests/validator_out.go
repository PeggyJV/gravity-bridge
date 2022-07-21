package integration_tests

// Validator out tests a validator that is not running the mandatory Ethereum node. This validator will be slashed and the bridge will remain functioning.

// Start the chain with validators
func (s *IntegrationTestSuite) TestValidatorOut() {
	s.Run("Bring up chain, and test the valset update", func() {
		s.dockerPool.RemoveContainerByName("orchestrator3")
	})
}

// Check that there's no slashing and validator orchestrator are working fine.

// Stop a validator's orchestrator.

// Check that the chain is still running

// Run an ethereum tx transaction

// Check that the validator has been slashed.
