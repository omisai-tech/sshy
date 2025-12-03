package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/spf13/cobra"
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Edit local configuration overrides",
	Long:  "Open local.yaml in your default editor to edit local overrides for servers.",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		localConfigPath := filepath.Join(cfg.ConfigPath, "local.yaml")
		// Ensure the file exists with template if empty or missing
		stat, err := os.Stat(localConfigPath)
		if os.IsNotExist(err) || (err == nil && stat.Size() == 0) {
			// Create file with template
			template := `# Local configuration for sshy
# This file allows you to override shared server configurations and define private servers.

# Overrides for shared servers:
# servers:
#   server_name:
#     host: new_host
#     user: new_user
#     port: 22
#     key: /path/to/key
#     options:
#       option1: value1

# Private servers (not shared in repo):
# private:
#   - name: private_server
#     host: private_host
#     user: private_user
#     port: 22
#     key: /path/to/key
#     tags: [tag1, tag2]
#     options:
#       option1: value1
`
			err = os.WriteFile(localConfigPath, []byte(template), 0644)
			if err != nil {
				fmt.Println("Error creating local.yaml:", err)
				return
			}
		}

		// Open in editor
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano" // default
		}
		editCmd := exec.Command(editor, localConfigPath)
		editCmd.Stdin = os.Stdin
		editCmd.Stdout = os.Stdout
		editCmd.Stderr = os.Stderr
		err = editCmd.Run()
		if err != nil {
			fmt.Println("Error opening editor:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(localCmd)
}
