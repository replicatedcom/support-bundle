package ginkgo

import (
	"path"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
)

var _ = Describe("Docker run command", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Successfully executes the docker run command", func() {

		WriteFile("config.yml", `
specs:
  - builtin: docker.run-command
    raw: /dockerext/run-command/
    config:
      image: replicated/support-bundle:latest
      command: echo
      args: ["Hello World!"]
      `)

		err := cmd.Generate(
			path.Join(tmpdir, "config.yml"),
			path.Join(tmpdir, "bundle.tar.gz"),
			true,
			60,
		)

		Expect(err).To(BeNil())

		contents := ReadFileFromBundle(
			path.Join("bundle.tar.gz"),
			"/dockerext/run-command/stdout",
		)

		Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))
	})

})
