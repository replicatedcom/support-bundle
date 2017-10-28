package ginkgo

import (
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
)

var _ = Describe("Given file paths", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("And the user runs the following commands within desired locations", func() {

		WriteFile("config.yml", `
specs:
  - builtin: core.read-command
    raw: /daemon/commands/date
    config:
      command: "date"
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
			"/daemon/commands/date",
		)

		Expect(contents).ToNot(BeEmpty())
	})
})
