package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicChain(t *testing.T) {
	chain := Chain{
		DataDir: "testdata",
		ID:        "testchain",
		Nodes:     nil,
	}

	err := chain.CreateAndInitializeNodes(4)
	require.NoError(t, err, "error initializing chain")
}
