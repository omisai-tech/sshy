package cmd

import (
	"fmt"
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
- servers.yaml: Shared server configuration
- local.yaml: Local overrides and private servers`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		options := []string{
			fmt.Sprintf("%s - Shared server configuration", cfg.ServersPath),
			"local.yaml - Local overrides and private servers",
		}

		idx, err := fuzzyfinder.Find(options, func(i int) string { return options[i] })
		if err != nil {
			fmt.Println("Selection cancelled")
			return
		}

		var filePath string
		var title string

		switch idx {
		case 0:
			filePath = fmt.Sprintf("%s/%s", cfg.ConfigPath, cfg.ServersPath)
			title = fmt.Sprintf("Shared Configuration (%s)", cfg.ServersPath)
		case 1:
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				return
			}
			filePath = fmt.Sprintf("%s/.sshy/local.yaml", home)
			title = "Local Configuration (local.yaml)"
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("%s does not exist or is empty\n", title)
				return
			}
			fmt.Printf("Error reading %s: %v\n", title, err)
			return
		}

		fmt.Printf("=== %s ===\n\n", title)

		// Pretty print YAML if possible
		var obj interface{}
		if err := yaml.Unmarshal(data, &obj); err == nil {
			prettyData, err := yaml.Marshal(obj)
			if err == nil {
				fmt.Print(string(prettyData))
				return
			}
		}

		// Fallback to raw output
		fmt.Print(string(data))
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)
}
