package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.service-logs", func() {

	var _ = Describe("swarm", func() {

		// FIXME: deploy test-stack

		BeforeEach(EnterNewTempDir)
		AfterEach(LogResultsFomBundle)
		AfterEach(CleanupDir)

		Context("When the spec is run", func() {

			It("should output the correct files in the bundle", func() {

				WriteBundleConfig(`
specs:
  - docker.service-logs:
      id: 02cxslno3h9liblktnr0o6li7
    output_dir: /docker/service-logs-by-id/
  - docker.service-logs:
      name: test-stack_visualizer
    output_dir: /docker/service-logs-by-name/
  - docker.service-logs:
      service_list_options:
        all: true
        filters:
          label:
            - com.docker.stack.namespace=test-stack
    output_dir: /docker/service-logs-by-labels/`)

				GenerateBundle()

				var contents string

				_ = GetResultFromBundle("docker/service-logs-by-id/02cxslno3h9liblktnr0o6li7.raw")
				contents = GetFileFromBundle("docker/service-logs-by-id/02cxslno3h9liblktnr0o6li7.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

				_ = GetResultFromBundle("docker/service-logs-by-name/test-stack_visualizer.raw")
				contents = GetFileFromBundle("docker/service-logs-by-name/test-stack_visualizer.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

				_ = GetResultFromBundle("docker/service-logs-by-labels/test-stack_visualizer.raw")
				contents = GetFileFromBundle("docker/service-logs-by-labels/test-stack_visualizer.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))
			})
		})
	})
})
