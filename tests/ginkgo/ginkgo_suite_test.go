package ginkgo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	jww "github.com/spf13/jwalterweatherman"
)

func TestGinkgo(t *testing.T) {
	jww.SetLogOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Support Bundle")
}
