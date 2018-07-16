package docker

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("docker.container-ls", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.container-ls:
      all: true
    output_dir: /docker/container-ls/
  - docker.ps:
      all: true
    output_dir: /docker/ps/`)

			GenerateBundle()

			var contents string
			var m interface{}
			var err error

			_ = GetResultFromBundle("docker/container-ls/container_ls.json")
			contents = GetFileFromBundle("docker/container-ls/container_ls.json")
			Expect(contents).To(ContainSubstring("Command"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())

			_ = GetResultFromBundle("docker/ps/container_ls.json")
			contents = GetFileFromBundle("docker/ps/container_ls.json")
			Expect(contents).To(ContainSubstring("Command"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
