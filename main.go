package main

import "github.com/omisai-tech/sshy/cmd"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func run() {
	cmd.SetVersionInfo(version, commit, date)
	cmd.Execute()
}

func main() {
	run()
}
