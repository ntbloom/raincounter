package main

import (
	"fmt"
	"os"

	"github.com/ntbloom/raincounter/pkg/gateway"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{} //nolint:gochecknoglobals

func addSubcommand(command string, short string, callable func()) {
	rootCmd.AddCommand(&cobra.Command{Use: command, Short: short, Run: func(_ *cobra.Command, _ []string) {
		callable()
	}})
}

func main() {
	addSubcommand("gateway", "shuffle data from sensor to mqtt on the gateway", gateway.Start)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
