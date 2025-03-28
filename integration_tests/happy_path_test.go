package integration_tests

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ory/dockertest/v3"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

func (s *IntegrationTestSuite) TestHappyPath() {
	s.Run("Bring up chain, and test the happy path", func() {
		s.T().Logf("approving Gravity to spend ERC 20")
		err := s.approveERC20()
		s.Require().NoError(err, "error approving spending balance for the gravity contract on behalf of the first validator")

		allowance, err := s.getERC20AllowanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), gravityContract)
		s.Require().NoError(err, "error getting allowance of gravity contract spending on behalf of first validator")
		s.Require().Equal(UInt256Max(), allowance.BigInt(), "spending allowance not set correctly, got: %s", allowance.String())

		initialBalance, err := s.getEthTokenBalanceOf(common.HexToAddress(s.chain.validators[0].ethereumKey.address), testERC20contract)
		s.Require().NoError(err, "error getting first validator balance")
		s.Require().Equal(sdk.NewUint(10000).BigInt(), initialBalance.BigInt(), "balance was %s, expected 10000", initialBalance.String())

		for _, val := range s.chain.validators {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].address().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.address().String(), res.Balances)
		}

		// send from val 0 on eth to val 1 on cosmos
		s.T().Logf("sending to cosmos")
		err = s.sendToCosmos(s.chain.validators[1].address(), sdk.NewInt(200))
		s.Require().NoError(err, "error sending test denom to cosmos")

		for _, val := range s.chain.validators {
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
			s.Require().NoError(err)

			queryClient := banktypes.NewQueryClient(clientCtx)
			res, err := queryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].address().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.address().String(), res.Balances)
		}

		var gravityDenom string
		s.Require().Eventuallyf(func() bool {
			val := s.chain.validators[0]
			kb, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
			s.Require().NoError(err)

			bankQueryClient := banktypes.NewQueryClient(clientCtx)
			res, err := bankQueryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].address().String(),
				})
			if err != nil {
				return false
			}

			gbQueryClient := types.NewQueryClient(clientCtx)
			denomRes, err := gbQueryClient.ERC20ToDenom(context.Background(),
				&types.ERC20ToDenomRequest{
					Erc20: testERC20contract.String(),
				})
			if err != nil {
				s.T().Logf("error querying ERC20 denom %s, %e", testERC20contract.String(), err)
				return false
			}
			s.Require().False(denomRes.CosmosOriginated, "ERC20-originated denom marked as cosmos originated")
			gravityDenom = denomRes.Denom

			for _, coin := range res.Balances {
				if coin.Denom == gravityDenom && coin.Amount.Equal(sdk.NewInt(200)) {
					return true
				}
			}

			s.T().Logf("balance not found, received %v", res.Balances)

			return false
		}, 105*time.Second, 10*time.Second, "balance never found on cosmos")

		s.T().Logf("sending to ethereum")
		sendToEthereumMsg := types.NewMsgSendToEthereum(
			s.chain.validators[1].address(),
			s.chain.validators[1].ethereumKey.address,
			sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(100)},
			sdk.Coin{Denom: gravityDenom, Amount: sdk.NewInt(1)},
		)

		s.Require().Eventuallyf(func() bool {
			val := s.chain.validators[1]
			keyring, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
			s.Require().NoError(err)

			response, err := s.chain.sendMsgs(*clientCtx, sendToEthereumMsg)
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			if response.Code != 0 {
				if response.Code != 32 {
					s.T().Log(response)
				}

				return false
			}
			return true
		}, 105*time.Second, 10*time.Second, "unable to send to ethereum")

		s.T().Logf("funding community pool")
		orch := s.chain.orchestrators[0]
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err)

		fundCommunityPoolMsg := distrtypes.NewMsgFundCommunityPool(
			sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(1000000000))),
			orch.address(),
		)

		s.Require().Eventuallyf(func() bool {
			response, err := s.chain.sendMsgs(*clientCtx, fundCommunityPoolMsg)
			if err != nil {
				s.T().Logf("error: %s", err)
				return false
			}
			if response.Code != 0 {
				if response.Code != 32 {
					s.T().Log(response)
				}
				return false
			}
			return true
		}, 105*time.Second, 10*time.Second, "unable to fund community pool")

		s.Require().Eventuallyf(func() bool {
			distrQueryClient := distrtypes.NewQueryClient(clientCtx)
			poolRes, err := distrQueryClient.CommunityPool(context.Background(),
				&distrtypes.QueryCommunityPoolRequest{},
			)
			s.Require().NoError(err, "error retrieving community pool")
			s.Require().Greater(poolRes.Pool.AmountOf(testDenom).BigInt().Int64(), sdk.NewDec(1000000000).BigInt().Int64())

			return true
		}, 20*time.Second, 2*time.Second, "community pool balance not high enough")

		s.T().Logf("deploying testgb as an ERC20")
		gbQueryClient := types.NewQueryClient(clientCtx)
		paramsRes, err := gbQueryClient.DenomToERC20Params(context.Background(),
			&types.DenomToERC20ParamsRequest{
				Denom: testDenom,
			})
		s.Require().NoError(err, "error retrieving ERC20 params for testgb denom")

		err = s.deployERC20(paramsRes.BaseDenom, paramsRes.Erc20Name, paramsRes.Erc20Symbol, uint8(paramsRes.Erc20Decimals))
		s.Require().NoError(err, "error deploying testgb as an ERC20")

		s.Require().Eventuallyf(func() bool {
			erc20Res, err := gbQueryClient.DenomToERC20(context.Background(),
				&types.DenomToERC20Request{
					Denom: testDenom,
				},
			)
			if err != nil {
				s.T().Logf("erc20 not deployed yet, waiting")
				return false
			}

			s.Require().True(erc20Res.CosmosOriginated)
			return true
		}, 180*time.Second, 10*time.Second, "unable to verify ERC20 deployment")

		erc20Res, err := gbQueryClient.DenomToERC20(context.Background(),
			&types.DenomToERC20Request{
				Denom: testDenom,
			},
		)
		s.Require().NoError(err, "error querying ERC20 for testgb denom")
		erc20Contract := common.HexToAddress(erc20Res.Erc20)
		initialBalance, err = s.getEthTokenBalanceOf(common.HexToAddress(s.chain.validators[2].ethereumKey.address), erc20Contract)
		s.Require().NoError(err, "error getting destination balance")

		s.T().Logf("create governance proposal to fund an ethereum address")
		orch = s.chain.orchestrators[0]
		clientCtx, err = s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.address())
		s.Require().NoError(err)

		sendAmount := int64(900)
		proposal := types.CommunityPoolEthereumSpendProposal{
			Title:       "community pool spend ethereum",
			Description: "community pool spend ethereum",
			Recipient:   s.chain.validators[2].ethereumKey.address,
			Amount:      sdk.NewCoin(testDenom, sdk.NewInt(sendAmount)),
			BridgeFee:   sdk.NewCoin(testDenom, sdk.NewInt(1000000)),
		}

		proposalMsg, err := govtypesv1beta1.NewMsgSubmitProposal(
			&proposal,
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: sdk.NewInt(2),
				},
			},
			orch.address(),
		)
		s.Require().NoError(err, "unable to create governance proposal")

		s.T().Log("submit proposal spending community pool funds")
		submitProposalResponse, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		govQueryClient := govtypesv1beta1.NewQueryClient(clientCtx)

		s.Require().Eventually(func() bool {
			s.T().Log("check proposal was submitted correctly")
			govQueryClient = govtypesv1beta1.NewQueryClient(clientCtx)
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypesv1beta1.QueryProposalsRequest{})
			s.Require().NoError(err)
			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypesv1beta1.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")
			return true
		}, 60*time.Second, 2*time.Second, "proposal not submitted correctly")

		s.T().Log("vote for community spend proposal")
		for _, val := range s.chain.validators {
			kr, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.address())
			s.Require().NoError(err)

			voteMsg := govtypesv1beta1.NewMsgVote(val.address(), 1, govtypesv1beta1.OptionYes)
			voteResponse, err := s.chain.sendMsgs(*clientCtx, voteMsg)
			s.Require().NoError(err)
			s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
		}

		s.T().Log("wait for community spend proposal to be approved")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypesv1beta1.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)
			return govtypesv1beta1.StatusPassed == proposalQueryResponse.Proposal.Status
		}, time.Second*30, time.Second*5, "proposal was never accepted")

		s.T().Logf("initial balance of %s of token %s is %v", s.chain.validators[2].ethereumKey.address, erc20Contract.Hex(), initialBalance)
		s.T().Log("waiting for community funds to reach destination")
		s.Require().Eventuallyf(func() bool {
			s.T().Logf("getting balance in %s of token %s", s.chain.validators[2].ethereumKey.address, erc20Contract.Hex())
			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(s.chain.validators[2].ethereumKey.address), erc20Contract)
			s.Require().NoError(err, "error getting destination balance")
			s.T().Logf("balance is %v", balance)

			if balance.LT(sdk.NewInt(sendAmount)) {
				s.T().Logf("funds not received yet, dest balance: %s", balance.String())
				return false
			}

			s.Require().Equal(balance.BigInt(), sdk.NewInt(sendAmount).BigInt(), "balance was %s, expected %v", balance.String(), sendAmount)
			return true
		}, time.Second*180, time.Second*10, "community funds did not reach destination")
	})

	s.Run("Test orchestrator EthereumTxConfirmation after transaction execution", func() {
		val := s.chain.validators[0]
		keyring, err := val.keyring()
		s.Require().NoError(err)
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
		s.Require().NoError(err)

		recipient := s.chain.validators[3].ethereumKey.address
		gravityQueryClient := types.NewQueryClient(clientCtx)
		gravityResponse, err := gravityQueryClient.DenomToERC20(context.Background(),
			&types.DenomToERC20Request{
				Denom: testDenom,
			},
		)
		s.Require().NoError(err, "error querying ERC20 for testgb denom")
		s.Require().True(gravityResponse.CosmosOriginated)
		testDenomERC20 := common.HexToAddress(gravityResponse.Erc20)

		initialBalance, err := s.getEthTokenBalanceOf(common.HexToAddress(recipient), testDenomERC20)
		s.Require().NoError(err, "error getting initial balance")
		// Turn off one orchestrator (let's use the second one)
		s.T().Log("Turning off orchestrator1")
		err = s.dockerPool.RemoveContainerByName("orchestrator1")
		s.Require().NoError(err, "Failed to remove orchestrator1 container")

		// Prepare SendToEthereum message
		s.T().Log("Preparing SendToEthereum message")
		// Query the balance of testERC20Denom for validator0
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		balanceRes, err := bankQueryClient.AllBalances(
			context.Background(),
			&banktypes.QueryAllBalancesRequest{
				Address: val.address().String(),
			},
		)
		s.Require().NoError(err)

		s.T().Logf("Validator0 balances: %s", balanceRes.Balances)
		sendAmount := sdk.NewInt(420)
		feeAmount := sdk.NewInt(1)
		sendToEthereumMsg := types.NewMsgSendToEthereum(
			s.chain.validators[0].address(),
			recipient,
			sdk.NewCoin(testDenom, sendAmount),
			sdk.NewCoin(testDenom, feeAmount),
		)
		s.T().Logf("Sent %s %s to %s", sendAmount, testDenom, recipient)

		// Getting latest Batch Nonce
		completed, err := gravityQueryClient.CompletedBatchTxs(context.Background(), &types.CompletedBatchTxsRequest{})
		s.Require().NoError(err, "error querying CompletedBatchTxs")
		// Search through completed batch txs to get the highest nonce
		var highestNonce uint64
		for _, batchTx := range completed.CompletedBatchTxs {
			if batchTx.BatchNonce > highestNonce {
				highestNonce = batchTx.BatchNonce
			}
		}
		s.T().Logf("Highest completed batch nonce: %d", highestNonce)

		// Send the message
		s.T().Logf("Sending SendToEthereum message with %s", val.address().String())
		response, err := s.chain.sendMsgs(*clientCtx, sendToEthereumMsg)
		s.Require().NoError(err)
		s.Require().Zero(response.Code, "SendToEthereum failed: %s", response.RawLog)

		// Wait for the transaction to complete on Ethereum
		s.T().Log("Waiting for transaction to complete on Ethereum")
		s.T().Logf("Recipient: %s, initial balance: %s", recipient, initialBalance.String())
		s.Require().Eventually(func() bool {
			balance, err := s.getEthTokenBalanceOf(common.HexToAddress(recipient), testDenomERC20)
			if err != nil {
				s.T().Logf("Error getting balance: %v", err)
				return false
			}
			s.T().Logf("Balance: %s", balance.String())
			expectedBalance := initialBalance.Add(sendAmount)
			return balance.Equal(expectedBalance)
		}, 5*time.Minute, 10*time.Second, "Transaction did not complete on Ethereum")

		// Wait for the CompletedOutgoingTx to be created
		s.T().Log("Waiting for CompletedOutgoingTx to be created")
		expectedNonce := highestNonce + 1
		s.T().Logf("Expected nonce: %d", expectedNonce)
		s.Require().Eventually(func() bool {
			res, err := gravityQueryClient.CompletedBatchTxs(context.Background(), &types.CompletedBatchTxsRequest{})
			if err != nil {
				s.T().Logf("Error querying CompletedBatchTxs: %v", err)
				return false
			}

			for _, batchTx := range res.CompletedBatchTxs {
				if batchTx.BatchNonce >= expectedNonce {
					return true
				}
			}

			return false
		}, 5*time.Minute, 3*time.Second, "CompletedBatchTx was not found")

		// Turn the orchestrator back on
		s.T().Log("Turning orchestrator1 back on")
		err = s.startOrchestrator1()
		s.Require().NoError(err, "Failed to restart orchestrator1")

		// Watch for EthereumTxConfirmation from the restarted orchestrator
		s.T().Log("Watching for EthereumTxConfirmation from restarted orchestrator")
		s.Require().Eventually(func() bool {
			val := s.chain.validators[1]
			keyring, err := val.keyring()
			s.Require().NoError(err)
			clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
			s.Require().NoError(err)

			// Get the validator address with cosmosvaloper prefix
			cosmosValOperAddr, err := getValOperatorAddress(val.address().String())
			s.Require().NoError(err)
			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BatchTxConfirmationsByValidator(context.Background(), &types.BatchTxConfirmationsByValidatorRequest{
				ValidatorAddress: cosmosValOperAddr,
			})
			if err != nil {
				s.T().Logf("Error querying EthereumTxConfirmation: %v", err)
				return false
			}

			// Check if the height has increased, indicating that the orchestrator has caught up
			return len(res.BatchTxConfirmations) > 0
		}, 5*time.Minute, 10*time.Second, "Orchestrator did not submit EthereumTxConfirmation after restart")

		s.T().Log("Orchestrator successfully submitted EthereumTxConfirmation after transaction execution")
	})
}

