package core

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/ginkgo"
)

var _ = Describe("os.read-file", func() {

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
  - os.read-file:
      filepath: /etc/os-release
    output_dir: /os/read-file/etc/os-release/
  - os.read-file:
      filepath: /tmp/notfound.txt
    output_dir: /os/read-file/notfound/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/etc/os-release/contents")
			contents := GetFileFromBundle("os/read-file/etc/os-release/contents")
			Expect(contents).To(ContainSubstring("VERSION="))

			ExpectBundleErrorToHaveOccured("os/read-file/notfound", "open /tmp/notfound.txt: no such file or directory")
		})
	})
})

var _ = Describe("os.read-file docker", func() {

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
  - os.read-file:
      filepath: /etc/profile
    output_dir: /os/read-file/etc/profile/
  - os.read-file:
      filepath: /tmp/notfound.txt
    output_dir: /os/read-file/notfound/`)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/etc/profile/contents")
			contents := GetFileFromBundle("os/read-file/etc/profile/contents")
			Expect(contents).To(ContainSubstring("profile.d"))

			ExpectBundleErrorToHaveOccured("os/read-file/notfound", "docker read file: file not found")
		})
	})
})