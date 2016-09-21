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
