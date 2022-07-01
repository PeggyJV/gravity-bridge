package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func GetTxCmd(storeKey string) *cobra.Command {
	gravityTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Gravity transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	gravityTxCmd.AddCommand(
		CmdSendToEVM(),
		CmdCancelSendToEVM(),
		CmdRequestBatchTx(),
		CmdSetDelegateKeys(),
	)

	return gravityTxCmd
}

func CmdSendToEVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "send-to-evm [chain-id] [evm-reciever] [send-coins] [fee-coins]",
		Aliases: []string{"send", "transfer"},
		Args:    cobra.ExactArgs(4),
		Short:   "Send tokens from cosmos chain to connected evm chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			if from == nil {
				return fmt.Errorf("must pass from flag")
			}

			chainID := uint32(sdk.NewUintFromString(args[0]).Uint64())

			if !common.IsHexAddress(args[1]) {
				return fmt.Errorf("must be a valid evm address got %s", args[1])
			}
			receiver := common.HexToAddress(args[1])

			// Get amount of coins
			sendCoin, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			feeCoin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgSendToEVM(chainID, from, receiver.Hex(), sendCoin, feeCoin)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdCancelSendToEVM() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-send-to-evm [id] [chain-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel evm send by id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			if from == nil {
				return fmt.Errorf("must pass from flag")
			}

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			chainID := uint32(sdk.NewUintFromString(args[0]).Uint64())

			msg := types.NewMsgCancelSendToEVM(uint32(chainID), id, from)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdRequestBatchTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "request-batch-tx [chain-id] [denom] [signer]",
		Args:  cobra.ExactArgs(3),
		Short: "Request batch transaction for denom by signer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainID := uint32(sdk.NewUintFromString(args[0]).Uint64())

			denom := args[1]
			signer, err := sdk.AccAddressFromHex(args[2])
			if err != nil {
				return err
			}

			msg := types.NewMsgRequestBatchTx(chainID, denom, signer)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdSetDelegateKeys() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-delegate-keys [validator-address] [orchestrator-address] [evm-address] [evm-signature]",
		Args:  cobra.ExactArgs(4),
		Short: "Set gravity delegate keys",
		Long: `Set a validator's EVM and orchestrator addresses. The validator must
sign over a binary Proto-encoded DelegateKeysSignMsg message. The message contains
the validator's address and operator account current nonce.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			orcAddr, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			ethAddr, err := parseContractAddress(args[2])
			if err != nil {
				return err
			}

			ethSig, err := hexutil.Decode(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgDelegateKeys(valAddr, orcAddr, ethAddr, ethSig)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
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
	"amount": "20000stake",
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

			proposal, err := ParseCommunityPoolEVMSpendProposal(clientCtx.Codec, args[0])
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
