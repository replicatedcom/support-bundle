package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("os.hostname", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.hostname: 
      output_dir: /os/hostname/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/hostname/stdout")
			contents := GetFileFromBundle("os/hostname/stdout")
			Expect(contents).NotTo(BeEmpty())

			ExpectFileNotInBundle("os/hostname/stderr")
		})
	})
})
