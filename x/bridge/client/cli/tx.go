package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

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
		Short: "mints erc20 tokens locked on ethereum to kava eth co-chain",
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