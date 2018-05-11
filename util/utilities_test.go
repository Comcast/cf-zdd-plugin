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

package util_test

import (
	"errors"

	"code.cloudfoundry.org/cli/plugin/pluginfakes"
	. "github.com/comcast/cf-zdd-plugin/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("util", func() {
	Describe(".GetDefaultDomain", func() {
		var (
			fakeConnection *pluginfakes.FakeCliConnection
			utility        *Utility
		)

		BeforeEach(func() {
			fakeConnection = &pluginfakes.FakeCliConnection{}
			utility = &Utility{}
		})
		Context("when called with a valid cli connection", func() {
			var returnOutput = []string{"{\n   \"name\": \"vcap\",\n   \"build\": \"2222\",\n   \"support\": \"https://support.pivotal.io\",\n   \"version\": 2,\n   \"description\": \"Cloud Foundry sponsored by Pivotal\",\n   \"authorization_endpoint\": \"https://login.cloud.net\",\n   \"token_endpoint\": \"https://uaa.cloud.net\",\n   \"min_cli_version\": \"6.7.0\",\n   \"min_recommended_cli_version\": \"6.11.2\",\n   \"api_version\": \"2.43.0\",\n   \"app_ssh_endpoint\": \"ssh.cloud.net:2222\",\n   \"app_ssh_host_key_fingerprint\": \"9d:5e:3b:40:26\",\n   \"app_ssh_oauth_client\": \"ssh-proxy\",\n   \"routing_endpoint\": \"https://api.cloud.net/routing\",\n   \"logging_endpoint\": \"wss://loggregator.cloud.net:443\",\n   \"doppler_logging_endpoint\": \"wss://doppler.cloud.net:443\",\n   \"user\": \"5c8a94ef-3fb8--bef6\"\n}\n"}
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(returnOutput, nil)
			})
			It("should return the default domain from cf", func() {
				domain := utility.GetDefaultDomain(fakeConnection)
				Expect(domain).Should(Equal("cloud.net"))
			})
		})
		Context("when a valid cli connection returns an error", func() {
			BeforeEach(func() {
				fakeConnection.CliCommandWithoutTerminalOutputReturns(nil, errors.New("cf error"))
			})
			It("should return the default domain from cf", func() {
				domain := utility.GetDefaultDomain(fakeConnection)
				Expect(domain).Should(Equal("unknown"))
			})
		})
	})
})
