package cmd

import (
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/cmd/network"
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
	rootCmd.AddCommand(network.GetNetworkCmd())

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

	// Set the default log level to info and allow config via viper, default level of go-log is debug
	viper.SetDefault("log_level", "info")
	logLevelStr := viper.GetString("log_level")

	// Set GOLOG_LOG_LEVEL to override the default log level, loggers are not
	// manually set since this allows for per subsystem log levels
	err := os.Setenv("GOLOG_LOG_LEVEL", logLevelStr)
	cobra.CheckErr(err)
}
