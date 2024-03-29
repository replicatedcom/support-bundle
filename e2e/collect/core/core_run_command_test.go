package core

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("os.run-command", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo $HI]
      env: [HI=hello!]
      output_dir: /os/run-command/echo/
  - os.run-command:
      name: blah
      args: [blah, blah]
      output_dir: /os/run-command/notfound/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("hello!\n"))
			ExpectFileNotInBundle("os/run-command/echo/stderr")

			ExpectBundleErrorToHaveOccurred("os/run-command/notfound/stdout", `exec: "blah": executable file not found in \$PATH`)
			ExpectBundleErrorToHaveOccurred("os/run-command/notfound/stderr", `exec: "blah": executable file not found in \$PATH`)
		})
	})
})
