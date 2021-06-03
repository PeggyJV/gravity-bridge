package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	types3 "github.com/cosmos/cosmos-sdk/crypto/types"
	"os"
	"path"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	crypto2 "github.com/cosmos/cosmos-sdk/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	signing2 "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	tx2 "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/go-bip39"
	"github.com/cosmos/gravity-bridge/module/app"
	types2 "github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	cfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

type Validator struct {
	Chain   *Chain
	Index   uint8
	Moniker string

	// Key management
	Mnemonic string
	KeyInfo  keyring.Info
	PrivateKey types3.PrivKey

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

func getGenDoc(path string) (doc *tmtypes.GenesisDoc, err error) {
	serverCtx := server.NewDefaultContext()
	config := serverCtx.Config
	config.SetRoot(path)

	genFile := config.GenesisFile()
	doc = &tmtypes.GenesisDoc{}
	if _, err = os.Stat(genFile); err != nil {
		if !os.IsNotExist(err) {
			return
		}
		err = nil
	} else {
		doc, err = tmtypes.GenesisDocFromFile(genFile)
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

	privKeyArmor, err := kb.ExportPrivKeyArmor(name, "testpassphrase")
	if err != nil {
		return err
	}
	privKey, _, err := crypto2.UnarmorDecryptPrivKey(privKeyArmor, "testpassphrase")
	if err != nil {
		return err
	}
	v.PrivateKey = privKey

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

// BuildCreateValidatorMsg makes a new MsgCreateValidator.
func (v *Validator) buildCreateValidatorMsg(amount sdk.Coin) (sdk.Msg, error) {
	description := types.NewDescription(
		v.Moniker,
		"",
		"",
		"",
		"",
	)

	commissionRates := types.CommissionRates{
		Rate:          sdk.MustNewDecFromStr("0.1"),
		MaxRate:       sdk.MustNewDecFromStr("0.2"),
		MaxChangeRate: sdk.MustNewDecFromStr("0.01"),
	}

	// get the initial validator min self delegation
	minSelfDelegation, _ := sdk.NewIntFromString("1")

	msg, err := types.NewMsgCreateValidator(
		sdk.ValAddress(v.KeyInfo.GetAddress()), v.KeyInfo.GetPubKey(), amount, description, commissionRates, minSelfDelegation,
	)
	return msg, err
}

func (v *Validator) nodeID() string {
	return hex.EncodeToString(v.KeyInfo.GetPubKey().Address())
}

func (v *Validator) signMsg(msg sdk.Msg) (signing.Tx, error) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &types.MsgCreateValidator{}, &types2.MsgDelegateKeys{})
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	signModes := []signing2.SignMode{signing2.SignMode_SIGN_MODE_DIRECT}
	txConfig := tx2.NewTxConfig(marshaler, signModes)
	txBuilder := txConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(msg); err != nil {
		return nil, err
	}

	txBuilder.SetMemo(fmt.Sprintf("%s@%s%d:26656", v.nodeID(), v.Moniker, v.Index))
	fees := sdk.Coins{sdk.Coin{}}
	txBuilder.SetFeeAmount(fees)
	txBuilder.SetGasLimit(200000)
	txBuilder.SetTimeoutHeight(0)

	signerData := authsigning.SignerData{
		ChainID:       v.Chain.ID,
		AccountNumber: 0,
		Sequence:      0,
	}

	// For SIGN_MODE_DIRECT, calling SetSignatures calls setSignerInfos on
	// TxBuilder under the hood, and SignerInfos is needed to generated the
	// sign bytes. This is the reason for setting SetSignatures here, with a
	// nil signature.
	//
	// Note: this line is not needed for SIGN_MODE_LEGACY_AMINO, but putting it
	// also doesn't affect its generated sign bytes, so for code's simplicity
	// sake, we put it here.
	sigData := signing2.SingleSignatureData{
		SignMode:  signing2.SignMode_SIGN_MODE_DIRECT,
		Signature: nil,
	}
	sig := signing2.SignatureV2{
		PubKey:   v.KeyInfo.GetPubKey(),
		Data:     &sigData,
		Sequence: 0,
	}

	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	bytesToSign, err := txConfig.SignModeHandler().GetSignBytes(signing2.SignMode_SIGN_MODE_DIRECT, signerData, txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	// Sign those bytes
	sigBytes, err := v.PrivateKey.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}

	// Construct the SignatureV2 struct
	sigData = signing2.SingleSignatureData{
		SignMode:  signing2.SignMode_SIGN_MODE_DIRECT,
		Signature: sigBytes,
	}
	sig = signing2.SignatureV2{
		PubKey:   v.KeyInfo.GetPubKey(),
		Data:     &sigData,
		Sequence: 0,
	}
	if err := txBuilder.SetSignatures(sig); err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}
