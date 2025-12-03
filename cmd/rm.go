package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [name]",
	Short: "Remove an SSH server",
	Long:  "Remove the specified SSH server from the configuration and commit the changes. If no name is provided, select from available servers.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}
		localConfig, err := config.LoadLocalConfig()
		if err != nil {
			fmt.Println("Error loading local config:", err)
			return
		}
		serversWithSource, err := config.LoadServersWithSourceAndPath(cfg.ConfigPath, cfg.ServersPath)
		if err != nil {
			fmt.Println("Error loading servers:", err)
			return
		}

		// Filter to only local and overridden servers
		var removableServers []models.ServerWithSource
		for _, sws := range serversWithSource {
			if sws.Source == models.SourceLocal || sws.Source == models.SourceOverride {
				removableServers = append(removableServers, sws)
			}
		}

		var name string
		if len(args) == 0 {
			if len(removableServers) == 0 {
				fmt.Println("No removable servers configured")
				return
			}
			names := make([]string, len(removableServers))
			for i, sws := range removableServers {
				names[i] = sws.Server.Name
			}
			idx, err := fuzzyfinder.Find(names, func(i int) string { return names[i] })
			if err != nil {
				fmt.Println("Selection cancelled")
				return
			}
			name = names[idx]
		} else if len(args) == 1 {
			name = args[0]
		} else {
			fmt.Println("Usage: sshy rm [name]")
			return
		}

		// Find the server with source in removableServers
		var sws *models.ServerWithSource
		for i, s := range removableServers {
			if s.Server.Name == name {
				sws = &removableServers[i]
				break
			}
		}
		if sws == nil {
			fmt.Println("Server not found or not removable")
			return
		}

		// Confirmation
		fmt.Printf("Are you sure you want to remove '%s'? (y/N): ", name)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if confirm != "y" && confirm != "yes" {
			fmt.Println("Removal cancelled")
			return
		}

		// Remove based on source
		switch sws.Source {
		case models.SourceLocal:
			var newPrivate models.Servers
			for _, s := range localConfig.Private {
				if s.Name != name {
					newPrivate = append(newPrivate, s)
				}
			}
			localConfig.Private = newPrivate
		case models.SourceOverride:
			delete(localConfig.Servers, name)
		default:
			fmt.Println("Cannot remove shared server")
			return
		}

		err = config.SaveLocalConfig(localConfig)
		if err != nil {
			fmt.Println("Error saving local config:", err)
			return
		}
		fmt.Println("Server removed successfully")
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rmCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
