.DEFAULT_GOAL := e2e_slow_loris

e2e_build_images:
	@docker build -t gravity:prebuilt -f module/Dockerfile module/
	@docker build -t ethereum:prebuilt -f integration_tests/ethereum/Dockerfile integration_tests/ethereum/
	@docker build -t orchestrator:prebuilt -f orchestrator/Dockerfile orchestrator/


e2e_slow_loris:
	@make -s e2e_happy_path
	@make -s e2e_orchestrator_keys
	@make -s e2e_arbitrary_logic
	@make -s e2e_validator_out
	@make -s e2e_batch_stress
	@make -s e2e_valset_stress

e2e_clean_slate:
	@docker rm --force \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker wait \
		$(shell docker ps -qa --filter="name=contract_deployer") \
		$(shell docker ps -qa --filter="name=ethereum") \
		$(shell docker ps -qa --filter="name=gravity") \
		$(shell docker ps -qa --filter="name=orchestrator") \
		1>/dev/null \
		2>/dev/null \
		|| true
	@docker network prune --force 1>/dev/null 2>/dev/null || true
	@cd integration_tests && go test -c

e2e_batch_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestBatchStress -test.failfast -test.v || make -s fail

e2e_validator_out: e2e_clean_slate
	@testnet/testnet.test -test.run TestValidatorOut -test.failfast -test.v || make -s fail

e2e_valset_stress: e2e_clean_slate
	@testnet/testnet.test -test.run TestValsetStress -test.failfast -test.v || make -s fail

e2e_arbitrary_logic: e2e_clean_slate
	@testnet/testnet.test -test.run TestArbitraryLogic -test.failfast -test.v || make -s fail

e2e_orchestrator_keys: e2e_clean_slate
	@testnet/testnet.test -test.run TestOrchestratorKeys -test.failfast -test.v || make -s fail

fail:
	@echo 'test failed; dumping container logs into ./testlogs for review'
	@docker logs contract_deployer > testlogs/contract_deployer.log 2>&1 || true
	@docker logs gravity0 > testlogs/gravity0.log 2>&1 || true
	@docker logs gravity1 > testlogs/gravity1.log 2>&1 || true
	@docker logs gravity2 > testlogs/gravity2.log 2>&1 || true
	@docker logs gravity3 > testlogs/gravity3.log 2>&1 || true
	@docker logs orchestrator0 > testlogs/orchestrator0.log 2>&1 || true
	@docker logs orchestrator1 > testlogs/orchestrator1.log 2>&1 || true
	@docker logs orchestrator2 > testlogs/orchestrator2.log 2>&1 || true
	@docker logs orchestrator3 > testlogs/orchestrator3.log 2>&1 || true
	@docker logs ethereum > testlogs/ethereum.log 2>&1 || true
	@false

e2e_happy_path: e2e_clean_slate
	@E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestHappyPath || make -s fail
