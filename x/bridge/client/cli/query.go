package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/server/config"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/kava-labs/kava-bridge/contract"
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
		QueryERC20BridgePairCmd(),
		QueryConversionPairsCmd(),
		QueryConversionPairCmd(),
		QueryERC20BalanceOfCmd(),
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
		Use:   "params",
		Short: "Query the bridge module parameters",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s params",
			version.AppName, types.ModuleName,
		),
		Args: cobra.NoArgs,
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

// -----------------------------------------------------------------------------
// Bridge pair queries

// QueryERC20BridgePairsCmd queries the bridge module bridged ERC20 pairs
func QueryERC20BridgePairsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "erc20-pairs",
		Short: "Query all bridged ERC20 pairs",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s erc20-pairs",
			version.AppName, types.ModuleName,
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ERC20BridgePairs(
				context.Background(),
				&types.QueryERC20BridgePairsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// QueryERC20BridgePairCmd queries the bridge module for a bridged ERC20 pair
func QueryERC20BridgePairCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "erc20-pair [Ethereum or Kava address]",
		Short: "Query a bridged ERC20 pair by Ethereum or Kava address",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s erc20-pair 0x404F9466d758eA33eA84CeBE9E444b06533b369e",
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid hex address: %v", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ERC20BridgePair(
				context.Background(),
				&types.QueryERC20BridgePairRequest{
					Address: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion pair queries

// QueryConversionPairsCmd queries the bridge module conversion ERC20/sdk.Coin pairs
func QueryConversionPairsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "conversion-pairs",
		Short: "Query all ERC20 / Coin conversion pairs",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s conversion-pairs",
			version.AppName, types.ModuleName,
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ConversionPairs(
				context.Background(),
				&types.QueryConversionPairsRequest{},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// QueryConversionPairCmd queries the bridge module for a conversion pair
func QueryConversionPairCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "conversion-pair [ERC20 contract address or denom]",
		Short: "Query a conversion pair by ERC20 contract address or denom",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s conversion-pair 0x404F9466d758eA33eA84CeBE9E444b06533b369e",
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) || sdk.ValidateDenom(args[0]) != nil {
				return fmt.Errorf("invalid hex address or denom: %v", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ConversionPair(
				context.Background(),
				&types.QueryConversionPairRequest{
					AddressOrDenom: args[0],
				},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
}

// -----------------------------------------------------------------------------
// ERC20 queries

// QueryERC20BridgePairsCmd queries the bridge module bridged ERC20 pairs
func QueryERC20BalanceOfCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "erc20-balance [contract address] [account address]",
		Short: "Query the balance of a ERC20 token",
		Example: fmt.Sprintf(
			"%[1]s q %[2]s erc20-balance 0x404F9466d758eA33eA84CeBE9E444b06533b369e 0x7Bbf300890857b8c241b219C6a489431669b3aFA",
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			evmQueryClient := evmtypes.NewQueryClient(clientCtx)

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("invalid contract address: %v", args[0])
			}

			if !common.IsHexAddress(args[1]) {
				return fmt.Errorf("invalid account address: %v", args[0])
			}

			contractAddr := common.HexToAddress(args[0])
			accountAddr := common.HexToAddress(args[1])

			data, err := contract.ERC20MintableBurnableContract.ABI.Pack(
				"balanceOf",
				accountAddr,
			)
			if err != nil {
				return fmt.Errorf("failed to pack balanceOf data: %w", err)
			}

			transactionArgs := evmtypes.TransactionArgs{
				To:   &contractAddr,
				Data: (*hexutil.Bytes)(&data),
			}

			ethCalArgs, err := json.Marshal(transactionArgs)
			if err != nil {
				return err
			}

			res, err := evmQueryClient.EthCall(context.Background(), &evmtypes.EthCallRequest{
				Args:   ethCalArgs,
				GasCap: uint64(config.DefaultGasCap),
			})
			if err != nil {
				return err
			}

			if res.Failed() {
				// No revert handling, not making a tx
				return fmt.Errorf(res.VmError)
			}

			anyOutput, err := contract.ERC20MintableBurnableContract.ABI.Unpack("balanceOf", res.Ret)
			if err != nil {
				return fmt.Errorf(
					"failed to unpack method %v response: %w",
					"balanceOf",
					err,
				)
			}

			bal, ok := anyOutput[0].(*big.Int)
			if !ok {
				return fmt.Errorf(
					"invalid ERC20 return type %T, expected %T",
					anyOutput[0],
					&big.Int{},
				)
			}

			return clientCtx.PrintString(fmt.Sprintf("%v\n", bal))
		},
	}
}
