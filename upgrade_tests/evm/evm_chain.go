package evm

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/network"
	"github.com/ory/dockertest/v3/docker/types/container"
	"io"

	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	"github.com/strangelove-ventures/ibctest/v6/internal/dockerutil"
	"go.uber.org/zap"
)

type EVMChain struct {
	log         *zap.Logger
	testName    string
	cfg         ibc.ChainConfig
	Image       ibc.DockerImage
	DockerClient *dockerclient.Client
	containerID string
	port        string
	chainID 	uint32
}

//func NewEVMChain(log *zap.Logger, testName string, chainConfig ibc.ChainConfig, port string) *EVMChain {
//	return &EVMChain{
//		log:      log,
//		testName: testName,
//		cfg:      chainConfig,
//		port:     port,
//	}
//}

func (e *EVMChain) Config() ibc.ChainConfig {
	return e.cfg
}

func (e *EVMChain) Initialize(ctx context.Context, testName string, cli *dockerclient.Client, networkID string) error {
	chainCfg := e.Config()

	for _, image := range chainCfg.Images {
		rc, err := cli.ImagePull(
			ctx,
			image.Repository+":"+image.Version,
			types.ImagePullOptions{},
		)
		if err != nil {
			e.log.Error("Failed to pull image",
				zap.Error(err),
				zap.String("repository", image.Repository),
				zap.String("tag", image.Version),
			)
		} else {
			_, _ = io.Copy(io.Discard, rc)
			_ = rc.Close()
		}
	}

	if len(chainCfg.Images) != 1 {
		e.log.Error("Multiple images for EVM")
	}

	e.Image = chainCfg.Images[0]
	e.DockerClient = cli

	return nil
}

func (e *EVMChain) Start(testName string, ctx context.Context, _ ...ibc.WalletAmount) error {
	cc, err := e.DockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image: e.Image.Ref(),
			ExposedPorts: e.port,

		},
		&container.HostConfig{

		},
		&network.NetworkingConfig{},
		nil,
		e.Name(),
	)

	if err := dockerutil.StartContainer(ctx, e.DockerClient, e.containerID)
}

func (e *EVMChain) Name() string {
	return fmt.Sprintf("evm-%d", e.chainID)
}

func (e *EVMChain) Exec(ctx context.Context, cmd []string, env []string) (stdout, stderr []byte, err error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) ExportState(ctx context.Context, height int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetRPCAddress() string {
	return e.port
}

func (e *EVMChain) GetGRPCAddress() string {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetHostRPCAddress() string {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetHostGRPCAddress() string {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) HomeDir() string {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) CreateKey(ctx context.Context, keyName string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) RecoverKey(ctx context.Context, name, mnemonic string) error {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetAddress(ctx context.Context, keyName string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) SendFunds(ctx context.Context, keyName string, amount ibc.WalletAmount) error {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) SendIBCTransfer(ctx context.Context, channelID, keyName string, amount ibc.WalletAmount, options ibc.TransferOptions) (ibc.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) Height(ctx context.Context) (uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetBalance(ctx context.Context, address string, denom string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) GetGasFeesInNativeDenom(gasPaid int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) Acknowledgements(ctx context.Context, height uint64) ([]ibc.PacketAcknowledgement, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EVMChain) Timeouts(ctx context.Context, height uint64) ([]ibc.PacketTimeout, error) {
	//TODO implement me
	panic("implement me")
}
