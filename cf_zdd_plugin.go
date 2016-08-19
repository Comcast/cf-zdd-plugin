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

	"github.com/cloudfoundry/cli/plugin"
	"github.com/comcast/cf-zdd-plugin/canarydeploy"
	"github.com/comcast/cf-zdd-plugin/canarypromote"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	"github.com/comcast/cf-zdd-plugin/scaleover"
	"github.com/comcast/cf-zdd-plugin/zdddeploy"
)

// constants
const (
	CanaryDeployHelpText  = "Deploys an application with a canary route"
	CanaryPromoteHelpText = "Performs a promotion on the canary"
	ZddDeployHelpText     = "ZDD deployment using scale-over plugin"
	ScaleoverHelpText     = "Scalesover one application version to another"
	PluginName            = "cf-zero-downtime-deployment"
)

// var - exported vars
var (
	CanaryDeployCmdName  = canarydeploy.CanaryDeployCmdName
	CanaryPromoteCmdName = canarypromote.CanaryPromoteCmdName
	ZddDeployCmdName     = zdddeploy.ZddDeployCmdName
	ScaleoverCmdName     = scaleover.ScaleoverCmdName
)

// CfZddCmd - struct to initialize.
type CfZddCmd struct{}

//GetMetadata - required method to implement plugin
func (CfZddCmd) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: PluginName,
		Version: plugin.VersionType{
			Major: 1,
			Minor: 2,
			Build: 1,
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
		},
	}
}

//GetPluginRunnable - function to return runnable.
func GetPluginRunnable(args []string) (pluginRunnable canaryrepo.PluginRunnable) {
	pluginRunnable = canaryrepo.GetRegistry()[args[0]]
	pluginRunnable.SetArgs(args)
	return
}

// main - entry point to the plugin
func main() {
	plugin.Start(&CfZddCmd{})
}

// Run - required method to implement plugin.
func (cmd CfZddCmd) Run(cliConnection plugin.CliConnection, args []string) {

	if err := GetPluginRunnable(args).Run(cliConnection); err != nil {
		fmt.Printf("Caught panic: %s", err.Error())
		panic(err)
	}
}

// ApplicationRepo - wrapper struct
type ApplicationRepo struct {
	conn plugin.CliConnection
}

// NewApplicationRepo - wrapper for cf cliConnection
func NewApplicationRepo(conn plugin.CliConnection) *ApplicationRepo {
	return &ApplicationRepo{
		conn: conn,
	}
}
