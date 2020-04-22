package meta

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("meta.files_redaction", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When a meta.redact.files spec is included", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo $HI]
      env: [HI=hello!]
      output_dir: /os/run-command/echohi/
  - os.run-command:
      name: sh
      args: [-c, echo $BYE]
      env: [BYE=goodbye!]
      output_dir: /os/run-command/echobye/
  - meta.redact:
      output_dir: /redact/
      files: ["**/echohi/*"]
`)

			GenerateBundle()

			hiResult := GetResultFromBundle("os/run-command/echohi/stdout")
			Expect(hiResult.Redacted).To(Equal(true))
			ExpectFileNotInBundle("os/run-command/echohi/stdout")

			byeResult := GetResultFromBundle("os/run-command/echobye/stdout")
			Expect(byeResult.Redacted).To(Equal(false))
			contents := GetFileFromBundle("os/run-command/echobye/stdout")
			Expect(contents).To(Equal("goodbye!\n"))

			_ = GetResultFromBundle("redact/file_redactions.json")
			redactions := GetFileFromBundle("redact/file_redactions.json")
			Expect(redactions).To(Equal(`[
  "**/echohi/*"
]
`))
		})
	})
})
