package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicChain(t *testing.T) {
	err := os.RemoveAll("testdata/")
	require.NoError(t, err, "unable to reset testdata directory")

	chain := Chain{
		DataDir:    "testdata",
		ID:         "testchain",
		Validators: nil,
	}

	err = chain.CreateAndInitializeValidators(4)
	require.NoError(t, err, "error initializing validators")

	err = chain.CreateAndInitializeOrchestrators(uint8(len(chain.Validators)))
	require.NoError(t, err, "error initializing orchestrators")

	// add validator accounts to genesis file
	path := chain.Validators[0].ConfigDir()
	for _, n := range chain.Validators {
		err = addGenesisAccount(path, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating validator accounts")
	}

	// add orchestrator accounts to genesis file
	for _, n := range chain.Orchestrators {
		err = addGenesisAccount(path, "", n.KeyInfo.GetAddress(), "100000000000stake,100000000000footoken")
		require.NoError(t, err, "error creating orchestrator accounts")
	}

	// copy around the genesis file with the accounts
	for _, v := range chain.Validators[1:] {
		_, err = copy(filepath.Join(path, "config", "genesis.json"), filepath.Join(v.ConfigDir(), "config", "genesis.json"))
		require.NoError(t, err, "error copying over genesis files")
	}

	// generate ethereum keys for validators,
	// add them to the ethereum genesis
	ethGenesis := EthereumGenesis{
		Difficulty: "0x400",
		GasLimit: "0xB71B00",
		Alloc: make(map[string]Allocation, len(chain.Validators)),
	}
	for _, v := range chain.Validators {
		err = v.generateEthereumKey()
		require.NoError(t, err, "error copying over genesis files")

		ethGenesis.Alloc[v.EthereumKey.Address] = Allocation{Balance: "0x1337000000000000000000"}
	}

	// write out the genesis file
	ethGenesisMarshal, err := json.MarshalIndent(ethGenesis, "", "  ")
	require.NoError(t, err, "error marshalling ethereum genesis file")

	err = ioutil.WriteFile(filepath.Join(chain.ConfigDir(), "ethGenesis.json"), ethGenesisMarshal, 0644)
	require.NoError(t, err, "error writing ethereum genesis file")
}
