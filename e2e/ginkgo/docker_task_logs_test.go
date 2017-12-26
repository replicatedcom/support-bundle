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
      id: x1xzye79vohgp2a9dly0iigqc
    output_dir: /docker/task-logs-by-id/
  - docker.task-logs:
      task_list_options:
        filters:
          label:
            - com.docker.stack.namespace=test-stack
    output_dir: /docker/task-logs-by-labels/
  - docker.stack-task-logs:
      namespace: test-stack
    output_dir: /docker/stack-task-logs/`)
				GenerateBundle()

				var contents string

				_ = GetResultFromBundle("docker/task-logs-by-id/x1xzye79vohgp2a9dly0iigqc.raw")
				contents = GetFileFromBundle("docker/task-logs-by-id/x1xzye79vohgp2a9dly0iigqc.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

				_ = GetResultFromBundle("docker/task-logs-by-labels/x1xzye79vohgp2a9dly0iigqc.raw")
				contents = GetFileFromBundle("docker/task-logs-by-labels/x1xzye79vohgp2a9dly0iigqc.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))

				_ = GetResultFromBundle("docker/stack-task-logs/x1xzye79vohgp2a9dly0iigqc.raw")
				contents = GetFileFromBundle("docker/stack-task-logs/x1xzye79vohgp2a9dly0iigqc.raw")
				Expect(contents).To(ContainSubstring("npm info it worked if it ends with ok"))
			})
		})
	})
})
