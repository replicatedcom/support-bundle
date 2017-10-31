package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.daemon", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Finds DriverStatus in docker_info.json", func() {

		WriteFile("config.yml", `
specs:
  - builtin: docker.daemon
    json: /daemon/docker/

      `)

		GenerateBundle()

		contents := GetFileFromBundle("/daemon/docker/docker_info.json")

		Expect(contents).To(ContainSubstring("DriverStatus"))
	})

	Describe("Container tests", func() {
		var containerID string
		BeforeEach(func() {
			containerID = MakeDockerContainer()
		})
		AfterEach(func() {
			RemoveDockerContainer(containerID)
		})

		It("Copies a file from the docker container", func() {
			WriteFile("config.yml", `
specs:
  - builtin: docker.read-file
    raw: /daemon/docker/readfile
    config:
      file_path: "/usr/lib/os-release"
      container_id: `+containerID)

			GenerateBundle()

			contents := GetFileFromBundle("/daemon/docker/readfile")

			Expect(contents).To(ContainSubstring("ubuntu"))
		})

		It("Runs a command on the docker container", func() {
			WriteFile("config.yml", `
specs:
  - builtin: docker.exec-command
    raw: /daemon/docker/command-succeed.
    config:
      command: "echo"
      args: ["testingEchoCommand"]
      container_id: `+containerID)

			GenerateBundle()

			contents := GetFileFromBundle("/daemon/docker/command-succeed.stdout")

			Expect(contents).To(ContainSubstring("testingEchoCommand"))
		})

		It("Runs a command on the docker container that generates output on stderr", func() {
			WriteFile("config.yml", `
specs:
  - builtin: docker.exec-command
    raw: /daemon/docker/command-fail.
    config:
      command: "cat"
      args: ["fileThatDoesNotExist"]
      container_id: `+containerID)

			GenerateBundle()

			contents := GetFileFromBundle("/daemon/docker/command-fail.stderr")

			Expect(contents).To(ContainSubstring("fileThatDoesNotExist"))
		})
	})

})
