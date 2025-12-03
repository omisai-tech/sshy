package config

import (
	"os"
	"path/filepath"

	"github.com/omisai-tech/sshy/internal/models"
	"gopkg.in/yaml.v3"
)

type GlobalConfig struct {
	ServersPath string `yaml:"servers_path"`
	ConfigPath  string `yaml:"config_path"`
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
	configPath := filepath.Join(configDir, "config.yaml")

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
			data, err := yaml.Marshal(defaultLocal)
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
	var cfg GlobalConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.ServersPath == "" {
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
	configPath := filepath.Join(configDir, "config.yaml")
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}
