package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/gravity-bridge/module/v2/x/gravity/client/cli"
)

// ProposalHandler is the community Ethereum spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitCommunityPoolEthereumSpendProposal)
)
