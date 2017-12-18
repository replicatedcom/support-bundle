package ginkgo

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("os.read-file", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteFile("/tmp/blah.txt", `
Hey there!
Let's take a peek into my file!`)
			defer os.RemoveAll("/tmp/blah.txt")

			WriteBundleConfig(`
specs:
  - os.read-file:
      filepath: /tmp/blah.txt
    output_dir: /os/read-file/blah/
  - os.read-file:
      filepath: /tmp/notfound.txt
    output_dir: /os/read-file/notfound/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/blah/contents")
			contents := GetFileFromBundle("os/read-file/blah/contents")
			Expect(contents).To(Equal(`
Hey there!
Let's take a peek into my file!`))

			ExpectBundleErrorToHaveOccured("os/read-file/notfound", "open /tmp/notfound.txt: no such file or directory")
		})
	})
})
