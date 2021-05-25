package testnet

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/types"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/gravity-bridge/module/app"
	"github.com/spf13/cobra"
	tmcfg "github.com/tendermint/tendermint/config"
)

const (
	flagChainID = "chain-id"
	flagHomeDir = "home"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "testnet",
		Short: "Gravity Testnet",
	}
	rootCmd.PersistentFlags().String("log_level", "info", "The logging level in the format of <module>:<level>,...")
	rootCmd.PersistentFlags().String("chain-id", "testchain", "The name of the chain generated for test")

	rootCmd.AddCommand(
		StartTestnet(),
		GenerateKeys(),
	)
}

func StartTestnet() *cobra.Command {
	cmd := &cobra.Command{
		Use: "start",
		Short: "bring up the testnet",

		RunE: runStartTestnetCmd,
	}

	cmd.SetOut(cmd.OutOrStdout())
	cmd.SetErr(cmd.ErrOrStderr())

	return cmd
}

func runStartTestnetCmd(cmd *cobra.Command, args []string) error {

	return nil
}

func InitializeTestnet() *cobra.Command {
	cmd := &cobra.Command{
		Use: "initialize",
		Short: "generate the files necessary for a cosmos chain",

		RunE: runInitializeTestnetCmd,
	}

	cmd.Flags().String(flagHomeDir, "./testdata", "test chain's home directory")
	cmd.Flags().String(flagNodePrefix, "gravity", "name prefix for initialized nodes")
	cmd.Flags().Uint8(flagNodeCount, 1, "number of nodes to initialize")

	cmd.SetOut(cmd.OutOrStdout())
	cmd.SetErr(cmd.ErrOrStderr())

	return cmd
}

func runInitializeTestnetCmd(cmd *cobra.Command, args []string) error {
	config := tmcfg.DefaultConfig()
	homeDir, err := cmd.Flags().GetString(flagHomeDir)
	if err != nil {
		return err
	}
	config.SetRoot(homeDir)

	chainID, err := cmd.Flags().GetString(flagChainID)
	if err != nil {
		return err
	}

	prefix, err := cmd.Flags().GetString(flagNodePrefix)
	if err != nil {
		return err
	}

	count, err := cmd.Flags().GetUint8(flagNodeCount)
	if err != nil {
		return err
	}

	return initializeTestnet(homeDir, chainID, prefix, count)
}

func initializeTestnet(homeDir string, chainID string, prefix string, count uint8) error {
	for i := uint8(0); i < count; i++ {
		name := fmt.Sprintf("%s%d", prefix, i)
		err := initializeTestnetNode(homeDir, chainID, name)
		if err != nil {
			return err
		}
	}

	return nil
}

func initializeTestnetNode(homeDir string, chainID string, name string) error {
	config := tmcfg.DefaultConfig()
	config.SetRoot(homeDir)

	nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, "")
	if err != nil {
		return err
	}

	config.Moniker = name

	ctx := client.Context{}
	cdc := ctx.JSONMarshaler
	genFile := config.GenesisFile()
	appState, err := json.MarshalIndent(app.ModuleBasics.DefaultGenesis(cdc), "", " ")

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

	genDoc.ChainID = chainID
	genDoc.Validators = nil
	genDoc.AppState = appState
	if err = genutil.ExportGenesisFile(genDoc, genFile); err != nil {
		return errors.Wrap(err, "Failed to export gensis file")
	}

	toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
	return displayInfo(toPrint)
}
