package commands_test

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	"github.com/comcast/cf-zdd-plugin/commands"
	"github.com/comcast/cf-zdd-plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("zddDeploy", func() {

	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := commands.GetRegistry()[commands.ZddDeployCmdName]
				Ω(ok).Should(BeTrue())
			})
		})
	})

	Describe("when calling a run method on zddDeploy which has been initialized with valid args", func() {
		var (
			fakeConnection *fakes.FakeCliConnection
			fakeCommands   *fakes.FakeCommonCmd
			fakeScaleover  *fakes.FakeScaleoverCommand
			zddDeploy      *commands.ZddDeploy
			cfZddCmd       *commands.CfZddCmd
			//ctrlAppName        = "myTestApp#1.2.3-abcde"
			//ctrlManifestPath   = "../fixtures/manifest.yml"
			//ctrlPathToArtifact = "application.jar"
			//ctrlArgs           = []string{ZddDeployCmdName, ctrlAppName, "-f", ctrlManifestPath, "-p", ctrlPathToArtifact}
			err error
		)
		BeforeEach(func() {
			fakeConnection = new(fakes.FakeCliConnection)
			zddDeploy = new(commands.ZddDeploy)
			fakeCommands = new(fakes.FakeCommonCmd)

			cfZddCmd = &commands.CfZddCmd{
				CmdName:         commands.ZddDeployCmdName,
				NewApp:          "myTestApp#1.2.3-abcde",
				ManifestPath:    "../fixtures/manifest.yml",
				ApplicationPath: "application.jar",
				Conn:            fakeConnection,
				Commands:        fakeCommands,
			}

			zddDeploy = new(commands.ZddDeploy)
			fakeScaleover = new(fakes.FakeScaleoverCommand)
			zddDeploy.SetArgs(cfZddCmd)
			zddDeploy.ScalerOverCmd = fakeScaleover
		})

		Context("when called with a valid connection object for a new application deploy", func() {
			BeforeEach(func() {
				fakeCommands.IsApplicationDeployedReturns("", false)
				fakeCommands.PushApplicationReturns(nil)
			})
			It("should issue a single push of the application", func() {
				err = zddDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("when called with a valid connection object for a new deploy", func() {
			BeforeEach(func() {
				fakeCommands.IsApplicationDeployedReturns("myTestApp1.2.3#abcd", true)
				fakeCommands.PushApplicationReturns(nil)
				fakeCommands.RemoveApplicationReturns(nil)
				fakeScaleover.DoScaleoverReturns(nil)
				cfZddCmd.BaseAppName = "myTestApp"
			})
			It("should not return an error and scale over", func() {
				err = zddDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
				//Expect(fakeConnection.CliCommandCallCount()).Should(Equal(4))
			})
		})

	})
	XDescribe("given: a valid run() method on a zdddeploy object which has been initialized with valid args", func() {
		var zddDeploy *commands.ZddDeploy
		var cfZddCmd *commands.CfZddCmd
		var fakeConnection *pluginfakes.FakeCliConnection
		var ctrlAppName = "myTestApp#1.2.3-abcde"
		var ctrlManifestPath = "../fixtures/manifest.yml"
		var ctrlPathToArtifact = "application.jar"
		var ctrlArgs = []string{commands.ZddDeployCmdName, ctrlAppName, "-f", ctrlManifestPath, "-p", ctrlPathToArtifact}

		XContext("when called with a valid connection object for a new application deploy", func() {

			var err error

			BeforeEach(func() {

				fakeConnection = new(pluginfakes.FakeCliConnection)

				cfZddCmd = &commands.CfZddCmd{
					CmdName:         commands.ZddDeployCmdName,
					NewApp:          "myTestApp#1.2.3-abcde",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
				}

				zddDeploy = new(commands.ZddDeploy)
				zddDeploy.SetArgs(cfZddCmd)

				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{{Name: "1234"}}, nil)

				err = zddDeploy.Run()
			})
			It("should issue a single push of the application", func() {
				Ω(err).ShouldNot(HaveOccurred())
				args := append([]string{"push"}, ctrlArgs[1:]...)
				Ω(fakeConnection.CliCommandArgsForCall(fakeConnection.CliCommandCallCount() - 1)).Should(Equal(args))
			})
		})
		XContext("when called with a valid connection object for a new deploy", func() {
			var (
				err error
			)
			BeforeEach(func() {

				fakeConnection = new(pluginfakes.FakeCliConnection)

				cfZddCmd = &commands.CfZddCmd{
					CmdName:         "deploy-zdd",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
					BaseAppName:     "myTestApp",
				}

				fakeConnection.GetAppsReturns([]plugin_models.GetAppsModel{{Name: "myTestApp1.2.3#abcd"}, {Name: "anotherApp1.2.3#abcd"}}, nil)

				zddDeploy = new(commands.ZddDeploy)
				zddDeploy.SetArgs(cfZddCmd)

				err = zddDeploy.Run()
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			//It("should deploy an application in a scaleover way", func() {
			//	Ω(fakeConnection.CliCommandCallCount()).Should(Equal(4))
			//})
		})
		XContext("when called with a valid connection object for a new deploy and with additional scaleover arg", func() {
			var (
				cfZddCmd          *commands.CfZddCmd
				zddDeploy         *commands.ZddDeploy
				err               error
				fakeCliConnection *pluginfakes.FakeCliConnection
			)
			BeforeEach(func() {
				fakeCliConnection = new(pluginfakes.FakeCliConnection)
				cfZddCmd = &commands.CfZddCmd{
					CmdName:         "deploy-zdd",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeCliConnection,
				}

				zddDeploy = new(commands.ZddDeploy)
				zddDeploy.SetArgs(cfZddCmd)

				err = zddDeploy.Run()
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application in a scaleover way", func() {
				Ω(fakeCliConnection.CliCommandCallCount()).Should(Equal(1))
			})
		})

		XContext("when called as a redeploy of the current app version", func() {
			var (
				fakeCliConnection *pluginfakes.FakeCliConnection
				returnModels      []plugin_models.GetAppsModel
				err               error
				callOrder         = map[string]int{
					"rename": 0,
					"push":   1,
					"delete": 2,
				}
				evaluateCallChainForAction = func(action string) {
					command := fakeCliConnection.CliCommandArgsForCall(callOrder[action])[0]
					Ω(command).Should(Equal(action))
				}
			)

			BeforeEach(func() {
				returnModels = make([]plugin_models.GetAppsModel, 0)
				fakeCliConnection = new(pluginfakes.FakeCliConnection)
				returnModels = append(returnModels,
					plugin_models.GetAppsModel{Name: ctrlAppName},
				)
				fakeCliConnection = new(pluginfakes.FakeCliConnection)
				cfZddCmd = &commands.CfZddCmd{
					CmdName:         "deploy-zdd",
					NewApp:          ctrlAppName,
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeCliConnection,
				}

				zddDeploy = new(commands.ZddDeploy)
				zddDeploy.SetArgs(cfZddCmd)
				fakeCliConnection.GetAppsReturns(returnModels, nil)
				fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{"requested state: stopped\ninstances: 0/10\nurls: app.cfapps.io"}, nil)
				err = zddDeploy.Run()
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have renamed the old version", func() {
				action := "rename"
				evaluateCallChainForAction(action)
			})

			It("should have pushed the new version", func() {
				action := "push"
				evaluateCallChainForAction(action)
			})

			It("should have deleted the old version", func() {
				action := "delete"
				evaluateCallChainForAction(action)
			})
		})
		XContext("when called as a newly deployed app version", func() {
			var (
				fakeCliConnection *pluginfakes.FakeCliConnection
				returnModels      []plugin_models.GetAppsModel
				ctrlAppNameV2     string
				err               error
				callOrder         = map[string]int{
					"push":   0,
					"delete": 1,
				}
				evaluateCallChainForAction = func(action string) {
					command := fakeCliConnection.CliCommandArgsForCall(callOrder[action])[0]
					Ω(command).Should(Equal(action))
				}
				cfZddCmd  *commands.CfZddCmd
				zddDeploy *commands.ZddDeploy
			)

			BeforeEach(func() {
				fakeConnection = new(pluginfakes.FakeCliConnection)

				cfZddCmd = &commands.CfZddCmd{
					CmdName:         "deploy-zdd",
					NewApp:          "myTestApp#1.2.3-abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeCliConnection,
				}

				returnModels = make([]plugin_models.GetAppsModel, 0)
				ctrlAppNameV2 = "myTestApp#1.2.2-abcde"
				returnModels = append(returnModels, plugin_models.GetAppsModel{Name: ctrlAppNameV2})
				zddDeploy = new(commands.ZddDeploy)
				fakeCliConnection = new(pluginfakes.FakeCliConnection)
				fakeCliConnection.GetAppsReturns(returnModels, nil)
				zddDeploy.SetArgs(cfZddCmd)
				err = zddDeploy.Run()
			})

			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})

			It("should have pushed the new version", func() {
				action := "push"
				evaluateCallChainForAction(action)
			})

			It("should have deleted the old version", func() {
				action := "delete"
				evaluateCallChainForAction(action)
			})
		})
	})
})
