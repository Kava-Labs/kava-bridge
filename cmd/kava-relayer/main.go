package main

import (
	"fmt"
	"os"

	logging "github.com/ipfs/go-log/v2"

	"github.com/kava-labs/kava-bridge/cmd/kava-relayer/cmd"
)

func main() {
	rootCmd, err := cmd.NewRootCmd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Set the default log level to info, default level of go-log is debug
	if _, present := os.LookupEnv("GOLOG_LOG_LEVEL"); !present {
		logging.SetAllLoggers(logging.LevelInfo)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
