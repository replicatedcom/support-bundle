package core

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

func GenerateLargeFileString() string {
	// Create 10MB string for text file
	data := make([]byte, int(1e7))
	for i := range data {
		data[i] = byte('f')
	}
	return string(data)
}

var _ = Describe("os.read-file", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {
		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
    - os.read-file:
        filepath: /etc/os-release
        output_dir: /os/read-file/etc/
    - os.read-file:
        filepath: /tmp/notfound.txt
        output_dir: /os/read-file/notfound/
`)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/etc/os-release")
			contents := GetFileFromBundle("os/read-file/etc/os-release")
			Expect(contents).To(ContainSubstring("VERSION="))

			ExpectBundleErrorToHaveOccurred("os/read-file/notfound", "stat /tmp/notfound.txt: no such file or directory")
		})

		Context("IncludeEmpty is set", func() {
			It("should output all files in the bundle", func() {
				WriteFile("/tmp/empty.txt", "")
				WriteBundleConfig(`
specs:
    - os.read-file:
        filepath: /tmp/empty.txt
        output_dir: /os/read-file/notincludedempty/
    - os.read-file:
        filepath: /tmp/empty.txt
        include_empty: true
        output_dir: /os/read-file/includedempty/
`)

				GenerateBundle()

				notIncludedEmptyFilePath := "os/read-file/notincludedempty/empty.txt"
				_ = GetResultFromBundle(notIncludedEmptyFilePath)
				ExpectFileNotInBundle(notIncludedEmptyFilePath)

				includedEmptyFilePath := "os/read-file/includedempty/empty.txt"
				_ = GetResultFromBundle(includedEmptyFilePath)
				contents := GetFileFromBundle(includedEmptyFilePath)
				Expect(contents).To(Equal(""))
			})
		})

		Context("large file provided", func() {
			It("is able to read the file", func() {
				expectedContents := GenerateLargeFileString()
				WriteFile("/tmp/large.txt", expectedContents)
				WriteBundleConfig(`
specs:
    - os.read-file:
        filepath: /tmp/large.txt
        output_dir: /os/read-file/largefolder/
`)

				GenerateBundle()

				outputPath := "os/read-file/largefolder/large.txt"
				_ = GetResultFromBundle(outputPath)
				contents := GetFileFromBundle(outputPath)
				Expect(contents).To(Equal(expectedContents))
			})
		})
	})
})

var _ = XDescribe("os.read-file docker", func() {

	inContainer := os.Getenv("IN_CONTAINER")
	BeforeEach(func() {
		os.Setenv("IN_CONTAINER", "1")
	})
	AfterEach(func() {
		os.Setenv("IN_CONTAINER", inContainer)
	})

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When the spec is run", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
    - os.read-file:
        filepath: /etc/profile
        output_dir: /os/read-file/etc/
    - os.read-file:
        filepath: /tmp/notfound.txt
        output_dir: /os/read-file/notfound/
`)

			GenerateBundle()

			_ = GetResultFromBundle("os/read-file/etc/profile")
			contents := GetFileFromBundle("os/read-file/etc/profile")
			Expect(contents).To(ContainSubstring("profile.d"))

			ExpectBundleErrorToHaveOccurred("os/read-file/notfound", "docker read file: file not found")
		})

		Context("large file provided", func() {
			It("is able to read the file", func() {
				expectedContents := GenerateLargeFileString()
				WriteFile("/tmp/large.txt", expectedContents)

				WriteBundleConfig(`
specs:
    - os.read-file:
        filepath: /tmp/large.txt
        output_dir: /os/read-file/largefolder/
`)

				GenerateBundle()

				outputPath := "os/read-file/largefolder/large.txt"
				_ = GetResultFromBundle(outputPath)
				contents := GetFileFromBundle(outputPath)
				Expect(contents).To(Equal(expectedContents))
			})
		})
	})
})
