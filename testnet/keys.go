package testnet

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

const (
	flagNodeCount = "count"
	flagNodePrefix = "prefix"
)

func GenerateKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use: "generate-keys",
		Short: "generate testnet keys",
		Long: "Generates the keys for the nodes and outputs them to the target directory",

		RunE: runGenerateKeysCmd,
	}

	cmd.Flags().Int(flagNodeCount, 1, "The number of nodes to generate keys for")
	cmd.Flags().String(flagNodePrefix, "gravity", "The prefix for the names of the nodes")

	cmd.SetOut(cmd.OutOrStdout())
	cmd.SetErr(cmd.ErrOrStderr())

	return cmd
}

func runGenerateKeysCmd(cmd *cobra.Command, args []string) error {
	nodeCount, err := cmd.Flags().GetInt(flagNodeCount)
	if err != nil {
		return err
	}

	return generateKeys(nodeCount)
}

type GeneratedKey struct {

}

func generateKeys(nodeCount int, nodePrefix string) error {
	nodeKeys := make(map[string]GeneratedKey)

	for i :=0; i < nodeCount; i++ {
		nodeName := fmt.Sprintf("%s%d", nodePrefix, i)
		nodeKeys[nodeName] = GeneratedKey{}
	}
}