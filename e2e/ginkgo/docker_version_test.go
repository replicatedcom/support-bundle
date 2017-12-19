package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.version", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.version: {}
    output_dir: /docker/version/`)

			GenerateBundle()

			var contents string
			_ = GetResultFromBundle("docker/version/docker_version.raw")
			contents = GetFileFromBundle("docker/version/docker_version.raw")
			Expect(contents).NotTo(BeEmpty())
			_ = GetResultFromBundle("docker/version/docker_version.json")
			contents = GetFileFromBundle("docker/version/docker_version.json")
			Expect(contents).NotTo(BeEmpty())
			var m interface{}
			err := json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
			_ = GetResultFromBundle("docker/version/docker_version.human")
			contents = GetFileFromBundle("docker/version/docker_version.human")
			Expect(contents).NotTo(BeEmpty())
		})
	})
})
