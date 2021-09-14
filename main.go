package main

import (
	"fmt"
	"os"

	raincloud "github.com/ntbloom/raincounter/pkg/server"

	"github.com/ntbloom/raincounter/cli"

	"github.com/ntbloom/raincounter/pkg/gateway"
)

func main() {
	cli.Configure()
	cli.AddSubcommand("gateway", "shuffle data from sensor to MQTT on the gateway", gateway.Start)
	cli.AddSubcommand("receiver", "receive data over MQTT on the server", raincloud.Start)
	if err := cli.RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

}
