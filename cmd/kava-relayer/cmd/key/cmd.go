package key

import "github.com/spf13/cobra"

func GetKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Manage relayer key",
	}

	cmds := []*cobra.Command{
		newPrecomputePreParamsCmd(),
	}

	cmd.AddCommand(cmds...)

	return cmd
}
