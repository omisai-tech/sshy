package models

import "testing"

func TestServerStruct(t *testing.T) {
	server := Server{
		Name:    "test-server",
		Host:    "192.168.1.1",
		User:    "admin",
		Port:    2222,
		Tags:    []string{"prod", "web"},
		Key:     "~/.ssh/id_rsa",
		Options: map[string]interface{}{"ForwardAgent": "yes"},
	}

	if server.Name != "test-server" {
		t.Errorf("Expected name 'test-server', got '%s'", server.Name)
	}
	if server.Host != "192.168.1.1" {
		t.Errorf("Expected host '192.168.1.1', got '%s'", server.Host)
	}
	if server.User != "admin" {
		t.Errorf("Expected user 'admin', got '%s'", server.User)
	}
	if server.Port != 2222 {
		t.Errorf("Expected port 2222, got %d", server.Port)
	}
	if len(server.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(server.Tags))
	}
	if server.Key != "~/.ssh/id_rsa" {
		t.Errorf("Expected key '~/.ssh/id_rsa', got '%s'", server.Key)
	}
	if server.Options["ForwardAgent"] != "yes" {
		t.Errorf("Expected ForwardAgent 'yes', got '%v'", server.Options["ForwardAgent"])
	}
}

func TestServersSlice(t *testing.T) {
	servers := Servers{
		{Name: "server1", Host: "host1"},
		{Name: "server2", Host: "host2"},
	}

	if len(servers) != 2 {
		t.Errorf("Expected 2 servers, got %d", len(servers))
	}
	if servers[0].Name != "server1" {
		t.Errorf("Expected first server name 'server1', got '%s'", servers[0].Name)
	}
}

func TestServerSource(t *testing.T) {
	if SourceShared != 0 {
		t.Errorf("Expected SourceShared to be 0, got %d", SourceShared)
	}
	if SourceLocal != 1 {
		t.Errorf("Expected SourceLocal to be 1, got %d", SourceLocal)
	}
	if SourceOverride != 2 {
		t.Errorf("Expected SourceOverride to be 2, got %d", SourceOverride)
	}
}

func TestServerWithSource(t *testing.T) {
	server := Server{Name: "test", Host: "localhost"}
	sws := ServerWithSource{
		Server: server,
		Source: SourceLocal,
	}

	if sws.Server.Name != "test" {
		t.Errorf("Expected server name 'test', got '%s'", sws.Server.Name)
	}
	if sws.Source != SourceLocal {
		t.Errorf("Expected source SourceLocal, got %d", sws.Source)
	}
}

func TestServerDefaultValues(t *testing.T) {
	server := Server{}

	if server.Name != "" {
		t.Errorf("Expected empty name, got '%s'", server.Name)
	}
	if server.Port != 0 {
		t.Errorf("Expected port 0, got %d", server.Port)
	}
	if server.Tags != nil {
		t.Errorf("Expected nil tags, got %v", server.Tags)
	}
	if server.Options != nil {
		t.Errorf("Expected nil options, got %v", server.Options)
	}
}
