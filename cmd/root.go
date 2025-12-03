package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

func SetVersionInfo(v, _, _ string) {
	version = v
	rootCmd.Version = v
}

var rootCmd = &cobra.Command{
	Use:     "sshy",
	Short:   "Manage SSH servers via YAML config",
	Version: version,
	Long:    `sshy is a CLI tool for managing SSH servers via YAML configuration.`,
}

var osExit = os.Exit

func Execute() {
	ExecuteWithArgs(os.Args)
}

func ExecuteWithArgs(args []string) {
	if len(args) == 1 || (len(args) > 1 && !isSubcommand(args[1]) && !isFlag(args[1])) {
		connectCmd.RunE(connectCmd, args[1:])
		return
	}

	if err := rootCmd.Execute(); err != nil {
		osExit(1)
	}
}

func isSubcommand(arg string) bool {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == arg {
			return true
		}
	}
	return false
}

func isFlag(arg string) bool {
	return len(arg) > 1 && arg[0] == '-'
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
