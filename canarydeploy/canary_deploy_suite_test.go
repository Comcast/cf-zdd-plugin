package canarydeploy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCanaryDeploy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CanaryDeploy Suite")
}
