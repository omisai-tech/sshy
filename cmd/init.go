package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func promptConfirm(reader *bufio.Reader, message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize sshy configuration",
	Long:  "Interactively set up the initial configuration for sshy by specifying the path and format for configuration files.",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		sshyDir := filepath.Join(home, ".sshy")

		existingYAMLConfig := fileExists(filepath.Join(sshyDir, "config.yaml"))
		existingJSONConfig := fileExists(filepath.Join(sshyDir, "config.json"))

		if existingYAMLConfig || existingJSONConfig {
			currentFormat := "YAML"
			if existingJSONConfig {
				currentFormat = "JSON"
			}
			fmt.Printf("Existing configuration found (%s format).\n", currentFormat)
			if !promptConfirm(reader, "Do you want to reconfigure?") {
				fmt.Println("Configuration unchanged.")
				return
			}
		}

		fmt.Println("\nChoose your preferred configuration format:")
		fmt.Println("  1) YAML (default)")
		fmt.Println("  2) JSON")
		fmt.Print("Enter choice [1/2]: ")
		formatChoice, _ := reader.ReadString('\n')
		formatChoice = strings.TrimSpace(formatChoice)

		var fileFormat config.FileFormat
		var fileExt, altExt string
		switch formatChoice {
		case "2", "json", "JSON":
			fileFormat = config.FormatJSON
			fileExt = ".json"
			altExt = ".yaml"
		default:
			fileFormat = config.FormatYAML
			fileExt = ".yaml"
			altExt = ".json"
		}

		fmt.Println("\nChoose how to configure shared servers:")
		fmt.Println("  1) Local file path (default)")
		fmt.Println("  2) Remote URL (HTTP/HTTPS)")
		fmt.Print("Enter choice [1/2]: ")
		sourceChoice, _ := reader.ReadString('\n')
		sourceChoice = strings.TrimSpace(sourceChoice)

		var serversPath string
		var serversURL string
		var serversDir string

		if sourceChoice == "2" {
			fmt.Print("Enter the URL for the shared servers configuration: ")
			serversURL, _ = reader.ReadString('\n')
			serversURL = strings.TrimSpace(serversURL)
			if serversURL == "" {
				fmt.Println("Error: URL cannot be empty")
				return
			}
			if err := config.ValidateURL(serversURL); err != nil {
				fmt.Println("Error:", err)
				return
			}
			serversDir = sshyDir
		} else {
			defaultServersPath := fmt.Sprintf("~/.sshy/servers%s", fileExt)
			fmt.Printf("Enter the path for the shared servers file (default: %s): ", defaultServersPath)
			serversPath, _ = reader.ReadString('\n')
			serversPath = strings.TrimSpace(serversPath)
			if serversPath == "" {
				serversPath = defaultServersPath
			}

			if strings.HasPrefix(serversPath, "~") {
				serversPath = strings.Replace(serversPath, "~", home, 1)
			}

			serversDir = filepath.Dir(serversPath)
			err = os.MkdirAll(serversDir, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
		}

		err = os.MkdirAll(sshyDir, 0755)
		if err != nil {
			fmt.Println("Error creating .sshy directory:", err)
			return
		}

		altConfigPath := filepath.Join(sshyDir, "config"+altExt)
		if fileExists(altConfigPath) {
			if promptConfirm(reader, fmt.Sprintf("Remove old config file (%s)?", altConfigPath)) {
				os.Remove(altConfigPath)
				fmt.Printf("Removed: %s\n", altConfigPath)
			}
		}

		altLocalPath := filepath.Join(sshyDir, "local"+altExt)
		if fileExists(altLocalPath) {
			if promptConfirm(reader, fmt.Sprintf("Remove old local config file (%s)?", altLocalPath)) {
				os.Remove(altLocalPath)
				fmt.Printf("Removed: %s\n", altLocalPath)
			}
		}

		cfg := config.DefaultConfig()
		cfg.ConfigPath = serversDir
		if serversURL != "" {
			cfg.ServersURL = serversURL
			cfg.ServersPath = ""
		} else {
			cfg.ServersPath = filepath.Base(serversPath)
		}

		err = config.SaveGlobalConfigWithFormat(cfg, fileFormat)
		if err != nil {
			fmt.Println("Error saving global config:", err)
			return
		}

		serversCreated := false
		if serversURL == "" && !fileExists(serversPath) {
			var content []byte
			if fileFormat == config.FormatJSON {
				content = []byte("[\n]\n")
			} else {
				content = []byte("# Shared SSH servers configuration\n# Add your servers here\n")
			}
			err = os.WriteFile(serversPath, content, 0644)
			if err != nil {
				fmt.Println("Error creating servers file:", err)
				return
			}
			serversCreated = true
		}

		localPath := filepath.Join(sshyDir, "local"+fileExt)
		localCreated := false
		if !fileExists(localPath) {
			defaultLocal := config.LocalConfig{
				Servers: make(map[string]models.Server),
				Private: make(models.Servers, 0),
			}
			data, err := config.Marshal(defaultLocal, fileFormat)
			if err != nil {
				fmt.Println("Error marshaling local config:", err)
				return
			}
			err = os.WriteFile(localPath, data, 0644)
			if err != nil {
				fmt.Println("Error creating local config file:", err)
				return
			}
			localCreated = true
		}

		formatName := "YAML"
		if fileFormat == config.FormatJSON {
			formatName = "JSON"
		}

		fmt.Println("\nConfiguration initialized successfully.")
		fmt.Printf("Format: %s\n", formatName)
		fmt.Printf("Global config: ~/.sshy/config%s\n", fileExt)
		if localCreated {
			fmt.Printf("Local config: %s (created)\n", localPath)
		} else {
			fmt.Printf("Local config: %s (exists, unchanged)\n", localPath)
		}
		if serversURL != "" {
			fmt.Printf("Servers URL: %s\n", serversURL)
		} else if serversCreated {
			fmt.Printf("Servers file: %s (created)\n", serversPath)
		} else {
			fmt.Printf("Servers file: %s (exists, unchanged)\n", serversPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
