package key

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newPrecomputePreParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "precompute-preparams",
		Short: "Generates the pre-params for keygen",
		Long:  "Pre-computes 2 safe primes and Paillier secret required for the protocol.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := viper.BindPFlags(cmd.Flags())
			cobra.CheckErr(err)

			homeDir := viper.GetString("home")
			if homeDir == "" {
				return fmt.Errorf("home directory must be set")
			}

			force := viper.GetBool(KeyFlagForce)

			preParamsFilePath := path.Join(homeDir, "pre-params.json")
			if !force {
				if _, err := os.Stat(preParamsFilePath); err == nil {
					return fmt.Errorf("pre-params file already exists: %s", preParamsFilePath)
				}
			}

			// When using the keygen party it is recommended that you pre-compute the
			// "safe primes" and Paillier secret beforehand because this can take some time.
			// This code will generate those parameters using a concurrency limit equal
			// to the number of available CPU cores.
			preParams, err := keygen.GeneratePreParams(1 * time.Minute)
			if err != nil {
				return fmt.Errorf("failed to generate pre-params: %s", err)
			}

			b, err := json.MarshalIndent(preParams, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to json marshal pre-params: %s", err)
			}

			// No perms for public
			err = os.WriteFile(preParamsFilePath, b, 0600)
			if err != nil {
				return fmt.Errorf("failed to write pre-params file: %s", err)
			}

			fmt.Printf("pre-params written to: %v\n", preParamsFilePath)

			return nil
		},
	}

	cmd.Flags().Bool(KeyFlagForce, false, "Overwrite existing pre-params file if it exists")

	return cmd
}
