package cli

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/kava-labs/kava-bridge/contract"
	"github.com/kava-labs/kava-bridge/x/bridge/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	bridgeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmds := []*cobra.Command{
		getCmdMsgBridgeEthereumToKava(),
		getCmdMsgConvertCoinToERC20(),
		getCmdBridgeKavaToEthereum(),
		getCmdConvertERC20ToCoin(),
	}

	for _, cmd := range cmds {
		flags.AddTxFlagsToCmd(cmd)
	}

	bridgeTxCmd.AddCommand(cmds...)

	return bridgeTxCmd
}

func getCmdMsgBridgeEthereumToKava() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-eth-to-kava [token] [receiver] [amount] [sequence]",
		Short: "mints ERC20 tokens locked on Ethereum to Kava EVM co-chain",
		Example: fmt.Sprintf(
			`%s tx %s bridge-eth-to-kava 0xc778417e063141139fce010982780140aa0cd5ab 0x6B1088f788b412Ad1280F95240d56B886A64bc05 1000000000000000 1 --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			token := args[0]
			receiver := args[1]

			amount, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("amount '%s' is an invalid int", args[2])
			}

			sequence, ok := sdk.NewIntFromString(args[3])
			if !ok {
				return fmt.Errorf("amount '%s' is an invalid int", args[2])
			}

			signer := clientCtx.GetFromAddress()
			msg := types.NewMsgBridgeEthereumToKava(signer.String(), token, amount, receiver, sequence)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

func getCmdMsgConvertCoinToERC20() *cobra.Command {
	return &cobra.Command{
		Use:   "convert-coin-to-erc20 [receiver] [coin]",
		Short: "converts sdk.Coin to erc20 tokens on Kava eth co-chain",
		Example: fmt.Sprintf(
			`%s tx %s convert-coin-to-erc20 0x6B1088f788b412Ad1280F95240d56B886A64bc05 100000000weth --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			receiver := args[0]
			if !common.IsHexAddress(receiver) {
				return fmt.Errorf("receiver '%s' is an invalid hex address", args[0])
			}

			coin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			signer := clientCtx.GetFromAddress()
			msg := types.NewMsgConvertCoinToERC20(signer.String(), receiver, coin)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}
}

func getCmdBridgeKavaToEthereum() *cobra.Command {
	return &cobra.Command{
		Use:   "bridge-kava-to-eth [Ethereum receiver address] [Kava ERC20 address] [amount]",
		Short: "burns ERC20 tokens on Kava EVM co-chain and unlocks on Ethereum",
		Example: fmt.Sprintf(
			`%s tx %s bridge-kava-to-eth 0x21E360e198Cde35740e88572B59f2CAdE421E6b1 0x8223259205A3E31C54469fCbfc9F7Cf83D515ff6 1000000000000000 --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if !common.IsHexAddress(args[0]) {
				return fmt.Errorf("receiver '%s' is an invalid hex address", args[1])
			}
			receiver := common.HexToAddress(args[0])

			if !common.IsHexAddress(args[1]) {
				return fmt.Errorf("contract address '%s' is an invalid hex address", args[0])
			}
			contractAddr := common.HexToAddress(args[1])

			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				return fmt.Errorf("amount '%s' is invalid", args[2])
			}

			data, err := PackContractCallData(
				contract.ERC20MintableBurnableContract.ABI,
				"withdraw",
				receiver,
				amount,
			)
			if err != nil {
				return err
			}

			ethTx, err := CreateEthCallContractTx(
				clientCtx,
				&contractAddr,
				data,
			)
			if err != nil {
				return err
			}

			return GenerateOrBroadcastTx(clientCtx, ethTx)
		},
	}
}

func getCmdConvertERC20ToCoin() *cobra.Command {
	return &cobra.Command{
		Use:   "convert-erc20-to-coin [Kava receiver address] [Kava ERC20 address or Denom] [amount]",
		Short: "burns ERC20 tokens on Kava EVM co-chain and unlocks on Ethereum",
		Example: fmt.Sprintf(
			`%s tx %s convert-erc20-to-coin 0x8223259205A3E31C54469fCbfc9F7Cf83D515ff6 0x21E360e198Cde35740e88572B59f2CAdE421E6b1 1000000000000000 --from <key>`,
			version.AppName, types.ModuleName,
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Support both bech32 and hex address, try to parse as bech32 if is
			// not a hex address
			var receiver common.Address
			if !common.IsHexAddress(args[0]) {
				accAddr, err := sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return fmt.Errorf("receiver '%s' is not a hex or bech32 address", args[1])
				}

				receiver = common.BytesToAddress(accAddr)
			} else {
				receiver = common.HexToAddress(args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)

			var contractAddr common.Address
			if !common.IsHexAddress(args[1]) {
				if err := sdk.ValidateDenom(args[1]); err != nil {
					return fmt.Errorf("Kava ERC20 '%s' is not a valid hex address or denom", args[0])
				}

				// Valid denom, try looking up as denom to get corresponding Kava ERC20 address
				paramsRes, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
				if err != nil {
					return err
				}

				found := false
				for _, pair := range paramsRes.Params.EnabledConversionPairs {
					if pair.Denom == args[1] {
						contractAddr = pair.GetAddress().Address
						found = true
						break
					}
				}

				if !found {
					return fmt.Errorf("Kava ERC20 '%s' is not a valid hex address or denom", args[0])
				}
			} else {
				contractAddr = common.HexToAddress(args[1])
			}

			amount, ok := new(big.Int).SetString(args[2], 10)
			if !ok {
				return fmt.Errorf("amount '%s' is invalid", args[2])
			}

			data, err := PackContractCallData(
				contract.ERC20MintableBurnableContract.ABI,
				"convertToCoin",
				receiver,
				amount,
			)
			if err != nil {
				return err
			}

			ethTx, err := CreateEthCallContractTx(
				clientCtx,
				&contractAddr,
				data,
			)
			if err != nil {
				return err
			}

			return GenerateOrBroadcastTx(clientCtx, ethTx)
		},
	}
}
