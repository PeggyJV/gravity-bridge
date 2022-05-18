package cli

import (
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

// ParseCommunityPoolEVMSpendProposal reads and parses a CommunityPoolEVMSpendProposalForCLI from a file.
func ParseCommunityPoolEVMSpendProposal(cdc codec.JSONCodec, proposalFile string) (types.CommunityPoolEVMSpendProposalForCLI, error) {
	proposal := types.CommunityPoolEVMSpendProposalForCLI{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
