package cmd

import (
	"testing"

	"github.com/omisai-tech/sshy/internal/models"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedServer string
		expectedPath   string
	}{
		{
			name:           "remote path with server",
			path:           "server1:/path/to/file",
			expectedServer: "server1",
			expectedPath:   "/path/to/file",
		},
		{
			name:           "local path",
			path:           "/local/path/file.txt",
			expectedServer: "",
			expectedPath:   "/local/path/file.txt",
		},
		{
			name:           "relative local path",
			path:           "relative/path/file.txt",
			expectedServer: "",
			expectedPath:   "relative/path/file.txt",
		},
		{
			name:           "server with empty path",
			path:           "server:",
			expectedServer: "server",
			expectedPath:   "",
		},
		{
			name:           "path with multiple colons",
			path:           "server:/path:with:colons",
			expectedServer: "server",
			expectedPath:   "/path:with:colons",
		},
		{
			name:           "empty path",
			path:           "",
			expectedServer: "",
			expectedPath:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, path := parsePath(tt.path)
			if server != tt.expectedServer {
				t.Errorf("parsePath(%q) server = %q, expected %q", tt.path, server, tt.expectedServer)
			}
			if path != tt.expectedPath {
				t.Errorf("parsePath(%q) path = %q, expected %q", tt.path, path, tt.expectedPath)
			}
		})
	}
}

func TestBuildScpTarget(t *testing.T) {
	tests := []struct {
		name       string
		server     models.Server
		remotePath string
		sshArgs    []string
		expected   string
	}{
		{
			name:       "basic server with user",
			server:     models.Server{Host: "example.com", User: "admin"},
			remotePath: "/path/to/file",
			sshArgs:    []string{},
			expected:   "admin@example.com:/path/to/file",
		},
		{
			name:       "server without user",
			server:     models.Server{Host: "example.com"},
			remotePath: "/path/to/file",
			sshArgs:    []string{},
			expected:   "example.com:/path/to/file",
		},
		{
			name:       "user override with -l flag",
			server:     models.Server{Host: "example.com", User: "original"},
			remotePath: "/path/to/file",
			sshArgs:    []string{"-l", "override"},
			expected:   "override@example.com:/path/to/file",
		},
		{
			name:       "user override with -l flag and other args",
			server:     models.Server{Host: "example.com", User: "original"},
			remotePath: "/path/to/file",
			sshArgs:    []string{"-v", "-l", "override", "-C"},
			expected:   "override@example.com:/path/to/file",
		},
		{
			name:       "-l flag at end without value",
			server:     models.Server{Host: "example.com", User: "admin"},
			remotePath: "/path",
			sshArgs:    []string{"-v", "-l"},
			expected:   "admin@example.com:/path",
		},
		{
			name:       "empty remote path",
			server:     models.Server{Host: "example.com", User: "user"},
			remotePath: "",
			sshArgs:    []string{},
			expected:   "user@example.com:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildScpTarget(tt.server, tt.remotePath, tt.sshArgs)
			if result != tt.expected {
				t.Errorf("buildScpTarget() = %q, expected %q", result, tt.expected)
			}
		})
	}
}
