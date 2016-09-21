package canarypromote

import (
	"fmt"
	"os"

	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	"github.com/comcast/cf-zdd-plugin/scaleover"
)

// CanaryPromote - struct
type CanaryPromote struct {
	CliConnection plugin.CliConnection
	args          []string
}

// CmdRunner - interface type
type CmdRunner interface {
	Run() error
}

// Route - to hold the hostname and domain
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
	s.CliConnection = conn
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

	app, err := s.CliConnection.GetApp(appName)
	canary, err := s.CliConnection.GetApp(canaryAppName)

	err = s.UpdateRoutes(app, canary)

	// Do the scaleover
	scaleovercmdArgs := []string{"scaleover", app.Name, canary.Name, scaleovertime}
	fmt.Printf("Scaleover args: %v\n", scaleovercmdArgs)
	scaleovercmd := &scaleover.ScaleoverCmd{
		CliConnection: s.CliConnection,
		Args:          scaleovercmdArgs,
	}
	if err = scaleovercmd.ScaleoverCommand(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	if err = s.RemoveApplication(app); err != nil {
		fmt.Printf(err.Error())
		os.Exit(1)
	}
	return
}

// UpdateRoutes - function to add or remove routes from the application
func (s *CanaryPromote) UpdateRoutes(app1 plugin_models.GetAppModel, app2 plugin_models.GetAppModel) (err error) {
	var (
		output  []string
		cliArgs []string
	)
	for _, route := range app1.Routes {

		fmt.Printf("Host: %s, Domain: %s\n ", route.Host, route.Domain.Name)

		cliArgs = []string{"map-route", app2.Name, route.Domain.Name, "-n", route.Host}

		output, err = s.CliConnection.CliCommand(cliArgs...)
		fmt.Printf("Add Routes output: %+v\n", output)
	}

	for _, route := range app2.Routes {
		fmt.Printf("Host: %s, Domain: %s\n ", route.Host, route.Domain.Name)
		cliArgs = []string{"delete-route", route.Domain.Name, "-n", route.Host, "-f"}

		output, err = s.CliConnection.CliCommand(cliArgs...)
		fmt.Printf("Delete Routes output: %+v\n", output)
	}
	return
}

//RemoveApplication - function to remove application
func (s *CanaryPromote) RemoveApplication(app plugin_models.GetAppModel) (err error) {
	// Remove the old application
	removeAppArgs := []string{"delete", app.Name, "-f"}
	_, err = s.CliConnection.CliCommand(removeAppArgs...)
	return
}
