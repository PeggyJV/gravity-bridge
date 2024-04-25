#!/bin/bash
set -eu

gravity init --chain-id=testing local

# Create delegate keys
CONFIG=$(echo "$(pwd)/gorc_config.toml")
gorc keys eth add signer 
gorc keys cosmos add orchestrator

SIGNER=$(gorc --config $CONFIG keys eth show signer)
ORCHESTRATOR=$(gorc --config $CONFIG keys cosmos show orchestrator)
# remove prefix from output
ORCHESTRATOR=$(echo $ORCHESTRATOR | sed s/"orchestrator "//)

gravity add-genesis-account validator 1000000000stake

# assumes key has been created
VALIDATOR=$(gravity keys show validator --bech val -a)
SIGNATURE="0x$(gorc sign-delegate-keys signer $VALIDATOR 0)"

echo $SIGNER $ORCHESTRATOR $SIGNATURE

gravity gentx validator 1000000000stake $SIGNER $ORCHESTRATOR $SIGNATURE --chain-id testing
gravity collect-gentxs
