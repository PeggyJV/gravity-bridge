package integration_tests

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func (s *IntegrationTestSuite) TestHappyPath() {
	s.Run("Bring up chain, and test the happy path", func() {
		var err error

		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)
		bankQueryClient := banktypes.NewQueryClient(clientCtx)
		gbQueryClient := types.NewQueryClient(clientCtx)
		distrQueryClient := distrtypes.NewQueryClient(clientCtx)
		govQueryClient := govtypes.NewQueryClient(clientCtx)

		for chainIndex, evm := range s.evms {
			s.T().Logf("approving Gravity to spend ERC 20 on %s", evm.Name)
			err := s.approveERC20(evm, testERC20Contracts[chainIndex], gravityContracts[chainIndex])
			s.Require().NoError(err, "error approving spending balance for the gravity contract on behalf of the first validator")

			allowance, err := s.getERC20AllowanceOf(evm, testERC20Contracts[chainIndex], common.HexToAddress(s.chain.validators[0].ethereumKey.address), gravityContracts[chainIndex])
			s.Require().NoError(err, "error getting allowance of gravity contract spending on behalf of first validator")
			s.Require().Equal(UInt256Max(), allowance.BigInt(), "spending allowance not set correctly, got: %s", allowance.String())

			balance, err := s.getEthTokenBalanceOf(evm, common.HexToAddress(s.chain.validators[0].ethereumKey.address), testERC20Contracts[chainIndex])
			s.Require().NoError(err, "error getting first validator balance")
			s.Require().Equal(sdk.NewUint(10000).BigInt(), balance.BigInt(), "balance was %s, expected 10000", balance.String())
		}
		for _, val := range s.chain.validators {
			res, err := bankQueryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].keyInfo.GetAddress().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.keyInfo.GetAddress().String(), res.Balances)
		}

		// send from val 0 on eth to val 1 on cosmos
		for chainIndex, evm := range s.evms {
			s.T().Logf("sending to cosmos from %s", evm.Name)
			err = s.sendToCosmos(evm, gravityContracts[chainIndex], testERC20Contracts[chainIndex], s.chain.validators[1].keyInfo.GetAddress(), sdk.NewInt(200))
			s.Require().NoError(err, "error sending test denom to cosmos from %s", ChainNames[chainIndex])
		}

		for i, evm := range s.evms {
			res, err := gbQueryClient.ERC20ToDenom(context.Background(), &types.ERC20ToDenomRequest{
				Erc20:   testERC20Contracts[i].String(),
				ChainId: evm.ID,
			})
			s.Require().NoError(err)
			gravityDenoms[i] = res.Denom
		}

		for _, val := range s.chain.validators {
			res, err := bankQueryClient.AllBalances(context.Background(),
				&banktypes.QueryAllBalancesRequest{
					Address: s.chain.validators[1].keyInfo.GetAddress().String(),
				})
			s.Require().NoError(err)
			s.T().Logf("balances for %s: %s", val.keyInfo.GetAddress().String(), res.Balances)
		}

		s.T().Logf("gravity contracts: %v", gravityContracts)
		s.T().Logf("gravity denoms: %v", gravityDenoms)
		s.T().Logf("test erc20 contracts: %v", testERC20Contracts)

		for chainIndex, evm := range s.evms {
			s.T().Logf("checking balances for %s", evm.Name)
			s.Require().Eventuallyf(func() bool {
				res, err := bankQueryClient.AllBalances(context.Background(),
					&banktypes.QueryAllBalancesRequest{
						Address: s.chain.validators[1].keyInfo.GetAddress().String(),
					})
				if err != nil {
					return false
				}

				denomRes, err := gbQueryClient.ERC20ToDenom(context.Background(),
					&types.ERC20ToDenomRequest{
						Erc20:   testERC20Contracts[chainIndex].String(),
						ChainId: evm.ID,
					})
				if err != nil {
					s.T().Logf("error querying ERC20 denom %s, %e", testERC20Contracts[chainIndex].String(), err)
					return false
				}
				s.Require().False(denomRes.CosmosOriginated, "ERC20-originated denom marked as cosmos originated")
				gravityDenoms[chainIndex] = denomRes.Denom

				for _, coin := range res.Balances {
					if coin.Denom == gravityDenoms[chainIndex] && coin.Amount.Equal(sdk.NewInt(200)) {
						return true
					}
				}

				s.T().Logf("balance not found, received %v", res.Balances)

				return false
			}, 105*time.Second, 10*time.Second, "balance from %s never found on cosmos", evm.Name)

		}
		for chainIndex, evm := range s.evms {
			receiver, err := generateEthereumKey()
			s.Require().NoError(err)
			testReceivers[chainIndex] = common.HexToAddress(receiver.address)
			s.T().Logf("created evm receiver wallet: %s", testReceivers[chainIndex])

			beforeBalance, err := s.getEthTokenBalanceOf(s.evms[chainIndex], testReceivers[chainIndex], testERC20Contracts[chainIndex])
			s.Require().NoError(err, "error getting eth balance")

			s.T().Logf("sending to %s", evm.Name)
			sendToEthereumMsg := types.NewMsgSendToEVM(
				evm.ID,
				s.chain.validators[1].keyInfo.GetAddress(),
				testReceivers[chainIndex].String(),
				sdk.Coin{Denom: gravityDenoms[chainIndex], Amount: sdk.NewInt(100)},
				sdk.Coin{Denom: gravityDenoms[chainIndex], Amount: sdk.NewInt(1)},
			)

			s.Require().Eventuallyf(func() bool {
				val := s.chain.validators[1]
				keyring, err := val.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.keyInfo.GetAddress())
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

			s.T().Logf("verifying send to %s", evm.Name)
			s.Require().Eventuallyf(func() bool {
				balance, err := s.getEthTokenBalanceOf(evm, testReceivers[chainIndex], testERC20Contracts[chainIndex])
				s.Require().NoError(err, "error getting destination balance")

				if balance.LTE(*beforeBalance) {
					s.T().Logf("funds not received yet, dest balance: %s, starting balance: %s", balance.String(), beforeBalance.String())
					unbatchedSendToEVMS, err := gbQueryClient.UnbatchedSendToEVMs(context.Background(),
						&types.UnbatchedSendToEVMsRequest{
							SenderAddress: s.chain.validators[1].keyInfo.GetAddress().String(),
							Pagination:    nil,
							ChainId:       evm.ID,
						})
					s.Require().NoError(err)
					s.T().Logf("unbatched send to evms: %v", unbatchedSendToEVMS.SendToEvms)

					batchedSendToEVMS, err := gbQueryClient.BatchedSendToEVMs(context.Background(),
						&types.BatchedSendToEVMsRequest{
							SenderAddress: s.chain.validators[1].keyInfo.GetAddress().String(),
							ChainId:       evm.ID,
						})
					s.Require().NoError(err)
					s.T().Logf("batched send to evms: %v", batchedSendToEVMS.SendToEvms)

					for _, val := range s.chain.validators {
						unsignedBatchTxs, err := gbQueryClient.UnsignedBatchTxs(context.Background(),
							&types.UnsignedBatchTxsRequest{
								Address: val.keyInfo.GetAddress().String(),
								ChainId: evm.ID,
							})
						s.Require().NoError(err)
						s.T().Logf("unsigned batches for val %s: %v", val.keyInfo.GetAddress().String(), unsignedBatchTxs.Batches)
					}

					batchTxs, err := gbQueryClient.BatchTxs(context.Background(),
						&types.BatchTxsRequest{
							ChainId: evm.ID,
						})
					s.Require().NoError(err)
					s.T().Logf("batches: %v", batchTxs.Batches)

					return false
				}

				s.Require().Equal(sdk.NewInt(100).BigInt(), balance.BigInt(), "balance was %s, expected 100", balance.String())
				return true
			}, time.Second*240, time.Second, "send to ethereum did not reach destination")
		}

		s.T().Logf("funding community pool")
		orch := s.chain.orchestrators[0]
		clientCtx, err = s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
		s.Require().NoError(err)

		fundCommunityPoolMsg := distrtypes.NewMsgFundCommunityPool(
			sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(1000000000))),
			orch.keyInfo.GetAddress(),
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

		poolRes, err := distrQueryClient.CommunityPool(context.Background(),
			&distrtypes.QueryCommunityPoolRequest{},
		)
		s.Require().NoError(err, "error retrieving community pool")
		s.Require().True(poolRes.Pool.AmountOf(testDenom).GT(sdk.NewDec(1000000000)))

		for chainIndex, evm := range s.evms {
			s.T().Logf("deploying testgb as an ERC20")
			paramsRes, err := gbQueryClient.DenomToERC20Params(context.Background(),
				&types.DenomToERC20ParamsRequest{
					Denom:   testDenom,
					ChainId: evm.ID,
				})
			s.Require().NoError(err, "error retrieving ERC20 params for testgb denom")

			err = s.deployERC20(evm, gravityContracts[chainIndex], paramsRes.BaseDenom, paramsRes.Erc20Name, paramsRes.Erc20Symbol, uint8(paramsRes.Erc20Decimals))
			s.Require().NoError(err, "error deploying testgb as an ERC20")

			s.Require().Eventuallyf(func() bool {
				erc20Res, err := gbQueryClient.DenomToERC20(context.Background(),
					&types.DenomToERC20Request{
						Denom:   testDenom,
						ChainId: evm.ID,
					},
				)
				if err != nil {
					s.T().Logf("erc20 not deployed yet, waiting. err: %e", err)
					return false
				}

				s.Require().True(erc20Res.CosmosOriginated)
				return true
			}, 180*time.Second, 10*time.Second, "unable to verify ERC20 deployment")

			erc20Res, err := gbQueryClient.DenomToERC20(context.Background(),
				&types.DenomToERC20Request{
					Denom:   testDenom,
					ChainId: evm.ID,
				},
			)
			s.Require().NoError(err, "error querying ERC20 for testgb denom")
			testERC20Contracts[chainIndex] = common.HexToAddress(erc20Res.Erc20)

			communitySpendReceiver, err := generateEthereumKey()
			s.Require().NoError(err)
			communitySpendReceivers[chainIndex] = common.HexToAddress(communitySpendReceiver.address)

			s.T().Logf("create governance proposal to fund an ethereum address")
			orch = s.chain.orchestrators[0]
			clientCtx, err = s.chain.clientContext("tcp://localhost:26657", orch.keyring, "orch", orch.keyInfo.GetAddress())
			s.Require().NoError(err)

			proposal := types.CommunityPoolEVMSpendProposal{
				Title:       fmt.Sprintf("community pool spend %s", evm.Name),
				Description: fmt.Sprintf("community pool spend %s", evm.Name),
				Recipient:   communitySpendReceivers[chainIndex].String(),
				Amount:      sdk.NewCoin(testDenom, sdk.NewInt(900000000)),
				BridgeFee:   sdk.NewCoin(testDenom, sdk.NewInt(1000000)),
				ChainId:     evm.ID,
			}

			proposalMsg, err := govtypes.NewMsgSubmitProposal(
				&proposal,
				sdk.Coins{
					{
						Denom:  testDenom,
						Amount: sdk.NewInt(2),
					},
				},
				orch.keyInfo.GetAddress(),
			)
			s.Require().NoError(err, "unable to create governance proposal")

			s.T().Log("submit proposal spending community pool funds")
			submitProposalResponse, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
			s.Require().NoError(err)
			s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

			s.T().Log("check proposal was submitted correctly")
			proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
			s.Require().NoError(err)
			s.Require().NotEmpty(proposalsQueryResponse.Proposals)
			s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "not proposal id 1")
			s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")

			s.T().Log("vote for community spend proposal")
			for _, val := range s.chain.validators {
				kr, err := val.keyring()
				s.Require().NoError(err)
				clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kr, "val", val.keyInfo.GetAddress())
				s.Require().NoError(err)

				voteMsg := govtypes.NewMsgVote(val.keyInfo.GetAddress(), 1, govtypes.OptionYes)
				voteResponse, err := s.chain.sendMsgs(*clientCtx, voteMsg)
				s.Require().NoError(err)
				s.Require().Zero(voteResponse.Code, "vote error: %s", voteResponse.RawLog)
			}
		}

		s.T().Log("wait for community spend proposal to be approved")
		for chainIndex, _ := range s.evms {
			s.Require().Eventuallyf(func() bool {
				proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: uint64(chainIndex + 1)})
				s.Require().NoError(err)
				return govtypes.StatusPassed == proposalQueryResponse.Proposal.Status
			}, time.Second*30, time.Second*5, "proposal was never accepted")
		}

		s.T().Log("waiting for community funds to reach destination")
		for chainIndex, evm := range s.evms {

			s.Require().Eventuallyf(func() bool {
				balance, err := s.getEthTokenBalanceOf(evm, communitySpendReceivers[chainIndex], testERC20Contracts[chainIndex])
				s.Require().NoError(err, "error getting destination balance")

				if balance.LT(sdk.NewInt(900000000)) {
					s.T().Logf("funds not received yet, dest balance: %s", balance.String())
					return false
				}

				s.Require().Equal(balance.BigInt(), sdk.NewInt(900000000).BigInt(), "balance was %s, expected 900000000", balance.String())
				return true
			}, time.Second*180, time.Second*10, "community funds did not reach destination")
		}
	})
}
