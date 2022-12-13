package upgrade_tests

import (
	"sync"

	dockerclient "github.com/docker/docker/client"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"go.uber.org/zap"
)

// OrchestratorNode represents orchestrator sidecar processes for a validator in the network
type OrchestratorNode struct {
	VolumeName   string
	Index        int
	Chain        ibc.Chain
	Validator    bool
	NetworkID    string
	DockerClient *dockerclient.Client
	Client       rpcclient.Client
	TestName     string
	Image        ibc.DockerImage

	lock sync.Mutex
	log  *zap.Logger

	containerID string
}

// OrchestratorNodes is a collection of OrchestratorNode
type OrchestratorNodes []*OrchestratorNode
