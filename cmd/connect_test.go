package cmd

import (
	"errors"
	"testing"

	"github.com/omisai-tech/sshy/internal/models"
)

type MockCommandRunner struct {
	LastCommand string
	LastArgs    []string
	ShouldError bool
}

func (m *MockCommandRunner) Run(name string, args []string) error {
	m.LastCommand = name
	m.LastArgs = args
	if m.ShouldError {
		return errors.New("mock error")
	}
	return nil
}

func TestConnectTo(t *testing.T) {
	oldRunner := cmdRunner
	defer func() { cmdRunner = oldRunner }()

	mock := &MockCommandRunner{}
	cmdRunner = mock

	server := models.Server{
		Name: "test",
		Host: "example.com",
		User: "admin",
		Key:  "~/.ssh/id_rsa",
	}

	connectTo(server, []string{"-v"}, "ls -la")

	if mock.LastCommand != "ssh" {
		t.Errorf("Expected command 'ssh', got '%s'", mock.LastCommand)
	}

	if !contains(mock.LastArgs, "-i") {
		t.Error("Expected -i flag in args")
	}
	if !contains(mock.LastArgs, "~/.ssh/id_rsa") {
		t.Error("Expected key path in args")
	}
	if !contains(mock.LastArgs, "-v") {
		t.Error("Expected -v flag in args")
	}
	if !contains(mock.LastArgs, "admin@example.com") {
		t.Error("Expected user@host in args")
	}
	if !contains(mock.LastArgs, "ls -la") {
		t.Error("Expected remote command in args")
	}
}

func TestConnectTo_Error(t *testing.T) {
	oldRunner := cmdRunner
	defer func() { cmdRunner = oldRunner }()

	mock := &MockCommandRunner{ShouldError: true}
	cmdRunner = mock

	server := models.Server{
		Name: "test",
		Host: "example.com",
		User: "admin",
	}

	connectTo(server, []string{}, "")
}

func TestDefaultCommandRunner(t *testing.T) {
	runner := &DefaultCommandRunner{}
	err := runner.Run("echo", []string{"test"})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestBuildSSHArgs(t *testing.T) {
	tests := []struct {
		name          string
		server        models.Server
		sshArgs       []string
		remoteCommand string
		checkFunc     func(t *testing.T, args []string)
	}{
		{
			name: "basic server with key",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
				Key:  "~/.ssh/id_rsa",
				Port: 22,
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				hasKey := containsSequence(args, "-i", "~/.ssh/id_rsa")
				if !hasKey {
					t.Error("Expected -i flag with key")
				}
				if !contains(args, "admin@example.com") {
					t.Error("Expected user@host in args")
				}
			},
		},
		{
			name: "server with custom port",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
				Port: 2222,
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				hasPort := containsSequence(args, "-p", "2222")
				if !hasPort {
					t.Error("Expected -p flag with port 2222")
				}
			},
		},
		{
			name: "server without user",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				Port: 22,
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if !contains(args, "example.com") {
					t.Error("Expected host without user prefix")
				}
				if contains(args, "@") {
					for _, arg := range args {
						if arg == "example.com" {
							return
						}
					}
					t.Error("Should not have @ in host when user is empty")
				}
			},
		},
		{
			name: "server with user override via -l flag",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
				Port: 22,
			},
			sshArgs:       []string{"-l", "root"},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if contains(args, "admin@example.com") {
					t.Error("Should not prepend user when -l flag is used")
				}
				if !contains(args, "example.com") {
					t.Error("Expected host in args")
				}
			},
		},
		{
			name: "server with remote command",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
			},
			sshArgs:       []string{},
			remoteCommand: "ls -la",
			checkFunc: func(t *testing.T, args []string) {
				if !contains(args, "ls -la") {
					t.Error("Expected remote command in args")
				}
			},
		},
		{
			name: "server with ForwardAgent option",
			server: models.Server{
				Name:    "test",
				Host:    "example.com",
				User:    "admin",
				Options: map[string]interface{}{"ForwardAgent": "yes"},
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if !contains(args, "-A") {
					t.Error("Expected -A flag for ForwardAgent")
				}
			},
		},
		{
			name: "server with ForwardAgent no",
			server: models.Server{
				Name:    "test",
				Host:    "example.com",
				User:    "admin",
				Options: map[string]interface{}{"ForwardAgent": "no"},
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if contains(args, "-A") {
					t.Error("Should not have -A flag when ForwardAgent is no")
				}
			},
		},
		{
			name: "server with RequestTTY force",
			server: models.Server{
				Name:    "test",
				Host:    "example.com",
				User:    "admin",
				Options: map[string]interface{}{"RequestTTY": "force"},
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				count := 0
				for _, arg := range args {
					if arg == "-t" {
						count++
					}
				}
				if count != 2 {
					t.Errorf("Expected two -t flags for RequestTTY force, got %d", count)
				}
			},
		},
		{
			name: "server with RequestTTY yes",
			server: models.Server{
				Name:    "test",
				Host:    "example.com",
				User:    "admin",
				Options: map[string]interface{}{"RequestTTY": "yes"},
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				count := 0
				for _, arg := range args {
					if arg == "-t" {
						count++
					}
				}
				if count != 1 {
					t.Errorf("Expected one -t flag for RequestTTY yes, got %d", count)
				}
			},
		},
		{
			name: "server with LocalForward option",
			server: models.Server{
				Name:    "test",
				Host:    "example.com",
				User:    "admin",
				Options: map[string]interface{}{"LocalForward": "8080:localhost:80"},
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if !containsSequence(args, "-L", "8080:localhost:80") {
					t.Error("Expected -L flag with LocalForward value")
				}
			},
		},
		{
			name: "server with additional ssh args",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
			},
			sshArgs:       []string{"-v", "-C"},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if !contains(args, "-v") || !contains(args, "-C") {
					t.Error("Expected additional ssh args in output")
				}
			},
		},
		{
			name: "server with port 0",
			server: models.Server{
				Name: "test",
				Host: "example.com",
				User: "admin",
				Port: 0,
			},
			sshArgs:       []string{},
			remoteCommand: "",
			checkFunc: func(t *testing.T, args []string) {
				if contains(args, "-p") {
					t.Error("Should not have -p flag when port is 0")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := buildSSHArgs(tt.server, tt.sshArgs, tt.remoteCommand)
			tt.checkFunc(t, args)
		})
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsSequence(slice []string, first, second string) bool {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i] == first && slice[i+1] == second {
			return true
		}
	}
	return false
}

func TestConnectCmdHelp(t *testing.T) {
	if connectCmd.Use != "connect [name] [ssh-flags...] [command]" {
		t.Errorf("Unexpected Use: %s", connectCmd.Use)
	}
	if connectCmd.Short != "Connect to an SSH server" {
		t.Errorf("Unexpected Short: %s", connectCmd.Short)
	}
}
