package ginkgo

import (
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
)

var _ = Describe("Checking contents of the file", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Validating text in blah.txt", func() {

		WriteFile("blah.txt", `
Hey there!
Let's take a peek into my file!`)

		WriteFile("config.yml", `
specs:
  - builtin: core.read-file
    raw: /daemon/etc/default/replicated
    config:
      file_path: blah.txt
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
			"/daemon/etc/default/replicated",
		)
		Expect(contents).To(ContainSubstring("Hey there!"))
		Expect(contents).To(ContainSubstring("Let's take a peek into my file!"))

	})

})
