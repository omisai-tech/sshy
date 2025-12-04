package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/omisai-tech/sshy/internal/models"
	"gopkg.in/yaml.v3"
)

func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "sshy-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir, func() { os.RemoveAll(tmpDir) }
}

func setupTestHomeDir(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "sshy-home-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	sshyDir := filepath.Join(tmpDir, ".sshy")
	if err := os.MkdirAll(sshyDir, 0755); err != nil {
		t.Fatalf("Failed to create .sshy dir: %v", err)
	}

	return tmpDir, func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}
}

func TestLoadServers_EmptyFile(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServers(configDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(servers) != 0 {
		t.Errorf("Expected 0 servers, got %d", len(servers))
	}
}

func TestLoadServers_WithSharedServers(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sharedServers := models.Servers{
		{Name: "server1", Host: "host1", User: "user1", Port: 22},
		{Name: "server2", Host: "host2", User: "user2", Port: 2222},
	}
	sharedData, _ := yaml.Marshal(sharedServers)
	os.WriteFile(filepath.Join(configDir, SharedConfigFile), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServers(configDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}
	if servers[0].Name != "server1" {
		t.Errorf("Expected first server name 'server1', got '%s'", servers[0].Name)
	}
}

func TestLoadServers_WithLocalOverrides(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sharedServers := models.Servers{
		{Name: "server1", Host: "host1", User: "user1", Port: 22},
	}
	sharedData, _ := yaml.Marshal(sharedServers)
	os.WriteFile(filepath.Join(configDir, SharedConfigFile), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: map[string]models.Server{
			"server1": {Host: "override-host", User: "override-user", Port: 3333, Key: "/path/to/key", Tags: []string{"new-tag"}, Options: map[string]interface{}{"opt": "val"}},
		},
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServers(configDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].Host != "override-host" {
		t.Errorf("Expected host 'override-host', got '%s'", servers[0].Host)
	}
	if servers[0].User != "override-user" {
		t.Errorf("Expected user 'override-user', got '%s'", servers[0].User)
	}
	if servers[0].Port != 3333 {
		t.Errorf("Expected port 3333, got %d", servers[0].Port)
	}
	if servers[0].Key != "/path/to/key" {
		t.Errorf("Expected key '/path/to/key', got '%s'", servers[0].Key)
	}
}

func TestLoadServers_WithPrivateServers(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{
			{Name: "private1", Host: "private-host", User: "private-user"},
		},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServers(configDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}
	if servers[0].Name != "private1" {
		t.Errorf("Expected server name 'private1', got '%s'", servers[0].Name)
	}
}

func TestLoadServersWithPath(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	customPath := "custom-servers.yaml"
	sharedServers := models.Servers{
		{Name: "custom-server", Host: "custom-host"},
	}
	sharedData, _ := yaml.Marshal(sharedServers)
	os.WriteFile(filepath.Join(configDir, customPath), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServersWithPath(configDir, customPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(servers))
	}
}

func TestLoadServers_InvalidYAML(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(configDir, SharedConfigFile), []byte("invalid: yaml: content: ["), 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	_, err := LoadServers(configDir)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadServers_ReadError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	serverFile := filepath.Join(configDir, SharedConfigFile)
	os.WriteFile(serverFile, []byte("- name: test"), 0644)
	os.Chmod(serverFile, 0000)
	defer os.Chmod(serverFile, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	_, err := LoadServers(configDir)
	if err == nil {
		t.Error("Expected error for unreadable file, got nil")
	}
}

func TestSaveServers(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers := models.Servers{
		{Name: "server1", Host: "host1", User: "user1", Key: "should-be-removed"},
	}

	err := SaveServers(configDir, servers)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(configDir, SharedConfigFile))
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	var savedServers models.Servers
	yaml.Unmarshal(data, &savedServers)

	if len(savedServers) != 1 {
		t.Errorf("Expected 1 saved server, got %d", len(savedServers))
	}
	if savedServers[0].Key != "" {
		t.Errorf("Expected key to be empty (not saved), got '%s'", savedServers[0].Key)
	}
}

func TestSaveServers_ExcludesPrivate(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{
			{Name: "private-server", Host: "private-host"},
		},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers := models.Servers{
		{Name: "shared-server", Host: "shared-host"},
		{Name: "private-server", Host: "private-host"},
	}

	err := SaveServers(configDir, servers)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(configDir, SharedConfigFile))
	var savedServers models.Servers
	yaml.Unmarshal(data, &savedServers)

	if len(savedServers) != 1 {
		t.Errorf("Expected 1 saved server (private excluded), got %d", len(savedServers))
	}
	if savedServers[0].Name != "shared-server" {
		t.Errorf("Expected 'shared-server', got '%s'", savedServers[0].Name)
	}
}

func TestSaveServersWithPath(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers := models.Servers{
		{Name: "server1", Host: "host1"},
	}

	customPath := "custom-save.yaml"
	err := SaveServersWithPath(configDir, customPath, servers)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	_, err = os.Stat(filepath.Join(configDir, customPath))
	if err != nil {
		t.Errorf("Expected file to exist: %v", err)
	}
}

func TestLoadLocalConfig(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: map[string]models.Server{
			"test": {Host: "test-host"},
		},
		Private: models.Servers{
			{Name: "private", Host: "private-host"},
		},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	loaded, err := LoadLocalConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(loaded.Servers) != 1 {
		t.Errorf("Expected 1 server override, got %d", len(loaded.Servers))
	}
	if len(loaded.Private) != 1 {
		t.Errorf("Expected 1 private server, got %d", len(loaded.Private))
	}
}

func TestLoadLocalConfig_NotExists(t *testing.T) {
	_, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	loaded, err := LoadLocalConfig()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(loaded.Servers) != 0 {
		t.Errorf("Expected empty servers, got %d", len(loaded.Servers))
	}
}

func TestLoadLocalConfig_InvalidYAML(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), []byte("invalid: [yaml"), 0644)

	_, err := LoadLocalConfig()
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestSaveLocalConfig(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: map[string]models.Server{
			"test": {Host: "test-host"},
		},
		Private: models.Servers{
			{Name: "private", Host: "private-host"},
		},
	}

	err := SaveLocalConfig(localConfig)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(homeDir, ".sshy", "local.yaml"))
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	var loaded LocalConfig
	yaml.Unmarshal(data, &loaded)

	if len(loaded.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(loaded.Servers))
	}
}

