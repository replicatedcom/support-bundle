package journald

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("journald.logs", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - journald.logs:
      unit: docker
      since: -7 days
      reverse: true
    output_dir: /journald/logs/docker/`)

			GenerateBundle("--journald")

			_ = GetResultFromBundle("journald/logs/docker/logs.raw")
			contents := GetFileFromBundle("journald/logs/docker/logs.raw")
			Expect(contents).To(ContainSubstring("dockerd"))
		})
	})
})

var _ = Describe("journald.logs docker", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "1")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - journald.logs:
      unit: docker
      since: -7 days
      reverse: true
    output_dir: /journald/logs/docker/`)

			GenerateBundle("--journald")

			_ = GetResultFromBundle("journald/logs/docker/logs.raw")
			contents := GetFileFromBundle("journald/logs/docker/logs.raw")
			Expect(contents).To(ContainSubstring("dockerd"))
		})
	})
})
