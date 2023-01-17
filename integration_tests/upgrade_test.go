package integration_tests

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"time"
)

func (s *IntegrationTestSuite) TestMultiEVMUpgrade() {
	s.Run("bring up chain, test upgrade path", func() {
		var err error

		val := s.chain.validators[0]
		kb, err := val.keyring()
		s.Require().NoError(err)
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.keyInfo.GetAddress())
		s.Require().NoError(err)

		upgradeQueryClient := upgradetypes.NewQueryClient(clientCtx)
		govQueryClient := govtypes.NewQueryClient(clientCtx)

		// validate that there is no upgrade planned yet
		currentPlan, err := upgradeQueryClient.CurrentPlan(context.Background(), &upgradetypes.QueryCurrentPlanRequest{})
		s.Require().NoError(err)
		s.Require().Nil(currentPlan.Plan)

		rpcClient, err := rpchttp.New("tcp://localhost:26657", "/websocket")
		s.Require().NoError(err)
		status, err := rpcClient.Status(context.Background())
		s.Require().NoError(err)

		upgradeBlockHeight := status.SyncInfo.LatestBlockHeight + 30

		proposalMsg, err := govtypes.NewMsgSubmitProposal(
			&upgradetypes.SoftwareUpgradeProposal{
				Title:       "test multi-evm upgrade",
				Description: "test multi-evm upgrade description",
				Plan: upgradetypes.Plan{
					Name:   "multi-evm-upgrade",
					Height: upgradeBlockHeight,
				},
			},
			sdk.Coins{
				{
					Denom:  testDenom,
					Amount: sdk.NewInt(2),
				},
			},
			val.keyInfo.GetAddress(),
		)

		s.T().Log("submit proposal upgrading software version")
		submitProposalResponse, err := s.chain.sendMsgs(*clientCtx, proposalMsg)
		s.Require().NoError(err)
		s.Require().Zero(submitProposalResponse.Code, "raw log: %s", submitProposalResponse.RawLog)

		s.T().Log("check proposal was submitted correctly")
		proposalsQueryResponse, err := govQueryClient.Proposals(context.Background(), &govtypes.QueryProposalsRequest{})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposalsQueryResponse.Proposals)
		s.Require().Equal(govtypes.StatusVotingPeriod, proposalsQueryResponse.Proposals[0].Status, "proposal not in voting period")
		s.Require().Equal(uint64(1), proposalsQueryResponse.Proposals[0].ProposalId, "proposal id is not 1")

		s.T().Log("vote for upgrade proposal")
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

		s.T().Log("waiting for proposal to complete")
		s.Require().Eventuallyf(func() bool {
			proposalQueryResponse, err := govQueryClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: 1})
			s.Require().NoError(err)

			return proposalQueryResponse.Proposal.Status == govtypes.StatusPassed
		}, time.Minute, time.Second*10, "proposal failed to be accepted")

		s.T().Logf("waiting for upgrade block height %d", upgradeBlockHeight)
		s.Require().Eventuallyf(func() bool {
			status, err = rpcClient.Status(context.Background())
			s.Require().NoError(err)
			return status.SyncInfo.LatestBlockHeight >= upgradeBlockHeight
		}, time.Minute, time.Second*10, "upgrade height failed to reach")
	})
}
