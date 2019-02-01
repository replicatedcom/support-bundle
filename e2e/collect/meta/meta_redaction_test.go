package meta

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
)

var _ = Describe("meta.redaction", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFromBundle)
	AfterEach(CleanupDir)

	Context("When it is the only redaction", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo $HI]
      env: [HI=hello!]
      output_dir: /os/run-command/echo/
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "[he][he]llo"
        replace: "goodbye"
`)

			GenerateBundle()

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("goodbye!"))

			_ = GetResultFromBundle("redact/scrubs.json")
			redactions := GetFileFromBundle("redact/scrubs.json")
			Expect(redactions).To(Equal(`[
  {
    "regex": "[he][he]llo",
    "replace": "goodbye"
  }
]
`))
		})
	})

	Context("When it has multiple redactions", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo 'abc 123!']
      output_dir: /os/run-command/echo/
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "abc"
        replace: "xyz"
      - regex: "123"
        replace: "789"
`)

			GenerateBundle()

			_ = GetResultFromBundle("redact/scrubs.json")
			redactions := GetFileFromBundle("redact/scrubs.json")
			Expect(redactions).To(Equal(`[
  {
    "regex": "xyz",
    "replace": "xyz"
  },
  {
    "regex": "789",
    "replace": "789"
  }
]
`))

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("xyz 789!\n"))

		})
	})

	Context("When it has multiple, dependent redactions", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo 'thetest']
      output_dir: /os/run-command/echo/
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "thetest"
        replace: "456456"
      - regex: "456"
        replace: "654"
`)

			GenerateBundle()

			_ = GetResultFromBundle("redact/scrubs.json")
			redactions := GetFileFromBundle("redact/scrubs.json")
			Expect(redactions).To(Equal(`[
  {
    "regex": "654654",
    "replace": "654654"
  },
  {
    "regex": "654",
    "replace": "654"
  }
]
`))

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("654654\n"))

		})
	})

	Context("When it has multiple meta.redact entries", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo 'first second!']
      output_dir: /os/run-command/echo/
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "first"
        replace: "third"
  - meta.redact:
      output_dir: /redact2/
      scrubs:
      - regex: "second"
        replace: "fourth"
`)

			GenerateBundle()

			_ = GetResultFromBundle("redact/scrubs.json")
			redactions := GetFileFromBundle("redact/scrubs.json")
			Expect(redactions).To(Equal(`[
  {
    "regex": "third",
    "replace": "third"
  }
]
`))

			_ = GetResultFromBundle("redact2/scrubs.json")
			redactions2 := GetFileFromBundle("redact2/scrubs.json")
			Expect(redactions2).To(Equal(`[
  {
    "regex": "fourth",
    "replace": "fourth"
  }
]
`))

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("third fourth!\n"))

		})
	})

	Context("When the command also has a redact entry", func() {

		It("should output the correct files in the bundle", func() {

			WriteBundleConfig(`
specs:
  - os.run-command:
      name: sh
      args: [-c, echo 'fifth sixth!']
      output_dir: /os/run-command/echo/
      scrub:
        regex: "sixth"
        replace: "eighth"
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "fifth"
        replace: "seventh"
`)

			GenerateBundle()

			_ = GetResultFromBundle("redact/scrubs.json")
			redactions := GetFileFromBundle("redact/scrubs.json")
			Expect(redactions).To(Equal(`[
  {
    "regex": "seventh",
    "replace": "seventh"
  }
]
`))

			_ = GetResultFromBundle("os/run-command/echo/stdout")
			contents := GetFileFromBundle("os/run-command/echo/stdout")
			Expect(contents).To(Equal("seventh eighth!"))

		})
	})

	Context("overwrites file contents", func() {
		inContainer := os.Getenv("IN_CONTAINER")
		BeforeEach(func() {
			os.Setenv("IN_CONTAINER", "")
		})
		AfterEach(func() {
			os.Setenv("IN_CONTAINER", inContainer)
		})

		It("is able to read the file", func() {
			WriteFile("/tmp/testfile.txt", "this_is_a_test_file")
			WriteBundleConfig(`
specs:
  - os.read-file:
      filepath: /tmp/testfile.txt
      output_dir: /file/
  - meta.redact:
      output_dir: /redact/
      scrubs:
      - regex: "_a_"
        replace: "_not_a_"
`)

			GenerateBundle()

			outputPath := "file/testfile.txt"
			_ = GetResultFromBundle(outputPath)
			contents := GetFileFromBundle(outputPath)
			Expect(contents).To(Equal("this_is_not_a_test_file"))
		})
	})

})
