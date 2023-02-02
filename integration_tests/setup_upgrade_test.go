package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	premultievmapp "github.com/peggyjv/gravity-bridge/module/v2/app"
	premultigravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	"github.com/stretchr/testify/suite"
	tmcfg "github.com/tendermint/tendermint/config"
	tmjson "github.com/tendermint/tendermint/libs/json"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type UpgradeTestSuite struct {
	IntegrationTestSuite
}

func (s *UpgradeTestSuite) SetupSuite() {
	s.T().Log("setting up e2e upgrade test suite...")

	var err error
	s.chain, err = newChain()
	s.Require().NoError(err)

	s.T().Logf("starting e2e infrastructure; chain-id: %s; datadir: %s", s.chain.id, s.chain.dataDir)

	// initialization
	mnemonics := MNEMONICS()
	s.Require().NoError(s.createAndInitValidatorsPreEvmMulti(len(mnemonics)))
	s.initNodesWithMnemonics(mnemonics...)

	// we only need to generate eth keys, there is no genesis for hardhat
	for i, val := range s.chain.validators {
		s.Require().NoError(val.generateEthereumKeyFromMnemonic(mnemonics[i]))
	}
	s.initPreMultiGenesis()
	s.initValidatorConfigs()

	s.dockerPool, err = dockertest.NewPool("")
	s.Require().NoError(err)

	s.dockerNetwork, err = s.dockerPool.CreateNetwork(fmt.Sprintf("%s-testnet", s.chain.id))
	s.Require().NoError(err)

	/// container infrastructure

	// evm containers don't change for the upgrade, and can run as normal
	s.runEVMContainers()

	// original validators and orchestrators have protos specific to Ethereum mainnet,
	// so the old versions must run genesis
	s.runPreMultiValidators()
	s.runPreMultiOrchestrators()
}

func (s *IntegrationTestSuite) createAndInitValidatorsPreEvmMulti(count int) error {
	for i := 0; i < count; i++ {
		// create node
		node := s.chain.createValidator(i)

		// generate genesis files
		if err := node.createConfig(); err != nil {
			return err
		}

		serverCtx := server.NewDefaultContext()
		config := serverCtx.Config

		config.SetRoot(node.configDir())
		config.Moniker = node.moniker

		genDoc, err := getGenDoc(node.configDir())
		if err != nil {
			return err
		}

		appState, err := json.MarshalIndent(premultievmapp.ModuleBasics.DefaultGenesis(cdc), "", " ")
		if err != nil {
			return fmt.Errorf("failed to JSON encode premultievmapp genesis state: %w", err)
		}

		genDoc.ChainID = node.chain.id
		genDoc.Validators = nil
		genDoc.AppState = appState

		if err = genutil.ExportGenesisFile(genDoc, config.GenesisFile()); err != nil {
			return fmt.Errorf("failed to export premultievmapp genesis state: %w", err)
		}

		tmcfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

		s.chain.validators = append(s.chain.validators, node)
	}

	return nil
}

func (s *IntegrationTestSuite) initPreMultiGenesis() {
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

	var govGenState govtypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[govtypes.ModuleName], &govGenState))

	// set short voting period to allow gov proposals in tests
	govGenState.VotingParams.VotingPeriod = time.Second * 20
	govGenState.DepositParams.MinDeposit = sdk.Coins{{Denom: testDenom, Amount: sdk.OneInt()}}
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

		delKeysMsg := val.buildPreEvmMultiDelegateKeysMsg()
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

	// create gravity genesis using only the v2 version, only for ethereum
	var gravityGenState premultigravitytypes.GenesisState
	s.Require().NoError(cdc.UnmarshalJSON(appGenState[gravitytypes.ModuleName], &gravityGenState))

	gravityGenState.Params = premultigravitytypes.DefaultParams()
	gravityGenState.Params.BridgeChainId = gravitytypes.EthereumChainID
	gravityGenState.Params.GravityId = "gravitytest"
	gravityGenState.Params.SignedBatchesWindow = 50

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

