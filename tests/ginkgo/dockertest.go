package ginkgo

import (
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
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

		err := cmd.Generate(
			path.Join(tmpdir, "config.yml"),
			"",
			path.Join(tmpdir, "bundle.tar.gz"),
			true,
			60,
		)

		Expect(err).To(BeNil())

		contents := ReadFileFromBundle(
			path.Join("bundle.tar.gz"),
			"/daemon/docker/docker_info.json",
		)

		Expect(contents).To(ContainSubstring("DriverStatus"))
	})

})
