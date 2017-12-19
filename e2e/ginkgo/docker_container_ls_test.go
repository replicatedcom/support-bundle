package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.container-ls", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.container-ls:
      all: true
    output_dir: /docker/container-ls/`)

			GenerateBundle()

			var contents string
			_ = GetResultFromBundle("docker/container-ls/container_ls.raw")
			contents = GetFileFromBundle("docker/container-ls/container_ls.raw")
			Expect(contents).To(ContainSubstring("Command"))
			_ = GetResultFromBundle("docker/container-ls/container_ls.json")
			contents = GetFileFromBundle("docker/container-ls/container_ls.json")
			Expect(contents).To(ContainSubstring("Command"))
			var m interface{}
			err := json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
			_ = GetResultFromBundle("docker/container-ls/container_ls.human")
			contents = GetFileFromBundle("docker/container-ls/container_ls.human")
			Expect(contents).To(ContainSubstring("Command"))
		})
	})
})
