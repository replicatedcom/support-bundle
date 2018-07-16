package docker

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("docker.version", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.version: {}
    output_dir: /docker/version/`)

			GenerateBundle()

			var contents string
			_ = GetResultFromBundle("docker/version/docker_version.json")
			contents = GetFileFromBundle("docker/version/docker_version.json")
			Expect(contents).To(ContainSubstring("ApiVersion"))
			var m interface{}
			err := json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
