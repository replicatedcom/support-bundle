package supportbundle

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("supportbundle.scrub streamsource", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteFile("pg.env", `
PGDATABASE=mydata
PGPASSWORD=mypass`)

			WriteBundleConfig(`
specs:
  - os.read-file:
      filepath: pg.env
    output_dir: /os/read-file/
    scrub:
      regex: (PGPASSWORD)=(.*)
      replace: $1=REDACTED
      `)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/pg.env")
			contents := GetFileFromBundle("os/read-file/pg.env")

			Expect(contents).To(ContainSubstring("PGDATABASE=mydata"))
			Expect(contents).NotTo(ContainSubstring("PGPASSWORD=mypass"))
			Expect(contents).To(ContainSubstring("PGPASSWORD=REDACTED"))

		})
	})
})
