package canaryrepo

import "github.com/cloudfoundry/cli/plugin"

var registry = make(map[string]PluginRunnable)

type PluginRunnable interface {
	Run(conn plugin.CliConnection) error
	SetArgs(args []string)
}

func Register(canaryName string, canaryRunnable PluginRunnable) {
	registry[canaryName] = canaryRunnable

}

func GetRegistry() map[string]PluginRunnable {
	return registry
}
