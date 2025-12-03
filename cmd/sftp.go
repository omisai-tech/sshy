package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var sftpCmd = &cobra.Command{
	Use:   "sftp [ssh-flags...] <name>",
	Short: "Start SFTP session with SSH server",
	Long:  "Start an interactive SFTP session with the specified SSH server. SSH flags can be passed through.",
	Args:  cobra.MinimumNArgs(0), // Allow any number of args, we'll parse them
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle help flag manually since we need to allow unknown flags
		for _, arg := range args {
			if arg == "--help" || arg == "-h" {
				return cmd.Help()
			}
		}

		if len(args) < 1 {
			return fmt.Errorf("usage: sshy sftp [ssh-flags...] <name>")
		}

		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}
		servers, err := config.LoadServers(cfg.ConfigPath)
		if err != nil {
			return fmt.Errorf("error loading servers: %w", err)
		}

		// Assume last arg is server name, rest are SSH flags
		name := args[len(args)-1]
		sshArgs := args[:len(args)-1]

		var selectedServer models.Server
		for _, s := range servers {
			if s.Name == name {
				selectedServer = s
				break
			}
		}
		if selectedServer.Name == "" {
			return fmt.Errorf("server not found: %s", name)
		}

		sftpArgs := []string{}

		if selectedServer.Key != "" {
			sftpArgs = append(sftpArgs, "-i", selectedServer.Key)
		}
		if selectedServer.Port != 0 && selectedServer.Port != 22 {
			sftpArgs = append(sftpArgs, "-P", fmt.Sprintf("%d", selectedServer.Port))
		}

		filteredArgs := []string{}
		for i := 0; i < len(sshArgs); i++ {
			if sshArgs[i] == "-l" && i+1 < len(sshArgs) {
				i++
				continue
			}
			filteredArgs = append(filteredArgs, sshArgs[i])
		}
		sftpArgs = append(sftpArgs, filteredArgs...)

		user := selectedServer.User
		for i, arg := range sshArgs {
			if arg == "-l" && i+1 < len(sshArgs) {
				user = sshArgs[i+1]
				break
			}
		}

		userHost := selectedServer.Host
		if user != "" {
			userHost = user + "@" + selectedServer.Host
		}
		sftpArgs = append(sftpArgs, userHost)

		sftpCmd := exec.Command("sftp", sftpArgs...)
		sftpCmd.Stdin = os.Stdin
		sftpCmd.Stdout = os.Stdout
		sftpCmd.Stderr = os.Stderr
		err = sftpCmd.Run()
		if err != nil {
			return fmt.Errorf("error starting SFTP: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sftpCmd)
	sftpCmd.DisableFlagParsing = true
}
