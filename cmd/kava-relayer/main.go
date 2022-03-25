package main

import (
	"fmt"
	"os"
	"strings"
)

var cmds = []string{}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s %s\n", os.Args[0], strings.Join(cmds, "|"))
		os.Exit(1)
	}

	switch cmd := os.Args[1]; cmd {
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
