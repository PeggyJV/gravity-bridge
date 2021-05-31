package main

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
	"path"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/go-bip39"
	"github.com/cosmos/gravity-bridge/module/app"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/types"

)

type Validator struct {
	Chain   *Chain
	Index   uint8
	Moniker string

	// Key management
	Mnemonic string
	KeyInfo  keyring.Info

	EthereumKey EthereumKey
}

type EthereumKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Address    string `json:"address"`
}

type Orchestrator struct {
	Chain *Chain
	Index uint8

	// Key management
	Mnemonic string
	KeyInfo  keyring.Info
}

// createMnemonic creates a new mnemonic
func createMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func (v *Validator) ConfigDir() string {
	return fmt.Sprintf("%s/%s%d", v.Chain.ConfigDir(), v.Moniker, v.Index)
}

// MkDir creates the directory for the testnode
func (v *Validator) MkDir() {
	p := path.Join(v.ConfigDir(), "config")
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}
}

func getGenDoc(path string) (doc *types.GenesisDoc, err error) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)

	genFile := config.GenesisFile()
	doc = &types.GenesisDoc{}
	if _, err = os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
	} else {
		doc, err = types.GenesisDocFromFile(genFile)
		if err != nil {
			err = errors.Wrap(err, "Failed to read genesis doc from file")
			return
		}
	}

	return
}

func (v *Validator) init() error {
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	v.MkDir()
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(v.ConfigDir())
	config.Moniker = v.Moniker

	_, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, "")
	if err != nil {
		return err
	}

	genDoc, err := getGenDoc(v.ConfigDir())
	if err != nil {
		return err
	}

	appState, err := json.MarshalIndent(app.ModuleBasics.DefaultGenesis(cdc), "", " ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshall default genesis state")
	}

	genDoc.ChainID = v.Chain.ID
	genDoc.Validators = nil
	genDoc.AppState = appState
	if err = genutil.ExportGenesisFile(genDoc, config.GenesisFile()); err != nil {
		return errors.Wrap(err, "Failed to export gensis file")
	}

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	return nil
}

// createMemoryKey creates a key but doesn't store it to any files
func createMemoryKey() (mnemonic string, info *keyring.Info, err error) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return
	}

	// Get bip39 mnemonic
	mnemonic, err = createMnemonic()
	if err != nil {
		return
	}

	account, err := kb.NewAccount("", mnemonic, "", "", algo)
	info = &account
	return
}

// createKey creates a new account and writes it to a validator's config directory
func (v *Validator) createKey(name string) (err error) {
	kb, err := keyring.New("testnet", keyring.BackendTest, v.ConfigDir(), nil)
	if err != nil {
		return err
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return err
	}

	// Get bip39 mnemonic
	mnemonic, err := createMnemonic()
	if err != nil {
		return err
	}
	v.Mnemonic = mnemonic

	info, err := kb.NewAccount(name, mnemonic, "", "", algo)
	if err != nil {
		return err
	}
	v.KeyInfo = info

	return nil
}

func addGenesisAccount(path string, moniker string, accAddr sdk.AccAddress, coinsStr string) (err error) {
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)
	config.Moniker = moniker

	coins, err := sdk.ParseCoinsNormalized(coinsStr)
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	balances := banktypes.Balance{Address: accAddr.String(), Coins: coins.Sort()}
	genAccount := authtypes.NewBaseAccount(accAddr, nil, 0, 0)

	genFile := config.GenesisFile()
	appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
	if err != nil {
		return fmt.Errorf("failed to unmarshal genesis state: %w", err)
	}

	authGenState := authtypes.GetGenesisStateFromAppState(cdc, appState)

	accs, err := authtypes.UnpackAccounts(authGenState.Accounts)
	if err != nil {
		return fmt.Errorf("failed to get accounts from any: %w", err)
	}

	if accs.Contains(accAddr) {
		return fmt.Errorf("cannot add account at existing address %s", accAddr)
	}

	// Add the new account to the set of genesis accounts and sanitize the
	// accounts afterwards.
	accs = append(accs, genAccount)
	accs = authtypes.SanitizeGenesisAccounts(accs)


	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		return fmt.Errorf("failed to convert accounts into any's: %w", err)
	}
	authGenState.Accounts = genAccs

	authGenStateBz, err := cdc.MarshalJSON(&authGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal auth genesis state: %w", err)
	}

	appState[authtypes.ModuleName] = authGenStateBz


	bankGenState := banktypes.GetGenesisStateFromAppState(cdc, appState)
	bankGenState.Balances = append(bankGenState.Balances, balances)
	bankGenState.Balances = banktypes.SanitizeGenesisBalances(bankGenState.Balances)

	bankGenStateBz, err := cdc.MarshalJSON(bankGenState)
	if err != nil {
		return fmt.Errorf("failed to marshal bank genesis state: %w", err)
	}

	appState[banktypes.ModuleName] = bankGenStateBz

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return fmt.Errorf("failed to marshal application genesis state: %w", err)
	}

	genDoc.AppState = appStateJSON
	return genutil.ExportGenesisFile(genDoc, genFile)
}

func (v *Validator) generateEthereumKey() (err error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	v.EthereumKey = EthereumKey{
		PrivateKey: hexutil.Encode(privateKeyBytes),
		PublicKey:  hexutil.Encode(publicKeyBytes),
		Address:    crypto.PubkeyToAddress(*publicKeyECDSA).Hex(),
	}
	return
}