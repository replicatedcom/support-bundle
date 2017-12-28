package journald

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	jww "github.com/spf13/jwalterweatherman"
)

func TestJournald(t *testing.T) {
	jww.SetLogOutput(GinkgoWriter)
	jww.SetLogThreshold(jww.LevelTrace)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Journald")
}
