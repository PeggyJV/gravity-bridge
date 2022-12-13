package upgrade_tests

import (
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	"go.uber.org/zap"
	"sync"
)

// GravityBridgeChain is a local docker testnet for a Gravity Bridge chain.
// Implements the ibc.Chain interface.
type GravityBridgeChain struct {
	testName         string
	cfg              ibc.ChainConfig
	numValidators    int
	numFullNodes     int
	numOrchestrators int
	Validators       cosmos.ChainNodes
	FullNodes        cosmos.ChainNodes
	Orchestrators    OrchestratorNodes

	log *zap.Logger

	findTxMu sync.Mutex
}
