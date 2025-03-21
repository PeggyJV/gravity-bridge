package integration_tests

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	sdkmath "cosmossdk.io/math"
	tmconfig "github.com/cometbft/cometbft/config"
	tmjson "github.com/cometbft/cometbft/libs/json"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/server"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	testDenom           = "testgb"
	initBalanceStr      = "1000000000000testgb"
	minGasPrice         = "2"
	ethChainID     uint = 15
)

func MNEMONICS() []string {
	return []string{
		"say monitor orient heart super local purse cricket caution primary bring insane road expect rather help two extend own execute throw nation plunge subject",
		"march carpet enact kiss tribe plastic wash enter index lift topic riot try juice replace supreme original shift hover adapt mutual holiday manual nut",
		"assault section bleak gadget venture ship oblige pave fabric more initial april dutch scene parade shallow educate gesture lunar match patch hawk member problem",
		"receive roof marine sure lady hundred sea enact exist place bean wagon kingdom betray science photo loop funny bargain floor suspect only strike endless",
	}
}

var (
	stakeAmount, _         = sdk.NewIntFromString("100000000000")
	stakeAmountCoin        = sdk.NewCoin(testDenom, stakeAmount)
	gravityContract        = common.HexToAddress("0x04C89607413713Ec9775E14b954286519d836FEf")
	testERC20contract      = common.HexToAddress("0x4C4a2f8c81640e47606d3fd77B353E87Ba015584")
	maliciousERC20contract = common.HexToAddress("0x4C4a2f8c81640e47606d3fd77B353E87Ba015584")
)

type IntegrationTestSuite struct {
	suite.Suite

	chain         *chain
	dockerPool    *dockertest.Pool
	dockerNetwork *dockertest.Network
	ethResource   *dockertest.Resource
	valResources  []*dockertest.Resource
	orchResources []*dockertest.Resource
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up e2e integration test suite...")

	var err error
	s.chain, err = newChain()
	s.Require().NoError(err)

	s.T().Logf("starting e2e infrastructure; chain-id: %s; datadir: %s", s.chain.id, s.chain.dataDir)

	// initialization
	mnemonics := MNEMONICS()
	s.initNodesWithMnemonics(mnemonics...)
	s.initEthereumFromMnemonics(mnemonics)
	s.initGenesis()
	s.initValidatorConfigs()

	s.dockerPool, err = dockertest.NewPool("")
	s.Require().NoError(err)

	s.dockerNetwork, err = s.dockerPool.CreateNetwork(fmt.Sprintf("%s-testnet", s.chain.id))
	s.Require().NoError(err)

	// container infrastructure
	s.runEthContainer()
	s.runValidators()
	s.runOrchestrators()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	if str := os.Getenv("E2E_SKIP_CLEANUP"); len(str) > 0 {
		skipCleanup, err := strconv.ParseBool(str)
		s.Require().NoError(err)

		if skipCleanup {
			s.T().Log("skipping teardown")
			return
		}
	}

	s.T().Log("tearing down e2e integration test suite...")

	s.Require().NoError(s.dockerPool.Purge(s.ethResource))

	for _, vc := range s.valResources {
		s.Require().NoError(s.dockerPool.RemoveContainerByName(vc.Container.Name))
	}

	for _, oc := range s.orchResources {
		s.Require().NoError(s.dockerPool.RemoveContainerByName(oc.Container.Name))
	}

	s.Require().NoError(s.dockerPool.RemoveNetwork(s.dockerNetwork))
	s.Require().NoError(os.RemoveAll(s.chain.dataDir))
}

