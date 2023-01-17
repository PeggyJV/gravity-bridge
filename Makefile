.DEFAULT_GOAL := e2e_slow_loris

e2e_build_images: e2e_clean_slate
	@docker build -t gravity:prebuilt -f module/Dockerfile module/
	@docker build -t evm:prebuilt -f integration_tests/evm/Dockerfile integration_tests/evm/
	@docker build -t orchestrator:prebuilt -f orchestrator/Dockerfile orchestrator/


e2e_slow_loris_deprecated:
	@make -s e2e_happy_path
	@make -s e2e_orchestrator_keys
	@make -s e2e_arbitrary_logic
	@make -s e2e_batch_stress
	@make -s e2e_valset_stress
	@make -s e2e_transaction_stress

e2e_clean_slate:
	@./clean_slate.sh

e2e_batch_stress_deprecated: e2e_clean_slate
	@testnet/testnet.test -test.run TestBatchStress -test.failfast -test.v || make -s fail

e2e_valset_stress_deprecated: e2e_clean_slate
	@testnet/testnet.test -test.run TestValsetStress -test.failfast -test.v || make -s fail

e2e_arbitrary_logic_deprecated: e2e_clean_slate
	@testnet/testnet.test -test.run TestArbitraryLogic -test.failfast -test.v || make -s fail

e2e_orchestrator_keys_deprecated: e2e_clean_slate
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

e2e_slow_loris:
	@make -s e2e_happy_path
	@make -s e2e_validator_out
	@make -s e2e_valset_update
	@make -s e2e_transaction_stress

e2e_happy_path: e2e_clean_slate
	E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestHappyPath || make -s fail

e2e_valset_update: e2e_clean_slate
	integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestValsetUpdate || make -s fail

e2e_validator_out: e2e_clean_slate
	integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestValidatorOut || make -s fail

e2e_transaction_stress: e2e_clean_slate
	integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestTransactionStress || make -s fail

e2e_evm_upgrade: e2e_clean_slate
	E2E_SKIP_CLEANUP=true integration_tests/integration_tests.test -test.failfast -test.v -test.run IntegrationTestSuite -testify.m TestMultiEVMUpgrade || make -s fail
