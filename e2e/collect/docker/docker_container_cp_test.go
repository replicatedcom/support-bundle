package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	uuid "github.com/satori/go.uuid"
)

var _ = Describe("docker.container-cp", func() {
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
  - docker.container-cp:
      container: %s
      src_path: /etc/default/halt
    output_dir: /docker/container-cp-file/
  - docker.container-cp:
      container: %s
      src_path: /etc/default/
    output_dir: /docker/container-cp-dir/`, containerID, containerName))

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

var _ = Describe("docker.container-cp-by-label", func() {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	var containerID1, containerID2, containerName string
	containerName = uuid.NewV4().String()
	BeforeEach(func() {
		containerID1 = MakeDockerContainer(dockerClient, containerName+"1", map[string]string{"container-cp-test": "", "container-cp": "value", "bothLabel": ""}, nil)
		containerID2 = MakeDockerContainer(dockerClient, containerName+"2", map[string]string{"bothLabel": ""}, nil)
	})
	AfterEach(func() {
		RemoveDockerContainer(dockerClient, containerID1)
		RemoveDockerContainer(dockerClient, containerID2)
	})

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.container-cp:
      labels:
        - container-cp-test
      src_path: /etc/default/halt
    output_dir: /docker/container-cp-file/
  - docker.container-cp:
      labels:
        - "container-cp=value"
      src_path: /etc/default/
    output_dir: /docker/container-cp-dir/
  - docker.container-cp:
      labels:
        - bothLabel
        - container-cp-test
      src_path: /etc/default/
    output_dir: /docker/container-cp-multilabel/
  - docker.container-cp:
      labels:
        - bothLabel
        - container-cp-test
        - notexist
      src_path: /etc/default/
    output_dir: /docker/container-cp-multilabel-fail/
  - docker.container-cp:
      labels:
        - bothLabel
      src_path: /etc/default/
    output_dir: /docker/container-cp-multiple/`)

			GenerateBundle()

			var contents string

			//copy a file from a simple labeled container
			_ = GetResultFromBundle("docker/container-cp-file/halt")
			contents = GetFileFromBundle("docker/container-cp-file/halt")
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))

			//copy a directory from a key=val labeled container
			_ = GetResultFromBundle("docker/container-cp-dir/default/halt")
			contents = GetFileFromBundle("docker/container-cp-dir/default/halt")
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))

			//copy a directory from a container specified by multiple labels
			_ = GetResultFromBundle("docker/container-cp-multilabel/default/halt")
			contents = GetFileFromBundle("docker/container-cp-multilabel/default/halt")
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))

			//when there isn't a container that matches *all* of the labels, the file should not be created
			ExpectFileNotInBundle("docker/container-cp-multilabel-fail/default/halt")

			//copy a directory from multiple containers at once
			container1path := fmt.Sprintf("docker/container-cp-multiple/%s/default/halt", containerID1)
			container2path := fmt.Sprintf("docker/container-cp-multiple/%s/default/halt", containerID2)
			_ = GetResultFromBundle(container1path)
			contents = GetFileFromBundle(container1path)
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))
			_ = GetResultFromBundle(container2path)
			contents = GetFileFromBundle(container2path)
			Expect(contents).To(Equal(`# Default behaviour of shutdown -h / halt. Set to "halt" or "poweroff".
HALT=poweroff
`))
		})
	})
})