func (s *IntegrationTestSuite) initNodes(nodeCount int) {
	s.Require().NoError(s.chain.createAndInitValidators(nodeCount))
	s.Require().NoError(s.chain.createAndInitOrchestrators(nodeCount))

	// initialize a genesis file for the first validator
	val0ConfigDir := s.chain.validators[0].configDir()
	for _, val := range s.chain.validators {
		s.Require().NoError(
			addGenesisAccount(val0ConfigDir, "", initBalanceStr, val.address()),
		)
	}

	// add orchestrator accounts to genesis file
	for _, orch := range s.chain.orchestrators {
		s.Require().NoError(
			addGenesisAccount(val0ConfigDir, "", initBalanceStr, orch.address()),
		)
	}

	// copy the genesis file to the remaining validators
	for _, val := range s.chain.validators[1:] {
		err := copyFile(
			filepath.Join(val0ConfigDir, "config", "genesis.json"),
			filepath.Join(val.configDir(), "config", "genesis.json"),
		)
		s.Require().NoError(err)
	}
}

func (s *IntegrationTestSuite) initNodesWithMnemonics(mnemonics ...string) {
	s.Require().NoError(s.chain.createAndInitValidatorsWithMnemonics(mnemonics))
	s.Require().NoError(s.chain.createAndInitOrchestratorsWithMnemonics(mnemonics))

	//initialize a genesis file for the first validator
	val0ConfigDir := s.chain.validators[0].configDir()
	for _, val := range s.chain.validators {
		s.Require().NoError(
			addGenesisAccount(val0ConfigDir, "", initBalanceStr, val.address()),
		)
	}

	// add orchestrator accounts to genesis file
	for _, orch := range s.chain.orchestrators {
		s.Require().NoError(
			addGenesisAccount(val0ConfigDir, "", initBalanceStr, orch.address()),
		)
	}

	// copy the genesis file to the remaining validators
	for _, val := range s.chain.validators[1:] {
		err := copyFile(
			filepath.Join(val0ConfigDir, "config", "genesis.json"),
			filepath.Join(val.configDir(), "config", "genesis.json"),
		)
		s.Require().NoError(err)
	}
}

func (s *IntegrationTestSuite) initEthereum() {
	// generate ethereum keys for validators add them to the ethereum genesis
	ethGenesis := EthereumGenesis{
		Difficulty: "0x400",
		GasLimit:   "0xB71B00",
		Config:     EthereumConfig{ChainID: ethChainID},
		Alloc:      make(map[string]Allocation, len(s.chain.validators)+1),
	}

	alloc := Allocation{
		Balance: "0x1337000000000000000000",
	}

	ethGenesis.Alloc["0xBf660843528035a5A4921534E156a27e64B231fE"] = alloc
	for _, val := range s.chain.validators {
		s.Require().NoError(val.generateEthereumKey())
		ethGenesis.Alloc[val.ethereumKey.address] = alloc
	}

	ethGenBz, err := json.MarshalIndent(ethGenesis, "", "  ")
	s.Require().NoError(err)

	// write out the genesis file
	s.Require().NoError(writeFile(filepath.Join(s.chain.configDir(), "eth_genesis.json"), ethGenBz))
}

func (s *IntegrationTestSuite) initEthereumFromMnemonics(mnemonics []string) {
	// generate ethereum keys for validators add them to the ethereum genesis
	ethGenesis := EthereumGenesis{
		Difficulty: "0x400",
		GasLimit:   "0xB71B00",
		Config:     EthereumConfig{ChainID: ethChainID},
		Alloc:      make(map[string]Allocation, len(s.chain.validators)+1),
	}

	alloc := Allocation{
		Balance: "0x1337000000000000000000",
	}

	ethGenesis.Alloc["0xBf660843528035a5A4921534E156a27e64B231fE"] = alloc
	for i, val := range s.chain.validators {
		s.Require().NoError(val.generateEthereumKeyFromMnemonic(mnemonics[i]))
		ethGenesis.Alloc[val.ethereumKey.address] = alloc
	}

	ethGenBz, err := json.MarshalIndent(ethGenesis, "", "  ")
	s.Require().NoError(err)

	// write out the genesis file
	s.Require().NoError(writeFile(filepath.Join(s.chain.configDir(), "eth_genesis.json"), ethGenBz))
}

