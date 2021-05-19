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
		CmdSignerSetTx(),
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

			queryClient := types.NewQueryClient(clientCtx)

			req := types.ParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				nonce uint64 // TODO(levi) init and validate from args[0]
			)

			req := types.SignerSetTxRequest{
				Nonce: nonce,
			}

			res, err := queryClient.SignerSetTx(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
