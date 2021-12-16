package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ntbloom/raincounter/cli"
	"github.com/ntbloom/raincounter/pkg/config"
	"github.com/ntbloom/raincounter/pkg/rainbase"
	"github.com/ntbloom/raincounter/pkg/raincloud"
)

func main() {
	cli.Configure()
	cli.AddSubcommand("rainbase", "shuffle data from sensor to MQTT on the rainbase", rainbase.Start)
	cli.AddSubcommand("receiver", "receive data over MQTT on the cloud", raincloud.Receive)
	cli.AddSubcommand("server", "serve the rest API on the cloud", raincloud.Serve)

	cli.RootCmd.PersistentFlags().StringVar(&config.RegularFile, "config", "", "config file")
	cobra.OnInitialize(config.Configure)

	if err := cli.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
