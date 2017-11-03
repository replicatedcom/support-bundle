package ginkgo

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Docker container run command", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogErrorsFomBundle)
	AfterEach(CleanupDir)

	It("Successfully executes the docker container run command", func() {

		WriteBundleConfig(`
specs:
  - builtin: docker.run-command
    raw: /dockerext/run-command/
    docker.run-command:
      ContainerCreateConfig:
        Config:
          Image: alpine:latest
          Cmd: ["echo", "Hello World!"]
      EnablePull: true
      `)

		GenerateBundle()

		contents := GetFileFromBundle("dockerext/run-command/stdout")

		Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))
	})

})
