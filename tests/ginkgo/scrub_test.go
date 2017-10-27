package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scrubbing secrets from file", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Scrubs any instances of PGPASSWORD=.*", func() {

		WriteFile("pg.env", `
PGDATABASE=mydata
PGPASSWORD=mypass`)

		WriteBundleConfig(`
specs:
  - builtin: core.read-file
    raw: /pg/pg.env
    config:
      file_path: pg.env
      scrub:
        regex: (PGPASSWORD)=(.*)
        replace: $1=REDACTED
      `)

		GenerateBundle()

		contents := GetFileFromBundle("/pg/pg.env")

		Expect(contents).To(ContainSubstring("PGDATABASE=mydata"))
		Expect(contents).NotTo(ContainSubstring("PGPASSWORD=mypass"))
		Expect(contents).To(ContainSubstring("PGPASSWORD=REDACTED"))
	})

})
