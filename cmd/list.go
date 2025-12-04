package cmd

import (
	"fmt"
	"strings"

	"github.com/omisai-tech/sshy/internal/config"
	"github.com/omisai-tech/sshy/internal/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured SSH servers",
	Long: `List all SSH servers from the configuration.

Server prefixes indicate source:
- [S] Shared servers from servers.yaml
- [L] Local private servers from local.yaml
- [O] Shared servers with local overrides in local.yaml

Use --tags to filter by tags.`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetStringSlice("tags")

		cfg, err := config.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		serversWithSource, err := config.LoadServersWithSourceAndPath(cfg.ConfigPath, cfg.GetServersSource())
		if err != nil {
			fmt.Println("Error loading servers:", err)
			return
		}

		for _, sws := range serversWithSource {
			s := sws.Server
			if len(tags) == 0 || hasAllTags(s.Tags, tags) {
				sourceFlag := ""
				switch sws.Source {
				case models.SourceShared:
					sourceFlag = "[S]"
				case models.SourceLocal:
					sourceFlag = "[L]"
				case models.SourceOverride:
					sourceFlag = "[O]"
				}
				fmt.Printf("%s %s: %s@%s [%s]\n", sourceFlag, s.Name, s.User, s.Host, strings.Join(s.Tags, ", "))
			}
		}
	},
}

func hasAllTags(serverTags, filterTags []string) bool {
	tagMap := make(map[string]bool)
	for _, t := range serverTags {
		tagMap[t] = true
	}
	for _, t := range filterTags {
		if !tagMap[t] {
			return false
		}
	}
	return true
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringSliceP("tags", "t", []string{}, "Filter servers by tags (comma-separated)")
}
