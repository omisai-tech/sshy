package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"http://example.com/servers.yaml", true},
		{"https://example.com/servers.json", true},
		{"HTTP://EXAMPLE.COM/servers.yaml", true},
		{"HTTPS://example.com/servers.yaml", true},
		{"/path/to/servers.yaml", false},
		{"servers.yaml", false},
		{"", false},
		{"ftp://example.com/servers.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsURL(tt.input)
			if result != tt.expected {
				t.Errorf("IsURL(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"http://example.com/servers.yaml", false},
		{"https://example.com/servers.json", false},
		{"https://internal.vpn.company.com/api/servers", false},
		{"", true},
		{"not-a-url", true},
		{"ftp://example.com/file", true},
		{"file:///local/file", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ValidateURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestFetchServersFromURL_YAML(t *testing.T) {
	servers := `
- name: server1
  host: 192.168.1.1
  user: admin
  port: 22
- name: server2
  host: 192.168.1.2
  user: root
  port: 2222
  tags:
    - production
    - web
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write([]byte(servers))
	}))
	defer server.Close()

	result, err := FetchServersFromURL(server.URL)
	if err != nil {
		t.Fatalf("FetchServersFromURL() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(result))
	}

	if result[0].Name != "server1" {
		t.Errorf("Expected first server name 'server1', got '%s'", result[0].Name)
	}

	if result[1].Port != 2222 {
		t.Errorf("Expected second server port 2222, got %d", result[1].Port)
	}

	if len(result[1].Tags) != 2 {
		t.Errorf("Expected second server to have 2 tags, got %d", len(result[1].Tags))
	}
}

func TestFetchServersFromURL_JSON(t *testing.T) {
	servers := []map[string]interface{}{
		{"name": "server1", "host": "192.168.1.1", "user": "admin", "port": 22},
		{"name": "server2", "host": "192.168.1.2", "user": "root", "port": 2222},
	}
	jsonData, _ := json.Marshal(servers)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer server.Close()

	result, err := FetchServersFromURL(server.URL)
	if err != nil {
		t.Fatalf("FetchServersFromURL() error = %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(result))
	}

	if result[0].Name != "server1" {
		t.Errorf("Expected first server name 'server1', got '%s'", result[0].Name)
	}
}

func TestFetchServersFromURL_AutoDetectFormat(t *testing.T) {
	yamlServers := `
- name: yaml-server
  host: 10.0.0.1
  user: admin
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(yamlServers))
	}))
	defer server.Close()

	result, err := FetchServersFromURL(server.URL)
	if err != nil {
		t.Fatalf("FetchServersFromURL() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 server, got %d", len(result))
	}
}

func TestFetchServersFromURL_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := FetchServersFromURL(server.URL)
	if err == nil {
		t.Error("Expected error for HTTP 500, got nil")
	}
}

func TestFetchServersFromURL_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := FetchServersFromURL(server.URL)
	if err == nil {
		t.Error("Expected error for HTTP 404, got nil")
	}
}

