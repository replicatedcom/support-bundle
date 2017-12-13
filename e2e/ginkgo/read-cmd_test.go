package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Given file paths", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("And the user runs the following commands within desired locations", func() {

		WriteBundleConfig(`
specs:
  - builtin: core.read-command
    raw: /daemon/commands/date
    config:
      command: "date"`)

		GenerateBundle()

		contents := GetFileFromBundle("daemon/commands/date")
		Expect(contents).ToNot(BeEmpty())
	})
})
