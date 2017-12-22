package ginkgo

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.task-ls", func() {

	var _ = Describe("swarm", func() {

		// FIXME: deploy test-stack

		BeforeEach(EnterNewTempDir)
		AfterEach(LogResultsFomBundle)
		AfterEach(CleanupDir)

		Context("When the spec is run", func() {

			It("should output the correct files in the bundle", func() {

				WriteBundleConfig(`
specs:
  - docker.task-ls: {}
    output_dir: /docker/task-ls/
  - docker.service-ps: {}
    output_dir: /docker/service-ps/
  - docker.stack-service-ps:
      namespace: test-stack
    output_dir: /docker/stack-service-ps/`)

				GenerateBundle()

				var contents string
				var m interface{}
				var err error

				_ = GetResultFromBundle("docker/task-ls/task_ls.raw")
				contents = GetFileFromBundle("docker/task-ls/task_ls.raw")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				_ = GetResultFromBundle("docker/task-ls/task_ls.json")
				contents = GetFileFromBundle("docker/task-ls/task_ls.json")
				Expect(contents).To(ContainSubstring("ContainerSpec"))
				err = json.Unmarshal([]byte(contents), &m)
				Expect(err).NotTo(HaveOccurred())
				_ = GetResultFromBundle("docker/task-ls/task_ls.human")
				contents = GetFileFromBundle("docker/task-ls/task_ls.human")
				Expect(contents).To(ContainSubstring("ContainerSpec"))

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
