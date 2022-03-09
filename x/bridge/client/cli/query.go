package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group bep3 queries under a subcommand
	bridgeQueryCmd := &cobra.Command{
		Use:                        "bridge",
		Short:                      "Querying commands for the bridge module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmds := []*cobra.Command{
		QueryParamsCmd(),
		QueryERC20BridgePairsCmd(),
	}

	for _, cmd := range cmds {
		flags.AddQueryFlagsToCmd(cmd)
	}

	bridgeQueryCmd.AddCommand(cmds...)

	return bridgeQueryCmd
}

// QueryParamsCmd queries the bridge module parameters
func QueryParamsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "params",
		Short:   "get the bridge module parameters",
		Example: "bridge params",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(&res.Params)
		},
	}
}

// QueryERC20BridgePairsCmd queries the bridge module bridged ERC20 pairs
func QueryERC20BridgePairsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "erc20-pairs",
		Short:   "get the bridge module bridged ERC20 pairs",
		Example: "bridge erc20-pairs",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.BridgedERC20Pairs(context.Background(), &types.QueryBridgedERC20PairsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}
