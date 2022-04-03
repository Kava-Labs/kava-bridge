package cmd

import (
	"github.com/kava-labs/kava-bridge/relayer"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts the relayer",
		Long:  "Starts processing blocks and relaying transactions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			// get config from env, flags, file, etc
			ethRpcUrl := viper.GetString("eth.rpc")
			ethBridgeAddr := viper.GetString("eth.bridge")
			kavaGrpcUrl := viper.GetString("kava.grpc")
			relayerMnemonic := viper.GetString("relayer-mnemonic")

			// create service, perform validations on parameters
			srv, err := relayer.NewService(ethRpcUrl, ethBridgeAddr, kavaGrpcUrl, relayerMnemonic)
			if err != nil {
				return err
			}

			// blocking unless an error occurs or context is cancelled, start syncing blocks and relaying outputs
			// TODO: add cancellation and signal handling, allowing for graceful termination
			//       (save state, broadcast, finishing signing, etc)
			return srv.Run()
		},
	}

	cmd.Flags().String("eth.rpc", "", "Ethereum JSON-RPC URL")
	cmd.Flags().String("eth.bridge", "", "Ethereum bridge contract address")
	cmd.Flags().String("kava.grpc", "", "Kava GRPC endpoint")
	cmd.Flags().String("relayer-mnemonic", "", "Relayer Mnemonic for signing transactions")

	return cmd
}
