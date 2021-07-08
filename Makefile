.DEFAULT_GOAL := e2e_slow_loris

# TODO(levi) e2e_help target

e2e_slow_loris:
	@make -s e2e_happy_path
	@make -s e2e_v2_happy_path
	@make -s e2e_orchestrator_keys
	@make -s e2e_batch_stress
	@make -s e2e_arbitrary_logic
	@make -s e2e_validator_out
	@make -s e2e_valset_stress

e2e_clean_slate:
	@#TODO(levi) use a tag in these filters: (network is insufficient / only works on running containers)
	@docker rm --force \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		$(shell docker ps -qa --filter="name=test_runner") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker wait \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		$(shell docker ps -qa --filter="name=test_runner") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@sudo rm -fr testdata
	@cd testnet && go test -c

e2e_batch_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestBatchStress -test.failfast -test.v

e2e_happy_path: e2e_clean_slate
	@testnet/testnet.test -test.run TestHappyPath -test.failfast -test.v

e2e_validator_out: e2e_clean_slate
	@testnet/testnet.test -test.run TestValidatorOut -test.failfast -test.v

e2e_valset_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestValsetStress -test.failfast -test.v

e2e_v2_happy_path: e2e_clean_slate
	@testnet/testnet.test -test.run TestV2HappyPath -test.failfast -test.v

e2e_arbitrary_logic: e2e_clean_slate
	@testnet/testnet.test -test.run TestArbitraryLogic -test.failfast -test.v

e2e_orchestrator_keys: e2e_clean_slate
	@testnet/testnet.test -test.run TestOrchestratorKeys -test.failfast -test.v
