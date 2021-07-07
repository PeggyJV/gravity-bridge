
# TODO(levi) ensure there are no stale containers running in our test network here:
e2e_setup:
	@sudo rm -fr testdata
	@cd testnet && go test -c

e2e_happy_path: e2e_setup
	@testnet/testnet.test -test.run HappyPath -test.failfast -test.v

e2e_validator_out: e2e_setup
	@testnet/testnet.test -test.run ValidatorOut -test.failfast -test.v

e2e_batch_stress: e2e_setup
	@testnet/testnet.test -test.run BatchStress -test.failfast -test.v

e2e_valset_stress: e2e_setup
	@testnet/testnet.test -test.run ValsetStress -test.failfast -test.v

e2e_v2_happy_path: e2e_setup
	@testnet/testnet.test -test.run V2HappyPath -test.failfast -test.v

e2e_arbitrary_logic: e2e_setup
	@testnet/testnet.test -test.run ArbitraryLogic -test.failfast -test.v

e2e_orchestrator_keys: e2e_setup
	@testnet/testnet.test -test.run OrchestratorKeys -test.failfast -test.v