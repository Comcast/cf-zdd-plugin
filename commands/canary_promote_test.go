package commands

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("canaryPromote", func() {
	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := GetRegistry()[CanaryPromoteCmdName]
				Ω(ok).Should(BeTrue())
			})
		})
	})

	Describe(".RemoveApplication", func() {
		Context("when called with an application that exists", func() {
			var (
				ctrlApp        plugin_models.GetAppModel
				fakeConnection *pluginfakes.FakeCliConnection
				canaryPromote  *CanaryPromote
				cfZddCmd       *CfZddCmd
			)
			BeforeEach(func() {
				fakeConnection = new(pluginfakes.FakeCliConnection)
				cfZddCmd = &CfZddCmd{
					CmdName:         "deploy-canary",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
				}

				canaryPromote = &CanaryPromote{}
				canaryPromote.SetArgs(cfZddCmd)

				ctrlApp = plugin_models.GetAppModel{
					Name: "ars-generic#0.1.1.8-5d1bef",
				}
				fakeConnection.CliCommandReturns([]string{"return"}, nil)
			})

			It("should remove the application", func() {
				err := canaryPromote.RemoveApplication(ctrlApp)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
	Describe(".UpdateRoutes", func() {
		var (
			fakeConnection *pluginfakes.FakeCliConnection
			canaryPromote  *CanaryPromote
			cfZddCmd       *CfZddCmd

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
				fakeConnection = new(pluginfakes.FakeCliConnection)
				cfZddCmd = &CfZddCmd{
					CmdName:         "deploy-canary",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
				}
				canaryPromote = &CanaryPromote{}
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
