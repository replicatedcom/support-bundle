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
Hello World!
What's in my Bundle`)

		WriteFile("config.yml", `
specs:
  - builtin: core.read-file
    raw: /daemon/etc/default/replicated
    config:
      file_path: blah.txt
      `)
		/*
		   - builtin: core.read-file
		       raw: /daemon/etc/default/replicated
		       config:
		         file_path: /etc/default/replicated */

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
		Expect(contents).To(ContainSubstring("Hello World!"))
		Expect(contents).To(ContainSubstring("What's in my Bundle"))

	})

})
