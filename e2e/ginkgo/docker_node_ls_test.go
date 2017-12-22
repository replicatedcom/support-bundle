package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.node-ls", func() {

	var _ = Describe("swarm", func() {

		BeforeEach(EnterNewTempDir)
		AfterEach(LogResultsFomBundle)
		AfterEach(CleanupDir)

		Context("When the spec is run", func() {

			It("should output the correct files in the bundle", func() {

				WriteBundleConfig(`
specs:
  - docker.node-ls: {}
    output_dir: /docker/node-ls/`)

				GenerateBundle()

				var contents string
				var m interface{}
				var err error

				_ = GetResultFromBundle("docker/node-ls/node_ls.raw")
				contents = GetFileFromBundle("docker/node-ls/node_ls.raw")
				Expect(contents).To(ContainSubstring("Availability"))
				_ = GetResultFromBundle("docker/node-ls/node_ls.json")
				contents = GetFileFromBundle("docker/node-ls/node_ls.json")
				Expect(contents).To(ContainSubstring("Availability"))
				err = json.Unmarshal([]byte(contents), &m)
				Expect(err).NotTo(HaveOccurred())
				_ = GetResultFromBundle("docker/node-ls/node_ls.human")
				contents = GetFileFromBundle("docker/node-ls/node_ls.human")
				Expect(contents).To(ContainSubstring("Availability"))
			})
		})
	})
})
