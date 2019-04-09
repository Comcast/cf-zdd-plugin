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
	"errors"
	"github.com/comcast/cf-zdd-plugin/commands"
	"github.com/comcast/cf-zdd-plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("common command", func() {

	var (
		fakeCliConnection *fakes.FakeCliConnection
		cmd               commands.CommonCmd
	)

	BeforeEach(func() {
		fakeCliConnection = new(fakes.FakeCliConnection)
		cmd = commands.NewCommonCmd(fakeCliConnection)
	})

	Describe(".GetDefaultDomain", func() {
		var (
			domain     string
			ctrlDomain string
		)
		BeforeEach(func() {
			ctrlDomain = "cloud.net"
			b, err := ioutil.ReadFile("../fixtures/curlResponse.txt")
			if err != nil {
				Skip(err.Error())
			}
			fakeCliConnection.CliCommandWithoutTerminalOutputReturns([]string{string(b[:])}, nil)
		})
		Context("when called", func() {
			It("should return the default domain", func() {
				domain = cmd.GetDefaultDomain()
				Expect(domain).Should(Equal(ctrlDomain))
			})
		})
	})

	Describe(".IsApplicationDeployed", func() {
		var (
			ctrlAppName = "demoApp-1.2.3"
			baseAppName = "demoApp"
		)
		Context("when called for an application that is deployed", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppsReturns([]plugin_models.GetAppsModel{{
					Name: ctrlAppName,
				}, {
					Name: "someOtherApp",
				}}, nil)
			})
			It("should return an app and return true", func() {
				app, dep := cmd.IsApplicationDeployed(baseAppName)
				Expect(app).Should(Equal(ctrlAppName))
				Expect(dep).Should(BeTrue())
			})
		})
		Context("when called for an application that is not deployed", func() {
			BeforeEach(func() {
				fakeCliConnection.GetAppsReturns([]plugin_models.GetAppsModel{{
					Name: "someOtherApp",
				}}, nil)
			})
			It("should not return an app and return false", func() {
				app, dep := cmd.IsApplicationDeployed(baseAppName)
				Expect(app).Should(Equal(""))
				Expect(dep).Should(BeFalse())
			})
		})
	})

	Describe(".PushApplication", func() {
		Context("when called with a valid appliction, path, and manifest", func() {
			var (
				err error
			)
			BeforeEach(func() {
				fakeCliConnection.CliCommandReturns(nil, nil)
			})
			It("should deploy the application and not return an error", func() {
				err = cmd.PushApplication("appname", "appPath", "manifestPath", "")
				Expect(err).ShouldNot(HaveOccurred())
			})

		})
		Context("when called with a valid appliction, path, and no manifest", func() {
			var (
				err error
			)
			BeforeEach(func() {
				fakeCliConnection.CliCommandReturns(nil, nil)
			})
			It("should deploy the application and not return an error", func() {
				err = cmd.PushApplication("appname", "appPath", "", "")
				Expect(err).ShouldNot(HaveOccurred())
			})

		})
		Context("when called with a valid appliction, no path, and manifest", func() {
			var (
				err error
			)
			BeforeEach(func() {
				fakeCliConnection.CliCommandReturns(nil, nil)
			})
			It("should deploy the application and not return an error", func() {
				err = cmd.PushApplication("appname", "", "manifestPath", "")
				Expect(err).ShouldNot(HaveOccurred())
			})

		})
		Context("when called with a invalid appliction, path, and manifest", func() {
			var (
				err error
			)
			BeforeEach(func() {
				fakeCliConnection.CliCommandReturns(nil, errors.New("appname must be specified"))
			})
			It("should not deploy the application and return an error", func() {
				err = cmd.PushApplication("", "appPath", "manifestPath", "")
				Expect(err).Should(HaveOccurred())
			})

		})
	})

	Describe(".RenameApplication", func() {
		var (
			err error
		)
		Context("when called with a valid application name and new name", func() {
			It("should rename the application and not return an error", func() {
				err = cmd.RenameApplication("appname", "newappname")
				Expect(err).ShouldNot(HaveOccurred())
			})

		})
		Context("when called with an empty application name and new name", func() {
			It("should not rename the application and return an error", func() {
				err = cmd.RenameApplication("", "newappname")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe(".RemoveApplication", func() {
		var (
			err error
		)
		Context("when called with a valid application name and new name", func() {
			It("should delete the application and not return an error", func() {
				err = cmd.RemoveApplication("appname")
				Expect(err).ShouldNot(HaveOccurred())
			})

		})
		Context("when called with an empty application name and new name", func() {
			It("should not delete the application and return an error", func() {
				err = cmd.RemoveApplication("")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Describe(".RemapRoutes", func() {
		Context("when called with a valid application", func() {
			var (
				err error
			)
			BeforeEach(func() {
				fakeCliConnection.GetAppReturns(plugin_models.GetAppModel{
					Routes: []plugin_models.GetApp_RouteSummary{
						{
							Domain: plugin_models.GetApp_DomainFields{
								Name: "adomain.com",
							},
							Host: "myapp",
						},
					},
				}, nil)

				fakeCliConnection.CliCommandReturnsOnCall(0, nil, nil)
				fakeCliConnection.CliCommandReturnsOnCall(1, nil, nil)
			})

			It("should remap the routes", func() {
				err = cmd.RemapRoutes("oldApp", "newApp")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