func (s *IntegrationTestSuite) initGenesis() {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config

	config.SetRoot(s.chain.validators[0].configDir())
	config.Moniker = s.chain.validators[0].moniker

	genFilePath := config.GenesisFile()
	appGenState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFilePath)
	s.Require().NoError(err)

	var bankGenState banktypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[banktypes.ModuleName], &bankGenState))

	bankGenState.DenomMetadata = append(bankGenState.DenomMetadata, banktypes.Metadata{
		Description: "The native staking token of the test gravity bridge network",
		Display:     testDenom,
		Base:        testDenom,
		Name:        testDenom,
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    testDenom,
				Exponent: 0,
				Aliases: []string{
					"tgb",
				},
			},
		},
	})

	bz, err := cdc.MarshalJSON(&bankGenState)
	s.Require().NoError(err)
	appGenState[banktypes.ModuleName] = bz

	var govGenState govtypesv1.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState))

	// set short voting period to allow gov proposals in tests
	*govGenState.Params.VotingPeriod = time.Second * 20
	govGenState.Params.MinDeposit = sdk.Coins{{Denom: testDenom, Amount: sdk.OneInt()}}
	bz, err = cdc.MarshalJSON(&govGenState)
	s.Require().NoError(err)
	appGenState[govtypes.ModuleName] = bz

	// set crisis denom
	var crisisGenState crisistypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[crisistypes.ModuleName], &crisisGenState))
	crisisGenState.ConstantFee.Denom = testDenom
	bz, err = cdc.MarshalJSON(&crisisGenState)
	s.Require().NoError(err)
	appGenState[crisistypes.ModuleName] = bz

	// set staking bond denom
	var stakingGenState stakingtypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[stakingtypes.ModuleName], &stakingGenState))
	stakingGenState.Params.BondDenom = testDenom
	bz, err = cdc.MarshalJSON(&stakingGenState)
	s.Require().NoError(err)
	appGenState[stakingtypes.ModuleName] = bz

	// set mint denom
	var mintGenState minttypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[minttypes.ModuleName], &mintGenState))
	mintGenState.Params.MintDenom = testDenom
	bz, err = cdc.MarshalJSON(&mintGenState)
	s.Require().NoError(err)
	appGenState[minttypes.ModuleName] = bz

	var genUtilGenState genutiltypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[genutiltypes.ModuleName], &genUtilGenState))

	// generate genesis txs
	genTxs := make([]json.RawMessage, len(s.chain.validators))
	for i, val := range s.chain.validators {
		createValmsg, err := val.buildCreateValidatorMsg(stakeAmountCoin)
		s.Require().NoError(err)

		delKeysMsg, err := val.buildDelegateKeysMsg()
		s.Require().NoError(err)

		signedTx, err := val.signMsg(createValmsg, delKeysMsg)
		s.Require().NoError(err)

		txRaw, err := cdc.MarshalJSON(signedTx)
		s.Require().NoError(err)

		genTxs[i] = txRaw
	}

	genUtilGenState.GenTxs = genTxs

	bz, err = cdc.MarshalJSON(&genUtilGenState)
	s.Require().NoError(err)
	appGenState[genutiltypes.ModuleName] = bz

	// set contract addr
	var gravityGenState gravitytypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[gravitytypes.ModuleName], &gravityGenState))
	gravityGenState.Params.GravityId = "gravitytest"
	gravityGenState.Params.BridgeEthereumAddress = gravityContract.String()
	gravityGenState.Params.ConfirmedOutgoingTxWindow = 350
	gravityGenState.Params.TargetEthTxTimeout = 3600000
	gravityGenState.Params.AverageBlockTime = 1000
	gravityGenState.Params.AverageEthereumBlockTime = 1000

	bz, err = cdc.MarshalJSON(&gravityGenState)
	s.Require().NoError(err)
	appGenState[gravitytypes.ModuleName] = bz

	// serialize genesis state
	bz, err = json.MarshalIndent(appGenState, "", "  ")
	s.Require().NoError(err)

	genDoc.AppState = bz

	bz, err = tmjson.MarshalIndent(genDoc, "", "  ")
	s.Require().NoError(err)

	// write the updated genesis file to each validator
	for _, val := range s.chain.validators {
		s.Require().NoError(writeFile(filepath.Join(val.configDir(), "config", "genesis.json"), bz))
	}
}