func (s *IntegrationTestSuite) startOrchestrator1() error {
	i := 1
	orch := s.chain.orchestrators[i]
	gorcCfg := fmt.Sprintf(`keystore = "/root/gorc/keystore/"

[gravity]
contract = "%s"
fees_denom = "%s"

[ethereum]
key_derivation_path = "m/44'/60'/0'/0/0"
rpc = "http://%s:8545"

[cosmos]
key_derivation_path = "m/44'/118'/1'/0/0"
grpc = "http://%s:9090"
gas_price = { amount = %s, denom = "%s" }
prefix = "cosmos"
gas_adjustment = 2.0
msg_batch_size = 5
`,
		gravityContract.String(),
		testDenom,
		// NOTE: container names are prefixed with '/'
		s.ethResource.Container.Name[1:],
		s.valResources[i].Container.Name[1:],
		minGasPrice,
		testDenom,
	)

	val := s.chain.validators[i]

	gorcCfgPath := path.Join(val.configDir(), "gorc")
	err := os.MkdirAll(gorcCfgPath, 0755)
	if err != nil {
		return err
	}

	filePath := path.Join(gorcCfgPath, "config.toml")
	err = writeFile(filePath, []byte(gorcCfg))
	if err != nil {
		return err
	}

	// We must first populate the orchestrator's keystore prior to starting
	// the orchestrator gorc process. The keystore must contain the Ethereum
	// and orchestrator keys. These keys will be used for relaying txs to
	// and from the test network and Ethereum. The gorc_bootstrap.sh scripts encapsulates
	// this entire process.
	//
	// NOTE: If the Docker build changes, the script might have to be modified
	// as it relies on busybox.
	err = copyFile(
		filepath.Join("integration_tests", "gorc_bootstrap.sh"),
		filepath.Join(gorcCfgPath, "gorc_bootstrap.sh"),
	)
	if err != nil {
		return err
	}

	resource, err := s.dockerPool.RunWithOptions(
		&dockertest.RunOptions{
			Name:       orch.instanceName(),
			NetworkID:  s.dockerNetwork.Network.ID,
			Repository: "orchestrator",
			Tag:        "prebuilt",
			Mounts: []string{
				fmt.Sprintf("%s/:/root/gorc", gorcCfgPath),
			},
			Env: []string{
				fmt.Sprintf("ORCH_MNEMONIC=%s", orch.mnemonic),
				fmt.Sprintf("ETH_PRIV_KEY=%s", val.ethereumKey.privateKey),
				"RUST_BACKTRACE=full",
				"RUST_LOG=debug",
			},
			Entrypoint: []string{
				"sh",
				"-c",
				"chmod +x /root/gorc/gorc_bootstrap.sh && /root/gorc/gorc_bootstrap.sh",
			},
		},
		noRestart,
	)
	if err != nil {
		return err
	}

	s.orchResources[i] = resource
	s.T().Logf("started orchestrator container: %s", resource.Container.ID)
	return nil
}

func getValOperatorAddress(address string) (string, error) {
	// Decode the old address
	_, bz, err := bech32.DecodeAndConvert(address)
	if err != nil {
		return "", fmt.Errorf("failed to decode address: %w", err)
	}

	// Encode with the new prefix
	newAddress, err := bech32.ConvertAndEncode("cosmosvaloper", bz)
	if err != nil {
		return "", fmt.Errorf("failed to encode address: %w", err)
	}

	return newAddress, nil
}
