package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("docker.container-exec", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	var containerID, containerName string
	containerName = uuid.NewV4().String()
	BeforeEach(func() {
		containerID = MakeDockerContainer(dockerClient, containerName, nil, nil)
	})
	AfterEach(func() {
		RemoveDockerContainer(dockerClient, containerID)
	})

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(fmt.Sprintf(`
specs:
  - docker.container-exec:
      container: %s
      exec_config:
        Cmd: ["echo", "Hello World!"]
    output_dir: /docker/container-exec/
  - docker.exec:
      container: %s
      exec_config:
        Cmd: ["echo", "foo bar"]
    output_dir: /docker/exec/
  - docker.container-exec:
      container: %s
      exec_config:
        Cmd: ["foobar", "bah"]
    output_dir: /docker/container-exec-notexist/`, containerID, containerName, containerID))
			GenerateBundle()

			var contents string

			_ = GetResultFromBundle("docker/container-exec/stdout.raw")
			_ = GetResultFromBundle("docker/container-exec/stderr.raw")
			contents = GetFileFromBundle("docker/container-exec/stdout.raw")
			Expect(contents).To(Equal("Hello World!\n"))

			_ = GetResultFromBundle("docker/exec/stdout.raw")
			_ = GetResultFromBundle("docker/exec/stderr.raw")
			contents = GetFileFromBundle("docker/exec/stdout.raw")
			Expect(contents).To(Equal("foo bar\n"))

			_ = GetResultFromBundle("docker/container-exec-notexist/stdout.raw")
			_ = GetResultFromBundle("docker/container-exec-notexist/stderr.raw")
			contents = GetFileFromBundle("docker/container-exec-notexist/stdout.raw") // FIXME: stdout and no error!
			Expect(contents).To(ContainSubstring("executable file not found in"))
		})
	})
})