func (s *IntegrationTestSuite) initValidatorConfigs() {
	for i, val := range s.chain.validators {
		tmCfgPath := filepath.Join(val.configDir(), "config", "config.toml")

		vpr := viper.New()
		vpr.SetConfigFile(tmCfgPath)
		s.Require().NoError(vpr.ReadInConfig())

		valConfig := &tmconfig.Config{}
		s.Require().NoError(vpr.Unmarshal(valConfig))

		valConfig.P2P.ListenAddress = "tcp://0.0.0.0:26656"
		valConfig.P2P.AddrBookStrict = false
		valConfig.P2P.ExternalAddress = fmt.Sprintf("%s:%d", val.instanceName(), 26656)
		valConfig.RPC.ListenAddress = "tcp://0.0.0.0:26657"
		valConfig.StateSync.Enable = false
		valConfig.LogLevel = "info"

		// speed up blocks
		valConfig.Consensus.TimeoutCommit = 1 * time.Second
		valConfig.Consensus.TimeoutPropose = 1 * time.Second

		var peers []string

		for j := 0; j < len(s.chain.validators); j++ {
			if i == j {
				continue
			}

			peer := s.chain.validators[j]
			peerID := fmt.Sprintf("%s@%s%d:26656", peer.nodeKey.ID(), peer.moniker, j)
			peers = append(peers, peerID)
		}

		valConfig.P2P.PersistentPeers = strings.Join(peers, ",")

		tmconfig.WriteConfigFile(tmCfgPath, valConfig)

		// set application configuration
		appCfgPath := filepath.Join(val.configDir(), "config", "app.toml")

		appConfig := srvconfig.DefaultConfig()
		appConfig.GRPC.Address = "0.0.0.0:9090"
		appConfig.API.Enable = true
		appConfig.Pruning = "nothing"
		appConfig.MinGasPrices = fmt.Sprintf("%s%s", minGasPrice, testDenom)

		srvconfig.WriteConfigFile(appCfgPath, appConfig)
	}
}

func (s *IntegrationTestSuite) runHardhatContainer() {
	s.T().Log("starting Ethereum Hardhat container...")

}

