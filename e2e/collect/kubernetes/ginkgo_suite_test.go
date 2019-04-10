package kubernetes

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

func TestCore(t *testing.T) {
	ginkgo.SetupLogger()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Kubernetes")
}
