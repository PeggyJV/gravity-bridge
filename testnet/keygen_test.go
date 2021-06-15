package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeygen(t *testing.T) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	require.NoError(t, err)

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	require.NoError(t, err)

	mnemonic := "weasel lunch attack blossom tone drum unfair worry risk level negative height sight nation inside task oyster client shiver aware neck mansion gun dune"

	account, err := kb.NewAccount("", mnemonic, "", "", algo)
	require.NoError(t, err)

	require.Equal(t, "cosmos18umn8nad5m8vcr3567v0ylu0m3q2ksrp5c3zf5", account.GetAddress().String(), "mismatch of g")
}