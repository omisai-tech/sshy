package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

type CommandRunner interface {
	Run(name string, args []string) error
}

type DefaultCommandRunner struct{}

func (r *DefaultCommandRunner) Run(name string, args []string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var cmdRunner CommandRunner = &DefaultCommandRunner{}

var fuzzyFind = func(names []string, itemFunc func(int) string) (int, error) {
	return fuzzyfinder.Find(names, itemFunc)
}

var connectCmd = &cobra.Command{
	Use:   "connect [name] [ssh-flags...] [command]",
	Short: "Connect to an SSH server",
	Long:  "Connect to the specified SSH server or select one interactively if no name is provided. SSH flags can be passed through. Use -- to separate SSH options from remote commands.",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if arg == "--help" || arg == "-h" {
				return cmd.Help()
			}
		}

		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %w", err)
		}

		servers, err := config.LoadServersWithPath(cfg.ConfigPath, cfg.ServersPath)
		if err != nil {
			return fmt.Errorf("error loading servers: %w", err)
		}

		var selectedServer models.Server
		var sshArgs []string
		var remoteCommand string

		if len(args) == 0 {
			names := make([]string, len(servers))
			for i, s := range servers {
				names[i] = s.Name
			}
			idx, err := fuzzyFind(names, func(i int) string { return names[i] })
			if err != nil {
				fmt.Println("No server selected")
				return nil
			}
			selectedServer = servers[idx]
		} else {
			name := args[0]
			for _, s := range servers {
				if s.Name == name {
					selectedServer = s
					break
				}
			}
			if selectedServer.Name == "" {
				return fmt.Errorf("server not found: %s", name)
			}

			remainingArgs := args[1:]
			commandStart := -1

			for i, arg := range remainingArgs {
				if arg == "--" {
					commandStart = i + 1
					break
				}
			}

			if commandStart >= 0 {
				sshArgs = remainingArgs[:commandStart-1]
				remoteCommand = strings.Join(remainingArgs[commandStart:], " ")
			} else {
				sshArgs = remainingArgs
			}
		}
		connectTo(selectedServer, sshArgs, remoteCommand)
		return nil
	},
}

func buildSSHArgs(s models.Server, sshArgs []string, remoteCommand string) []string {
	args := []string{}

	if s.Key != "" {
		args = append(args, "-i", s.Key)
	}

	if s.Port != 0 && s.Port != 22 {
		args = append(args, "-p", fmt.Sprintf("%d", s.Port))
	}

	for key, value := range s.Options {
		switch key {
		case "ForwardAgent":
			if value == "yes" {
				args = append(args, "-A")
			}
		case "RequestTTY":
			switch value {
			case "force":
				args = append(args, "-t", "-t")
			case "yes":
				args = append(args, "-t")
			}
		case "LocalForward":
			args = append(args, "-L", fmt.Sprintf("%v", value))
		}
	}

	args = append(args, sshArgs...)

	hasUserFlag := false
	for _, arg := range sshArgs {
		if arg == "-l" {
			hasUserFlag = true
			break
		}
	}

	userHost := s.Host
	if s.User != "" && !hasUserFlag {
		userHost = s.User + "@" + s.Host
	}
	args = append(args, userHost)

	if remoteCommand != "" {
		args = append(args, remoteCommand)
	}

	return args
}

func connectTo(s models.Server, sshArgs []string, remoteCommand string) {
	args := buildSSHArgs(s, sshArgs, remoteCommand)
	err := cmdRunner.Run("ssh", args)
	if err != nil {
		fmt.Println("Error connecting:", err)
	}
}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.DisableFlagParsing = true
}
