package docker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("docker.service-logs swarm", func() {

	// FIXME: deploy test-stack

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.service-logs:
      service: uxa1uzb3wd6khtu3qv4nbqpbe
    output_dir: /docker/service-logs-by-id/
  - docker.service-logs:
      service: test-stack_visualizer
    output_dir: /docker/service-logs-by-name/
  - docker.service-logs:
      service_list_options:
        filters:
          label:
            - com.docker.stack.namespace=test-stack
    output_dir: /docker/service-logs-by-labels/
  - docker.stack-service-logs:
      namespace: test-stack
    output_dir: /docker/stack-service-logs/`)

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle("docker/service-logs-by-id/uxa1uzb3wd6khtu3qv4nbqpbe")
			contents = GetFileFromBundle("docker/service-logs-by-id/uxa1uzb3wd6khtu3qv4nbqpbe.stdout")
			Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

			_ = GetResultFromBundle("docker/service-logs-by-name/test-stack_visualizer")
			contents = GetFileFromBundle("docker/service-logs-by-name/test-stack_visualizer.stdout")
			Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

			_ = GetResultFromBundle("docker/service-logs-by-labels/test-stack_visualizer")
			contents = GetFileFromBundle("docker/service-logs-by-labels/test-stack_visualizer.stdout")
			Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

			_ = GetResultFromBundle("docker/stack-service-logs/test-stack_visualizer")
			contents = GetFileFromBundle("docker/stack-service-logs/test-stack_visualizer.stdout")
			Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))
		})
	})
})
