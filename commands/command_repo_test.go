package commands_test

import (
	. "github.com/comcast/cf-zdd-plugin/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeRunnable struct{ CommandRunnable }

var _ = Describe("commandRunnable", func() {
	Describe(".Register", func() {
		Context("When called with a valid CommandRunnable", func() {
			It("should register and set a CommandRunnable", func() {
				Register("canary", new(fakeRunnable))
				_, ok := GetRegistry()["canary"]
				Expect(ok).Should(BeTrue())
			})
		})
	})
	Describe(".GetRegistry", func() {
		Context("when containing a registered CommandRunnable", func() {
			It("should return a registry", func() {
				Register("canary", new(fakeRunnable))
				_, ok := GetRegistry()["canary"]
				Expect(ok).Should(BeTrue())
			})
		})
	})
})
