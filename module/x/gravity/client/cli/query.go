package cli

import (
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
		CmdBatchTx(),
		CmdContractCallTx(),
		CmdSignerSetTxs(),
		CmdBatchTxs(),
		CmdContractCallTxs(),
		CmdSignerSetTxEthereumSignatures(),
		CmdBatchTxEthereumSignatures(),
		CmdContractCallTxEthereumSignatures(),
		CmdPendingSignerSetTxEthereumSignatures(),
		CmdPendingBatchTxEthereumSignatures(),
		CmdPendingContractCallTxEthereumSignatures(),
		CmdLastSubmittedEthereumEvent(),
		CmdBatchTxFees(),
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

func CmdBatchTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx [contract-address] [nonce]",
		Args:  cobra.ExactArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				contractAddress string // TODO(levi) init and validate from args[0]
				nonce           uint64 // TODO(levi) init and validate from args[1]
			)

			req := types.BatchTxRequest{
				ContractAddress: contractAddress,
				Nonce:           nonce,
			}

			res, err := queryClient.BatchTx(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-tx [invalidation-scope] [invalidation-nonce]",
		Args:  cobra.ExactArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				invalidationScope []byte // TODO(levi) init and validate from args[0]
				invalidationNonce uint64 // TODO(levi) init and validate from args[1]
			)

			req := types.ContractCallTxRequest{
				InvalidationScope: invalidationScope,
				InvalidationNonce: invalidationNonce,
			}

			res, err := queryClient.ContractCallTx(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTxs() *cobra.Command {
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

func CmdBatchTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-txs",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := types.BatchTxsRequest{}

			res, err := queryClient.BatchTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-txs",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := types.ContractCallTxsRequest{}

			res, err := queryClient.ContractCallTxs(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdSignerSetTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-set-tx-ethereum-signatures [nonce] [validator-or-orchestrator-address]",
		Args:  cobra.ExactArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				nonce   uint64 // TODO(levi) init and validate from args[0]
				address string // TODO(levi) init and validate from args[1]
			)

			req := types.SignerSetTxEthereumSignaturesRequest{
				Nonce:   nonce,
				Address: address,
			}

			res, err := queryClient.SignerSetTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx-ethereum-signatures [nonce] [contract-address] [validator-or-orchestrator-address]",
		Args:  cobra.MinimumNArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				nonce           uint64 // TODO(levi) init and validate from args[0]
				contractAddress string // TODO(levi) init and validate from args[1]
				address         string // TODO(levi) init and validate from args[2]
			)

			req := types.BatchTxEthereumSignaturesRequest{
				Nonce:           nonce,
				ContractAddress: contractAddress,
				Address:         address,
			}

			res, err := queryClient.BatchTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdContractCallTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "contract-call-tx-ethereum-signatures [invalidation-scope] [invalidation-nonce] [validator-or-orchestrator-address]",
		Args:  cobra.MinimumNArgs(2),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				invalidationScope []byte // TODO(levi) init and validate from args[0]
				invalidationNonce uint64 // TODO(levi) init and validate from args[1]
				address           string // TODO(levi) init and validate from args[2]
			)

			req := types.ContractCallTxEthereumSignaturesRequest{
				InvalidationNonce: invalidationNonce,
				InvalidationScope: invalidationScope,
				Address:           address,
			}

			res, err := queryClient.ContractCallTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdPendingSignerSetTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-signer-set-tx-ethereum-signatures [validator-or-orchestrator-address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				address string // TODO(levi) init and validate from args[0]
			)

			req := types.PendingSignerSetTxEthereumSignaturesRequest{
				Address: address,
			}

			res, err := queryClient.PendingSignerSetTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdPendingBatchTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-batch-tx-ethereum-signatures [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				address string // TODO(levi) init and validate from args[0]
			)

			req := types.PendingBatchTxEthereumSignaturesRequest{
				Address: address,
			}

			res, err := queryClient.PendingBatchTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdPendingContractCallTxEthereumSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-contract-call-tx-ethereum-signatures [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				address string // TODO(levi) init and validate from args[0]
			)

			req := types.PendingContractCallTxEthereumSignaturesRequest{
				Address: address,
			}

			res, err := queryClient.PendingContractCallTxEthereumSignatures(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdLastSubmittedEthereumEvent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-submitted-ethereum-event [address]",
		Args:  cobra.ExactArgs(1),
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			var ( // args
				address string // TODO(levi) init and validate from args[0]
			)

			req := types.LastSubmittedEthereumEventRequest{
				Address: address,
			}

			res, err := queryClient.LastSubmittedEthereumEvent(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdBatchTxFees() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch-tx-fees",
		Args:  cobra.NoArgs,
		Short: "", // TODO(levi) provide short description
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := types.BatchTxFeesRequest{}

			res, err := queryClient.BatchTxFees(cmd.Context(), &req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
