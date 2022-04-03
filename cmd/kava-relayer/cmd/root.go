package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	appDir string
)

const (
	defaultAppDirName = ".kava-relayer"
	defaultConfigFile = "config"
	defaultConfigType = "yaml"
)

func NewRootCmd() (*cobra.Command, error) {
	cobra.OnInitialize(initConfig)

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	defaultAppDir := path.Join(userHomeDir, defaultAppDirName)

	rootCmd := &cobra.Command{
		Use:   "kava-relayer",
		Short: "Kava-relayer relays funds between ethereum and kava",
		Long:  `The kava relayer processes ethereum and kava blocks to transfer ERC20 tokens between chains.`,
	}

	// TODO: allow configuration of config and data separately
	rootCmd.PersistentFlags().StringVar(&appDir, "home", defaultAppDir, "Directory for config and data")

	rootCmd.AddCommand(newStartCmd())

	return rootCmd, nil
}

func initConfig() {
	viper.AddConfigPath(appDir)
	viper.SetConfigName(defaultConfigFile)
	viper.SetConfigType(defaultConfigType)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			cobra.CheckErr(err)
		}
	}
}