func (s *IntegrationTestSuite) runEthContainer() {
	s.T().Log("starting Ethereum container...")
	var err error

	nodeURL := os.Getenv("ARCHIVE_NODE_URL")
	s.Require().NotEmptyf(nodeURL, "ARCHIVE_NODE_URL env variable must be set")

	runOpts := dockertest.RunOptions{
		Name:       "ethereum",
		Repository: "ethereum",
		Tag:        "prebuilt",
		NetworkID:  s.dockerNetwork.Network.ID,
		PortBindings: map[docker.Port][]docker.PortBinding{
			"8545/tcp": {{HostIP: "", HostPort: "8545"}},
		},
		ExposedPorts: []string{"8545/tcp"},
		Env:          []string{fmt.Sprintf("ARCHIVE_NODE_URL=%s", nodeURL)},
	}

	s.ethResource, err = s.dockerPool.RunWithOptions(
		&runOpts,
		noRestart,
	)
	s.Require().NoError(err)

	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	s.Require().NoError(err)

	// Wait for the Ethereum node to respond to a request
	s.Require().Eventually(
		func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()

			balance, err := ethClient.BalanceAt(ctx, common.HexToAddress(s.chain.validators[0].ethereumKey.address), nil)
			if err != nil {
				s.T().Logf("error querying balance: %e", err)
				return false
			}

			if balance == nil {
				s.T().Logf("balance for first validator is nil")
			}

			if balance.Cmp(big.NewInt(0)) == 0 {
				s.T().Logf("balance for first validator is %s", balance.String())
				return false
			}

			return true
		},
		5*time.Minute,
		10*time.Second,
		"ethereum node failed to respond",
	)

	s.T().Logf("waiting for contract to deploy")
	ethereumLogOutput := bytes.Buffer{}
	err = s.dockerPool.Client.Logs(docker.LogsOptions{
		Container:    s.ethResource.Container.ID,
		OutputStream: &ethereumLogOutput,
		Stdout:       true,
	})
	s.Require().NoError(err, "error getting contract deployer logs")

	s.Require().Eventuallyf(func() bool {
		for _, s := range strings.Split(ethereumLogOutput.String(), "\n") {
			if strings.HasPrefix(s, "gravity contract deployed at") {
				strSpl := strings.Split(s, "-")
				gravityContract = common.HexToAddress(strings.ReplaceAll(strSpl[1], " ", ""))
				return true
			}
		}
		return false
	}, time.Minute*5, time.Second*10, "unable to retrieve gravity address from logs")
	s.T().Logf("gravity contract deployed at %s", gravityContract.String())

	s.Require().Eventuallyf(func() bool {
		for _, s := range strings.Split(ethereumLogOutput.String(), "\n") {
			if strings.HasPrefix(s, "test ERC20 TestGB TGB deployed at") {
				strSpl := strings.Split(s, "-")
				testERC20contract = common.HexToAddress(strings.ReplaceAll(strSpl[1], " ", ""))
				return true
			}
		}
		return false
	}, time.Minute*5, time.Second*10, "unable to retrieve test erc20 address from logs")
	s.T().Logf("test erc20 contract deployed at %s", testERC20contract.String())

	s.Require().Eventuallyf(func() bool {
		for _, s := range strings.Split(ethereumLogOutput.String(), "\n") {
			if strings.HasPrefix(s, "MaliciousERC20 deployed at") {
				strSpl := strings.Split(s, "-")
				maliciousERC20contract = common.HexToAddress(strings.ReplaceAll(strSpl[1], " ", ""))
				return true
			}
		}
		return false
	}, time.Minute*5, time.Second*10, "unable to retrieve test erc20 address from logs")
	s.T().Logf("malicious erc20 contract deployed at %s", maliciousERC20contract.String())

	s.T().Logf("started Ethereum container: %s", s.ethResource.Container.ID)
}

func (s *IntegrationTestSuite) runValidators() {
	s.T().Log("starting validator containers...")

	s.valResources = make([]*dockertest.Resource, len(s.chain.validators))
	for i, val := range s.chain.validators {
		runOpts := &dockertest.RunOptions{
			Name:       val.instanceName(),
			NetworkID:  s.dockerNetwork.Network.ID,
			Repository: "gravity",
			Tag:        "prebuilt",
			Mounts: []string{
				fmt.Sprintf("%s/:/root/.gravity", val.configDir()),
			},
			Entrypoint: []string{"gravity", "start", "--trace=true"},
		}

		// expose the first validator for debugging and communication
		if val.index == 0 {
			runOpts.PortBindings = map[docker.Port][]docker.PortBinding{
				"1317/tcp":  {{HostIP: "", HostPort: "1317"}},
				"9090/tcp":  {{HostIP: "", HostPort: "9090"}},
				"26656/tcp": {{HostIP: "", HostPort: "26656"}},
				"26657/tcp": {{HostIP: "", HostPort: "26657"}},
			}
			runOpts.ExposedPorts = []string{"1317/tcp", "9090/tcp", "26656/tcp", "26657/tcp"}
		}

		resource, err := s.dockerPool.RunWithOptions(runOpts, noRestart)
		s.Require().NoError(err)

		s.valResources[i] = resource
		s.T().Logf("started validator container: %s", resource.Container.ID)
	}

	rpcClient, err := rpchttp.New("tcp://localhost:26657", "/websocket")
	s.Require().NoError(err)

	s.Require().Eventually(
		func() bool {
			status, err := rpcClient.Status(context.Background())
			if err != nil {
				s.T().Logf("can't get container status: %s", err.Error())
			}
			if status == nil {
				container, ok := s.dockerPool.ContainerByName("gravity0")
				if !ok {
					s.T().Logf("no container by 'gravity0'")
				} else {
					if container.Container.State.Status == "exited" {
						s.Fail("validators exited", "state: %s logs: \n%s", container.Container.State.String(), s.logsByContainerID(container.Container.ID))
						s.T().FailNow()
					}
					s.T().Logf("state: %v, health: %v", container.Container.State.Status, container.Container.State.Health)
				}
				return false
			}

			// let the node produce a few blocks
			if status.SyncInfo.CatchingUp {
				s.T().Logf("catching up: %t", status.SyncInfo.CatchingUp)
				return false
			}
			if status.SyncInfo.LatestBlockHeight < 2 {
				s.T().Logf("block height %d", status.SyncInfo.LatestBlockHeight)
				return false
			}

			return true
		},
		10*time.Minute,
		15*time.Second,
		"validator node failed to produce blocks",
	)
}

