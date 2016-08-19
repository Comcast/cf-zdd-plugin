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

package canarydeploy_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	. "github.com/comcast/cf-zdd-plugin/canarydeploy"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	"github.com/comcast/cf-zdd-plugin/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("canaryDeploy", func() {
	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := canaryrepo.GetRegistry()[CanaryDeployCmdName]
				Ω(ok).Should(BeTrue())
			})
		})
	})
	Describe("given: a run() method on a canarydeploy object which has been initialized with valid args", func() {
		var canaryDeploy *CanaryDeploy
		var ctrlAppName = "myTestApp1.2.3#abcd"
		var ctrlManifestPath = "../fixtures/manifest.yml"
		var fakeUtilities *fakes.FakeUtilities

		BeforeEach(func() {
			fakeUtilities = new(fakes.FakeUtilities)

			canaryDeploy = &CanaryDeploy{
				Utils: fakeUtilities,
			}

			canaryDeploy.SetArgs([]string{CanaryDeployCmdName, ctrlAppName, "-f", ctrlManifestPath})
		})
		Context("when called with a valid connection object and a domain defined in the manifest", func() {
			var err error
			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlArgsNoRoute = []string{"push", ctrlAppName, "-f", ctrlManifestPath, "-i", "1", "--no-route", "--no-start"}
			var ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "mylocaldomain.com", "-n", CreateCanaryRouteName(ctrlAppName)}
			BeforeEach(func() {
				fakeConnection = new(pluginfakes.FakeCliConnection)
				err = canaryDeploy.Run(fakeConnection)
			})
			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(3))
				Ω(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsNoRoute))
				Ω(fakeConnection.CliCommandArgsForCall(1)).Should(Equal(ctrlArgsMapRoute))
			})
		})
		Context("when called with a valid connection object and multiple domains defined in the manifest", func() {
			var err error

			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlArgsNoRoute []string
			var ctrlArgsMapRoute []string
			BeforeEach(func() {
				ctrlArgsNoRoute = []string{"push", ctrlAppName, "-f", ctrlManifestPath, "-i", "1", "--no-route", "--no-start"}
				ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "mylocaldomain.com", "-n", CreateCanaryRouteName(ctrlAppName)}
				ctrlManifestPath = "../fixtures/manifest-multidomain.yml"
				fakeConnection = new(pluginfakes.FakeCliConnection)
				err = canaryDeploy.Run(fakeConnection)
			})
			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(4))
				Ω(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsNoRoute))
				Ω(fakeConnection.CliCommandArgsForCall(1)).Should(Equal(ctrlArgsMapRoute))
			})
		})
		Context("when called with a valid connection object and both domain and domains defined in the manifest", func() {
			var err error

			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlArgsNoRoute []string
			var ctrlArgsMapRoute []string
			BeforeEach(func() {
				ctrlArgsNoRoute = []string{"push", ctrlAppName, "-f", ctrlManifestPath, "-i", "1", "--no-route", "--no-start"}
				ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "mylocaldomain.com", "-n", CreateCanaryRouteName(ctrlAppName)}
				ctrlManifestPath = "../fixtures/manifest-bothdomain.yml"
				fakeConnection = new(pluginfakes.FakeCliConnection)
				err = canaryDeploy.Run(fakeConnection)
			})
			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(5))
				Ω(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsNoRoute))
				Ω(fakeConnection.CliCommandArgsForCall(1)).Should(Equal(ctrlArgsMapRoute))
			})
		})
		Context("when called with a valid connection object and no domain defined in the manifest", func() {
			var err error

			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlArgsNoRoute []string
			var ctrlArgsMapRoute []string
			BeforeEach(func() {
				ctrlManifestPath = "../fixtures/manifest-nodomain.yml"
				ctrlArgsNoRoute = []string{"push", ctrlAppName, "-f", ctrlManifestPath, "-i", "1", "--no-route", "--no-start"}
				ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "cloud.net", "-n", CreateCanaryRouteName(ctrlAppName)}

				fakeConnection = new(pluginfakes.FakeCliConnection)
				fakeUtilities.GetDefaultDomainReturns("cloud.net")
				err = canaryDeploy.Run(fakeConnection)
			})
			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(3))
				Ω(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsNoRoute))
				Ω(fakeConnection.CliCommandArgsForCall(1)).Should(Equal(ctrlArgsMapRoute))
			})
		})
	})
	Describe(".CreateCanaryRouteName string", func() {
		Context("when given an appname with dots", func() {
			var ctrlAppname = "ctrlAppName-1.2.3"
			It("should remove dots and return a valid canary routename", func() {
				routename := CreateCanaryRouteName(ctrlAppname)
				canaryRoute := fmt.Sprintf("%s-%s", ctrlAppname, CanaryRouteSuffix)
				canaryRoute = strings.Replace(canaryRoute, ".", CanaryRouteSeparator, -1)
				Ω(routename).Should(Equal(canaryRoute))
			})
		})
		Context("when given an appname containing #", func() {
			var ctrlAppname = "ctrlAppName#45"
			It("should remove hashes and return a valid canary routename", func() {
				routename := CreateCanaryRouteName(ctrlAppname)
				canaryRoute := fmt.Sprintf("%s-%s", ctrlAppname, CanaryRouteSuffix)
				canaryRoute = strings.Replace(canaryRoute, "#", CanaryRouteSeparator, -1)
				Ω(routename).Should(Equal(canaryRoute))
			})
		})
		Context("when given an appname", func() {
			var ctrlAppname = "ctrlAppName"
			It("should return a valid canary routename", func() {
				routename := CreateCanaryRouteName(ctrlAppname)
				Ω(routename).Should(Equal(fmt.Sprintf("%s-%s", ctrlAppname, CanaryRouteSuffix)))
			})
		})
	})
})
