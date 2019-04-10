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
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// CanaryDeploy - struct
type CanaryDeploy struct {
	args *CfZddCmd
}

// DomainList - struct
type DomainList struct {
	Routes []string `yaml:"routes,omitempty"`
}

// CanaryDeployCmdName - constants
const (
	CanaryDeployCmdName = "deploy-canary"
)

func init() {
	Register(CanaryDeployCmdName, new(CanaryDeploy))
}

// Run - Run method
func (s *CanaryDeploy) Run() (err error) {
	err = s.deploy()
	return
}

// SetArgs - set command args
func (s *CanaryDeploy) SetArgs(args *CfZddCmd) {
	s.args = args
}

// DeployCanary - function to create and push a canary deployment
func (s *CanaryDeploy) deploy() (err error) {
	appName := s.args.NewApp

	//Deploy an initial canary version
	deployArgs := []string{"-i", "1", "--no-route", "--no-start"}

	fmt.Printf("Calling with deploy args: %v\n", deployArgs)
	err = s.args.Commands.PushApplication(appName, s.args.ApplicationPath, s.args.ManifestPath, deployArgs...)

	deployArgsMapRoute := []string{"map-route", appName, s.getDomain(), "-n", CreateCanaryRouteName(appName)}
	fmt.Printf("Calling with deploy args: %v\n", deployArgsMapRoute)
	_, err = s.args.Conn.CliCommand(deployArgsMapRoute...)

	if err != nil {
		fmt.Println(err.Error())
	}

	startArgs := []string{"start", appName}
	_, err = s.args.Conn.CliCommand(startArgs...)

	return
}

// CreateCanaryRouteName - function to create a properly formatted routename from an appname.
func CreateCanaryRouteName(appname string) (routename string) {
	appname = strings.Replace(appname, ".", CanaryRouteSeparator, -1)
	appname = strings.Replace(appname, "#", CanaryRouteSeparator, -1)
	routename = fmt.Sprintf("%s-%s", appname, CanaryRouteSuffix)
	return
}

func (s *CanaryDeploy) getDomain() (domain string) {

	var yamlFile []byte
	var err error

	// Check manifest file for a route
	if s.args.ManifestPath != "" {
		fmt.Printf("Reading manifest file at %s\n", s.args.ManifestPath)
		yamlFile, err = ioutil.ReadFile(s.args.ManifestPath)
		if err != nil {
			fmt.Printf("###ERROR reading file: %s\n", err.Error())
			domain = s.args.Commands.GetDefaultDomain()
			return
		}
	} else if _, err = os.Stat("manifest.yml"); err == nil {
		fmt.Println("Reading default manifest file")
		yamlFile, err = ioutil.ReadFile("manifest.yml")

		if err != nil {
			fmt.Printf("###ERROR reading file: %s\n", err.Error())
			domain = s.args.Commands.GetDefaultDomain()
			return
		}
	}

	var domainList DomainList
	if len(yamlFile) > 0 {
		err = yaml.Unmarshal(yamlFile, &domainList)
		if err != nil {
			fmt.Printf("YAML UNMARSHAL ERROR: %s\n", err.Error())
			domain = s.args.Commands.GetDefaultDomain()
			return
		}
		if len(domainList.Routes) > 0 {
			domain = strings.Join(strings.Split(domainList.Routes[0], ".")[1:], ".")
			return
		} else {
			domain = s.args.Commands.GetDefaultDomain()
			return
		}
	}
	return
}
