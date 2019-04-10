///*
//* Copyright 2016 Comcast Cable Communications Management, LLC
//*
//* Licensed under the Apache License, Version 2.0 (the "License");
//* you may not use this file except in compliance with the License.
//* You may obtain a copy of the License at
//*
//* http://www.apache.org/licenses/LICENSE-2.0
//*
//* Unless required by applicable law or agreed to in writing, software
//* distributed under the License is distributed on an "AS IS" BASIS,
//* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//* See the License for the specific language governing permissions and
//* limitations under the License.
// */
//
package commands_test

import (
	"fmt"
	"github.com/comcast/cf-zdd-plugin/commands"
	"github.com/comcast/cf-zdd-plugin/fakes"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("canaryDeploy", func() {
	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := commands.GetRegistry()[commands.CanaryDeployCmdName]
				Expect(ok).Should(BeTrue())
			})
		})
	})
	Describe("given: a run() method on a canarydeploy object which has been initialized with valid args", func() {
		var (
			canaryDeploy     *commands.CanaryDeploy
			ctrlAppName      = "myTestApp1.2.3#abcd"
			ctrlManifestPath = "../fixtures/manifest.yml"
			cfZddCmd         *commands.CfZddCmd
			fakeConnection   *fakes.FakeCliConnection
			fakeCommand      *fakes.FakeCommonCmd
		)

		BeforeEach(func() {
			fakeConnection = new(fakes.FakeCliConnection)
			fakeCommand = new(fakes.FakeCommonCmd)

			cfZddCmd = &commands.CfZddCmd{
				CmdName:         "deploy-canary",
				NewApp:          "myTestApp1.2.3#abcd",
				ManifestPath:    "../fixtures/manifest.yml",
				ApplicationPath: "application.jar",
				Conn:            fakeConnection,
				Commands:        fakeCommand,
			}
			canaryDeploy = &commands.CanaryDeploy{}

			canaryDeploy.SetArgs(cfZddCmd)
		})
		Context("when called with a valid connection object and a domain defined in the manifest", func() {
			var err error

			var ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "mylocaldomain.com", "-n", commands.CreateCanaryRouteName(ctrlAppName)}
			BeforeEach(func() {
				err = canaryDeploy.Run()
			})
			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Expect(fakeConnection.CliCommandCallCount()).Should(Equal(2))
				Expect(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsMapRoute))
			})
		})
		Context("when called with a valid connection object and multiple domains defined in the manifest", func() {
			var err error
			var ctrlArgsMapRoute []string
			BeforeEach(func() {
				ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "mylocaldomain.com", "-n", commands.CreateCanaryRouteName(ctrlAppName)}
				ctrlManifestPath = "../fixtures/manifest-multidomain.yml"
				cfZddCmd.ManifestPath = ctrlManifestPath
				err = canaryDeploy.Run()
			})
			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Expect(fakeConnection.CliCommandCallCount()).Should(Equal(2))
				Expect(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsMapRoute))
			})
		})
		Context("when called with a valid connection object and no domain defined in the manifest", func() {
			var err error
			var ctrlArgsMapRoute []string

			BeforeEach(func() {
				ctrlManifestPath = "../fixtures/manifest-nodomain.yml"
				cfZddCmd.ManifestPath = ctrlManifestPath
				fakeCommand.GetDefaultDomainReturns("u1.app.cloud.comcast.net")
				ctrlArgsMapRoute = []string{"map-route", ctrlAppName, "u1.app.cloud.comcast.net", "-n", commands.CreateCanaryRouteName(ctrlAppName)}
				err = canaryDeploy.Run()
			})
			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should deploy an application with a canary route", func() {
				Expect(fakeConnection.CliCommandCallCount()).Should(Equal(2))
				Expect(fakeConnection.CliCommandArgsForCall(0)).Should(Equal(ctrlArgsMapRoute))
			})
		})
	})
	Describe(".CreateCanaryRouteName string", func() {
		Context("when given an appname with dots", func() {
			var ctrlAppname = "ctrlAppName-1.2.3"
			It("should remove dots and return a valid canary routename", func() {
				routename := commands.CreateCanaryRouteName(ctrlAppname)
				canaryRoute := fmt.Sprintf("%s-%s", ctrlAppname, commands.CanaryRouteSuffix)
				canaryRoute = strings.Replace(canaryRoute, ".", commands.CanaryRouteSeparator, -1)
				Expect(routename).Should(Equal(canaryRoute))
			})
		})
		Context("when given an appname containing #", func() {
			var ctrlAppname = "ctrlAppName#45"
			It("should remove hashes and return a valid canary routename", func() {
				routename := commands.CreateCanaryRouteName(ctrlAppname)
				canaryRoute := fmt.Sprintf("%s-%s", ctrlAppname, commands.CanaryRouteSuffix)
				canaryRoute = strings.Replace(canaryRoute, "#", commands.CanaryRouteSeparator, -1)
				Expect(routename).Should(Equal(canaryRoute))
			})
		})
		Context("when given an appname", func() {
			var ctrlAppname = "ctrlAppName"
			It("should return a valid canary routename", func() {
				routename := commands.CreateCanaryRouteName(ctrlAppname)
				Expect(routename).Should(Equal(fmt.Sprintf("%s-%s", ctrlAppname, commands.CanaryRouteSuffix)))
			})
		})
	})
})
