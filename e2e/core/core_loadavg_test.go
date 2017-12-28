package core

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("os.loadavg", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.loadavg: {}
    output_dir: /os/loadavg/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/loadavg/contents")
			contents := GetFileFromBundle("os/loadavg/contents")
			Expect(contents).NotTo(BeEmpty())
		})
	})
})
