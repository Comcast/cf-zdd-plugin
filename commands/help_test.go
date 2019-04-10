package commands_test

import (
	"github.com/comcast/cf-zdd-plugin/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".Help", func() {

	Describe(".init", func() {
		Context("when the package is imported", func() {
			It("should then be registered with the command repo", func() {
				_, ok := commands.GetRegistry()[commands.HelpCommandName]
				Expect(ok).Should(BeTrue())
			})
		})
	})

	Describe("with a valid arg and run method", func() {
		var (
			err      error
			helpCmd  *commands.HelpCmd
			cfZddCmd *commands.CfZddCmd
		)
		BeforeEach(func() {
			cfZddCmd = &commands.CfZddCmd{
				CmdName: commands.HelpCommandName,
			}
			helpCmd = new(commands.HelpCmd)
			helpCmd.SetArgs(cfZddCmd)
		})
		Context("when called without a command", func() {
			It("should return the default help", func() {
				err = helpCmd.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("when called with a valid subcommand", func() {
			BeforeEach(func() {
				cfZddCmd.HelpTopic = commands.BlueGreenCmdName
			})
			It("should return the appropriate help", func() {
				err = helpCmd.Run()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
