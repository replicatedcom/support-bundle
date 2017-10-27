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
	AfterEach(LogErrors("bundle.tar.gz"))
	AfterEach(CleanupDir)

	It("Successfully executes the docker run command", func() {

		cfgDoc := `
specs:
  - builtin: docker.run-command
    raw: /dockerext/run-command/
    config:
      image: circleci/golang:1.8.1
      command: echo
      args: ["Hello World!"]`

		err := cmd.Generate(
			"",
			cfgDoc,
			path.Join(tmpdir, "bundle.tar.gz"),
			true,
			60,
		)

		Expect(err).To(BeNil())

		contents := ReadFileFromBundle(
			"bundle.tar.gz",
			"/dockerext/run-command/stdout",
		)

		Expect(strings.TrimSpace(contents)).To(Equal("Hello World!"))
	})

})
