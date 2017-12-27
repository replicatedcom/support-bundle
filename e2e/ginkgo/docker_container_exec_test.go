package ginkgo

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("docker.container-exec", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
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
    output_dir: /docker/exec/`, containerID, containerName))

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
		})
	})
})
