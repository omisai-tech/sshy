package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ServersPath != "servers.yaml" {
		t.Errorf("Expected ServersPath 'servers.yaml', got '%s'", cfg.ServersPath)
	}
	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, ".sshy")
	if cfg.ConfigPath != expectedPath {
		t.Errorf("Expected ConfigPath '%s', got '%s'", expectedPath, cfg.ConfigPath)
	}
}

func TestLoadGlobalConfig_CreatesDefault(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.Remove(configPath)

	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cfg.ServersPath != "servers.yaml" {
		t.Errorf("Expected ServersPath 'servers.yaml', got '%s'", cfg.ServersPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	localPath := filepath.Join(homeDir, ".sshy", "local.yaml")
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		t.Error("Expected local.yaml to be created")
	}
}

func TestLoadGlobalConfig_ExistingConfig(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	cfg := &GlobalConfig{
		ServersPath: "custom-servers.yaml",
		ConfigPath:  "/custom/path",
	}
	data, _ := yaml.Marshal(cfg)
	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.WriteFile(configPath, data, 0644)

	loaded, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if loaded.ServersPath != "custom-servers.yaml" {
		t.Errorf("Expected ServersPath 'custom-servers.yaml', got '%s'", loaded.ServersPath)
	}
	if loaded.ConfigPath != "/custom/path" {
		t.Errorf("Expected ConfigPath '/custom/path', got '%s'", loaded.ConfigPath)
	}
}

func TestLoadGlobalConfig_EmptyFields(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	cfg := &GlobalConfig{
		ServersPath: "",
		ConfigPath:  "",
	}
	data, _ := yaml.Marshal(cfg)
	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.WriteFile(configPath, data, 0644)

	loaded, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if loaded.ServersPath != "servers.yaml" {
		t.Errorf("Expected ServersPath 'servers.yaml' (default), got '%s'", loaded.ServersPath)
	}
	if loaded.ConfigPath != "." {
		t.Errorf("Expected ConfigPath '.' (default), got '%s'", loaded.ConfigPath)
	}
}

func TestLoadGlobalConfig_InvalidYAML(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.WriteFile(configPath, []byte("invalid: [yaml"), 0644)

	_, err := LoadGlobalConfig()
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadGlobalConfig_ReadError(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.WriteFile(configPath, []byte("servers_path: test"), 0644)
	os.Chmod(configPath, 0000)
	defer os.Chmod(configPath, 0644)

	_, err := LoadGlobalConfig()
	if err == nil {
		t.Error("Expected error for unreadable file, got nil")
	}
}

func TestSaveGlobalConfig(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	cfg := &GlobalConfig{
		ServersPath: "test-servers.yaml",
		ConfigPath:  "/test/path",
	}

	err := SaveGlobalConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	var loaded GlobalConfig
	yaml.Unmarshal(data, &loaded)

	if loaded.ServersPath != "test-servers.yaml" {
		t.Errorf("Expected ServersPath 'test-servers.yaml', got '%s'", loaded.ServersPath)
	}
}

func TestSaveGlobalConfig_CreatesDirectory(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sshyDir := filepath.Join(homeDir, ".sshy")
	os.RemoveAll(sshyDir)

	cfg := &GlobalConfig{
		ServersPath: "servers.yaml",
		ConfigPath:  ".",
	}

	err := SaveGlobalConfig(cfg)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, err := os.Stat(sshyDir); os.IsNotExist(err) {
		t.Error("Expected .sshy directory to be created")
	}
}

func TestSaveGlobalConfig_WriteError(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sshyDir := filepath.Join(homeDir, ".sshy")
	os.Chmod(sshyDir, 0444)
	defer os.Chmod(sshyDir, 0755)

	cfg := &GlobalConfig{
		ServersPath: "servers.yaml",
		ConfigPath:  ".",
	}

	err := SaveGlobalConfig(cfg)
	if err == nil {
		t.Error("Expected error for unwritable directory, got nil")
	}
}

func TestLoadGlobalConfig_SaveError(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	configPath := filepath.Join(homeDir, ".sshy", "config.yaml")
	os.Remove(configPath)

	sshyDir := filepath.Join(homeDir, ".sshy")
	os.Chmod(sshyDir, 0555)
	defer os.Chmod(sshyDir, 0755)

	_, err := LoadGlobalConfig()
	if err == nil {
		t.Error("Expected error when cannot save default config, got nil")
	}
}

func TestSaveGlobalConfig_UserHomeDirError(t *testing.T) {
	oldFunc := globalUserHomeDir
	globalUserHomeDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { globalUserHomeDir = oldFunc }()

	cfg := &GlobalConfig{ServersPath: "test.yaml", ConfigPath: "."}
	err := SaveGlobalConfig(cfg)
	if err == nil {
		t.Error("Expected error when UserHomeDir fails, got nil")
	}
}

func TestLoadGlobalConfig_UserHomeDirError(t *testing.T) {
	oldFunc := globalUserHomeDir
	globalUserHomeDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { globalUserHomeDir = oldFunc }()

	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Errorf("Expected no error (should return default), got %v", err)
	}
	if cfg == nil {
		t.Error("Expected default config, got nil")
	}
}

func TestDefaultConfig_UserHomeDirError(t *testing.T) {
	oldFunc := globalUserHomeDir
	globalUserHomeDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { globalUserHomeDir = oldFunc }()

	cfg := DefaultConfig()
	if cfg.ServersPath != "servers.yaml" {
		t.Errorf("Expected ServersPath 'servers.yaml', got '%s'", cfg.ServersPath)
	}
}
