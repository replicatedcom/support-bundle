package supportbundle

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("supportbundle.version", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - version: {}
    output_dir: /`)

			GenerateBundle()

			_ = GetResultFromBundle("VERSION.raw")
			contents := GetFileFromBundle("VERSION.raw")
			Expect(contents).To(Equal(`{"Version":"","GitSHA":"","BuildTime":"0001-01-01T00:00:00Z"}`))
			_ = GetResultFromBundle("VERSION.json")
			contents = GetFileFromBundle("VERSION.json")
			Expect(contents).To(Equal(`{
  "Version": "",
  "GitSHA": "",
  "BuildTime": "0001-01-01T00:00:00Z"
}
`))
			_ = GetResultFromBundle("VERSION.human")
			contents = GetFileFromBundle("VERSION.human")
			Expect(contents).To(Equal(`BuildTime: 0001-01-01T00:00:00Z
GitSHA: ""
Version: ""
`))
		})
	})
})
