package docker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	"github.com/satori/go.uuid"
)

var _ = Describe("docker.container-inspect", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient.NegotiateAPIVersion(context.Background())

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	now := strconv.FormatInt(time.Now().UnixNano(), 20)
	container1Name, container2Name := "c1-name-"+now, "c2-name-"+now
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
      container: %s
      scrub:
        regex: "(?m)(\"(?:ENVSCRUBBED|ENVSCRUBBEDANOTHER)=)([^\"]*)(\",?)"
        replace: "${1}***HIDDEN***${3}"
    output_dir: /docker/container-inspect-scrubbed/
  - docker.container-inspect:
      container_list_options:
        all: true
        filters:
          label:
            - foo=bar
    output_dir: /docker/container-inspect-by-labels/`,
				container1ID, container2Name, container1ID))

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-id/%s.json", container1ID))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-id/%s.json", container1ID))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-name/%s.json", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-name/%s.json", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-scrubbed/%s.json", container1ID))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-scrubbed/%s.json", container1ID))
			Expect(contents).To(ContainSubstring("Hello World!"))
			Expect(contents).To(ContainSubstring("ENVNORMAL=normal"))
			Expect(contents).NotTo(ContainSubstring("ENVSCRUBBED=secret"))
			Expect(contents).NotTo(ContainSubstring("ENVSCRUBBEDANOTHER=anothersecret"))
			Expect(contents).To(ContainSubstring("ENVNORMALTWO=normaltwo"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.json", container1Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.json", container1Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
			_ = GetResultFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.json", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-inspect-by-labels/%s.json", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
		})
	})
})
