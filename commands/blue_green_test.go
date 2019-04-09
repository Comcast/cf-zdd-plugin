/*
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package commands_test

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"github.com/comcast/cf-zdd-plugin/commands"
	"github.com/comcast/cf-zdd-plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".BlueGreenDeploy", func() {

	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the command repo", func() {
				_, ok := commands.GetRegistry()[commands.BlueGreenCmdName]
				Expect(ok).Should(BeTrue())
			})
		})
	})

	Describe("with a valid arg and run method", func() {
		var (
			err            error
			bgDeploy       *commands.BlueGreenDeploy
			cfZddCmd       *commands.CfZddCmd
			fakeConnection *fakes.FakeCliConnection
			fakeCommon     *fakes.FakeCommonCmd
		)

		BeforeEach(func() {
			fakeConnection = new(fakes.FakeCliConnection)
			fakeCommon = new(fakes.FakeCommonCmd)

		})
		Context("when called with an application not previously deployed", func() {
			BeforeEach(func() {
				fakeCommon.IsApplicationDeployedReturns("", false)
				fakeCommon.PushApplicationReturns(nil)

				cfZddCmd = &commands.CfZddCmd{
					CmdName:         commands.BlueGreenCmdName,
					NewApp:          "myTestApp#1.2.3-abcde",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
					Commands:        fakeCommon,
				}

				bgDeploy = new(commands.BlueGreenDeploy)
				bgDeploy.SetArgs(cfZddCmd)
			})
			It("should push the application and not return an error", func() {
				err = bgDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when called with a new version of the application", func() {
			BeforeEach(func() {
				fakeCommon.IsApplicationDeployedReturns("myTestApp#1.2.2-abcde", true)

				cfZddCmd = &commands.CfZddCmd{
					CmdName:         commands.BlueGreenCmdName,
					NewApp:          "myTestApp#1.2.3-abcde",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Conn:            fakeConnection,
					Commands:        fakeCommon,
					BaseAppName:     "mytestapp",
				}
				bgDeploy = new(commands.BlueGreenDeploy)
				bgDeploy.SetArgs(cfZddCmd)

				fakeConnection.GetAppReturnsOnCall(0, plugin_models.GetAppModel{
					RunningInstances: 1,
					InstanceCount:    2,
				}, nil)

				fakeConnection.GetAppReturnsOnCall(1, plugin_models.GetAppModel{
					RunningInstances: 2,
					InstanceCount:    2,
				}, nil)

			})
			It("should complete a blue-green deployment", func() {
				err = bgDeploy.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
