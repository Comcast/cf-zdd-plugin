package commands

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".BlueGreenDeploy", func() {

	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the command repo", func() {
				_, ok := GetRegistry()[BlueGreenCmdName]
				Expect(ok).Should(BeTrue())
			})
		})
	})

	Describe("with a valid arg and run method", func() {
		var err error
		var bgDeploy *BlueGreenDeploy
		var cfZddCmd *CfZddCmd
		var fakeConnection *pluginfakes.FakeCliConnection

		BeforeEach(func() {
			fakeConnection = new(pluginfakes.FakeCliConnection)
			cfZddCmd = &CfZddCmd{
				CmdName:         BlueGreenCmdName,
				NewApp:          "myTestApp#1.2.3-abcde",
				ManifestPath:    "../fixtures/manifest.yml",
				ApplicationPath: "application.jar",
				Conn:            fakeConnection,
			}

			bgDeploy = new(BlueGreenDeploy)
			bgDeploy.SetArgs(cfZddCmd)
		})
		Context("when called with an application not previously deployed", func() {
			BeforeEach(func() {
				returnModels := append([]plugin_models.GetAppsModel{}, plugin_models.GetAppsModel{
					Name: "1234",
				})
				returnModel := plugin_models.GetAppModel{
					Routes: []plugin_models.GetApp_RouteSummary{
						{
							Host: "myHost",
							Domain: plugin_models.GetApp_DomainFields{
								Name: "somedomain.com",
							},
						},
					},
				}
				fakeConnection.GetAppsReturns(returnModels, nil)
				fakeConnection.GetAppReturns(returnModel, nil)
			})
			It("should push the application and not return an error", func() {
				err = bgDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeConnection.CliCommandCallCount()).Should(Equal(1))
				Expect(fakeConnection.GetAppsCallCount()).Should(Equal(1))
			})
		})
		Context("when called with a new version of the application", func() {
			BeforeEach(func() {
				returnModels := append([]plugin_models.GetAppsModel{}, plugin_models.GetAppsModel{
					Name: "myTestApp#1.2.3-abcde",
				})
				fakeConnection.GetAppsReturns(returnModels, nil)
			})
			It("should complete a blue-green deployment", func() {
				err = bgDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
