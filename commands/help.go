package commands

import "fmt"

// HelpCmd - struct for HelpCmd
type HelpCmd struct {
	args *CfZddCmd
}

// HelpCommandName - constant to set command name text
const HelpCommandName = "zdd-help"

func init() {
	Register(HelpCommandName, new(HelpCmd))
}

// Run - run method for help command
func (h *HelpCmd) Run() (err error) {
	err = h.help()
	return
}

// SetArgs - set args command for help command
func (h *HelpCmd) SetArgs(args *CfZddCmd) {
	h.args = args
}

func (h *HelpCmd) help() (err error) {

	var helpString string
	switch h.args.HelpTopic {
	case ZddDeployCmdName:
		helpString = "zdd-deploy help" +
			"\n\t--oldapp = The name of the existing application" +
			"\n\t--newapp = The name of the new application" +
			"\n\t--duration = The time for scaling over the application, default is 480s" +
			"\n\t--p = The path to the application file" +
			"\n\t--f = The path to the application manifest"
	case CanaryDeployCmdName:
		helpString = "deploy-canary help" +
			"\n\t--newapp = The name of the new application" +
			"\n\t--p = The path to the application file" +
			"\n\t--f = The path to the application manifest"
	case CanaryPromoteCmdName:
		helpString = "promote-canary help" +
			"\n\t--oldapp = The name of the existing application" +
			"\n\t--newapp = The name of the new application" +
			"\n\t--duration = The time for scaling over the application, default is 480s" +
			"\n\t--p = The path to the application file" +
			"\n\t--f = The path to the application manifest"
	case ScaleoverCmdName:
		helpString = "scaleover help" +
			"\n\t--oldapp = The name of the existing application" +
			"\n\t--newapp = The name of the new application" +
			"\n\t--duration = The time for scaling over the application, default is 480s"
	case blueGreenCmdName:
		helpString = "blue-green help" +
			"\n\t--newapp = The name of the new application" +
			"\n\t--p = The path to the application file" +
			"\n\t--f = The path to the application manifest"
	default:
		helpString = "Help is available for the deployment types: \n\t - deploy-canary \n\t - promote-canary \n\t - blue-green \n\t - deploy-zdd \nUse the command help <deploy command> for command specific help"
	}

	fmt.Println(helpString)
	return
}
