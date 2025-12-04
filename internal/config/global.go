package config

import (
	"os"
	"path/filepath"

	"github.com/omisai-tech/sshy/internal/models"
)

const (
	GlobalConfigFile = "config.yaml"
)

type GlobalConfig struct {
	ServersPath string `yaml:"servers_path" json:"servers_path"`
	ServersURL  string `yaml:"servers_url,omitempty" json:"servers_url,omitempty"`
	ConfigPath  string `yaml:"config_path" json:"config_path"`
}

func (c *GlobalConfig) GetServersSource() string {
	if c.ServersURL != "" {
		return c.ServersURL
	}
	return c.ServersPath
}

func (c *GlobalConfig) IsRemoteSource() bool {
	return c.ServersURL != "" && IsURL(c.ServersURL)
}

var globalUserHomeDir = os.UserHomeDir

func DefaultConfig() *GlobalConfig {
	home, _ := globalUserHomeDir()
	return &GlobalConfig{
		ServersPath: "servers.yaml",
		ConfigPath:  filepath.Join(home, ".sshy"),
	}
}

func LoadGlobalConfig() (*GlobalConfig, error) {
	home, err := globalUserHomeDir()
	if err != nil {
		return DefaultConfig(), nil
	}
	configDir := filepath.Join(home, ".sshy")
	configPath, format := findConfigFile(configDir, GlobalConfigFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := DefaultConfig()
		err = SaveGlobalConfig(cfg)
		if err != nil {
			return nil, err
		}
		localPath := filepath.Join(configDir, "local.yaml")
		if _, err := os.Stat(localPath); os.IsNotExist(err) {
			defaultLocal := LocalConfig{
				Servers: make(map[string]models.Server),
				Private: make(models.Servers, 0),
			}
			data, err := Marshal(defaultLocal, FormatYAML)
			if err != nil {
				return nil, err
			}
			err = os.WriteFile(localPath, data, 0644)
			if err != nil {
				return nil, err
			}
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if format == FormatUnknown {
		format = DetectFormatFromContent(data)
	}

	var cfg GlobalConfig
	err = Unmarshal(data, format, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.ServersPath == "" && cfg.ServersURL == "" {
		cfg.ServersPath = "servers.yaml"
	}
	if cfg.ConfigPath == "" {
		cfg.ConfigPath = "."
	}
	return &cfg, nil
}

func SaveGlobalConfig(cfg *GlobalConfig) error {
	home, err := globalUserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".sshy")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	configPath, format := findConfigFile(configDir, GlobalConfigFile)
	if format == FormatUnknown {
		format = FormatYAML
	}

	data, err := Marshal(cfg, format)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func SaveGlobalConfigWithFormat(cfg *GlobalConfig, format FileFormat) error {
	home, err := globalUserHomeDir()
	if err != nil {
		return err
	}
	configDir := filepath.Join(home, ".sshy")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}

	var ext string
	switch format {
	case FormatJSON:
		ext = ".json"
	default:
		ext = ".yaml"
	}

	configPath := filepath.Join(configDir, "config"+ext)
	data, err := Marshal(cfg, format)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
