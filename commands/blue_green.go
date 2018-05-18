package commands

import (
	"code.cloudfoundry.org/cli/plugin/models"
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

	if !bg.isApplicationDeployed(applicationToDeploy) {
		fmt.Println("Application is not deployed.... pushing.")
		err = bg.pushApplication(applicationToDeploy, artifactPath, manifestPath, false)
	} else {
		fmt.Println("Application is deployed, renaming existing version")
		venerable := strings.Join([]string{applicationToDeploy, "venerable"}, "-")
		renameArgs := []string{"rename", applicationToDeploy, venerable}
		if _, err = bg.args.Conn.CliCommand(renameArgs...); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Pushing new version with name: %s\n", applicationToDeploy)
		if err = bg.pushApplication(applicationToDeploy, artifactPath, manifestPath, true); err != nil {
			fmt.Println(err.Error())
			return
		}

		for !bg.areAllInstancesStarted(applicationToDeploy) {
			time.Sleep(20 * time.Second)
		}
		fmt.Println("All instances started, remapping route.")
		if err = bg.remapRoute(applicationToDeploy, venerable); err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Removing old version")
		if err = bg.removeApplication(venerable); err != nil {
			fmt.Println(err.Error())
		}
	}

	return
}

func (bg *BlueGreenDeploy) isApplicationDeployed(appName string) bool {
	if output, err := bg.args.Conn.GetApps(); err == nil {
		for _, app := range output {
			if appName == app.Name {
				fmt.Println("Application is deployed")
				return true
			}
		}
	}
	return false
}

func (bg *BlueGreenDeploy) pushApplication(appName string, artifactPath string, manifestPath string, isDeployed bool) error {
	deployArgs := []string{"push", appName, "-f", manifestPath, "-p", artifactPath}

	if isDeployed {
		deployArgs = append(deployArgs, "--no-route")
	}

	if _, err := bg.args.Conn.CliCommand(deployArgs...); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (bg *BlueGreenDeploy) areAllInstancesStarted(appName string) bool {
	if output, err := bg.args.Conn.GetApp(appName); err == nil {
		if output.InstanceCount == output.RunningInstances {
			return true
		}
	}
	return false
}

func (bg *BlueGreenDeploy) removeApplication(appName string) error {
	removeArgs := []string{"delete", appName, "-f"}

	if _, err := bg.args.Conn.CliCommand(removeArgs...); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (bg *BlueGreenDeploy) remapRoute(appName string, venerable string) error {
	//Get the routes associated with the old version
	var (
		err            error
		venerableModel plugin_models.GetAppModel
	)

	venerableModel, err = bg.args.Conn.GetApp(venerable)
	if err != nil {
		return err
	}

	// Map the routes to the new application version
	for _, route := range venerableModel.Routes {
		mrArgs := []string{"map-route", appName, route.Domain.Name, "-n", route.Host}
		_, err = bg.args.Conn.CliCommand(mrArgs...)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// Remove the route from the old app version
	for _, route := range venerableModel.Routes {
		drArgs := []string{"unmap-route", venerable, route.Domain.Name, "-n", route.Host}
		_, err = bg.args.Conn.CliCommand(drArgs...)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return err
}
