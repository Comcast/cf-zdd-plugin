package commands

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin/models"
)

// CanaryPromote - struct
type CanaryPromote struct {
	args         *CfZddCmd
	ScaleoverCmd ScaleoverCommand
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
)

func init() {
	Register(CanaryPromoteCmdName, new(CanaryPromote))
}

// Run method
func (s *CanaryPromote) Run() (err error) {
	err = s.promote()
	return
}

// SetArgs - Setter for the args
func (s *CanaryPromote) SetArgs(args *CfZddCmd) {
	s.args = args
}

func (s *CanaryPromote) promote() (err error) {

	if s.ScaleoverCmd == nil {
		s.ScaleoverCmd = NewScaleoverCmd(s.args)
	}

	appName := s.args.OldApp
	canaryAppName := s.args.NewApp

	app, err := s.args.Conn.GetApp(appName)
	canary, err := s.args.Conn.GetApp(canaryAppName)

	err = s.UpdateRoutes(app, canary)

	if err = s.ScaleoverCmd.DoScaleover(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err = s.args.Commands.RemoveApplication(appName); err != nil {
		fmt.Printf("Error removing application: %s\n", err.Error())
		return
	}
	return
}

// UpdateRoutes - function to add or remove routes from the application. Apply the existing application routes to the
// canary version of the application.
func (s *CanaryPromote) UpdateRoutes(oldApp plugin_models.GetAppModel, canary plugin_models.GetAppModel) (err error) {
	var (
		output  []string
		cliArgs []string
	)

	for _, route := range oldApp.Routes {
		fmt.Printf("Host: %s, Domain: %s\n ", route.Host, route.Domain.Name)

		cliArgs = []string{"map-route", canary.Name, route.Domain.Name, "-n", route.Host}

		output, err = s.args.Conn.CliCommand(cliArgs...)
		fmt.Printf("Add Routes output: %+v\n", output)
	}

	for _, route := range canary.Routes {
		fmt.Printf("Host: %s, Domain: %s\n ", route.Host, route.Domain.Name)
		cliArgs = []string{"delete-route", route.Domain.Name, "-n", route.Host, "-f"}

		output, err = s.args.Conn.CliCommand(cliArgs...)
		fmt.Printf("Delete Routes output: %+v\n", output)
	}
	return
}
