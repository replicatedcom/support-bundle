package ginkgo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker.container-run", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - docker.container-run:
      Config:
        Image: alpine:latest
        Cmd: ["echo", "Hello World!"]
      EnablePull: true
    output_dir: /docker/container-run/
  - docker.run:
      Config:
        Image: alpine:latest
        Cmd: ["echo", "foo bar"]
      EnablePull: true
    output_dir: /docker/run/`)

			GenerateBundle()

			var contents string

			_ = GetResultFromBundle("docker/container-run/stdout.raw")
			_ = GetResultFromBundle("docker/container-run/stderr.raw")
			contents = GetFileFromBundle("docker/container-run/stdout.raw")
			Expect(contents).To(Equal("Hello World!\n"))

			_ = GetResultFromBundle("docker/run/stdout.raw")
			_ = GetResultFromBundle("docker/run/stderr.raw")
			contents = GetFileFromBundle("docker/run/stdout.raw")
			Expect(contents).To(Equal("foo bar\n"))
		})
	})
})
