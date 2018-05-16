package commands

import (
	"fmt"
	"os"
	"strings"
)

// ZddDeploy - struct
type ZddDeploy struct {
	args *CfZddCmd
}

// ZddDeployCmdName - constants
const (
	ZddDeployCmdName = "deploy-zdd"
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
		venerable string
		apps      []string
	)

	appName := s.args.NewApp
	manifestPath := s.args.ManifestPath
	artifactPath := s.args.ApplicationPath

	fmt.Printf("Calling zdd-deploy with args App2=%s, manifestPath=%s, artifactPath=%s\n", appName, manifestPath, artifactPath)
	if s.args.Duration == "" {
		s.args.Duration = "480s"
	}

	//Get the application list from cf
	apps = s.getDeployedApplications(appName)
	fmt.Printf("Found application versions: %v\n", apps)

	// Check if new deployment and deploy
	if apps == nil {
		fmt.Printf("Initial deployment of %s\n", appName)
		deployArgs := []string{"push", s.args.NewApp, "-f", s.args.ManifestPath, "-p", s.args.ApplicationPath}
		_, err = s.args.Conn.CliCommand(deployArgs...)
	} else {
		//Check if redeployment and rename old app.
		if isAppDeployed(apps, appName) {
			venerable = strings.Join([]string{appName, "venerable"}, "-")
			renameArgs := []string{"rename", appName, venerable}
			_, err = s.args.Conn.CliCommand(renameArgs...)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else {
			venerable = apps[0]
		}
		fmt.Printf("Venerable version assigned to %s\n", venerable)
		//Push new copy of the app with no-start
		deployArgs := []string{"push", appName, "-f", manifestPath, "-p", artifactPath, "-i", "1", "--no-start"}
		_, err = s.args.Conn.CliCommand(deployArgs...)
		if err != nil {
			fmt.Println(err.Error())
		}
		// Do the scaleover
		s.args.OldApp = venerable
		scaleovercmd := &ScaleoverCmd{
			Args: s.args,
		}
		if err = scaleovercmd.ScaleoverCommand(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Printf("Removing app: %s\n", venerable)
		removeOldAppArgs := []string{"delete", venerable, "-f"}
		_, err = s.args.Conn.CliCommand(removeOldAppArgs...)
	}

	return
}

func (s *ZddDeploy) getDeployedApplications(appName string) []string {
	var applist []string
	shortName := strings.Split(appName, "#")[0]
	if output, err := s.args.Conn.GetApps(); err == nil {
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
