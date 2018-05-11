package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/comcast/cf-zdd-plugin/util"

	"gopkg.in/yaml.v2"
)

// CanaryDeploy - struct
type CanaryDeploy struct {
	Utils         util.Utilities
	cliConnection plugin.CliConnection
	args          []string
}

// DomainList - struct
type DomainList struct {
	Domain  string   `yaml:"domain,omitempty"`
	Domains []string `yaml:"domains,omitempty"`
}

// CanaryDeployCmdName - constants
const (
	CanaryDeployCmdName  = "deploy-canary"
	CanaryRouteSuffix    = "canary"
	CanaryRouteSeparator = "-"
)

func init() {
	Register(CanaryDeployCmdName, new(CanaryDeploy))
}

// Run - Run method
func (s *CanaryDeploy) Run(conn plugin.CliConnection) (err error) {
	s.cliConnection = conn
	if s.Utils == nil {
		s.Utils = new(util.Utility)
	}
	err = s.deploy()

	return
}

// SetArgs - set command args
func (s *CanaryDeploy) SetArgs(args []string) {
	s.args = args
}

// DeployCanary - function to create and push a canary deployment
func (s *CanaryDeploy) deploy() (err error) {
	appName := s.args[1]

	deployArgs := append([]string{"push"}, s.args[1:]...)
	deployArgsNoroute := append(deployArgs, "-i", "1", "--no-route", "--no-start")

	fmt.Printf("Calling with deploy args: %v\n", deployArgsNoroute)
	_, err = s.cliConnection.CliCommand(deployArgsNoroute...)

	domains := s.getDomain()
	for _, val := range domains {
		deployArgsMapRoute := []string{"map-route", appName, val, "-n", CreateCanaryRouteName(appName)}
		fmt.Printf("Calling with deploy args: %v\n", deployArgsMapRoute)
		_, err = s.cliConnection.CliCommand(deployArgsMapRoute...)
	}
	if err != nil {

	}
	startArgs := []string{"start", appName}
	_, err = s.cliConnection.CliCommand(startArgs...)

	return
}

// CreateCanaryRouteName - function to create a properly formatted routename from an appname.
func CreateCanaryRouteName(appname string) (routename string) {
	appname = strings.Replace(appname, ".", CanaryRouteSeparator, -1)
	appname = strings.Replace(appname, "#", CanaryRouteSeparator, -1)
	routename = fmt.Sprintf("%s-%s", appname, CanaryRouteSuffix)
	return
}

func (s *CanaryDeploy) getDomain() (domains []string) {

	// check for file exists
	var manifestIndex = -1
	for idx, val := range s.args {
		if val == "-f" {
			manifestIndex = idx + 1
			fmt.Printf("Found manifest at index +%v\n", manifestIndex)
			break
		}
	}
	var yamlFile []byte
	var err error
	if manifestIndex != -1 {
		if _, err = os.Stat(s.args[manifestIndex]); err == nil {
			fmt.Printf("Reading manifest file at %s\n", s.args[manifestIndex])
			yamlFile, err = ioutil.ReadFile(s.args[manifestIndex])
			if err != nil {
				fmt.Printf("###ERROR: %s\n", err.Error())
				domains = []string{s.Utils.GetDefaultDomain(s.cliConnection)}
			}
		}
	} else {
		fmt.Println("Reading default manifest file")
		yamlFile, err = ioutil.ReadFile("manifest.yml")

		if err != nil {
			fmt.Printf("###ERROR: %s\n", err.Error())
			domains = []string{s.Utils.GetDefaultDomain(s.cliConnection)}
		}
	}

	var domainList DomainList
	if len(yamlFile) > 0 {
		err = yaml.Unmarshal(yamlFile, &domainList)
		if err != nil {
			fmt.Printf("###ERROR: %s\n", err.Error())
		}
		fmt.Printf("DomainList: %+v\n", domainList)
		if len(domainList.Domains) > 0 {
			if len(domainList.Domain) > 0 {
				domains = append(domainList.Domains, domainList.Domain)
			} else {
				domains = domainList.Domains
			}
		} else {
			fmt.Println("Domains are empty creating with domain tag")
			if len(domainList.Domain) > 0 {
				domains = []string{domainList.Domain}
			}
		}
	}
	if len(domains) <= 0 {
		fmt.Println("Domains are empty, calling default")
		domains = []string{s.Utils.GetDefaultDomain(s.cliConnection)}
	}
	fmt.Printf("Domains: %v", domains)
	return
}
