package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.service-ps", func() {

	var _ = Describe("swarm", func() {

		// FIXME: deploy test-stack

		BeforeEach(EnterNewTempDir)
		AfterEach(LogResultsFomBundle)
		AfterEach(CleanupDir)

		Context("When the spec is run", func() {

			It("should output the correct files in the bundle", func() {

				WriteBundleConfig(`
specs:
  - docker.service-ps: {}
    output_dir: /docker/service-ps/
  - docker.stack-service-ps:
      namespace: test-stack
    output_dir: /docker/stack-service-ps/`)

				GenerateBundle()

				var contents string
				var m interface{}
				var err error

				_ = GetResultFromBundle("docker/service-ps/service_ps.raw")
				contents = GetFileFromBundle("docker/service-ps/service_ps.raw")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				_ = GetResultFromBundle("docker/service-ps/service_ps.json")
				contents = GetFileFromBundle("docker/service-ps/service_ps.json")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				err = json.Unmarshal([]byte(contents), &m)
				Expect(err).NotTo(HaveOccurred())
				_ = GetResultFromBundle("docker/service-ps/service_ps.human")
				contents = GetFileFromBundle("docker/service-ps/service_ps.human")
				Expect(contents).To(ContainSubstring("ContainerSpec"))

				_ = GetResultFromBundle("docker/stack-service-ps/service_ps.raw")
				contents = GetFileFromBundle("docker/stack-service-ps/service_ps.raw")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				_ = GetResultFromBundle("docker/stack-service-ps/service_ps.json")
				contents = GetFileFromBundle("docker/stack-service-ps/service_ps.json")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				err = json.Unmarshal([]byte(contents), &m)
				Expect(err).NotTo(HaveOccurred())
				_ = GetResultFromBundle("docker/stack-service-ps/service_ps.human")
				contents = GetFileFromBundle("docker/stack-service-ps/service_ps.human")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
			})
		})
	})
})
