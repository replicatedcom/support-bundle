package supportbundle

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("supportbundle.logs", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - version: {}
    output_dir: /
  - logs:
      defer: true
    output_dir: /`)

			GenerateBundle()
			_ = GetResultFromBundle("logs")
			contents := GetFileFromBundle("logs")
			Expect(contents).To(MatchRegexp("Task with spec.+\"version\":{}"))
		})
	})
})
