package main

import (
	"encoding/json"
	"fmt"
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

type Node struct {
	Chain   *Chain
	Index   uint8
	Moniker string

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

func (n *Node) ConfigDir() string {
	return fmt.Sprintf("%s/%s%d", n.Chain.ConfigDir(), n.Moniker, n.Index)
}

// MkDir creates the directory for the testnode
func (n *Node) MkDir() {
	p := path.Join(n.ConfigDir(), "config")
	fmt.Println(p)
	if err := os.MkdirAll(p, 0755); err != nil {
		panic(err)
	}
}

func (n *Node) init() error {
	encodingConfig := app.MakeEncodingConfig()
	cdc := encodingConfig.Marshaler

	n.MkDir()
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(n.ConfigDir())
	config.Moniker = n.Moniker

	_, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, "")
	if err != nil {
		return err
	}

	appState, err := json.MarshalIndent(app.ModuleBasics.DefaultGenesis(cdc), "", " ")
	if err != nil {
		return errors.Wrap(err, "Failed to marshall default genesis state")
	}

	genFile := config.GenesisFile()
	genDoc := &types.GenesisDoc{}
	if _, err := os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		genDoc, err = types.GenesisDocFromFile(genFile)
		if err != nil {
			return errors.Wrap(err, "Failed to read genesis doc from file")
		}
	}

	genDoc.ChainID = n.Chain.ID
	genDoc.Validators = nil
	genDoc.AppState = appState
	if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
		return errors.Wrap(err, "Failed to export gensis file")
	}
	//
	//toPrint := newPrintInfo(config.Moniker, n.Chain.ID, nodeID, "", appState)
	//
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	//return displayInfo(toPrint)
	return nil
}

func (n *Node) createKey(name string) (err error) {
	kb, err := keyring.New("testnet", keyring.BackendTest, n.ConfigDir(), nil)
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
	n.Mnemonic = mnemonic

	info, err := kb.NewAccount(name, mnemonic, "", "", algo)
	if err != nil {
		return err
	}
	n.KeyInfo = info

	return nil
}
