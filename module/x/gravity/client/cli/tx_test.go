package cli

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseCommunityPoolEthereumSpendProposal(t *testing.T) {
	encodingConfig := params.MakeTestEncodingConfig()

	okJSON := testutil.WriteToNewTempFile(t, `
{
  "title": "Community Pool Ethereum Spend",
  "description": "Bridge me some tokens to Ethereum!",
  "recipient": "0x0000000000000000000000000000000000000000",
  "amount": "20000stake",
  "bridge_fee": "1000stake",
  "deposit": "1000stake"
}
`)

	proposal, err := ParseCommunityPoolEthereumSpendProposal(encodingConfig.Marshaler, okJSON.Name())
	require.NoError(t, err)

	require.Equal(t, "Community Pool Ethereum Spend", proposal.Title)
	require.Equal(t, "Bridge me some tokens to Ethereum!", proposal.Description)
	require.Equal(t, "0x0000000000000000000000000000000000000000", proposal.Recipient)
	require.Equal(t, "20000stake", proposal.Amount)
	require.Equal(t, "1000stake", proposal.BridgeFee)
	require.Equal(t, "1000stake", proposal.Deposit)
}
