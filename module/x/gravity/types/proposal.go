package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ProposalTypeCommunityPoolEVMSpend defines the type for a CommunityPoolEVMSpendProposal
	ProposalTypeCommunityPoolEVMSpend = "CommunityPoolEVMSpend"
)

// Assert CommunityPoolEVMSpendProposal implements govtypes.Content at compile-time
var _ govtypes.Content = &CommunityPoolEVMSpendProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeCommunityPoolEVMSpend)
	govtypes.RegisterProposalTypeCodec(&CommunityPoolEVMSpendProposal{}, "gravity/CommunityPoolEVMSpendProposal")
}

// NewCommunityPoolEVMSpendProposal creates a new community pool spend proposal.
//nolint:interfacer
func NewCommunityPoolEVMSpendProposal(title string, chainID uint32, description string, recipient string, amount sdk.Coin, bridgeFee sdk.Coin) *CommunityPoolEVMSpendProposal {
	return &CommunityPoolEVMSpendProposal{title, description, recipient, amount, bridgeFee, chainID}
}

// GetTitle returns the title of a community pool EVM spend proposal.
func (csp *CommunityPoolEVMSpendProposal) GetTitle() string { return csp.Title }

// GetDescription returns the description of a community pool EVM spend proposal.
func (csp *CommunityPoolEVMSpendProposal) GetDescription() string { return csp.Description }

// GetDescription returns the routing key of a community pool EVM spend proposal.
func (csp *CommunityPoolEVMSpendProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool EVM spend proposal.
func (csp *CommunityPoolEVMSpendProposal) ProposalType() string {
	return ProposalTypeCommunityPoolEVMSpend
}

// ValidateBasic runs basic stateless validity checks
func (csp *CommunityPoolEVMSpendProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}

	if !common.IsHexAddress(csp.Recipient) {
		return ErrInvalidEVMProposalRecipient
	}

	if !csp.Amount.IsValid() || csp.Amount.IsZero() {
		return ErrInvalidEVMProposalAmount
	}

	if !csp.BridgeFee.IsValid() {
		return ErrInvalidEVMProposalBridgeFee
	}

	if csp.Amount.Denom != csp.BridgeFee.Denom {
		return ErrEVMProposalDenomMismatch
	}

	return nil
}

// String implements the Stringer interface.
func (csp CommunityPoolEVMSpendProposal) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`Community Pool Spend Proposal:
  Title:       %s
  Description: %s
  Recipient:   %s
  Amount:      %s
  Bridge Fee:  %s
`, csp.Title, csp.Description, csp.Recipient, csp.Amount, csp.BridgeFee))
	return b.String()
}
