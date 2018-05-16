// Copyright 2016 Joshua Kruck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"errors"
	"time"

	"code.cloudfoundry.org/cli/plugin/models"
	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scaleover", func() {
	var scaleoverCmdPlugin *ScaleoverCmd
	var fakeCliConnection *pluginfakes.FakeCliConnection
	domain := plugin_models.GetApp_DomainFields{Name: "cfapps.io"}
	Describe("getAppStatus", func() {

		BeforeEach(func() {
			fakeCliConnection = &pluginfakes.FakeCliConnection{}
			scaleoverCmdPlugin = &ScaleoverCmd{
				Args: &CfZddCmd{
					Conn: fakeCliConnection,
				},
			}
		})

		It("should Fail Without App 1", func() {
			app := plugin_models.GetAppModel{}
			fakeCliConnection.GetAppReturns(app, errors.New("App app1 not found"))
			var err error
			_, err = scaleoverCmdPlugin.GetAppStatus("app1")
			Expect(err.Error()).To(Equal("App app1 not found"))
		})

		It("should Fail Without App 2", func() {
			app := plugin_models.GetAppModel{}
			fakeCliConnection.GetAppReturns(app, errors.New("App app2 not found"))
			var err error
			_, err = scaleoverCmdPlugin.GetAppStatus("app2")
			Expect(err.Error()).To(Equal("App app2 not found"))
		})

		It("should not start a stopped target with 1 instance", func() {
			app := plugin_models.GetAppModel{State: "stopped"}
			fakeCliConnection.GetAppReturns(app, nil)

			var status *AppStatus
			status, _ = scaleoverCmdPlugin.GetAppStatus("app1")

			Expect(status.Name).To(Equal("app1"))
			Expect(status.CountRequested).To(Equal(0))
			Expect(status.CountRunning).To(Equal(0))
			Expect(status.State).To(Equal("stopped"))
		})

		It("should start a started app with 10 instances", func() {

			cfApp := plugin_models.GetAppModel{
				InstanceCount:    10,
				RunningInstances: 10,
				State:            "started",
			}

			fakeCliConnection.GetAppReturns(cfApp, nil)

			var status *AppStatus
			status, _ = scaleoverCmdPlugin.GetAppStatus("app1")

			Expect(status.Name).To(Equal("app1"))
			Expect(status.CountRequested).To(Equal(10))
			Expect(status.CountRunning).To(Equal(10))
			Expect(status.State).To(Equal("started"))
		})

		It("should keep a stop app stopped with 10 instances", func() {
			cfApp := plugin_models.GetAppModel{
				InstanceCount:    10,
				RunningInstances: 0,
				State:            "stopped",
			}

			fakeCliConnection.GetAppReturns(cfApp, nil)

			var status *AppStatus
			status, _ = scaleoverCmdPlugin.GetAppStatus("app1")

			Expect(status.Name).To(Equal("app1"))
			Expect(status.CountRequested).To(Equal(0))
			Expect(status.CountRunning).To(Equal(0))
			Expect(status.State).To(Equal("stopped"))
		})

		It("should populate the routes for an app with one url", func() {
			routes := []plugin_models.GetApp_RouteSummary{
				{
					Host:   "app",
					Domain: domain,
				},
			}
			var status *AppStatus
			app := plugin_models.GetAppModel{Routes: routes}
			fakeCliConnection.GetAppReturns(app, nil)
			status, _ = scaleoverCmdPlugin.GetAppStatus("app1")
			Expect(len(status.Routes)).To(Equal(1))
			Expect(status.Routes[0]).To(Equal("app.cfapps.io"))
		})

		It("should populate the routes for an app with three urls", func() {
			routes := []plugin_models.GetApp_RouteSummary{
				{
					Host:   "app",
					Domain: domain,
				},
				{
					Host:   "foo-app",
					Domain: domain,
				},
				{
					Host:   "foo-app-b",
					Domain: domain,
				},
			}
			var status *AppStatus
			app := plugin_models.GetAppModel{Routes: routes}
			fakeCliConnection.GetAppReturns(app, nil)

			status, _ = scaleoverCmdPlugin.GetAppStatus("app1")
			Expect(len(status.Routes)).To(Equal(3))
			Expect(status.Routes[0]).To(Equal("app.cfapps.io"))
			Expect(status.Routes[1]).To(Equal("foo-app.cfapps.io"))
			Expect(status.Routes[2]).To(Equal("foo-app-b.cfapps.io"))
		})

	})

	Describe("It should handle weird time inputs", func() {
		BeforeEach(func() {
			scaleoverCmdPlugin = &ScaleoverCmd{}
		})

		It("like a negative number", func() {
			var err error
			_, err = scaleoverCmdPlugin.ParseTime("-1m")

			Expect(err.Error()).To(Equal("Duration must be a positive number in the format of 1m"))
		})

		It("like zero", func() {
			var t time.Duration
			t, _ = scaleoverCmdPlugin.ParseTime("0m")
			one, _ := time.ParseDuration("0s")
			Expect(t).To(Equal(one))
		})
	})

	Describe("scale up", func() {
		var appStatus *AppStatus

		BeforeEach(func() {
			appStatus = &AppStatus{
				Name:           "foo",
				CountRequested: 1,
				CountRunning:   1,
				State:          "stopped",
			}
			fakeCliConnection = &pluginfakes.FakeCliConnection{}

		})

		It("Starts a stopped app", func() {
			appStatus.ScaleUp(fakeCliConnection)
			Expect(appStatus.State).To(Equal("started"))
		})

		It("It increments the amount requested", func() {
			running := appStatus.CountRunning
			appStatus.ScaleUp(fakeCliConnection)
			Expect(appStatus.CountRequested).To(Equal(running + 1))
		})

		It("Leaves a started app started", func() {
			appStatus.State = "started"
			appStatus.ScaleUp(fakeCliConnection)
			Expect(appStatus.State).To(Equal("started"))
		})

	})

	Describe("scale down", func() {
		var appStatus *AppStatus

		BeforeEach(func() {
			appStatus = &AppStatus{
				Name:           "foo",
				CountRequested: 1,
				CountRunning:   1,
				State:          "started",
			}
			fakeCliConnection = &pluginfakes.FakeCliConnection{}

		})

		It("Stops a started app going to zero instances", func() {
			appStatus.ScaleDown(fakeCliConnection)
			Expect(appStatus.State).To(Equal("stopped"))
		})

		It("It decrements the amount requested", func() {
			running := appStatus.CountRunning
			appStatus.ScaleDown(fakeCliConnection)
			Expect(appStatus.CountRequested).To(Equal(running - 1))
		})

		It("Leaves a stopped app stopped", func() {
			appStatus.State = "stopped"
			appStatus.ScaleDown(fakeCliConnection)
			Expect(appStatus.State).To(Equal("stopped"))
		})

		It("Scales down the app", func() {
			appStatus.CountRequested = 2
			appStatus.ScaleDown(fakeCliConnection)
			Expect(appStatus.CountRunning).To(Equal(1))
			Expect(fakeCliConnection.CliCommandWithoutTerminalOutputCallCount()).To(Equal(1))
		})
	})

	Describe("Usage", func() {
		BeforeEach(func() {
			scaleoverCmdPlugin = &ScaleoverCmd{
				Args: &CfZddCmd{
					CmdName:         "scaleover",
					NewApp:          "myTestApp1.2.3#abcd",
					ManifestPath:    "../fixtures/manifest.yml",
					ApplicationPath: "application.jar",
					Duration:        "480s",
					Conn:            new(pluginfakes.FakeCliConnection),
				},
			}
		})

		It("shows usage for too few arguments", func() {
			Expect(scaleoverCmdPlugin.Usage(scaleoverCmdPlugin.Args)).NotTo(BeNil())
		})

		It("is just right", func() {
			scaleoverCmdPlugin.Args.OldApp = "myTestApp1.2.3#abcd"
			Expect(scaleoverCmdPlugin.Usage(scaleoverCmdPlugin.Args)).To(BeNil())
		})

		It("is okay with duration set", func() {
			scaleoverCmdPlugin.Args.OldApp = "myTestApp1.2.3#abcd"
			Expect(scaleoverCmdPlugin.Usage(scaleoverCmdPlugin.Args)).To(BeNil())
		})

		It("is not ok with duration and custom url not set", func() {
			scaleoverCmdPlugin.Args.Duration = ""
			Expect(scaleoverCmdPlugin.Usage(scaleoverCmdPlugin.Args)).ToNot(BeNil())
		})

	})

	Describe("Routes", func() {
		BeforeEach(func() {
			cfZddCmd := &CfZddCmd{
				RouteCheck: false,
			}

			scaleoverCmdPlugin = &ScaleoverCmd{
				Args: cfZddCmd,
			}
			var app1 = &AppStatus{
				Routes: []string{"a.b.c", "b.c.d"},
			}
			var app2 = &AppStatus{
				Routes: []string{"c.d.e", "d.e.f"},
			}
			scaleoverCmdPlugin.App1 = app1
			scaleoverCmdPlugin.App2 = app2
		})

		It("should return false if the apps don't share a route", func() {
			Expect(scaleoverCmdPlugin.AppsShareARoute()).To(BeFalse())
		})

		It("should return true when they share a route", func() {
			scaleoverCmdPlugin.App2 = scaleoverCmdPlugin.App1
			Expect(scaleoverCmdPlugin.AppsShareARoute()).To(BeTrue())
		})

		It("Should warn when apps don't share a route", func() {
			Expect(scaleoverCmdPlugin.ErrorIfNoSharedRoute().Error()).To(Equal("Apps do not share a route!"))
		})

		It("Should be just fine if apps share a route", func() {
			scaleoverCmdPlugin.App2.Routes = append(scaleoverCmdPlugin.App2.Routes, "a.b.c")
			Expect(scaleoverCmdPlugin.ErrorIfNoSharedRoute()).To(BeNil())
		})

		It("Should ignore route sanity if --no-route-check is at the end of args", func() {
			enforceRoutes := scaleoverCmdPlugin.ShouldEnforceRoutes()
			Expect(enforceRoutes).To(BeFalse())
		})

		It("Should carfuly consider routes if --no-route-check is not in the args", func() {
			scaleoverCmdPlugin.Args.RouteCheck = true
			enforceRoutes := scaleoverCmdPlugin.ShouldEnforceRoutes()
			Expect(enforceRoutes).To(BeTrue())
		})

	})
})
