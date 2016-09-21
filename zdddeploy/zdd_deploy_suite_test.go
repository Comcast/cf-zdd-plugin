package zdddeploy

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestZddDeploy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ZddDeploy Suite")
}
