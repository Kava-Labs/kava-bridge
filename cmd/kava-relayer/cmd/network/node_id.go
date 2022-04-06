package network

import (
	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/p2p"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newShowNodeIdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-node-id",
		Short: "Shows the node id",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			opts := p2p.ParseOptions()
			node, err := p2p.NewNode(opts...)
			if err != nil {
				return err
			}

			node.Close()

			return nil
		},
	}

	return cmd
}