func (v *validator) buildPreEvmMultiDelegateKeysMsg() sdk.Msg {
	privKeyBz, err := hexutil.Decode(v.ethereumKey.privateKey)
	if err != nil {
		panic(fmt.Sprintf("failed to HEX decode private key: %s", err))
	}

	privKey, err := crypto.ToECDSA(privKeyBz)
	if err != nil {
		panic(fmt.Sprintf("failed to convert private key: %s", err))
	}

	signMsg := premultigravitytypes.DelegateKeysSignMsg{
		ValidatorAddress: sdk.ValAddress(v.keyInfo.GetAddress()).String(),
		Nonce:            0,
	}

	signMsgBz := cdc.MustMarshal(&signMsg)
	hash := crypto.Keccak256Hash(signMsgBz).Bytes()
	ethSig, err := gravitytypes.NewEVMSignature(hash, privKey)
	if err != nil {
		panic(fmt.Sprintf("failed to create Ethereum signature: %s", err))
	}

	return premultigravitytypes.NewMsgDelegateKeys(
		sdk.ValAddress(v.keyInfo.GetAddress()),
		v.chain.orchestrators[v.index].keyInfo.GetAddress(),
		v.ethereumKey.address,
		ethSig,
	)
}

func (s *IntegrationTestSuite) runPreMultiValidators() {
	s.T().Log("starting validator containers...")

	s.valResources = make([]*dockertest.Resource, len(s.chain.validators))
	for i, val := range s.chain.validators {
		runOpts := &dockertest.RunOptions{
			Name:       val.instanceName(),
			NetworkID:  s.dockerNetwork.Network.ID,
			Repository: "gravity",
			Tag:        "pre-multi-evm",
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

// the only different in this function is the image tag and the chain IDs tagged to only Ethereum
func (s *IntegrationTestSuite) runPreMultiOrchestrators() {
	s.T().Log("starting orchestrator containers...")

	s.orchResources = make([]*dockertest.Resource, len(s.chain.orchestrators))
	for i, orch := range s.chain.orchestrators {
		val := s.chain.validators[i]

		gorcCfgPath := path.Join(val.configDir(), "gorc")
		s.Require().NoError(os.MkdirAll(gorcCfgPath, 0755))

		for j, chainID := range ChainIds {
			gorcCfg := fmt.Sprintf(`keystore = "/root/gorc/%d/keystore/"

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

[metrics]
listen_addr = "127.0.0.1:300%d"
`,
				chainID,
				gravityContracts[j].String(),
				testDenom,
				// NOTE: container names are prefixed with '/'
				s.evms[j].Resource.Container.Name[1:],
				s.valResources[i].Container.Name[1:],
				minGasPrice,
				testDenom,
				j,
			)

			gorcChainCfgPath := path.Join(gorcCfgPath, strconv.Itoa(int(chainID)))
			s.Require().NoError(os.MkdirAll(gorcChainCfgPath, 0755))

			filePath := path.Join(gorcChainCfgPath, "config.toml")
			s.Require().NoError(writeFile(filePath, []byte(gorcCfg)))
		}

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
				Tag:        "pre-multi-evm",
				Mounts: []string{
					fmt.Sprintf("%s/:/root/gorc", gorcCfgPath),
				},
				Env: []string{
					fmt.Sprintf("CHAIN_IDS=%d", gravitytypes.EthereumChainID),
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
			20*time.Second,
			"orchestrator %s not healthy",
			resource.Container.ID,
		)
	}
}

func (s *IntegrationTestSuite) TestUpgrade() {
	// chain is halted due to pending upgrade, take down existing validator and orchestrator nodes
	for _, oc := range s.orchResources {
		s.Require().NoError(s.dockerPool.RemoveContainerByName(oc.Container.Name))
	}
	for _, vc := range s.valResources {
		s.Require().NoError(s.dockerPool.RemoveContainerByName(vc.Container.Name))
	}

	// this will bring up the "normal" current validators and orchestrators
	s.runValidators()
	s.runOrchestrators()

	rpcClient, err := rpchttp.New("tcp://localhost:26657", "/websocket")
	s.Require().NoError(err)

	// check for creation of blocks
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

func (s *UpgradeTestSuite) TearDownSuite() {
	s.IntegrationTestSuite.TearDownSuite()
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}
