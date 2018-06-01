package docker

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("docker.image-ls", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.image-ls:
      all: true
    output_dir: /docker/image-ls/
  - docker.images:
      all: true
    output_dir: /docker/images/`)

			GenerateBundle()

			var contents string
			var m interface{}
			var err error
			_ = GetResultFromBundle("docker/image-ls/image_ls.json")
			contents = GetFileFromBundle("docker/image-ls/image_ls.json")
			Expect(contents).To(ContainSubstring("RepoTags"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())

			_ = GetResultFromBundle("docker/images/image_ls.json")
			contents = GetFileFromBundle("docker/images/image_ls.json")
			Expect(contents).To(ContainSubstring("RepoTags"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
