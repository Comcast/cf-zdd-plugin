# cf-zdd-plugin
Zero downtime deployment plugin for Cloud Foundry

### Assets  
Binaries are available for download in the [releases](https://github.com/Comcast/cf-zdd-plugin/releases) section of this respository.

### Requirements
Currently built with Go [1.10.2](https://golang.org/dl/), and [CloudFoundry CLI v6.36.1](https://github.com/cloudfoundry/cli/releases). The below dependencies are also required.
```
go get code.cloudfoundry.org/cli/plugin
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/andrew-d/go-termutil
```

### Build
```
go build cf_zdd_plugin.go
```

### Install
```
cf install-plugin cf_zdd_plugin
cf plugins
```

### deploy-zdd
Zero downtime deployments for applications. The plugin is designed to deploy applications without impact to the current version. Current scenarios covered:  
  - For a new deployment, application is pushed and started according to manifest contents  
  - For a update deployment, the application is pushed and then scaled over to the new version  
  - For a redeployment of the same version, the old application is renamed and the new application is pushed, the new application is then scaled over to  
The method then leverages scaleover mechanism as described below.  

**Usage**
```sh
cf deploy-zdd myapplication -base-name myapp -f path/to/manifest.yml -p path/to/application 15s
```
**myapplication** - my application name  
**-base-name** - [Optional] base name of application if you are using versioned application names, i.e. `myapplication-1.2.3` the base name would be `myapplication`  
**--f** - path to application manifest  
**--p** - path to deployable artifact
**15s** - duration in which to deploy application
### deploy-canary
The deploy-canary method deploys a single instance of an application under a custom route that can then be tested against before calling the counterpart promote-canary below.  
**Usage**
```sh
cf deploy-canary mycanaryapp -f path/to/manifest.yml -p path/to/application
```
**mycanaryapp** - my application name  
**-f** - path to application manifest  
**-p** - path to deployable artifact

### promote-canary
The promote-canary method takes the deployed canary application and deploys it to become the live application, as before this utilizes the scaleover method.  
**Usage**
```sh
cf promote-canary myapplication mycanaryapp 15s
```
**myapplication** - my application name  
**mycanaryapp** - my canary name  
**15s** - scaleover duration

### scaleover
Scaleover rolls an application from one version to another without the extra capacity needed for blue/green deployments. The duration argument is the total time taken to roll the application from the old to the new version.  
**Usage**
```sh
cf scaleover app1 app2 15s
```
**app1** - my old application  
**app2** - my new application  
**15s** - scaleover duration

### blue-green
Blue green deployment deploys your code and then remaps the route from the old to the new versions.  
**Usage**  
```jshelllanguage
cf blue-green myapplication -f path/to/manifest.yml -p path/to/application
```
**--f** - path to application manifest  
**--p** - path to deployable artifact  

##TODO
1. Remove opinions surrounding filenames and versions.
2. Another pass at refactoring.
