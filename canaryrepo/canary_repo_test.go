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

package canaryrepo_test

import (
	. "github.com/comcast/cf-zdd-plugin/canaryrepo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeRunnable struct{ PluginRunnable }

var _ = Describe("canaryrepo", func() {
	Describe(".Register", func() {
		Context("When called with a valid canaryRunnable", func() {
			It("should register and set a canaryRunnable", func() {
				Register("canary", new(fakeRunnable))
				_, ok := GetRegistry()["canary"]
				Ω(ok).Should(BeTrue())
			})
		})
	})
	Describe(".GetRegistry", func() {
		Context("when containing a registered canary runnable", func() {
			It("should return a registry", func() {
				Register("canary", new(fakeRunnable))
				_, ok := GetRegistry()["canary"]
				Ω(ok).Should(BeTrue())
			})
		})
	})
})
