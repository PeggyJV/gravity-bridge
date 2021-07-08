package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/require"
)

func TestValidatorOut(t *testing.T) {
	withPristineE2EEnvironment(t, func(
		wd string,
		pool *dockertest.Pool,
		network *dockertest.Network,
	) {
		err := pool.RemoveContainerByName("gravity0")
		require.NoError(t, err, "error removing gravity0")

		// bring up the test runner
		t.Log("building and deploying test runner")
		testRunner, err := pool.BuildAndRunWithBuildOptions(
			&dockertest.BuildOptions{
				Dockerfile: "testnet.Dockerfile",
				ContextDir: "./orchestrator",
			},
			&dockertest.RunOptions{
				Name:      "test_runner",
				NetworkID: network.Network.ID,
				PortBindings: map[docker.Port][]docker.PortBinding{
					"8545/tcp": {{HostIP: "", HostPort: "8545"}},
				},
				Mounts: []string{
					fmt.Sprintf("%s/testdata:/testdata", wd),
				},
				Env: []string{
					"RUST_BACKTRACE=1",
					"RUST_LOG=INFO",
					"TEST_TYPE=VALIDATOR_OUT",
				},
			},
			func(config *docker.HostConfig) {},
		)
		require.NoError(t, err, "error bringing up test runner")

		t.Logf("deployed test runner at %s", testRunner.Container.ID)
		defer func() {
			testRunner.Close()
		}()

		container := testRunner.Container
		for container.State.Running {
			time.Sleep(10 * time.Second)
			container, err = pool.Client.InspectContainer(testRunner.Container.ID)
			require.NoError(t, err, "error inspecting test runner")
		}
		require.Equal(t, 0, container.State.ExitCode, "container exited with error")
	})
}