func TestFetchServersFromURL_InvalidURL(t *testing.T) {
	_, err := FetchServersFromURL("not-a-valid-url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestFetchServersFromURL_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result, err := FetchServersFromURL(server.URL)
	if err != nil {
		t.Fatalf("FetchServersFromURL() error = %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected 0 servers for empty response, got %d", len(result))
	}
}

func TestFetchServersFromURLWithSource(t *testing.T) {
	servers := `
- name: server1
  host: 192.168.1.1
  user: admin
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write([]byte(servers))
	}))
	defer server.Close()

	result, err := FetchServersFromURLWithSource(server.URL)
	if err != nil {
		t.Fatalf("FetchServersFromURLWithSource() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 server, got %d", len(result))
	}

	if result[0].Server.Name != "server1" {
		t.Errorf("Expected server name 'server1', got '%s'", result[0].Server.Name)
	}
}

func TestDetectFormatFromContentType(t *testing.T) {
	tests := []struct {
		contentType string
		expected    FileFormat
	}{
		{"application/json", FormatJSON},
		{"application/json; charset=utf-8", FormatJSON},
		{"text/json", FormatJSON},
		{"application/x-yaml", FormatYAML},
		{"text/yaml", FormatYAML},
		{"application/x-yml", FormatYAML},
		{"text/plain", FormatUnknown},
		{"", FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.contentType, func(t *testing.T) {
			result := detectFormatFromContentType(tt.contentType)
			if result != tt.expected {
				t.Errorf("detectFormatFromContentType(%q) = %v, want %v", tt.contentType, result, tt.expected)
			}
		})
	}
}

func TestDetectFormatFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected FileFormat
	}{
		{"https://example.com/servers.json", FormatJSON},
		{"https://example.com/servers.yaml", FormatYAML},
		{"https://example.com/servers.yml", FormatYAML},
		{"https://example.com/servers.JSON", FormatJSON},
		{"https://example.com/servers.YAML", FormatYAML},
		{"https://example.com/api/servers", FormatUnknown},
		{"https://example.com/servers", FormatUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := detectFormatFromURL(tt.url)
			if result != tt.expected {
				t.Errorf("detectFormatFromURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestLoadServersWithPath_URL(t *testing.T) {
	servers := `
- name: remote-server
  host: 10.0.0.1
  user: admin
  port: 22
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write([]byte(servers))
	}))
	defer server.Close()

	originalUserHomeDir := userHomeDir
	userHomeDir = func() (string, error) {
		return t.TempDir(), nil
	}
	defer func() { userHomeDir = originalUserHomeDir }()

	result, err := LoadServersWithPath("", server.URL)
	if err != nil {
		t.Fatalf("LoadServersWithPath() error = %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 server, got %d", len(result))
	}

	if result[0].Name != "remote-server" {
		t.Errorf("Expected server name 'remote-server', got '%s'", result[0].Name)
	}
}

func TestGlobalConfig_GetServersSource(t *testing.T) {
	tests := []struct {
		name     string
		cfg      GlobalConfig
		expected string
		isRemote bool
	}{
		{
			name:     "URL takes precedence",
			cfg:      GlobalConfig{ServersPath: "servers.yaml", ServersURL: "https://example.com/servers.yaml"},
			expected: "https://example.com/servers.yaml",
			isRemote: true,
		},
		{
			name:     "Path when no URL",
			cfg:      GlobalConfig{ServersPath: "servers.yaml", ServersURL: ""},
			expected: "servers.yaml",
			isRemote: false,
		},
		{
			name:     "Empty URL returns path",
			cfg:      GlobalConfig{ServersPath: "custom.json", ServersURL: ""},
			expected: "custom.json",
			isRemote: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cfg.GetServersSource()
			if result != tt.expected {
				t.Errorf("GetServersSource() = %q, want %q", result, tt.expected)
			}
			isRemote := tt.cfg.IsRemoteSource()
			if isRemote != tt.isRemote {
				t.Errorf("IsRemoteSource() = %v, want %v", isRemote, tt.isRemote)
			}
		})
	}
}

func TestFetchServersFromURL_InvalidYAML(t *testing.T) {
	invalidYAML := `
- name: server1
  host: 192.168.1.1
  - invalid yaml structure
`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		w.Write([]byte(invalidYAML))
	}))
	defer server.Close()

	_, err := FetchServersFromURL(server.URL)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestFetchServersFromURL_DetectFromURLPath(t *testing.T) {
	type serverData struct {
		Name string `json:"name" yaml:"name"`
		Host string `json:"host" yaml:"host"`
		User string `json:"user" yaml:"user"`
	}

	jsonServers := []serverData{
		{Name: "json-server", Host: "10.0.0.1", User: "admin"},
	}
	jsonData, _ := json.Marshal(jsonServers)

	yamlServers := []serverData{
		{Name: "yaml-server", Host: "10.0.0.2", User: "root"},
	}
	yamlData, _ := yaml.Marshal(yamlServers)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/servers.json" {
			w.Write(jsonData)
		} else if r.URL.Path == "/servers.yaml" {
			w.Write(yamlData)
		}
	}))
	defer server.Close()

	jsonResult, err := FetchServersFromURL(server.URL + "/servers.json")
	if err != nil {
		t.Fatalf("FetchServersFromURL(json) error = %v", err)
	}
	if len(jsonResult) != 1 || jsonResult[0].Name != "json-server" {
		t.Errorf("Expected json-server, got %+v", jsonResult)
	}

	yamlResult, err := FetchServersFromURL(server.URL + "/servers.yaml")
	if err != nil {
		t.Fatalf("FetchServersFromURL(yaml) error = %v", err)
	}
	if len(yamlResult) != 1 || yamlResult[0].Name != "yaml-server" {
		t.Errorf("Expected yaml-server, got %+v", yamlResult)
	}
}
