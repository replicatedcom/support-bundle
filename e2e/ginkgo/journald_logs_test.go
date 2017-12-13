package ginkgo

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The journald.logs spec", func() {
	BeforeEach(EnterNewTempDir)
	AfterEach(CleanupDir)

	It("Collects logs from journalctl", func() {
		WriteBundleConfig(`
specs:
  - builtin: journald.logs
    raw: /journalctl/docker
    journald.logs:
      unit: docker`)

		GenerateBundle()

		Expect(err).NotTo(HaveOccurred())
		errors := GetFileFromBundle("error.json")
		if strings.Contains(errors, "No journal files were found") {
			// the host doesn't have journald installed
			return
		}

		Expect(errors).To(Equal("null"))

		stdout := GetFileFromBundle("journalctl/docker")
		Expect(stdout).ToNot(BeEmpty())
		Expect(stdout).To(ContainSubstring("Starting Docker Application Container Engine"))
	})
})
