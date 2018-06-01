package docker

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("docker.info", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.info: {}
    output_dir: /docker/info/`)

			GenerateBundle()

			var contents string
			_ = GetResultFromBundle("docker/info/docker_info.json")
			contents = GetFileFromBundle("docker/info/docker_info.json")
			Expect(contents).To(ContainSubstring("DriverStatus"))
			var m interface{}
			err := json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
