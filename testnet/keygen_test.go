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

func TestGeneratedKeygen(t *testing.T) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	require.NoError(t, err)

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	require.NoError(t, err)

	mnemonic := "weekend subject armor cupboard erosion cattle liquid rack play friend twenty magnet seminar right devote shoot sand design donor cigar employ outer debate rubber"

	account, err := kb.NewAccount("", mnemonic, "", "", algo)
	require.NoError(t, err)

	require.Equal(t, "cosmos1wwwfkhgmsvdr6wcdkq27w3kdmpye2s9khdhh20", account.GetAddress().String(), "mismatch of g")
}

func TestCompareKeygen(t *testing.T) {
	mnemonic := "weekend subject armor cupboard erosion cattle liquid rack play friend twenty magnet seminar right devote shoot sand design donor cigar employ outer debate rubber"

	memKeyring, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	require.NoError(t, err)

	memKeyringAlgos, _ := memKeyring.SupportedAlgorithms()
	memAlgo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), memKeyringAlgos)
	require.NoError(t, err)

	memAccount, err := memKeyring.NewAccount("", mnemonic, "", "", memAlgo)
	require.NoError(t, err)

	fileKeyring, err := keyring.New("testnet", keyring.BackendTest, "testnet/testdata", nil)
	require.NoError(t, err)

	fileKeyringAlgos, _ := fileKeyring.SupportedAlgorithms()
	fileAlgo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), fileKeyringAlgos)
	require.NoError(t, err)

	require.Equal(t, memAlgo, fileAlgo)

	fileAccount, err := fileKeyring.NewAccount("testname", mnemonic, "", "", fileAlgo)
	require.NoError(t, err)

	require.Equal(t, fileAccount.GetAddress(), memAccount.GetAddress(), "mismatch of g")
}