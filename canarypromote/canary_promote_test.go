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

package canarypromote_test

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/plugin/pluginfakes"
	. "github.com/comcast/cf-zdd-plugin/canarypromote"
	"github.com/comcast/cf-zdd-plugin/canaryrepo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("canaryPromote", func() {
	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the canary repo", func() {
				_, ok := canaryrepo.GetRegistry()[CanaryPromoteCmdName]
				Ω(ok).Should(BeTrue())
			})
		})
	})
	Describe("given: a run() method on a canarypromote object which has been initialized with valid args", func() {
		var canaryPromote *CanaryPromote
		var ctrlAppName = "ars-generic#0.1.1.8-5d1bef"
		var ctrlCanaryAppName = "ars-generic#0.1.1.9-5rtbef"
		var ctrlScaleoverTime = "15s"

		BeforeEach(func() {
			canaryPromote = new(CanaryPromote)
			canaryPromote.SetArgs([]string{CanaryPromoteCmdName, ctrlAppName, ctrlCanaryAppName, ctrlScaleoverTime})
		})
		Context("when called with a valid connection object and multiple routes", func() {
			var err error
			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlMapRouteArgs = []string{"map-route", ctrlCanaryAppName, "ula.app.cloud.comcast.net", "-n", "ars-dev"}
			BeforeEach(func() {
				fakeConnection = new(pluginfakes.FakeCliConnection)
				fakeConnection.CliCommandWithoutTerminalOutputReturns([]string{"requested state: stopped\ninstances: 0/10\nurls: app.cfapps.io"}, nil)
				fakeConnection.CliCommandReturns([]string{"requested state: started",
					"instances: 1/1", "usage: 1G x 1 instances",
					"urls: ars-dev.cloud.net,ars-int.cloud.net",
					"last uploaded: Thu Dec 17 15:34:33 UTC 2015", "stack: cflinuxfs2"}, nil)
				err = canaryPromote.Run(fakeConnection)
			})
			It("should not return an error", func() {
				Ω(err).ShouldNot(HaveOccurred())
			})
			It("should scaleover an application and remove the canary route", func() {
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(7))
				Ω(fakeConnection.CliCommandArgsForCall(2)).Should(Equal(ctrlMapRouteArgs))
			})
		})

		Context("when called with a valid connection object and pooled application", func() {
			var err error
			var fakeConnection *pluginfakes.FakeCliConnection
			var ctrlMapRouteArgs = []string{"map-route", ctrlCanaryAppName, "cloud.net", "-n", "ars-dev"}
			BeforeEach(func() {
				fakeConnection = new(pluginfakes.FakeCliConnection)
				fakeConnection.CliCommandReturns([]string{"requested state: started",
					"instances: 1/1", "usage: 1G x 1 instances",
					"urls: ars-dev.cloud.net",
					"last uploaded: Thu Dec 17 15:34:33 UTC 2015", "stack: cflinuxfs2"}, nil)

				err = canaryPromote.Run(fakeConnection)
			})
			It("should apply the correct route to the application", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fakeConnection.CliCommandCallCount()).Should(Equal(5))
				Ω(fakeConnection.CliCommandArgsForCall(2)).Should(Equal(ctrlMapRouteArgs))
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