func (s *IntegrationTestSuite) runOrchestrators() {
	s.T().Log("starting orchestrator containers...")

	s.orchResources = make([]*dockertest.Resource, len(s.chain.orchestrators))
	for i, orch := range s.chain.orchestrators {
		gorcCfg := fmt.Sprintf(`keystore = "/root/gorc/keystore/"

[gravity]
contract = "%s"
fees_denom = "%s"

[ethereum]
key_derivation_path = "m/44'/60'/0'/0/0"
rpc = "http://%s:8545"

[cosmos]
key_derivation_path = "m/44'/118'/1'/0/0"
grpc = "http://%s:9090"
gas_price = { amount = %s, denom = "%s" }
prefix = "cosmos"
gas_adjustment = 2.0
msg_batch_size = 5
`,
			gravityContract.String(),
			testDenom,
			// NOTE: container names are prefixed with '/'
			s.ethResource.Container.Name[1:],
			s.valResources[i].Container.Name[1:],
			minGasPrice,
			testDenom,
		)

		val := s.chain.validators[i]

		gorcCfgPath := path.Join(val.configDir(), "gorc")
		s.Require().NoError(os.MkdirAll(gorcCfgPath, 0755))

		filePath := path.Join(gorcCfgPath, "config.toml")
		s.Require().NoError(writeFile(filePath, []byte(gorcCfg)))

		// We must first populate the orchestrator's keystore prior to starting
		// the orchestrator gorc process. The keystore must contain the Ethereum
		// and orchestrator keys. These keys will be used for relaying txs to
		// and from the test network and Ethereum. The gorc_bootstrap.sh scripts encapsulates
		// this entire process.
		//
		// NOTE: If the Docker build changes, the script might have to be modified
		// as it relies on busybox.
		err := copyFile(
			filepath.Join("integration_tests", "gorc_bootstrap.sh"),
			filepath.Join(gorcCfgPath, "gorc_bootstrap.sh"),
		)
		s.Require().NoError(err)

		resource, err := s.dockerPool.RunWithOptions(
			&dockertest.RunOptions{
				Name:       orch.instanceName(),
				NetworkID:  s.dockerNetwork.Network.ID,
				Repository: "orchestrator",
				Tag:        "prebuilt",
				Mounts: []string{
					fmt.Sprintf("%s/:/root/gorc", gorcCfgPath),
				},
				Env: []string{
					fmt.Sprintf("ORCH_MNEMONIC=%s", orch.mnemonic),
					fmt.Sprintf("ETH_PRIV_KEY=%s", val.ethereumKey.privateKey),
					"RUST_BACKTRACE=full",
					"RUST_LOG=debug",
				},
				Entrypoint: []string{
					"sh",
					"-c",
					"chmod +x /root/gorc/gorc_bootstrap.sh && /root/gorc/gorc_bootstrap.sh",
				},
			},
			noRestart,
		)
		s.Require().NoError(err)

		s.orchResources[i] = resource
		s.T().Logf("started orchestrator container: %s", resource.Container.ID)
	}

	// TODO(mvid) Determine if there is a way to check the health or status of
	// the gorc orchestrator processes. For now, we search the logs to determine
	// when each orchestrator resource has synced all batches
	match := "No unsigned batches! Everything good!"
	for _, resource := range s.orchResources {
		resource := resource
		s.T().Logf("waiting for orchestrator to be healthy: %s", resource.Container.ID)

		s.Require().Eventuallyf(
			func() bool {
				var containerLogsBuf bytes.Buffer
				s.Require().NoError(s.dockerPool.Client.Logs(
					docker.LogsOptions{
						Container:    resource.Container.ID,
						OutputStream: &containerLogsBuf,
						Stdout:       true,
						Stderr:       true,
					},
				))

				return strings.Contains(containerLogsBuf.String(), match)
			},
			3*time.Minute,
			1*time.Second,
			"orchestrator %s not healthy",
			resource.Container.ID,
		)
	}
}

