package docker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
	"github.com/satori/go.uuid"
)

var _ = Describe("docker.container-inspect", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	now := strconv.FormatInt(time.Now().Unix(), 10)
	container1Name, container2Name := "container1-name-"+now, "container2-name-"+now
	labels := map[string]string{
		"foo": "bar",
	}
	cmd := []string{"echo", "Hello World!"}
	var container1ID, container2ID, container3ID string
	BeforeEach(func() {
		container1ID = MakeDockerContainer(dockerClient, container1Name, labels, cmd)
		container2ID = MakeDockerContainer(dockerClient, container2Name, labels, cmd)
		container3ID = MakeDockerContainer(dockerClient, uuid.NewV4().String(), nil, cmd)
	})
	AfterEach(func() {
		RemoveDockerContainer(dockerClient, container1ID)
		RemoveDockerContainer(dockerClient, container2ID)
		RemoveDockerContainer(dockerClient, container3ID)
	})

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(fmt.Sprintf(`
specs:
  - docker.container-inspect:
      container: %s
    output_dir: /docker/container-inspect-by-id/
  - docker.container-inspect:
      container: %s
    output_dir: /docker/container-inspect-by-name/
  - docker.container-inspect:
      container_list_options:
        all: true
        filters:
          label:
            - foo=bar
    output_dir: /docker/container-inspect-by-labels/`,
				container1ID, container2Name))

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-id/%s.raw", container1ID))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-id/%s.raw", container1ID))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-name/%s.raw", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-name/%s.raw", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.raw", container1Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.raw", container1Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.raw", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.raw", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
		})
	})
})
