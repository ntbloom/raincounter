package cli

import "github.com/spf13/cobra"

// RootCmd is the base of all command-line arguments
var RootCmd = &cobra.Command{} //nolint:gochecknoglobals

// Configure sets global CLI configs
func Configure() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
}

// AddSubcommand adds a subcommand to the CLI
func AddSubcommand(command string, short string, callable func()) {
	RootCmd.AddCommand(&cobra.Command{Use: command, Short: short, Run: func(_ *cobra.Command, _ []string) {
		callable()
	}})
}