func TestLoadServersWithSource(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sharedServers := models.Servers{
		{Name: "shared", Host: "shared-host"},
		{Name: "overridden", Host: "original-host"},
	}
	sharedData, _ := yaml.Marshal(sharedServers)
	os.WriteFile(filepath.Join(configDir, SharedConfigFile), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: map[string]models.Server{
			"overridden": {Host: "override-host", User: "override-user", Port: 3333, Tags: []string{"tag"}, Key: "key", Options: map[string]interface{}{"opt": "val"}},
		},
		Private: models.Servers{
			{Name: "private", Host: "private-host"},
		},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	serversWithSource, err := LoadServersWithSource(configDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(serversWithSource) != 3 {
		t.Errorf("Expected 3 servers, got %d", len(serversWithSource))
	}

	for _, sws := range serversWithSource {
		switch sws.Server.Name {
		case "shared":
			if sws.Source != models.SourceShared {
				t.Errorf("Expected SourceShared for 'shared', got %d", sws.Source)
			}
		case "overridden":
			if sws.Source != models.SourceOverride {
				t.Errorf("Expected SourceOverride for 'overridden', got %d", sws.Source)
			}
			if sws.Server.Host != "override-host" {
				t.Errorf("Expected override host, got '%s'", sws.Server.Host)
			}
		case "private":
			if sws.Source != models.SourceLocal {
				t.Errorf("Expected SourceLocal for 'private', got %d", sws.Source)
			}
		}
	}
}

func TestLoadServersWithSourceAndPath(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	customPath := "custom.yaml"
	sharedServers := models.Servers{
		{Name: "custom", Host: "custom-host"},
	}
	sharedData, _ := yaml.Marshal(sharedServers)
	os.WriteFile(filepath.Join(configDir, customPath), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	serversWithSource, err := LoadServersWithSourceAndPath(configDir, customPath)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(serversWithSource) != 1 {
		t.Errorf("Expected 1 server, got %d", len(serversWithSource))
	}
}

func TestLoadServersWithSource_InvalidYAML(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(configDir, SharedConfigFile), []byte("invalid: [yaml"), 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	_, err := LoadServersWithSource(configDir)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestLoadServersWithSource_ReadError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	serverFile := filepath.Join(configDir, SharedConfigFile)
	os.WriteFile(serverFile, []byte("- name: test"), 0644)
	os.Chmod(serverFile, 0000)
	defer os.Chmod(serverFile, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	_, err := LoadServersWithSource(configDir)
	if err == nil {
		t.Error("Expected error for unreadable file, got nil")
	}
}

func TestLoadServersWithSource_LocalConfigError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), []byte("invalid: [yaml"), 0644)

	_, err := LoadServersWithSource(configDir)
	if err == nil {
		t.Error("Expected error for invalid local config, got nil")
	}
}

func TestLoadServers_LocalConfigError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), []byte("invalid: [yaml"), 0644)

	_, err := LoadServers(configDir)
	if err == nil {
		t.Error("Expected error for invalid local config, got nil")
	}
}

func TestSaveServers_LocalConfigError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), []byte("invalid: [yaml"), 0644)

	servers := models.Servers{
		{Name: "test", Host: "test-host"},
	}
	err := SaveServers(configDir, servers)
	if err == nil {
		t.Error("Expected error for invalid local config, got nil")
	}
}

