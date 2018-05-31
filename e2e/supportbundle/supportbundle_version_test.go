package supportbundle

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("supportbundle.version", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - version: {}
    output_dir: /`)

			GenerateBundle()
			_ = GetResultFromBundle("VERSION.json")
			contents := GetFileFromBundle("VERSION.json")
			Expect(contents).To(Equal(`{
  "Version": "",
  "GitSHA": "",
  "BuildTime": "0001-01-01T00:00:00Z"
}
`))
		})
	})
})
