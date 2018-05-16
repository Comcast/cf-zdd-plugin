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

package main_test

import (
	"code.cloudfoundry.org/cli/plugin"
	. "github.com/comcast/cf-zdd-plugin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CfZddPlugin", func() {

	Describe(".GetMetaData", func() {
		var pluginMetadata plugin.PluginMetadata
		Context("when calling GetMetaData", func() {
			BeforeEach(func() {
				pluginMetadata = new(CfZddPlugin).GetMetadata()
			})
			It("should return the correct help text", func() {
				success := false
				for _, v := range pluginMetadata.Commands {
					if v.HelpText == CanaryDeployHelpText {
						success = true
					}
				}
				Ω(success).Should(BeTrue())
			})
			It("should return the correct cmd name", func() {
				success := false
				for _, v := range pluginMetadata.Commands {
					if v.Name == CanaryDeployCmdName {
						success = true
					}
				}
				Ω(success).Should(BeTrue())
			})
		})
	})
})
