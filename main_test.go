package main

import (
	"os"
	"testing"

	"github.com/omisai-tech/sshy/cmd"
)

func TestSetVersionInfo(t *testing.T) {
	cmd.SetVersionInfo("1.0.0", "abc123", "2024-01-01")
}

func TestRun(t *testing.T) {
	run()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
