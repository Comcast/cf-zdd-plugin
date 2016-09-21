package canarypromote

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCanaryPromote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CanaryPromote Suite")
}
