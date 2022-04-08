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

	// TODO: Default logger is debug level, we want to have info default but also
	// allow changing via GOLOG_LOG_LEVEL env var.
	logging.SetAllLoggers(logging.LevelInfo)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
