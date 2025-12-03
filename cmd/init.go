package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sshy configuration",
	Long:  "Interactively set up the initial configuration for sshy by specifying the path to the servers YAML file.",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter the path for the shared yaml file (eg: ~/.sshy/servers.yaml): ")
		serversPath, _ := reader.ReadString('\n')
		serversPath = strings.TrimSpace(serversPath)

		// Expand ~ to home directory
		if strings.HasPrefix(serversPath, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Println("Error getting home directory:", err)
				return
			}
			serversPath = strings.Replace(serversPath, "~", home, 1)
		}

		cfg := config.DefaultConfig()
		cfg.ConfigPath = filepath.Dir(serversPath)
		cfg.ServersPath = filepath.Base(serversPath)

		err := config.SaveGlobalConfig(cfg)
		if err != nil {
			fmt.Println("Error saving config:", err)
			return
		}

		// Create the directory for the servers file if it doesn't exist
		serversDir := filepath.Dir(serversPath)
		err = os.MkdirAll(serversDir, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		// Create an empty servers file if it doesn't exist
		if _, err := os.Stat(serversPath); os.IsNotExist(err) {
			err = os.WriteFile(serversPath, []byte("# Shared SSH servers configuration\n# Add your servers here\n"), 0644)
			if err != nil {
				fmt.Println("Error creating servers file:", err)
				return
			}
			fmt.Printf("Created empty servers file: %s\n", serversPath)
		}

		fmt.Println("Configuration initialized successfully.")
		fmt.Printf("Global config saved to: ~/.sshy/config.yaml\n")
		fmt.Printf("Servers file: %s\n", serversPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
