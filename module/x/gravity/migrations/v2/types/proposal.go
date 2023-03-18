package types

import (
	"fmt"
	"strings"
)

// String implements the Stringer interface.
func (csp CommunityPoolEthereumSpendProposal) String() string {
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
