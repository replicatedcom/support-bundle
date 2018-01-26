package ginkgo

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.daemon", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("Finds DriverStatus in docker_info.json", func() {

		WriteBundleConfig(`
specs:
  - builtin: docker.daemon
    json: /daemon/docker/

      `)

		GenerateBundle()

		index := GetFileFromBundle("index.json")
		Expect(index).To(ContainSubstring("docker_info.json"))
		Expect(index).To(ContainSubstring("docker_ps_all.json"))
		Expect(index).To(ContainSubstring("docker_images_all.json"))

		infoContents := GetFileFromBundle("daemon/docker/docker_info.json")
		Expect(infoContents).To(ContainSubstring("DriverStatus"))

		psContents := GetFileFromBundle("daemon/docker/docker_ps_all.json")
		Expect(psContents).ToNot(Equal(""))

		imagesContents := GetFileFromBundle("daemon/docker/docker_images_all.json")
		Expect(imagesContents).ToNot(Equal(""))
	})

	Describe("Container tests", func() {
		var containerID string
		BeforeEach(func() {
			containerID = MakeDockerContainer("", nil)
		})
		AfterEach(func() {
			RemoveDockerContainer(containerID)
		})

		Describe("docker container-ls-logs", func() {
			name := fmt.Sprintf("labeled-logs-container-%d", time.Now().UnixNano())
			name2 := fmt.Sprintf("labeled-inspect-container-2-%d", time.Now().UnixNano())
			labels := map[string]string{
				"foo": "bar",
			}
			var labeledContainerID, labeledContainer2ID string
			BeforeEach(func() {
				labeledContainerID = MakeDockerContainer(name, labels)
				labeledContainer2ID = MakeDockerContainer(name2, nil)
			})
			AfterEach(func() {
				RemoveDockerContainer(labeledContainerID)
				RemoveDockerContainer(labeledContainer2ID)
			})

			It("Gets logs from docker containers with matching labels", func() {
				WriteBundleConfig(`
specs:
- builtin: docker.container-ls-logs
  raw: /containers/foo
  docker.container-logs:
    container_list_options:
      filters:
        label:
          - foo=bar`)
				GenerateBundle()

				path := fmt.Sprintf("containers/foo/%s.log.stdout", name)
				_ = GetFileFromBundle(path)
				path2 := fmt.Sprintf("containers/foo/%s.log.stdout", name2)
				ExpectFileNotInBundle(path2)
			})
		})

		Describe("docker container-ls-inspect", func() {
			name := fmt.Sprintf("labeled-inspect-container-%d", time.Now().UnixNano())
			name2 := fmt.Sprintf("labeled-inspect-container-2-%d", time.Now().UnixNano())
			labels := map[string]string{
				"foo": "bar",
			}
			var labeledContainerID, labeledContainer2ID string
			BeforeEach(func() {
				labeledContainerID = MakeDockerContainer(name, labels)
				labeledContainer2ID = MakeDockerContainer(name2, nil)
			})
			AfterEach(func() {
				RemoveDockerContainer(labeledContainerID)
				RemoveDockerContainer(labeledContainer2ID)
			})

			It("Inspects docker containers with matching labels", func() {
				WriteBundleConfig(`
specs:
- builtin: docker.container-ls-inspect
  raw: /containers/foo
  docker.container-inspect:
    container_list_options:
      filters:
        label:
          - foo=bar`)
				GenerateBundle()

				path := fmt.Sprintf("containers/foo/%s.json", name)
				_ = GetFileFromBundle(path)
				path2 := fmt.Sprintf("containers/foo/%s.json", name2)
				ExpectFileNotInBundle(path2)
			})
		})

		It("Copies a file from the docker container", func() {
			WriteBundleConfig(`
specs:
  - builtin: docker.read-file
    raw: /daemon/docker/readfile
    config:
      file_path: "/usr/lib/os-release"
      container_id: ` + containerID)

			GenerateBundle()

			contents := GetFileFromBundle("daemon/docker/readfile/os-release")

			Expect(contents).To(ContainSubstring("ubuntu"))
		})

		It("Runs a command on the docker container", func() {
			WriteBundleConfig(`
specs:
  - builtin: docker.exec-command
    raw: /daemon/docker/command-succeed
    config:
      command: "echo"
      args: ["testingEchoCommand"]
      container_id: ` + containerID)

			GenerateBundle()

			contents := GetFileFromBundle("daemon/docker/command-succeed.stdout")

			Expect(contents).To(ContainSubstring("testingEchoCommand"))
		})

		It("Runs a command on the docker container that generates output on stderr", func() {
			WriteBundleConfig(`
specs:
  - builtin: docker.exec-command
    raw: /daemon/docker/command-fail
    config:
      command: "cat"
      args: ["fileThatDoesNotExist"]
      container_id: ` + containerID)

			GenerateBundle()

			contents := GetFileFromBundle("daemon/docker/command-fail.stderr")

			Expect(contents).To(ContainSubstring("fileThatDoesNotExist"))
		})
	})

})
