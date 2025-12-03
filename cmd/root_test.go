package cmd

import (
	"errors"
	"testing"
)

func TestIsSubcommand(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected bool
	}{
		{"valid command list", "list", true},
		{"valid command add", "add", true},
		{"valid command rm", "rm", true},
		{"valid command edit", "edit", true},
		{"valid command connect", "connect", true},
		{"valid command init", "init", true},
		{"valid command local", "local", true},
		{"valid command scp", "scp", true},
		{"valid command sftp", "sftp", true},
		{"valid command view", "view", true},
		{"invalid command", "invalid", false},
		{"empty string", "", false},
		{"random string", "foobar", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSubcommand(tt.arg)
			if result != tt.expected {
				t.Errorf("isSubcommand(%q) = %v, expected %v", tt.arg, result, tt.expected)
			}
		})
	}
}

func TestIsFlag(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected bool
	}{
		{"short flag", "-h", true},
		{"long flag", "--help", true},
		{"version flag", "-v", true},
		{"long version flag", "--version", true},
		{"not a flag", "command", false},
		{"empty string", "", false},
		{"single dash", "-", false},
		{"double dash only", "--", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isFlag(tt.arg)
			if result != tt.expected {
				t.Errorf("isFlag(%q) = %v, expected %v", tt.arg, result, tt.expected)
			}
		})
	}
}

func TestSetVersionInfo(t *testing.T) {
	SetVersionInfo("1.0.0", "abc123", "2024-01-01")

	if version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", version)
	}
	if rootCmd.Version != "1.0.0" {
		t.Errorf("Expected rootCmd.Version '1.0.0', got '%s'", rootCmd.Version)
	}
}

func TestRootCmdConfiguration(t *testing.T) {
	if rootCmd.Use != "sshy" {
		t.Errorf("Expected Use 'sshy', got '%s'", rootCmd.Use)
	}
	if rootCmd.Short != "Manage SSH servers via YAML config" {
		t.Errorf("Unexpected Short description: %s", rootCmd.Short)
	}
}

func TestExecuteWithArgs_SubcommandVersion(t *testing.T) {
	exitCalled := false
	exitCode := 0
	oldExit := osExit
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { osExit = oldExit }()

	ExecuteWithArgs([]string{"sshy", "--version"})

	if exitCalled && exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}
}

func TestExecuteWithArgs_SubcommandList(t *testing.T) {
	ExecuteWithArgs([]string{"sshy", "list"})
}

func TestExecuteWithArgs_NoArgs(t *testing.T) {
	oldRunner := cmdRunner
	defer func() { cmdRunner = oldRunner }()

	oldFuzzy := fuzzyFind
	defer func() { fuzzyFind = oldFuzzy }()

	fuzzyFind = func(names []string, itemFunc func(int) string) (int, error) {
		if len(names) > 0 {
			return 0, nil
		}
		return -1, errors.New("cancelled")
	}

	mock := &MockCommandRunner{}
	cmdRunner = mock

	ExecuteWithArgs([]string{"sshy"})
}

func TestExecuteWithArgs_UnknownArg(t *testing.T) {
	oldRunner := cmdRunner
	defer func() { cmdRunner = oldRunner }()

	oldFuzzy := fuzzyFind
	defer func() { fuzzyFind = oldFuzzy }()

	fuzzyFind = func(names []string, itemFunc func(int) string) (int, error) {
		return -1, errors.New("cancelled")
	}

	mock := &MockCommandRunner{}
	cmdRunner = mock

	ExecuteWithArgs([]string{"sshy", "someserver"})
}

func TestExecuteWithArgs_InvalidCommand(t *testing.T) {
	exitCalled := false
	exitCode := 0
	oldExit := osExit
	osExit = func(code int) {
		exitCalled = true
		exitCode = code
	}
	defer func() { osExit = oldExit }()

	rootCmd.SetArgs([]string{"--invalid-flag-that-does-not-exist"})
	ExecuteWithArgs([]string{"sshy", "--invalid-flag-that-does-not-exist"})
	rootCmd.SetArgs([]string{})

	if !exitCalled {
		t.Error("Expected exit to be called for invalid command")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
}

func TestExecute(t *testing.T) {
	oldRunner := cmdRunner
	defer func() { cmdRunner = oldRunner }()

	oldFuzzy := fuzzyFind
	defer func() { fuzzyFind = oldFuzzy }()

	fuzzyFind = func(names []string, itemFunc func(int) string) (int, error) {
		return -1, errors.New("cancelled")
	}

	mock := &MockCommandRunner{}
	cmdRunner = mock

	Execute()
}
