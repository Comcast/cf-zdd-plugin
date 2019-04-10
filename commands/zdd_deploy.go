package commands

import (
	"fmt"
	"os"
)

// ZddDeploy - struct
type ZddDeploy struct {
	args          *CfZddCmd
	ScalerOverCmd ScaleoverCommand
}

// ZddDeployCmdName - constants
const (
	ZddDeployCmdName = "deploy-zdd"
	DefaultDuration  = "480s"
)

func init() {
	Register(ZddDeployCmdName, new(ZddDeploy))
}

// Run method
func (s *ZddDeploy) Run() (err error) {
	err = s.deploy()
	return
}

// SetArgs - arg setter
func (s *ZddDeploy) SetArgs(args *CfZddCmd) {
	s.args = args
}

func (s *ZddDeploy) deploy() (err error) {
	var (
		oldApplication string
		venerable      string
		isAppDeployed  bool
		searchAppName  string
	)

	if s.ScalerOverCmd == nil {
		s.ScalerOverCmd = NewScaleoverCmd(s.args)
	}

	applicationToDeploy := s.args.NewApp
	manifestPath := s.args.ManifestPath
	artifactPath := s.args.ApplicationPath

	fmt.Printf("Calling zdd-deploy with args App2=%s, manifestPath=%s, artifactPath=%s\n", applicationToDeploy, manifestPath, artifactPath)
	if s.args.Duration == "" {
		s.args.Duration = DefaultDuration
	}

	// Verify if application is currently deployed and current name if app name is versioned
	if s.args.BaseAppName != "" {
		searchAppName = s.args.BaseAppName
	} else {
		searchAppName = applicationToDeploy
	}

	//Get the application list from cf
	oldApplication, isAppDeployed = s.args.Commands.IsApplicationDeployed(searchAppName)

	if !isAppDeployed {
		fmt.Printf("Initial deployment of %s\n", applicationToDeploy)
		if err = s.args.Commands.PushApplication(applicationToDeploy, artifactPath, manifestPath); err != nil {
			fmt.Printf("Error occurred pushing application: %s\n", err.Error())
		}
	} else {
		//Check if redeployment and rename old app.
		if oldApplication == applicationToDeploy {
			venerable = oldApplication + "-venerable"
			err = s.args.Commands.RenameApplication(oldApplication, venerable)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			venerable = oldApplication
		}
		fmt.Printf("Venerable version assigned to %s\n", venerable)

		if err = s.args.Commands.PushApplication(applicationToDeploy, artifactPath, manifestPath, "-i", "1", "--no-start"); err != nil {
			fmt.Println(err.Error())
		}

		// Do the scaleover
		s.args.OldApp = venerable

		if err = s.ScalerOverCmd.DoScaleover(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Removing app: %s\n", venerable)
		if err = s.args.Commands.RemoveApplication(venerable); err != nil {
			fmt.Printf("Unable to remove old application: %s, error: %s\n", venerable, err.Error())
		}
	}

	return
}
