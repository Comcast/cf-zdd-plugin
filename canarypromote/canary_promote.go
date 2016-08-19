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

package canarypromote

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	"github.com/comcast/cf-zdd-plugin/scaleover"
)

// CanaryPromote - struct
type CanaryPromote struct {
	cliConnection plugin.CliConnection
	args          []string
}

// CmdRunner - interface type
type CmdRunner interface {
	Run() error
}

//Route -  OBject to hold the hostname and domain
type Route struct {
	host   string
	domain string
}

// Constants - constants for objects
const (
	CanaryPromoteCmdName = "promote-canary"
	CanaryRouteSuffix    = "canary"
	CanaryRouteSeparator = "-"
)

func init() {
	canaryrepo.Register(CanaryPromoteCmdName, new(CanaryPromote))
}

// Run method
func (s *CanaryPromote) Run(conn plugin.CliConnection) (err error) {
	s.cliConnection = conn
	err = s.promote()
	return
}

// SetArgs - Setter for the args
func (s *CanaryPromote) SetArgs(args []string) {
	s.args = args
}

func (s *CanaryPromote) promote() (err error) {
	appName := s.args[1]
	canaryAppName := s.args[2]
	scaleovertime := s.args[3]

	appURLS := strings.Split(s.getUrls(appName), ",")
	canaryURLS := strings.Split(s.getUrls(canaryAppName), ",")

	for _, url := range appURLS {
		host, domain := getHostAndDomainFromURL(strings.TrimSpace(url))
		fmt.Printf("Existing Host: #%s#, Domains: #%s#\n ", host, domain)
		mapRouteArgs := []string{"map-route", canaryAppName, domain, "-n", host}
		_, err = s.cliConnection.CliCommand(mapRouteArgs...)
	}

	// Remove the canary routes
	for _, url := range canaryURLS {
		host, domain := getHostAndDomainFromURL(strings.TrimSpace(url))
		fmt.Printf("Canary Host: #%s#, Domains: #%s#\n ", host, domain)
		removeRouteArgs := []string{"delete-route", domain, "-n", host, "-f"}
		_, err = s.cliConnection.CliCommand(removeRouteArgs...)
	}

	// Do the scaleover
	scaleovercmdArgs := []string{"scaleover", appName, canaryAppName, scaleovertime}
	fmt.Printf("Scaleover args: %v\n", scaleovercmdArgs)
	scaleovercmd := &scaleover.ScaleoverCmd{
		CliConnection: s.cliConnection,
		Args:          scaleovercmdArgs,
	}
	if err = scaleovercmd.ScaleoverCommand(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}

	// Remove the old application
	removeOldAppArgs := []string{"delete", appName, "-f"}
	_, err = s.cliConnection.CliCommand(removeOldAppArgs...)

	return
}

// CreateCanaryRouteName - function to create a properly formatted routename from an appname.
func CreateCanaryRouteName(appname string) (routename string) {
	appname = strings.Replace(appname, ".", CanaryRouteSeparator, -1)
	appname = strings.Replace(appname, "#", CanaryRouteSeparator, -1)
	routename = fmt.Sprintf("%s-%s", appname, CanaryRouteSuffix)
	return
}

func (s *CanaryPromote) getUrls(appName string) (urls string) {
	getAppArgs := []string{"app", appName}
	cfAppOutput, err := s.cliConnection.CliCommand(getAppArgs...)
	if err != nil {
		panic(err)
	}
	for _, line := range strings.Split(cfAppOutput[0], "\n") {
		if strings.Contains(line, "urls:") {
			urls = strings.Trim(line, "urls:")
			//urls: h1.d1, h2.d2, h3.d3
			fmt.Printf("URLs: %s\n ", urls)
			break
		}
	}
	return
}

func getHostAndDomainFromURL(url string) (host string, domain string) {
	arlStrArray := strings.Split(url, ".")
	host = arlStrArray[0]
	domain = strings.Join(arlStrArray[1:], ".")
	fmt.Printf("Host: %s, Domains: %s\n ", host, domain)
	return
}

func (s *CanaryPromote) getHostAndDomain(appName string) (host string, domain string) {
	getAppArgs := []string{"app", appName}
	cfAppOutput, err := s.cliConnection.CliCommand(getAppArgs...)
	if err != nil {
		panic(err)
	}
	for _, line := range cfAppOutput {
		if strings.Contains(line, "urls:") {
			line = strings.TrimSpace(line)
			arlStrArray := strings.Split(strings.Fields(line)[1], ".")
			host = arlStrArray[0]
			domain = strings.Join(arlStrArray[1:], ".")
			fmt.Printf("Host: %s, Domains: %s\n ", host, domain)
			break
		}
	}
	return
}
