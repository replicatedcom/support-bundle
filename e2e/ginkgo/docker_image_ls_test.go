package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.image-ls", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.image-ls:
      all: true
    output_dir: /docker/image-ls/`)

			GenerateBundle()

			var contents string
			_ = GetResultFromBundle("docker/image-ls/image_ls.raw")
			contents = GetFileFromBundle("docker/image-ls/image_ls.raw")
			Expect(contents).To(ContainSubstring("RepoTags"))
			_ = GetResultFromBundle("docker/image-ls/image_ls.json")
			contents = GetFileFromBundle("docker/image-ls/image_ls.json")
			Expect(contents).To(ContainSubstring("RepoTags"))
			var m interface{}
			err := json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
			_ = GetResultFromBundle("docker/image-ls/image_ls.human")
			contents = GetFileFromBundle("docker/image-ls/image_ls.human")
			Expect(contents).To(ContainSubstring("RepoTags"))
		})
	})
})
