package cmd

import (
	"os"
	"path"
	"strings"

	logging "github.com/ipfs/go-log/v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/cmd/network"
)

var (
	appDir   string
	logLevel string
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
	rootCmd.PersistentFlags().StringVar(&logLevel, "log_level", "info", "The logging level (trace|debug|info|warn|error|fatal|panic)")

	// Bind log_level flag so that it shows in config
	err = viper.BindPFlag("log_level", rootCmd.Flag("log_level"))
	cobra.CheckErr(err)

	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newInitCmd())
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

	// Set the default log level to info, default level of go-log is debug
	// TODO: Does not support log level per subsystem such as with `GOLOG_LOG_LEVEL`
	logLevel, err := logging.LevelFromString(logLevel)
	cobra.CheckErr(err)

	logging.SetAllLoggers(logLevel)
}
