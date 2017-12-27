package ginkgo

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("os.uptime", func() {

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
  - os.uptime: {}
    output_dir: /os/uptime/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/uptime/contents")
			contents := GetFileFromBundle("os/uptime/contents")
			Expect(contents).NotTo(BeEmpty())
		})
	})
})
