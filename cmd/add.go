package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [<name> <host> <user> [<port>] [<key>] [<tags>]]",
	Short: "Add a new SSH server",
	Long:  "Add a new SSH server to the configuration and commit the changes. If no arguments are provided, prompts interactively. The port defaults to 22, key to ~/.ssh/id_rsa, and tags are optional.",
	Run: func(cmd *cobra.Command, args []string) {
		var name, host, user, key string
		var port int = 22
		var tags []string
		if len(args) == 0 {
			// Interactive mode
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Server name: ")
			scanner.Scan()
			name = strings.TrimSpace(scanner.Text())
			fmt.Print("Host: ")
			scanner.Scan()
			host = strings.TrimSpace(scanner.Text())
			fmt.Print("User (default: root): ")
			scanner.Scan()
			user = strings.TrimSpace(scanner.Text())
			if user == "" {
				user = "root"
			}
			fmt.Print("Port (default: 22): ")
			scanner.Scan()
			portStr := strings.TrimSpace(scanner.Text())
			if portStr != "" {
				if p, err := strconv.Atoi(portStr); err == nil {
					port = p
				}
			}
			fmt.Print("Key path (default: ~/.ssh/id_rsa): ")
			scanner.Scan()
			key = strings.TrimSpace(scanner.Text())
			if key == "" {
				key = "~/.ssh/id_rsa"
			}
			fmt.Print("Tags (comma separated, optional): ")
			scanner.Scan()
			tagsStr := strings.TrimSpace(scanner.Text())
			if tagsStr != "" {
				tags = strings.Split(tagsStr, ",")
				for i, t := range tags {
					tags[i] = strings.TrimSpace(t)
				}
			}
		} else {
			if len(args) < 3 {
				fmt.Println("Usage: sshy add [<name> <host> <user> [<port>] [<key>] [<tags>]]")
				return
			}
			name, host, user = args[0], args[1], args[2]
			if len(args) >= 4 {
				if p, err := strconv.Atoi(args[3]); err == nil {
					port = p
				}
			}
			if len(args) >= 5 {
				key = args[4]
			} else {
				key = "~/.ssh/id_rsa"
			}
			if len(args) >= 6 {
				tagsStr := args[5]
				if tagsStr != "" {
					tags = strings.Split(tagsStr, ",")
					for i, t := range tags {
						tags[i] = strings.TrimSpace(t)
					}
				}
			}
		}

		localConfig, err := config.LoadLocalConfig()
		if err != nil {
			fmt.Println("Error loading local config:", err)
			return
		}
		localConfig.Private = append(localConfig.Private, models.Server{Name: name, Host: host, User: user, Port: port, Key: key, Tags: tags})
		err = config.SaveLocalConfig(localConfig)
		if err != nil {
			fmt.Println("Error saving local config:", err)
			return
		}
		fmt.Println("Server added successfully")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
