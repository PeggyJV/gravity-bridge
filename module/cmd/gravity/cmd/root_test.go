package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/cli"
	"io"
	"os"
	"strings"
	"testing"
)


func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}


func executeCommandWithContext(ctx context.Context, root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.ExecuteContext(ctx)

	return buf.String(), err
}

type KeyOutput struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Address string `json:"address"`
	PubKey string `json:"pubkey"`
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestKeyGen(t *testing.T) {
	//rootCmd, config := NewRootCmd()
	//t.Log(config)

	buf := bytes.NewBuffer(nil)

	mnemonic := "weasel lunch attack blossom tone drum unfair worry risk level negative height sight nation inside task oyster client shiver aware neck mansion gun dune"

	// generate key from binary
	keyCmd := keys.AddKeyCommand()
	keyCmd.Flags().String(cli.OutputFlag, "json", "output flag")
	keyCmd.SetArgs([]string{"--dry-run=true", "--output=json", "--recover=true", "orch"})
	keyCmd.SetIn(strings.NewReader(mnemonic + "\n"))
	keyCmd.SetOut(buf)
	keyCmd.SetErr(buf)

	err := Execute(keyCmd)
	require.NoError(t, err)

	var key KeyOutput
	output := buf.Bytes()
	t.Log("outputs: ", string(output))
	err = json.Unmarshal(output, &key)
	require.NoError(t, err)

	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return
	}

	account, err := kb.NewAccount("", mnemonic, "", "", algo)
	require.NoError(t, err)

	require.Equal(t, account.GetAddress().String(), key.Address)
}
