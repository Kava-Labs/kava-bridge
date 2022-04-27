package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes the configuration file",
		Long:  "Creates the configuration file if it doesn't already exist.",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(appDir)
			if !os.IsExist(err) {
				if err := os.Mkdir(appDir, os.ModePerm); err != nil {
					return err
				}
			}

			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("could not write config: %w", err)
			}

			return nil
		},
	}

	return cmd
}
