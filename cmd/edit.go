package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [name]",
	Short: "Edit an existing SSH server configuration",
	Long:  "Edit the configuration of an existing SSH server. If no name is provided, select from available servers. Prompts interactively for each field.",
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
		serversWithSource, err := config.LoadServersWithSourceAndPath(cfg.ConfigPath, cfg.GetServersSource())
		if err != nil {
			fmt.Println("Error loading servers:", err)
			return
		}

		var name string
		if len(args) == 0 {
			if len(serversWithSource) == 0 {
				fmt.Println("No servers configured")
				return
			}
			names := make([]string, len(serversWithSource))
			for i, sws := range serversWithSource {
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
			fmt.Println("Usage: sshy edit [name]")
			return
		}

		// Find the server
		var serverIndex int = -1
		for i, sws := range serversWithSource {
			if sws.Server.Name == name {
				serverIndex = i
				break
			}
		}
		if serverIndex == -1 {
			fmt.Println("Server not found")
			return
		}

		server := serversWithSource[serverIndex].Server

		// Interactive edit
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Printf("Name (%s): ", server.Name)
		scanner.Scan()
		newName := strings.TrimSpace(scanner.Text())
		if newName != "" {
			// Check if name conflicts
			for _, sws := range serversWithSource {
				if sws.Server.Name == newName && sws.Server.Name != server.Name {
					fmt.Println("Server name already exists")
					return
				}
			}
			server.Name = newName
		}

		fmt.Printf("Host (%s): ", server.Host)
		scanner.Scan()
		newHost := strings.TrimSpace(scanner.Text())
		if newHost != "" {
			server.Host = newHost
		}

		userPrompt := server.User
		if userPrompt == "" {
			userPrompt = "root"
		}
		fmt.Printf("User (%s): ", userPrompt)
		scanner.Scan()
		newUser := strings.TrimSpace(scanner.Text())
		if newUser != "" {
			server.User = newUser
		}

		portPrompt := server.Port
		if portPrompt == 0 {
			portPrompt = 22
		}
		fmt.Printf("Port (%d): ", portPrompt)
		scanner.Scan()
		portStr := strings.TrimSpace(scanner.Text())
		if portStr != "" {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				fmt.Println("Invalid port number")
				return
			}
			server.Port = port
		}

		keyPrompt := server.Key
		if keyPrompt == "" {
			keyPrompt = "~/.ssh/id_rsa"
		}
		fmt.Printf("Key (%s): ", keyPrompt)
		scanner.Scan()
		newKey := strings.TrimSpace(scanner.Text())
		if newKey != "" {
			server.Key = newKey
		}

		fmt.Printf("Tags (%s): ", strings.Join(server.Tags, ", "))
		scanner.Scan()
		tagsStr := strings.TrimSpace(scanner.Text())
		if tagsStr != "" {
			server.Tags = strings.Split(tagsStr, ",")
			for i, t := range server.Tags {
				server.Tags[i] = strings.TrimSpace(t)
			}
		}

		// Update the server
		serversWithSource[serverIndex].Server = server

		// Save based on source
		sws := serversWithSource[serverIndex]
		if sws.Source == models.SourceLocal {
			// Update localConfig.Private
			found := false
			for i, s := range localConfig.Private {
				if s.Name == server.Name {
					localConfig.Private[i] = server
					found = true
					break
				}
			}
			if !found {
				localConfig.Private = append(localConfig.Private, server)
			}
		} else {
			// Save to servers map (override)
			localConfig.Servers[server.Name] = server
		}

		// Save
		err = config.SaveLocalConfig(localConfig)
		if err != nil {
			fmt.Println("Error saving local config:", err)
			return
		}
		fmt.Println("Server updated successfully")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
