package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checking contents of the file", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("Validating text in blah.txt", func() {

		WriteFile("blah.txt", `
Hey there!
Let's take a peek into my file!`)

		WriteBundleConfig(`
specs:
  - builtin: core.read-file
    raw: /daemon/etc/default/replicated
    config:
      file_path: blah.txt
      `)

		GenerateBundle()

		contents := GetFileFromBundle("daemon/etc/default/replicated")
		Expect(contents).To(ContainSubstring("Hey there!"))
		Expect(contents).To(ContainSubstring("Let's take a peek into my file!"))

	})

})
