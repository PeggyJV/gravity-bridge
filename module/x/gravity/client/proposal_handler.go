package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/client/cli"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/client/rest"
)

// ProposalHandler is the community EVM spend proposal handler.
var (
	ProposalHandler = govclient.NewProposalHandler(cli.CmdSubmitCommunityPoolEVMSpendProposal, rest.ProposalRESTHandler)
)
