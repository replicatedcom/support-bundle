package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.task-logs", func() {

	var _ = Describe("swarm", func() {

		// FIXME: deploy test-stack

		BeforeEach(EnterNewTempDir)
		AfterEach(LogResultsFomBundle)
		AfterEach(CleanupDir)

		Context("When the spec is run", func() {

			It("should output the correct files in the bundle", func() {

				WriteBundleConfig(`
specs:
  - docker.task-logs:
      id: manrfdu40g5sy6vd9ygpjx8tw
    output_dir: /docker/task-logs-by-id/
  - docker.task-logs:
      task_list_options:
        all: true
        filters:
          label:
            - com.docker.stack.namespace=test-stack
    output_dir: /docker/task-logs-by-labels/`)

				GenerateBundle()

				var contents string

				_ = GetResultFromBundle("docker/task-logs-by-id/manrfdu40g5sy6vd9ygpjx8tw.raw")
				contents = GetFileFromBundle("docker/task-logs-by-id/manrfdu40g5sy6vd9ygpjx8tw.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

				_ = GetResultFromBundle("docker/task-logs-by-labels/manrfdu40g5sy6vd9ygpjx8tw.raw")
				contents = GetFileFromBundle("docker/task-logs-by-labels/manrfdu40g5sy6vd9ygpjx8tw.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))
			})
		})
	})
})
