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

package zdddeploy

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	"github.com/comcast/cf-zdd-plugin/scaleover"
)

// ZddDeploy - struct
type ZddDeploy struct {
	cliConnection plugin.CliConnection
	args          []string
}

// CmdRunner - interface type
type CmdRunner interface {
	Run() error
}

// ZddDeployCmdName - constants
const (
	ZddDeployCmdName = "deploy-zdd"
)

func init() {
	canaryrepo.Register(ZddDeployCmdName, new(ZddDeploy))
}

// Run method
func (s *ZddDeploy) Run(conn plugin.CliConnection) (err error) {
	s.cliConnection = conn
	err = s.deploy()
	return
}

// SetArgs - arg setter
func (s *ZddDeploy) SetArgs(args []string) {
	s.args = args
}

func (s *ZddDeploy) deploy() (err error) {
	var (
		venerable     string
		scaleovertime string
		apps          []string
	)

	appName := s.args[1]
	manifestPath := s.args[3]
	artifactPath := s.args[5]

	if len(s.args) == 7 {
		scaleovertime = s.args[6]
	} else {
		scaleovertime = "480s"
	}

	//Get the application list from cf
	apps = s.getDeployedApplications(appName)
	fmt.Printf("Found application versions: %v\n", apps)

	// Check if new deployment and deploy
	if apps == nil {
		fmt.Printf("Initial deployment of %s\n", appName)
		deployArgs := append([]string{"push"}, s.args[1:6]...)
		_, err = s.cliConnection.CliCommand(deployArgs...)
	} else {
		//Check if redeployment and rename old app.
		if isAppDeployed(apps, appName) == true {
			venerable = strings.Join([]string{appName, "venerable"}, "-")
			renameArgs := []string{"rename", appName, venerable}
			_, err = s.cliConnection.CliCommand(renameArgs...)
		} else {
			venerable = apps[0]
		}
		fmt.Printf("Venerable version assigned to %s\n", venerable)
		//Push new copy of the app with no-start
		deployArgs := []string{"push", appName, "-f", manifestPath, "-p", artifactPath, "-i", "1", "--no-start"}
		_, err = s.cliConnection.CliCommand(deployArgs...)

		// Do the scaleover
		scaleovercmdArgs := []string{"scaleover", venerable, appName, scaleovertime}
		fmt.Printf("Scaleover args: %v\n", scaleovercmdArgs)
		scaleovercmd := &scaleover.ScaleoverCmd{
			CliConnection: s.cliConnection,
			Args:          scaleovercmdArgs,
		}
		if err = scaleovercmd.ScaleoverCommand(); err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Removing app: %s\n", venerable)
		removeOldAppArgs := []string{"delete", venerable, "-f"}
		_, err = s.cliConnection.CliCommand(removeOldAppArgs...)
	}

	return
}

func (s *ZddDeploy) getDeployedApplications(appName string) []string {
	var applist []string
	shortName := strings.Split(appName, "#")[0]
	if output, err := s.cliConnection.GetApps(); err == nil {
		for _, entry := range output {
			if strings.HasPrefix(entry.Name, shortName) {
				applist = append(applist, entry.Name)
			}
		}
	}
	return applist
}

func isAppDeployed(apps []string, app string) bool {
	for _, s := range apps {
		if s == app {
			return true
		}
	}
	return false
}