func TestLoadLocalConfig_ReadError(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localPath := filepath.Join(homeDir, ".sshy", "local.yaml")
	os.WriteFile(localPath, []byte("servers: {}"), 0644)
	os.Chmod(localPath, 0000)
	defer os.Chmod(localPath, 0644)

	_, err := LoadLocalConfig()
	if err == nil {
		t.Error("Expected error for unreadable file, got nil")
	}
}

func TestSaveServersWithPath_WriteError(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	os.Chmod(configDir, 0444)
	defer os.Chmod(configDir, 0755)

	servers := models.Servers{
		{Name: "test", Host: "test-host"},
	}
	err := SaveServersWithPath(configDir, "servers.yaml", servers)
	if err == nil {
		t.Error("Expected error for unwritable directory, got nil")
	}
}

func TestSaveLocalConfig_WriteError(t *testing.T) {
	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sshyDir := filepath.Join(homeDir, ".sshy")
	os.Chmod(sshyDir, 0444)
	defer os.Chmod(sshyDir, 0755)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	err := SaveLocalConfig(localConfig)
	if err == nil {
		t.Error("Expected error for unwritable directory, got nil")
	}
}

func TestSaveLocalConfig_UserHomeDirError(t *testing.T) {
	oldFunc := userHomeDir
	userHomeDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { userHomeDir = oldFunc }()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	err := SaveLocalConfig(localConfig)
	if err == nil {
		t.Error("Expected error when UserHomeDir fails, got nil")
	}
}

func TestLoadLocalConfig_UserHomeDirError(t *testing.T) {
	oldFunc := userHomeDir
	userHomeDir = func() (string, error) {
		return "", os.ErrNotExist
	}
	defer func() { userHomeDir = oldFunc }()

	_, err := loadLocalConfig()
	if err == nil {
		t.Error("Expected error when UserHomeDir fails, got nil")
	}
}

func TestLoadServers_WithJSONConfig(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sharedServers := models.Servers{
		{Name: "json-server1", Host: "json-host1", User: "json-user1", Port: 22},
		{Name: "json-server2", Host: "json-host2", User: "json-user2", Port: 2222},
	}
	sharedData, _ := Marshal(sharedServers, FormatJSON)
	os.WriteFile(filepath.Join(configDir, "servers.json"), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := Marshal(localConfig, FormatJSON)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.json"), localData, 0644)

	servers, err := LoadServersWithPath(configDir, "servers.json")
	if err != nil {
		t.Fatalf("Unexpected error loading JSON config: %v", err)
	}
	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}
	if servers[0].Name != "json-server1" {
		t.Errorf("Expected first server name 'json-server1', got '%s'", servers[0].Name)
	}
}

func TestLoadServers_FallbackToJSON(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	sharedServers := models.Servers{
		{Name: "fallback-server", Host: "fallback-host", User: "fallback-user", Port: 22},
	}
	sharedData, _ := Marshal(sharedServers, FormatJSON)
	os.WriteFile(filepath.Join(configDir, "servers.json"), sharedData, 0644)

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	servers, err := LoadServersWithPath(configDir, "servers.yaml")
	if err != nil {
		t.Fatalf("Unexpected error with JSON fallback: %v", err)
	}
	if len(servers) != 1 {
		t.Errorf("Expected 1 server (fallback to JSON), got %d", len(servers))
	}
	if servers[0].Name != "fallback-server" {
		t.Errorf("Expected server name 'fallback-server', got '%s'", servers[0].Name)
	}
}

func TestSaveServers_JSONFormat(t *testing.T) {
	configDir, cleanup := setupTestDir(t)
	defer cleanup()

	homeDir, homeCleanup := setupTestHomeDir(t)
	defer homeCleanup()

	localConfig := LocalConfig{
		Servers: make(map[string]models.Server),
		Private: models.Servers{},
	}
	localData, _ := yaml.Marshal(localConfig)
	os.WriteFile(filepath.Join(homeDir, ".sshy", "local.yaml"), localData, 0644)

	existingJSON := models.Servers{
		{Name: "existing", Host: "existing-host", User: "user", Port: 22},
	}
	existingData, _ := Marshal(existingJSON, FormatJSON)
	os.WriteFile(filepath.Join(configDir, "servers.json"), existingData, 0644)

	servers := models.Servers{
		{Name: "new-server", Host: "new-host", User: "new-user", Port: 3333},
	}
	err := SaveServersWithPath(configDir, "servers.json", servers)
	if err != nil {
		t.Fatalf("SaveServersWithPath failed: %v", err)
	}

	savedData, err := os.ReadFile(filepath.Join(configDir, "servers.json"))
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	format := DetectFormatFromContent(savedData)
	if format != FormatJSON {
		t.Errorf("Expected saved file to be JSON format, got %v", format)
	}

	var savedServers models.Servers
	err = Unmarshal(savedData, FormatJSON, &savedServers)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved JSON: %v", err)
	}
	if len(savedServers) != 1 || savedServers[0].Name != "new-server" {
		t.Errorf("Saved servers don't match expected: %+v", savedServers)
	}
}
