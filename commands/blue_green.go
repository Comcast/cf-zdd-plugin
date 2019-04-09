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

package commands

import (
	"fmt"
	"strings"
	"time"
)

// BlueGreenDeploy - struct for deployment
type BlueGreenDeploy struct {
	args *CfZddCmd
}

const BlueGreenCmdName = "blue-green"

func init() {
	Register(BlueGreenCmdName, new(BlueGreenDeploy))
}

// Run - run method as required by the interface
func (bg *BlueGreenDeploy) Run() (err error) {
	err = bg.deploy()
	return
}

// SetArgs - function to set the arguments for the deployment
func (bg *BlueGreenDeploy) SetArgs(args *CfZddCmd) {
	bg.args = args
}

func (bg *BlueGreenDeploy) deploy() (err error) {

	applicationToDeploy := bg.args.NewApp
	manifestPath := bg.args.ManifestPath
	artifactPath := bg.args.ApplicationPath

	fmt.Printf("Calling blue green deploy with args: Application=%s, manifestPath=%s, artifactPath=%s\n", applicationToDeploy, manifestPath, artifactPath)

	var (
		isAppDeployed bool
		searchAppName string
		oldAppName    string
	)

	// Verify if application is currently deployed and current name if app name is versioned
	if bg.args.BaseAppName != "" {
		searchAppName = bg.args.BaseAppName
	} else {
		searchAppName = applicationToDeploy
	}

	oldAppName, isAppDeployed = bg.args.Commands.IsApplicationDeployed(searchAppName)

	if !isAppDeployed {
		fmt.Println("Application is not deployed.... pushing.")
		err = bg.args.Commands.PushApplication(applicationToDeploy, artifactPath, manifestPath, "--no-route")
	} else {
		fmt.Println("Application is deployed, renaming existing version")
		venerable := strings.Join([]string{oldAppName, "venerable"}, "-")
		if err = bg.args.Commands.RenameApplication(oldAppName, venerable); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Pushing new version with name: %s\n", applicationToDeploy)
		if err = bg.args.Commands.PushApplication(applicationToDeploy, artifactPath, manifestPath, "--no-route"); err != nil {
			fmt.Println(err.Error())
			return
		}

		for !bg.areAllInstancesStarted(applicationToDeploy) {
			time.Sleep(20 * time.Second)
		}
		fmt.Println("All instances started, remapping route.")
		if err = bg.args.Commands.RemapRoutes(applicationToDeploy, venerable); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Removing old version")
		if err = bg.args.Commands.RemoveApplication(venerable); err != nil {
			fmt.Println(err.Error())
		}
	}

	return
}

func (bg *BlueGreenDeploy) areAllInstancesStarted(appName string) bool {
	if output, err := bg.args.Conn.GetApp(appName); err == nil {
		if output.InstanceCount == output.RunningInstances {
			return true
		}
	}
	return false
}
