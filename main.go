package main

import (
	"fmt"
	"os"

	"github.com/ntbloom/raincounter/pkg/config"

	raincloud "github.com/ntbloom/raincounter/pkg/server"

	"github.com/ntbloom/raincounter/cli"

	"github.com/ntbloom/raincounter/pkg/rainbase"
)

func main() {
	config.Configure()
	cli.Configure()
	cli.AddSubcommand("rainbase", "shuffle data from sensor to MQTT on the rainbase", rainbase.Start)
	cli.AddSubcommand("receiver", "receive data over MQTT on the cloud", raincloud.Receive)
	cli.AddSubcommand("server", "serve the rest API on the cloud", raincloud.Serve)
	if err := cli.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

}
