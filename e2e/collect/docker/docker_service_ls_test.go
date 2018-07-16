package docker

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("docker.service-ls swarm", func() {

	// FIXME: deploy test-stack

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.service-ls: {}
    output_dir: /docker/service-ls/
  - docker.stack-service-ls:
      namespace: test-stack
    output_dir: /docker/stack-service-ls/`)

			GenerateBundle()

			var contents string
			var m interface{}
			var err error

			_ = GetResultFromBundle("docker/service-ls/service_ls.json")
			contents = GetFileFromBundle("docker/service-ls/service_ls.json")
			Expect(contents).To(ContainSubstring("TaskTemplate"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())

			_ = GetResultFromBundle("docker/stack-service-ls/service_ls.json")
			contents = GetFileFromBundle("docker/stack-service-ls/service_ls.json")
			Expect(contents).To(ContainSubstring("TaskTemplate"))
			err = json.Unmarshal([]byte(contents), &m)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
