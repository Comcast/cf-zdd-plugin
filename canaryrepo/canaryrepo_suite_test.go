package canaryrepo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCanaryRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CanaryRepo Suite")
}
