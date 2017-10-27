package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
	"path"
)

var _ = Describe("Scrubbing secrets from file", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Scrubs any instances of PGPASSWORD=.*", func() {

		WriteFile("pg.env", `
PGDATABASE=mydata
PGPASSWORD=mypass`)

		WriteFile("config.yml", `
specs:
  - builtin: core.read-file
    raw: /pg/pg.env
    config:
      file_path: `+ path.Join(tmpdir, "pg.env")+ `
      scrub:
        regex: (PGPASSWORD)=(.*)
        replace: $1=REDACTED
      `)

		err := cmd.Generate(
			path.Join(tmpdir, "config.yml"),
			path.Join(tmpdir, "bundle.tar.gz"),
			true,
			60,
		)

		Expect(err).To(BeNil())

		contents := ReadFileFromBundle(
			path.Join("bundle.tar.gz"),
			"/pg/pg.env",
		)

		Expect(contents).To(ContainSubstring("PGDATABASE=mydata"))
		Expect(contents).NotTo(ContainSubstring("PGPASSWORD=mypass"))
		Expect(contents).To(ContainSubstring("PGPASSWORD=REDACTED"))
	})

})
