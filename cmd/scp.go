package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var scpCmd = &cobra.Command{
	Use:   "scp [ssh-flags...] <source> <destination>",
	Short: "Copy files to/from SSH servers",
	Long:  "Copy files between local machine and SSH servers using scp. Use server:path for remote paths. SSH flags can be passed through.",
	Args:  cobra.MinimumNArgs(0), // Allow any number of args, we'll parse them
	RunE: func(cmd *cobra.Command, args []string) error {
		// Handle help flag manually since we need to allow unknown flags
		for _, arg := range args {
			if arg == "--help" || arg == "-h" {
				return cmd.Help()
			}
		}

		if len(args) < 2 {
			return fmt.Errorf("usage: sshy scp [ssh-flags...] <source> <destination>")
		}

		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}
		servers, err := config.LoadServersWithPath(cfg.ConfigPath, cfg.GetServersSource())
		if err != nil {
			return fmt.Errorf("error loading servers: %w", err)
		}

		// Assume last two args are source and destination, rest are SSH flags
		source := args[len(args)-2]
		destination := args[len(args)-1]
		sshArgs := args[:len(args)-2]

		// Parse source and destination for server names
		serverName1, remotePath1 := parsePath(source)
		serverName2, remotePath2 := parsePath(destination)

		var server1, server2 *models.Server
		if serverName1 != "" {
			for i, s := range servers {
				if s.Name == serverName1 {
					server1 = &servers[i]
					break
				}
			}
			if server1 == nil {
				return fmt.Errorf("server %s not found", serverName1)
			}
		}
		if serverName2 != "" {
			for i, s := range servers {
				if s.Name == serverName2 {
					server2 = &servers[i]
					break
				}
			}
			if server2 == nil {
				return fmt.Errorf("server %s not found", serverName2)
			}
		}

		scpArgs := []string{}

		var server *models.Server
		if server1 != nil {
			server = server1
		} else if server2 != nil {
			server = server2
		}
		if server != nil {
			if server.Key != "" {
				scpArgs = append(scpArgs, "-i", server.Key)
			}
			if server.Port != 0 && server.Port != 22 {
				scpArgs = append(scpArgs, "-P", fmt.Sprintf("%d", server.Port))
			}
		}

		filteredArgs := []string{}
		for i := 0; i < len(sshArgs); i++ {
			if sshArgs[i] == "-l" && i+1 < len(sshArgs) {
				i++
				continue
			}
			filteredArgs = append(filteredArgs, sshArgs[i])
		}
		scpArgs = append(scpArgs, filteredArgs...)

		if server1 != nil {
			scpArgs = append(scpArgs, buildScpTarget(*server1, remotePath1, sshArgs))
		} else {
			scpArgs = append(scpArgs, source)
		}
		if server2 != nil {
			scpArgs = append(scpArgs, buildScpTarget(*server2, remotePath2, sshArgs))
		} else {
			scpArgs = append(scpArgs, destination)
		}

		scpCmd := exec.Command("scp", scpArgs...)
		scpCmd.Stdin = os.Stdin
		scpCmd.Stdout = os.Stdout
		scpCmd.Stderr = os.Stderr
		err = scpCmd.Run()
		if err != nil {
			return fmt.Errorf("error copying: %w", err)
		}
		return nil
	},
}

func parsePath(path string) (serverName, remotePath string) {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return "", path
}

func buildScpTarget(s models.Server, remotePath string, sshArgs []string) string {
	user := s.User

	for i, arg := range sshArgs {
		if arg == "-l" && i+1 < len(sshArgs) {
			user = sshArgs[i+1]
			break
		}
	}

	userHost := s.Host
	if user != "" {
		userHost = user + "@" + s.Host
	}
	return userHost + ":" + remotePath
}

func init() {
	rootCmd.AddCommand(scpCmd)
	scpCmd.DisableFlagParsing = true
}
