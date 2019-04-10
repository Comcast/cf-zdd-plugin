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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry/cli/plugin"
	"strings"
)

type CommonCmd interface {
	IsApplicationDeployed(string) (string, bool)
	PushApplication(string, string, string, ...string) error
	RenameApplication(string, string) error
	RemapRoutes(string, string) error
	RemoveApplication(string) error
	GetDefaultDomain() string
}

type commonCmd struct {
	cli plugin.CliConnection
}

func NewCommonCmd(conn plugin.CliConnection) CommonCmd {
	return &commonCmd{
		cli: conn,
	}
}

//GetDefaultDomain - function to get the default domain from cloud foundry
func (c *commonCmd) GetDefaultDomain() (domain string) {
	infoArgs := []string{"curl", "/v2/info"}
	infoOutput, err := c.cli.CliCommandWithoutTerminalOutput(infoArgs...)

	fmt.Println(infoOutput)

	if err != nil {
		domain = "unknown"
		return
	}
	infoBytes := []byte(strings.Join(infoOutput, ""))

	var info map[string]interface{}
	err = json.Unmarshal(infoBytes, &info)

	if err == nil {
		apiurl := info["authorization_endpoint"].(string)
		domain = apiurl[strings.IndexRune(apiurl, '.')+1:]
	}
	return
}

func (c *commonCmd) IsApplicationDeployed(appName string) (string, bool) {
	if output, err := c.cli.GetApps(); err == nil {
		for _, app := range output {
			if strings.HasPrefix(app.Name, appName) {
				fmt.Println("Application is deployed")
				return app.Name, true
			}
		}
	}
	return "", false
}

func (c *commonCmd) PushApplication(appName string, artifactPath string, manifestPath string, extraArgs ...string) error {

	if appName == "" {
		return errors.New("appname must be specified")
	}
	pushArgs := []string{"push", appName, "-f", manifestPath}

	if artifactPath != "" {
		pushArgs = append(pushArgs, "-p", artifactPath)
	}

	if manifestPath != "" {
		pushArgs = append(pushArgs, "-f", manifestPath)
	}

	pushArgs = append(pushArgs, extraArgs...)

	if _, err := c.cli.CliCommand(pushArgs...); err != nil {
		fmt.Printf("error pushing app - %s", err.Error())
		return err
	}
	return nil
}

func (c *commonCmd) RemoveApplication(appName string) error {
	if appName == "" {
		return errors.New("appname must be specified")
	}
	removeArgs := []string{"delete", appName, "-f"}

	if _, err := c.cli.CliCommand(removeArgs...); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (c *commonCmd) RemapRoutes(from string, to string) error {
	// Get the app model for the old version of the application
	fromModel, err := c.cli.GetApp(from)

	if err != nil {
		return err
	}
	// Map the routes to the new application version
	for _, r := range fromModel.Routes {
		mapArgs := []string{"map-route", to, r.Domain.Name, "-n", r.Host}
		if _, err = c.cli.CliCommand(mapArgs...); err != nil {
			fmt.Println(err.Error())
		}
	}
	// Remove the route from the old app version
	for _, r := range fromModel.Routes {
		unmapArgs := []string{"unmap-route", from, r.Domain.Name, "-n", r.Host}
		if _, err = c.cli.CliCommand(unmapArgs...); err != nil {
			fmt.Println(err.Error())
		}
	}

	return err
}

func (c *commonCmd) RenameApplication(from string, to string) (err error) {
	if from == "" || to == "" {
		return errors.New("appname and new appname must be specified")
	}
	renameArgs := []string{"rename", from, to}
	if _, err = c.cli.CliCommand(renameArgs...); err != nil {
		fmt.Println(err.Error())
	}
	return err
}
