package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	dt "github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
	tmjson "github.com/tendermint/tendermint/libs/json"
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
		GasLimit:   "0xB71B00",
		Alloc:      make(map[string]Allocation, len(chain.Validators)+1),
	}
	ethGenesis.Alloc["0xBf660843528035a5A4921534E156a27e64B231fE"] = Allocation{Balance: "0x1337000000000000000000"}
	for _, v := range chain.Validators {
		err = v.generateEthereumKey()
		require.NoError(t, err, "error copying over genesis files")

		ethGenesis.Alloc[v.EthereumKey.Address] = Allocation{Balance: "0x1337000000000000000000"}
	}

	// write out the genesis file
	ethGenesisMarshal, err := json.MarshalIndent(ethGenesis, "", "  ")
	require.NoError(t, err, "error marshalling ethereum genesis file")

	err = ioutil.WriteFile(filepath.Join(chain.ConfigDir(), "ETHGenesis.json"), ethGenesisMarshal, 0644)
	require.NoError(t, err, "error writing ethereum genesis file")

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = chain.Validators[0].Moniker

	genFilePath := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	require.NoError(t, err, "error reading genesis file")

	var genUtil GenUtil
	err = json.Unmarshal(appState["genutil"], &genUtil)
	require.NoError(t, err, "error unmarshalling genesis state")

	// generate gentxs
	amount, _ := sdktypes.NewIntFromString("100000000000")
	coin := sdktypes.Coin{Denom: "stake", Amount: amount}
	genTxs := make([]json.RawMessage, len(chain.Validators))

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdktypes.Msg)(nil), &types.MsgCreateValidator{}, &gravitytypes.MsgDelegateKeys{})
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{})
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	for i, v := range chain.Validators {
		cvm, err := v.buildCreateValidatorMsg(coin)
		require.NoError(t, err, "error building create validator msg")

		dm := v.buildDelegateKeysMsg()
		require.NoError(t, err, "error building delegate keys msg")

		signedTx, err := v.signMsg(cvm, dm)
		require.NoError(t, err, "error signing create validator msg")

		txRaw, err := marshaler.MarshalJSON(signedTx)
		require.NoError(t, err, "error marshalling tx")

		genTxs[i] = txRaw
	}
	genUtil.GenTxs = genTxs

	bz, err := json.Marshal(genUtil)
	require.NoError(t, err, "error marshalling gen_util state")
	appState["genutil"] = bz

	bz, err = json.Marshal(appState)
	require.NoError(t, err, "error marshalling app state")

	genDoc.AppState = bz
	out, err := tmjson.MarshalIndent(genDoc, "", "  ")
	require.NoError(t, err, "error marshalling genesis doc")

	err = ioutil.WriteFile(genFilePath, out, fs.ModePerm)
	require.NoError(t, err, "error writing out genesis file")

	// update config.toml files
	for i, v := range chain.Validators {
		var configToml ValidatorConfig
		path := filepath.Join(v.ConfigDir(), "config", "config.toml")
		_, err = toml.DecodeFile(path, &configToml)
		require.NoError(t, err, "error decoding config toml")

		configToml.P2P.Laddr = "tcp://0.0.0.0:26656"
		configToml.P2P.AddrBookStrict = false
		configToml.P2P.ExternalAddress = fmt.Sprintf("%s:%d", v.instanceName(), 26656)
		configToml.RPC.Laddr = "tcp://0.0.0.0:26657"
		configToml.StateSync.Enable = false

		if i > 0 {
			configToml.LogLevel = "info"
		}

		var peers []string

		for j := 0; j < len(chain.Validators); j++ {
			if i == j {
				continue
			}
			peer := chain.Validators[j]
			peerID := fmt.Sprintf("%s@%s%d:26656", peer.nodeID(), peer.Moniker, j)
			peers = append(peers, peerID)
		}

		configToml.P2P.PersistentPeers = strings.Join(peers, ",")

		var b bytes.Buffer
		encoder := toml.NewEncoder(&b)
		err = encoder.Encode(configToml)
		require.NoError(t, err, "error encoding config toml")

		err = os.WriteFile(path, b.Bytes(), fs.ModePerm)
		require.NoError(t, err, "error writing config toml")

		startupPath := filepath.Join(v.ConfigDir(), "startup.sh")
		err = os.WriteFile(startupPath, []byte(fmt.Sprintf("gravity --home home start --pruning=nothing > home.n%d.log", v.Index)), fs.ModePerm)
	}

	// bring up docker network
	pool, err := dt.NewPool("")
	require.NoError(t, err, "error creating docker pool")
	network, err := pool.CreateNetwork("testnet")
	require.NoError(t, err, "error creating testnet network")
	//defer network.Close()

	hostConfig := func(config *dc.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		//config.AutoRemove = true
		//config.RestartPolicy = dc.RestartPolicy{
		//	Name: "no",
		//}
		//config.LogConfig.Type = "json-file"
	}

	// bring up ethereum
	t.Log("building and running ethereum")
	ethereum, err := pool.BuildAndRunWithBuildOptions(&dt.BuildOptions{
		Dockerfile: "ethereum/Dockerfile",
		ContextDir: "./",
	},
		&dt.RunOptions{
			Name:      "ethereum",
			NetworkID: network.Network.ID,
			PortBindings: map[dc.Port][]dc.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		}, hostConfig)
	require.NoError(t, err, "error bringing up ethereum")
	t.Logf("deployed ethereum at %s", ethereum.Container.ID)

	// build validators
	for _, validator := range chain.Validators {
		t.Logf("building %s", validator.instanceName())
		err = pool.Client.BuildImage(dc.BuildImageOptions{
			Name: validator.instanceName(),
			Dockerfile: "Dockerfile",
			ContextDir: "./module",
			OutputStream: ioutil.Discard,
		})
		require.NoError(t, err, "error building %s", validator.instanceName())
		t.Logf("built %s", validator.instanceName())
	}

	for _, validator := range chain.Validators {
		runOpts := &dt.RunOptions{
			Name:      validator.instanceName(),
			NetworkID: network.Network.ID,
			Mounts: []string{fmt.Sprintf("/Users/mvid/Development/crypto/cosmos-gravity-bridge/testdata/testchain/%s/:/root/home", validator.instanceName())},
			Repository: validator.instanceName(),
		}

		// expose the first validator for debugging and communication
		if validator.Index == 0 {
			runOpts.PortBindings = map[dc.Port][]dc.PortBinding{
				"1317/tcp":  {{HostIP: "", HostPort: "1317"}},
				"6060/tcp":  {{HostIP: "", HostPort: "6060"}},
				"6061/tcp":  {{HostIP: "", HostPort: "6061"}},
				"6062/tcp":  {{HostIP: "", HostPort: "6062"}},
				"6063/tcp":  {{HostIP: "", HostPort: "6063"}},
				"6064/tcp":  {{HostIP: "", HostPort: "6064"}},
				"6065/tcp":  {{HostIP: "", HostPort: "6065"}},
				"9090/tcp":  {{HostIP: "", HostPort: "9090"}},
				"26656/tcp": {{HostIP: "", HostPort: "26656"}},
				"26657/tcp": {{HostIP: "", HostPort: "26657"}},
			}
		}

		resource, err := pool.RunWithOptions(runOpts, hostConfig)
		require.NoError(t, err, "error bringing up %s", validator.instanceName())
		t.Logf("deployed %s at %s", validator.instanceName(), resource.Container.ID)

	}

		// bring up the contract deployer and deploy contract
	t.Log("building contract_deployer")
	contractDeployer, err := pool.BuildAndRunWithBuildOptions(&dt.BuildOptions{
		Dockerfile: "Dockerfile",
		ContextDir: "./solidity",
	},
		&dt.RunOptions{
			Name:      "contract_deployer",
			NetworkID: network.Network.ID,
			PortBindings: map[dc.Port][]dc.PortBinding{
				"8545/tcp": {{HostIP: "", HostPort: "8545"}},
			},
			Env: []string{},
		}, hostConfig)
	require.NoError(t, err, "error bringing up contract deployer")
	t.Logf("deployed contract deployer at %s", contractDeployer.Container.ID)
}