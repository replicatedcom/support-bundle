package ginkgo

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/satori/go.uuid"
)

var _ = Describe("docker.container-logs", func() {

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
		container1ID = MakeDockerContainer(container1Name, labels, cmd)
		container2ID = MakeDockerContainer(container2Name, labels, cmd)
		container3ID = MakeDockerContainer(uuid.NewV4().String(), nil, cmd)
	})
	AfterEach(func() {
		RemoveDockerContainer(container1ID)
		RemoveDockerContainer(container2ID)
		RemoveDockerContainer(container3ID)
	})

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(fmt.Sprintf(`
specs:
  - docker.container-logs:
      id: %s
    output_dir: /docker/container-logs-by-id/
  - docker.container-logs:
      name: %s
    output_dir: /docker/container-logs-by-name/
  - docker.container-logs:
      container_list_options:
        all: true
        filters:
          label:
            - foo=bar
    output_dir: /docker/container-logs-by-labels/`,
				container1ID, container2Name))

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-logs-by-id/%s.raw", container1ID))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-logs-by-id/%s.raw", container1ID))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-logs-by-name/%s.raw", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-logs-by-name/%s.raw", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))

			_ = GetResultFromBundle(fmt.Sprintf("docker/container-logs-by-labels/%s.raw", container1Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-logs-by-labels/%s.raw", container1Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
			_ = GetResultFromBundle(fmt.Sprintf("docker/container-logs-by-labels/%s.raw", container2Name))
			contents = GetFileFromBundle(fmt.Sprintf("docker/container-logs-by-labels/%s.raw", container2Name))
			Expect(contents).To(ContainSubstring("Hello World!"))
		})
	})
})
