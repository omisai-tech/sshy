package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/omisai-tech/sshy/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View configuration files",
	Long: `View the contents of configuration files.

You can view:
- servers.yaml: Shared server configuration (local file or remote URL)
- local.yaml: Local overrides and private servers`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		serversLabel := cfg.GetServersSource()
		if cfg.IsRemoteSource() {
			serversLabel = fmt.Sprintf("URL: %s", cfg.ServersURL)
		}

		options := []string{
			fmt.Sprintf("%s - Shared server configuration", serversLabel),
			"local.yaml - Local overrides and private servers",
		}

		idx, err := fuzzyfinder.Find(options, func(i int) string { return options[i] })
		if err != nil {
			fmt.Println("Selection cancelled")
			return
		}

		var data []byte
		var title string

		switch idx {
		case 0:
			if cfg.IsRemoteSource() {
				title = fmt.Sprintf("Shared Configuration (URL: %s)", cfg.ServersURL)
				resp, err := http.Get(cfg.ServersURL)
				if err != nil {
					fmt.Printf("Error fetching from URL: %v\n", err)
					return
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Error: server returned status %d\n", resp.StatusCode)
					return
				}
				data, err = io.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("Error reading response: %v\n", err)
					return
				}
			} else {
				filePath := fmt.Sprintf("%s/%s", cfg.ConfigPath, cfg.ServersPath)
				title = fmt.Sprintf("Shared Configuration (%s)", cfg.ServersPath)
				data, err = os.ReadFile(filePath)
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Printf("%s does not exist or is empty\n", title)
						return
					}
					fmt.Printf("Error reading %s: %v\n", title, err)
					return
				}
			}
		case 1:
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				return
			}
			filePath := fmt.Sprintf("%s/.sshy/local.yaml", home)
			title = "Local Configuration (local.yaml)"
			data, err = os.ReadFile(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("%s does not exist or is empty\n", title)
					return
				}
				fmt.Printf("Error reading %s: %v\n", title, err)
				return
			}
		}

		fmt.Printf("=== %s ===\n\n", title)

		var obj interface{}
		if err := yaml.Unmarshal(data, &obj); err == nil {
			prettyData, err := yaml.Marshal(obj)
			if err == nil {
				fmt.Print(string(prettyData))
				return
			}
		}

		fmt.Print(string(data))
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