func noRestart(config *docker.HostConfig) {
	// in this case we don't want the nodes to restart on failure
	config.RestartPolicy = docker.RestartPolicy{
		Name: "no",
	}
}

func (s *IntegrationTestSuite) logsByContainerID(id string) string {
	var containerLogsBuf bytes.Buffer
	s.Require().NoError(s.dockerPool.Client.Logs(
		docker.LogsOptions{
			Container:    id,
			OutputStream: &containerLogsBuf,
			Stdout:       true,
			Stderr:       true,
		},
	))

	return containerLogsBuf.String()
}

func (s *IntegrationTestSuite) TestBasicChain() {
	// this test verifies that the setup functions all operate as expected
	s.Run("bring up basic chain", func() {
	})
}

func (s *IntegrationTestSuite) deployERC20(denom string, name string, symbol string, decimals uint8) error {
	return s.SendEthTransaction(&s.chain.validators[0].ethereumKey, gravityContract, PackDeployERC20(denom, name, symbol, decimals))
}

func (s *IntegrationTestSuite) approveERC20() error {
	return s.SendEthTransaction(&s.chain.validators[0].ethereumKey, testERC20contract, PackApproveERC20(gravityContract))
}

func (s *IntegrationTestSuite) sendToCosmos(destination sdk.AccAddress, amount sdkmath.Int) error {
	return s.SendEthTransaction(&s.chain.validators[0].ethereumKey, gravityContract, PackSendToCosmos(testERC20contract, destination, amount))
}

func (s *IntegrationTestSuite) getEthTokenBalanceOf(account common.Address, erc20contract common.Address) (*sdkmath.Int, error) {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return nil, err
	}

	data := PackBalanceOf(account)

	response, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress(s.chain.validators[0].ethereumKey.address),
		To:   &erc20contract,
		Gas:  0,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	balance := UnpackEthUInt(response)

	return &balance, err
}

func (s *IntegrationTestSuite) getERC20AllowanceOf(owner common.Address, spender common.Address) (*sdkmath.Int, error) {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return nil, err
	}

	data := PackAllowance(owner, spender)

	response, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress(s.chain.validators[0].ethereumKey.address),
		To:   &testERC20contract,
		Gas:  0,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	allowance := UnpackEthUInt(response)

	return &allowance, err
}

func (s *IntegrationTestSuite) SendEthTransaction(ethereumKey *ethereumKey, toAddress common.Address, data []byte) error {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return err
	}

	privateKey, err := crypto.HexToECDSA(ethereumKey.privateKey[2:])
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	value := big.NewInt(0)
	gasLimit := uint64(1000000)
	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	err = ethClient.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}

	return nil
}

func (s *IntegrationTestSuite) getLastValsetNonce(erc20contract common.Address) (*sdkmath.Int, error) {
	ethClient, err := ethclient.Dial(fmt.Sprintf("http://%s", s.ethResource.GetHostPort("8545/tcp")))
	if err != nil {
		return nil, err
	}

	data := PackLastValsetNonce()

	response, err := ethClient.CallContract(context.Background(), ethereum.CallMsg{
		To:   &erc20contract,
		Gas:  0,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	nonce := UnpackEthUInt(response)

	return &nonce, err
}
