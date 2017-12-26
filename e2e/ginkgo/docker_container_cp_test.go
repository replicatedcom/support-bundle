package ginkgo

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("docker.container-cp", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	var containerID string
	BeforeEach(func() {
		containerID = MakeDockerContainer(dockerClient, uuid.NewV4().String(), nil, nil)
	})
	AfterEach(func() {
		RemoveDockerContainer(dockerClient, containerID)
	})

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(fmt.Sprintf(`
specs:
  - docker.container-cp:
      container: %s
      src_path: /etc/default/halt
    output_dir: /docker/container-cp-file/
  - docker.container-cp:
      container: %s
      src_path: /etc/default/
    output_dir: /docker/container-cp-dir/`, containerID, containerID))

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle("docker/container-cp-file/halt")
			contents = GetFileFromBundle("docker/container-cp-file/halt")
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))

			_ = GetResultFromBundle("docker/container-cp-dir/default/halt")
			contents = GetFileFromBundle("docker/container-cp-dir/default/halt")
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))
		})
	})
})
