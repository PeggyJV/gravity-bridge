package gravity

import (
	"fmt"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/client/cli"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/keeper"
	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	"github.com/spf13/cobra"
)

// NewHandler returns a handler for "Gravity" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgSendToEVM:
			res, err := msgServer.SendToEVM(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgCancelSendToEVM:
			res, err := msgServer.CancelSendToEVM(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgRequestBatchTx:
			res, err := msgServer.RequestBatchTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSubmitEVMTxConfirmation:
			res, err := msgServer.SubmitEVMTxConfirmation(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgSubmitEVMEvent:
			res, err := msgServer.SubmitEVMEvent(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgDelegateKeys:
			res, err := msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		case *types.MsgEVMHeightVote:
			res, err := msgServer.SubmitEVMHeightVote(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func CmdSubmitCommunityPoolEVMSpendProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "community-pool-evm-spend [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool EVM spend proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a community pool EVM spend proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. The funds from the community pool
will be bridged to the EVM chain to the supplied recipient EVM address. Only one denomination
of Cosmos token can be sent, and the bridge fee supplied along with the amount must be of the
same denomination.

Example:
$ %s tx gov submit-proposal community-pool-evm-spend <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
	"title": "Community Pool EVM Spend",
	"description": "Bridge me some tokens to the EVM chain!",
	"recipient": "0x0000000000000000000000000000000000000000",
	"amount": "20000stake"
	"bridge_fee": "1000stake",
	"deposit": "1000stake",
	"chain_id": 1
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, err := cli.ParseCommunityPoolEVMSpendProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			if len(proposal.Title) == 0 {
				return fmt.Errorf("title is empty")
			}

			if len(proposal.Description) == 0 {
				return fmt.Errorf("description is empty")
			}

			if !common.IsHexAddress(proposal.Recipient) {
				return fmt.Errorf("recipient is not a valid EVM address")
			}

			amount, err := sdk.ParseCoinNormalized(proposal.Amount)
			if err != nil {
				return err
			}

			bridgeFee, err := sdk.ParseCoinNormalized(proposal.BridgeFee)
			if err != nil {
				return err
			}

			if amount.Denom != bridgeFee.Denom {
				return fmt.Errorf("amount and bridge fee denominations must match")
			}

			deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			content := types.NewCommunityPoolEVMSpendProposal(proposal.Title, proposal.ChainId, proposal.Description, proposal.Recipient, amount, bridgeFee)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
