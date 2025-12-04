package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/omisai-tech/sshy/internal/models"
)

const (
	SharedConfigFile = "servers.yaml"
	LocalConfigFile  = "local.yaml"
)

type LocalConfig struct {
	Servers map[string]models.Server `yaml:"servers" json:"servers"`
	Private models.Servers           `yaml:"private" json:"private"`
}

func findConfigFile(basePath, primaryFile string) (string, FileFormat) {
	primaryPath := filepath.Join(basePath, primaryFile)
	if _, err := os.Stat(primaryPath); err == nil {
		return primaryPath, DetectFormat(primaryFile)
	}
	alternatePath := filepath.Join(basePath, GetAlternateFilename(primaryFile))
	if _, err := os.Stat(alternatePath); err == nil {
		return alternatePath, DetectFormat(GetAlternateFilename(primaryFile))
	}
	return primaryPath, DetectFormat(primaryFile)
}

func loadLocalConfig() (LocalConfig, error) {
	home, err := userHomeDir()
	if err != nil {
		return LocalConfig{}, err
	}
	sshyDir := filepath.Join(home, ".sshy")
	localPath, format := findConfigFile(sshyDir, LocalConfigFile)
	localData, err := os.ReadFile(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return LocalConfig{}, nil
		}
		return LocalConfig{}, err
	}

	if format == FormatUnknown {
		format = DetectFormatFromContent(localData)
	}

	var config LocalConfig
	err = Unmarshal(localData, format, &config)
	return config, err
}

func LoadServers(configPath string) (models.Servers, error) {
	return LoadServersWithPath(configPath, SharedConfigFile)
}

func LoadServersWithPath(configPath, serversPath string) (models.Servers, error) {
	sharedPath, format := findConfigFile(configPath, serversPath)
	sharedData, err := os.ReadFile(sharedPath)
	if err != nil {
		if os.IsNotExist(err) {
			sharedData = []byte{}
		} else {
			return nil, err
		}
	}

	if format == FormatUnknown && len(sharedData) > 0 {
		format = DetectFormatFromContent(sharedData)
	}

	var sharedServers models.Servers
	if len(sharedData) > 0 {
		if err := Unmarshal(sharedData, format, &sharedServers); err != nil {
			return nil, err
		}
	}

	localConfig, err := loadLocalConfig()
	if err != nil {
		return nil, err
	}

	mergedServers := make(models.Servers, 0, len(sharedServers)+len(localConfig.Private))
	serverMap := make(map[string]*models.Server)

	for i := range sharedServers {
		server := sharedServers[i]
		serverMap[server.Name] = &server
		if override, ok := localConfig.Servers[server.Name]; ok {
			if override.Host != "" {
				server.Host = override.Host
			}
			if override.User != "" {
				server.User = override.User
			}
			if override.Port != 0 {
				server.Port = override.Port
			}
			if len(override.Tags) > 0 {
				server.Tags = override.Tags
			}
			if override.Key != "" {
				server.Key = override.Key
			}
			if override.Options != nil {
				server.Options = override.Options
			}
		}
		mergedServers = append(mergedServers, server)
	}

	mergedServers = append(mergedServers, localConfig.Private...)

	return mergedServers, nil
}

func SaveServers(configPath string, servers models.Servers) error {
	return SaveServersWithPath(configPath, SharedConfigFile, servers)
}

func SaveServersWithPath(configPath, serversPath string, servers models.Servers) error {
	localConfig, err := loadLocalConfig()
	if err != nil {
		return err
	}
	privateServerMap := make(map[string]struct{})
	for _, p := range localConfig.Private {
		privateServerMap[p.Name] = struct{}{}
	}

	serversToSave := make(models.Servers, 0, len(servers))
	for _, server := range servers {
		if _, isPrivate := privateServerMap[server.Name]; !isPrivate {
			serverToSave := server
			serverToSave.Key = ""
			serversToSave = append(serversToSave, serverToSave)
		}
	}

	existingPath, format := findConfigFile(configPath, serversPath)
	if format == FormatUnknown {
		format = DetectFormat(serversPath)
	}
	if format == FormatUnknown {
		format = FormatYAML
	}

	data, err := Marshal(serversToSave, format)
	if err != nil {
		return err
	}

	if strings.TrimSuffix(filepath.Base(existingPath), filepath.Ext(existingPath)) == strings.TrimSuffix(serversPath, filepath.Ext(serversPath)) {
		return os.WriteFile(existingPath, data, 0644)
	}
	fullPath := filepath.Join(configPath, serversPath)
	return os.WriteFile(fullPath, data, 0644)
}

func LoadLocalConfig() (LocalConfig, error) {
	return loadLocalConfig()
}

var userHomeDir = os.UserHomeDir

func SaveLocalConfig(config LocalConfig) error {
	home, err := userHomeDir()
	if err != nil {
		return err
	}
	sshyDir := filepath.Join(home, ".sshy")
	existingPath, format := findConfigFile(sshyDir, LocalConfigFile)
	if format == FormatUnknown {
		format = FormatYAML
	}

	data, err := Marshal(config, format)
	if err != nil {
		return err
	}
	return os.WriteFile(existingPath, data, 0644)
}

func LoadServersWithSource(configPath string) ([]models.ServerWithSource, error) {
	return LoadServersWithSourceAndPath(configPath, SharedConfigFile)
}

func LoadServersWithSourceAndPath(configPath, serversPath string) ([]models.ServerWithSource, error) {
	sharedPath, format := findConfigFile(configPath, serversPath)
	sharedData, err := os.ReadFile(sharedPath)
	if err != nil {
		if os.IsNotExist(err) {
			sharedData = []byte{}
		} else {
			return nil, err
		}
	}

	if format == FormatUnknown && len(sharedData) > 0 {
		format = DetectFormatFromContent(sharedData)
	}

	var sharedServers models.Servers
	if len(sharedData) > 0 {
		if err := Unmarshal(sharedData, format, &sharedServers); err != nil {
			return nil, err
		}
	}

	localConfig, err := loadLocalConfig()
	if err != nil {
		return nil, err
	}

	mergedServers := make([]models.ServerWithSource, 0, len(sharedServers)+len(localConfig.Private))
	serverMap := make(map[string]*models.Server)

	for i := range sharedServers {
		server := sharedServers[i]
		serverMap[server.Name] = &server
		if override, ok := localConfig.Servers[server.Name]; ok {
			if override.Host != "" {
				server.Host = override.Host
			}
			if override.User != "" {
				server.User = override.User
			}
			if override.Port != 0 {
				server.Port = override.Port
			}
			if len(override.Tags) > 0 {
				server.Tags = override.Tags
			}
			if override.Key != "" {
				server.Key = override.Key
			}
			if override.Options != nil {
				server.Options = override.Options
			}
			mergedServers = append(mergedServers, models.ServerWithSource{Server: server, Source: models.SourceOverride})
		} else {
			mergedServers = append(mergedServers, models.ServerWithSource{Server: server, Source: models.SourceShared})
		}
	}

	for _, server := range localConfig.Private {
		mergedServers = append(mergedServers, models.ServerWithSource{Server: server, Source: models.SourceLocal})
	}

	return mergedServers, nil
}
