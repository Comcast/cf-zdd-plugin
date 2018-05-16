package commands

import "code.cloudfoundry.org/cli/plugin"

// CfZddCmd - struct to initialize.
type CfZddCmd struct {
	Conn            plugin.CliConnection
	CmdName         string
	OldApp          string
	NewApp          string
	ManifestPath    string
	ApplicationPath string
	Duration        string
	CustomURL       string
	BatchSize       int
	RouteCheck      bool
	HelpTopic       string
}

// const - exported constants
const (
	CanaryRouteSuffix    = "canary"
	CanaryRouteSeparator = "-"
)
