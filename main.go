package main

import (
	"fmt"
	"os"

	"github.com/ntbloom/raincounter/cli"

	"github.com/ntbloom/raincounter/pkg/gateway"
)

func main() {
	cli.Configure()
	cli.AddSubcommand("gateway", "shuffle data from sensor to mqtt on the gateway", gateway.Start)
	if err := cli.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	//cli.AddSubcommand("--insecure", "disable TLS and connect over 1883 (for development only)", func())
}
