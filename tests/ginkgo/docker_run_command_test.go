package ginkgo

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Docker run command", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogErrors("bundle.tar.gz"))
	AfterEach(CleanupDir)

	It("Successfully executes the docker run command", func() {
		WriteFile("config.yml", `
specs:
  - builtin: docker.run-command
    raw: /dockerext/run-command/
    config:
      image: alpine:latest
      command: echo
      args: ["Hello World!"]
      enable_pull: true`)

		GenerateBundle()

		contents := GetFileFromBundle("/dockerext/run-command/stdout")

		Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))
	})

})
