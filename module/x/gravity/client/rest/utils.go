package rest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type (
	// CommunityPoolEthereumSpendProposalReq defines a community pool Ethereum spend proposal request body.
	CommunityPoolEthereumSpendProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Recipient   string         `json:"recipient" yaml:"recipient"`
		Amount      sdk.Coin       `json:"amount" yaml:"amount"`
		BridgeFee   sdk.Coin       `json:"bridge_fee" yaml:"bridge_fee"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.Coins      `json:"deposit" yaml:"deposit"`
	}
)
