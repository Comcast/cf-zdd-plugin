package commands_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestScaleoverPlugin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scaleover Suite")
}
