// Copyright 2016 Joshua Kruck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/andrew-d/go-termutil"
	"net/http"
)

type clientDoer interface {
	Do(*http.Request) (*http.Response, error)
}

//AppStatus represents the sattus of a app in CF
type AppStatus struct {
	Name           string
	GUID           string
	CountRunning   int
	CountRequested int
	State          string
	Routes         []string
}

//ScaleoverCmd is this plugin
type ScaleoverCmd struct {
	App1   *AppStatus
	App2   *AppStatus
	Args   *CfZddCmd
	Client clientDoer
}

// ScaleoverCmdName - constants
const (
	ScaleoverCmdName = "scaleover"
)

func init() {
	Register(ScaleoverCmdName, new(ScaleoverCmd))
}

// SetArgs - arg setter
func (cmd *ScaleoverCmd) SetArgs(args *CfZddCmd) {
	cmd.Args = args
}

func (cmd *ScaleoverCmd) ShouldEnforceRoutes() bool {
	return cmd.Args.RouteCheck
}

func (cmd *ScaleoverCmd) ParseTime(duration string) (time.Duration, error) {
	var (
		rolloverTime time.Duration
		err          error
	)
	if rolloverTime, err = time.ParseDuration(duration); err != nil {
		return rolloverTime, err
	}

	if 0 > rolloverTime {
		return rolloverTime, errors.New("Duration must be a positive number in the format of 1m")
	}
	return rolloverTime, nil
}

//Run runs the plugin
func (cmd *ScaleoverCmd) Run() (err error) {
	cmd.ScaleoverCommand()
	return
}

func (cmd *ScaleoverCmd) Usage(args *CfZddCmd) error {

	if args.OldApp == "" || args.NewApp == "" {
		return errors.New("App 1 and App2 are required")
	}

	if args.CustomURL == "" && args.Duration == "" {
		return errors.New("Custom URL or Duration is required")
	}

	return nil
}

//ScaleoverCommand creates a new instance of this plugin
func (cmd *ScaleoverCmd) ScaleoverCommand() (err error) {
	enforceRoutes := cmd.ShouldEnforceRoutes()

	if err = cmd.Usage(cmd.Args); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	rolloverTime, err := cmd.ParseTime(cmd.Args.Duration)
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	// The getAppStatus calls will exit with an error if the named apps don't exist
	if cmd.App1, err = cmd.GetAppStatus(cmd.Args.OldApp); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	if cmd.App2, err = cmd.GetAppStatus(cmd.Args.NewApp); nil != err {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("App1: %+v\nApp2: %+v\n", cmd.App1, cmd.App2)
	if enforceRoutes {
		if err = cmd.ErrorIfNoSharedRoute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd.showStatus()

	count := cmd.App1.CountRequested
	if count == 0 {
		fmt.Println("There are no instances of the source app to scale over")
		os.Exit(0)
	}

	//Standard scaleover using duration
	if cmd.Args.Duration != "" && cmd.Args.CustomURL == "" {
		sleepInterval := time.Duration(rolloverTime.Nanoseconds() / int64(count))

		for count > 0 {
			count--
			cmd.App2.ScaleUp(cmd.Args.Conn)
			cmd.App1.ScaleDown(cmd.Args.Conn)
			cmd.showStatus()
			if count > 0 {
				time.Sleep(sleepInterval)
			}
		}
	}

	return
}

func (cmd *ScaleoverCmd) GetAppStatus(name string) (*AppStatus, error) {
	app, err := cmd.Args.Conn.GetApp(name)

	if nil != err {
		return nil, err
	}
	status := &AppStatus{
		Name:           name,
		GUID:           app.Guid,
		CountRunning:   0,
		CountRequested: 0,
		State:          "unknown",
		Routes:         make([]string, len(app.Routes)),
	}

	status.State = app.State
	if app.State != "stopped" {
		status.CountRequested = app.InstanceCount
	}
	status.CountRunning = app.RunningInstances
	for idx, route := range app.Routes {
		status.Routes[idx] = route.Host + "." + route.Domain.Name
	}
	return status, nil
}

func (app *AppStatus) ScaleUp(cliConnection plugin.CliConnection) {
	// If not already started, start it
	if app.State != "started" {
		cliConnection.CliCommandWithoutTerminalOutput("start", app.Name)
		app.State = "started"
	}
	app.CountRequested++
	cliConnection.CliCommandWithoutTerminalOutput("scale", "-i", strconv.Itoa(app.CountRequested), app.Name)
}

func (app *AppStatus) ScaleDown(cliConnection plugin.CliConnection) {
	app.CountRequested--
	// If going to zero, stop the app
	if app.CountRequested == 0 {
		cliConnection.CliCommandWithoutTerminalOutput("stop", app.Name)
		app.State = "stopped"
	} else {
		cliConnection.CliCommandWithoutTerminalOutput("scale", "-i", strconv.Itoa(app.CountRequested), app.Name)
	}
}

func (cmd *ScaleoverCmd) showStatus() {
	if termutil.Isatty(os.Stdout.Fd()) {
		fmt.Printf("%s (%s) %s %s %s (%s) \r",
			cmd.App1.Name,
			cmd.App1.State,
			strings.Repeat("<", cmd.App1.CountRequested),
			strings.Repeat(">", cmd.App2.CountRequested),
			cmd.App2.Name,
			cmd.App2.State,
		)
	} else {
		fmt.Printf("%s (%s) %d instances, %s (%s) %d instances\n",
			cmd.App1.Name,
			cmd.App1.State,
			cmd.App1.CountRequested,
			cmd.App2.Name,
			cmd.App2.State,
			cmd.App2.CountRequested,
		)
	}
}

func (cmd *ScaleoverCmd) AppsShareARoute() bool {
	for _, r1 := range cmd.App1.Routes {
		for _, r2 := range cmd.App2.Routes {
			if r1 == r2 {
				return true
			}
		}
	}
	return false
}

func (cmd *ScaleoverCmd) ErrorIfNoSharedRoute() error {
	if cmd.AppsShareARoute() {
		return nil
	}
	return errors.New("Apps do not share a route!")
}
