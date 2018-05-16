package commands

var registry = make(map[string]CommandRunnable)

// CommandRunnable - interface type for other commands
type CommandRunnable interface {
	Run() error
	SetArgs(cmd *CfZddCmd)
}

// Register - function to add CommandRunnable to the registry map
func Register(cmdName string, canaryRunnable CommandRunnable) {
	registry[cmdName] = canaryRunnable

}

// GetRegistry - function to return the current registry
func GetRegistry() map[string]CommandRunnable {
	return registry
}
