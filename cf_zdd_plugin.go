/*
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	"flag"
	"github.com/comcast/cf-zdd-plugin/commands"
	"strconv"
)

// constants
const (
	CanaryDeployHelpText  = "Deploys an application with a canary route"
	CanaryPromoteHelpText = "Performs a promotion on the canary"
	ZddDeployHelpText     = "ZDD deployment using scale-over plugin"
	ScaleoverHelpText     = "Scalesover one application version to another"
	HelpText              = "Help is available for each of the commands in the form 'help <command name>'"
	BlueGreenHelpText     = "Deploys an application and then flips the route to the new application"
	PluginName            = "cf-zero-downtime-deployment"
)

// var - exported vars
var (
	CanaryDeployCmdName  = commands.CanaryDeployCmdName
	CanaryPromoteCmdName = commands.CanaryPromoteCmdName
	ZddDeployCmdName     = commands.ZddDeployCmdName
	ScaleoverCmdName     = commands.ScaleoverCmdName
	HelpCmdName          = commands.HelpCommandName
	BlueGreenCmdName     = commands.BlueGreenCmdName
	Major                string
	Minor                string
	Patch                string
)

// CfZddPlugin - struct to initialize.
type CfZddPlugin struct {
	cmd *commands.CfZddCmd
}

//GetMetadata - required method to implement plugin
func (c *CfZddPlugin) GetMetadata() plugin.PluginMetadata {

	major, _ := strconv.Atoi(Major)
	minor, _ := strconv.Atoi(Minor)
	patch, _ := strconv.Atoi(Patch)

	return plugin.PluginMetadata{
		Name: PluginName,
		Version: plugin.VersionType{
			Major: major,
			Minor: minor,
			Build: patch,
		},
		Commands: []plugin.Command{
			{
				Name:     CanaryDeployCmdName,
				HelpText: CanaryDeployHelpText,
			},
			{
				Name:     CanaryPromoteCmdName,
				HelpText: CanaryPromoteHelpText,
			},
			{
				Name:     ZddDeployCmdName,
				HelpText: ZddDeployHelpText,
			},
			{
				Name:     ScaleoverCmdName,
				HelpText: ScaleoverHelpText,
			},
			{
				Name:     BlueGreenCmdName,
				HelpText: BlueGreenHelpText,
			},
			{
				Name:     HelpCmdName,
				HelpText: HelpText,
			},
		},
	}
}

//GetPluginRunnable - function to return runnable.
func (c *CfZddPlugin) GetPluginRunnable() (pluginRunnable commands.CommandRunnable) {
	pluginRunnable = commands.GetRegistry()[c.cmd.CmdName]
	if pluginRunnable != nil {
		pluginRunnable.SetArgs(c.cmd)
	}
	return
}

// main - entry point to the plugin
func main() {
	plugin.Start(new(CfZddPlugin))
}

// Run - required method to implement plugin.
func (c *CfZddPlugin) Run(cliConnection plugin.CliConnection, args []string) {

	fs := flag.NewFlagSet(args[0], flag.ExitOnError)

	fmt.Printf("ARGS: %+v\n", args)

	app1Flag := fs.String("old-app", "", "current application name")
	app2Flag := fs.String("new-app", "", "new application being deployed")
	durationflag := fs.String("duration", "", "time between scalovers")
	applicationPathflag := fs.String("p", "", "path to applcation file")
	manifestPathFlag := fs.String("f", "", "path to application manifest")
	customURLFlag := fs.String("custom-health-url", "", "path to custom healthcheck page")
	batchSizeFlag := fs.Int("batch-size", 1, "number to restart/deploy at a time")
	routeCheckFlag := fs.Bool("no-route-check", false, "check to ensure a common route")

	fs.Parse(args[1:])

	c.cmd = &commands.CfZddCmd{
		OldApp:          *app1Flag,
		NewApp:          *app2Flag,
		CmdName:         args[0],
		Conn:            cliConnection,
		Duration:        *durationflag,
		ApplicationPath: *applicationPathflag,
		ManifestPath:    *manifestPathFlag,
		CustomURL:       *customURLFlag,
		BatchSize:       *batchSizeFlag,
		RouteCheck:      *routeCheckFlag,
	}

	if args[0] == HelpCmdName && len(args) > 1 {
		c.cmd.HelpTopic = args[1]
	}

	if pr := c.GetPluginRunnable(); pr != nil {
		pr.Run()
	}
}
