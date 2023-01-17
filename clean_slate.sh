#!/bin/bash

set -e

echo Cleaning up test environment
docker rm --force \
    $(docker ps -qa --filter="name=contract_deployer") \
    $(docker ps -qa --filter="name=evm") \
    $(docker ps -qa --filter="name=gravity") \
    $(docker ps -qa --filter="name=orchestrator") \
    1>/dev/null \
    2>/dev/null \
    || true
echo Waiting for container removal to complete
docker wait \
    $(docker ps -qa --filter="name=contract_deployer") \
    $(docker ps -qa --filter="name=evm") \
    $(docker ps -qa --filter="name=gravity") \
    $(docker ps -qa --filter="name=orchestrator") \
    1>/dev/null \
    2>/dev/null \
    || true
echo Removing network
docker network prune --force 1>/dev/null 2>/dev/null || true
echo Recompiling integration tests
cd integration_tests && go test -c
cd -
