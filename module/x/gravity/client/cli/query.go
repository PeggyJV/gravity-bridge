package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/gravity-bridge/module/x/gravity/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	gravityQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the gravity module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	gravityQueryCmd.AddCommand(
		CmdParams(),
		// CmdSignerSetTx(),
		// CmdBatchTx(),
		// CmdContractCallTx(),
		// CmdSignerSetTxs(),
		// CmdBatchTxs(),
		// CmdContractCallTxs(),
		// CmdSignerSetTxEthereumSignatures(),
		// CmdBatchTxEthereumSignatures(),
		// CmdContractCallTxEthereumSignatures(),
		// CmdPendingSignerSetTxEthereumSignatures(),
		// CmdPendingBatchTxEthereumSignatures(),
		// CmdPendingContractCallTxEthereumSignatures(),
		// CmdLastSubmittedEthereumEvent(),
		// CmdBatchTxFees(),
		// CmdERC20ToDenom(),
		// CmdDenomToERC20(),
		// CmdPendingSendToEthereums(),
		// CmdDelegateKeysByValidator(),
		// CmdDelegateKeysByEthereumSigner(),
		// CmdDelegateKeysByOrchestrator(),
	)

	return gravityQueryCmd
}

func CmdParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query votes on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(clientCtx).Params(context.Background(), &types.ParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
