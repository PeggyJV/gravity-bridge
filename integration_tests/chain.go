package integration_tests

import (
	"fmt"
	"os"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/mint"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	"github.com/peggyjv/gravity-bridge/module/v4/app"
	"github.com/peggyjv/gravity-bridge/module/v4/app/params"
	gravityclient "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/client"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = app.MakeEncodingConfig()

	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
		&stakingtypes.MsgBeginRedelegate{},
		&gravitytypes.MsgDelegateKeys{},
	)
	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*cryptotypes.PubKey)(nil),
		&secp256k1.PubKey{},
		&ed25519.PubKey{},
	)

	cdc = encodingConfig.Marshaler
}

type chain struct {
	dataDir       string
	id            string
	validators    []*validator
	orchestrators []*orchestrator
}

func newChain() (*chain, error) {
	var dir string
	var err error
	if _, found := os.LookupEnv("CI"); found {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	tmpDir, err := os.MkdirTemp(dir, "gravity-bridge-e2e-testnet")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.NewRand().Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error {
	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitValidatorsWithMnemonics(mnemonics []string) error {
	for i := 0; i < len(mnemonics); i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i], ""); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitOrchestrators(count int) error {
	mnemonics := make([]string, count)
	for i := 0; i < count; i++ {
		mnemonic, err := createMnemonic()
		if err != nil {
			return err
		}
		mnemonics = append(mnemonics, mnemonic)
	}

	return c.createAndInitOrchestratorsWithMnemonics(mnemonics)
}

func (c *chain) createAndInitOrchestratorsWithMnemonics(mnemonics []string) error {
	hdPath := hd.CreateHDPath(sdk.CoinType, 1, 0)

	for i := 0; i < len(mnemonics); i++ {
		// create orchestrator
		orchestrator := c.createOrchestrator(i)

		// create keys
		info, kb, err := createMemoryKeyFromMnemonic("orch", mnemonics[i], "", hdPath)
		if err != nil {
			return err
		}

		orchestrator.keyRecord = *info
		orchestrator.mnemonic = mnemonics[i]
		orchestrator.keyring = kb

		c.orchestrators = append(c.orchestrators, orchestrator)
	}

	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: "gravity",
	}
}

func (c *chain) createOrchestrator(index int) *orchestrator {
	return &orchestrator{
		index: index,
	}
}

func (c *chain) clientContext(nodeURI string, kb *keyring.Keyring, fromName string, fromAddr sdk.AccAddress) (*client.Context, error) {
	encodingConfig := moduletestutil.MakeTestEncodingConfig(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		consensus.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distribution.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				gravityclient.ProposalHandler,
			},
		),
		//params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
	)
	encodingConfig.InterfaceRegistry.RegisterImplementations((*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
		&gravitytypes.MsgDelegateKeys{},
	)
	encodingConfig.InterfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})

	rpcClient, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	clientContext := client.Context{}.
		WithChainID(c.id).
		WithCodec(encodingConfig.Codec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithNodeURI(nodeURI).
		WithClient(rpcClient).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(*kb).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithOutputFormat("json").
		WithFrom(fromName).
		WithFromName(fromName).
		WithFromAddress(fromAddr).
		WithSkipConfirmation(true)

	return &clientContext, nil
}

func (c *chain) sendMsgs(clientCtx client.Context, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithChainID(c.id).
		WithTxConfig(clientCtx.TxConfig).
		WithGasAdjustment(1.2).
		WithKeybase(clientCtx.Keyring).
		WithGas(12345678).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	fromAddr := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, fromAddr); err != nil {
		return nil, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, fromAddr)
		if err != nil {
			return nil, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	txf = txf.WithFees("246913560testgb")

	err := tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
	if err != nil {
		return nil, err
	}

	resBytes := []byte{}
	_, err = clientCtx.Input.Read(resBytes)
	if err != nil {
		return nil, err
	}

	var res sdk.TxResponse
	err = cdc.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
