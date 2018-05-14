package commands

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/cli/plugin/models"
)

// CanaryPromote - struct
type CanaryPromote struct {
	args *CfZddCmd
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
	appName := s.args.OldApp
	canaryAppName := s.args.NewApp

	app, err := s.args.Conn.GetApp(appName)
	canary, err := s.args.Conn.GetApp(canaryAppName)

	err = s.UpdateRoutes(app, canary)

	// Do the scaleover
	scaleovercmd := &ScaleoverCmd{
		Args: s.args,
	}
	if err = scaleovercmd.ScaleoverCommand(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err = s.RemoveApplication(app); err != nil {
		fmt.Println(err.Error())
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

		output, err = s.args.Conn.CliCommand(cliArgs...)
		fmt.Printf("Add Routes output: %+v\n", output)
	}

	for _, route := range app2.Routes {
		fmt.Printf("Host: %s, Domain: %s\n ", route.Host, route.Domain.Name)
		cliArgs = []string{"delete-route", route.Domain.Name, "-n", route.Host, "-f"}

		output, err = s.args.Conn.CliCommand(cliArgs...)
		fmt.Printf("Delete Routes output: %+v\n", output)
	}
	return
}

//RemoveApplication - function to remove application
func (s *CanaryPromote) RemoveApplication(app plugin_models.GetAppModel) (err error) {
	// Remove the old application
	removeAppArgs := []string{"delete", app.Name, "-f"}
	_, err = s.args.Conn.CliCommand(removeAppArgs...)
	return
}
