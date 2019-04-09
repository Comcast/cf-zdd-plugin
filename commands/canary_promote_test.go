package commands_test

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"github.com/comcast/cf-zdd-plugin/commands"
	"github.com/comcast/cf-zdd-plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("canaryPromote", func() {
	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := commands.GetRegistry()[commands.CanaryPromoteCmdName]
				Ω(ok).Should(BeTrue())
			})
		})
	})

	Describe("Run", func() {
		var (
			fakeConnection *fakes.FakeCliConnection
			canaryPromote  *commands.CanaryPromote
			cfZddCmd       *commands.CfZddCmd
			fakeScaleover  *fakes.FakeScaleoverCommand
			fakeCommand    *fakes.FakeCommonCmd
			err            error
		)
		Context("when called with a valid set of args", func() {
			BeforeEach(func() {
				fakeConnection = new(fakes.FakeCliConnection)
				fakeScaleover = new(fakes.FakeScaleoverCommand)
				fakeCommand = new(fakes.FakeCommonCmd)

				fakeScaleover.DoScaleoverReturns(nil)

				cfZddCmd = &commands.CfZddCmd{
					OldApp:   "app1",
					NewApp:   "canary",
					Conn:     fakeConnection,
					Commands: fakeCommand,
				}

				canaryPromote = &commands.CanaryPromote{
					ScaleoverCmd: fakeScaleover,
				}

				canaryPromote.SetArgs(cfZddCmd)
			})
			It("should execute the promotion of the canary and not return an error", func() {
				err = canaryPromote.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
	Describe(".UpdateRoutes", func() {
		var (
			fakeConnection *fakes.FakeCliConnection
			canaryPromote  *commands.CanaryPromote
			cfZddCmd       *commands.CfZddCmd

			app1 plugin_models.GetAppModel
			app2 plugin_models.GetAppModel
		)
		Context("given 2 app models with different routes", func() {
			BeforeEach(func() {
				app1 = plugin_models.GetAppModel{
					Name: "app1",
					Routes: []plugin_models.GetApp_RouteSummary{
						plugin_models.GetApp_RouteSummary{
							Host: "app1",
							Domain: plugin_models.GetApp_DomainFields{
								Name: "cf.app.io",
							},
						},
					},
				}
				app2 = plugin_models.GetAppModel{
					Name: "app2",
					Routes: []plugin_models.GetApp_RouteSummary{
						plugin_models.GetApp_RouteSummary{
							Host: "app2",
							Domain: plugin_models.GetApp_DomainFields{
								Name: "cf.app.io",
							},
						},
					},
				}
				fakeConnection = new(fakes.FakeCliConnection)
				cfZddCmd = &commands.CfZddCmd{
					CmdName:         "deploy-canary",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
				}
				canaryPromote = &commands.CanaryPromote{}
				canaryPromote.SetArgs(cfZddCmd)
			})
			It("map route should map app1 routes to app2", func() {
				err := canaryPromote.UpdateRoutes(app1, app2)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeConnection.CliCommandCallCount()).Should(Equal(2))
			})
		})
	})
})
