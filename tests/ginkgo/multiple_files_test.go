package ginkgo

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/replicatedcom/support-bundle/pkg/cli/commands"
)

var _ = Describe("Multiple files", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("Successfully collects specs from multiple files", func() {
		WriteFile("file1.txt", `File 1`)
		WriteFile("file2.txt", `File 2`)

		WriteFile("config1.yml", `
specs:
  - builtin: core.read-file
    raw: /custom/file1
    config:
      file_path: file1.txt`)
		WriteFile("config2.yml", `
specs:
  - builtin: core.read-file
    raw: /custom/file2
    config:
      file_path: file2.txt`)

		cmd := commands.NewSupportBundleCommand(cli.NewCli())
		buf := new(bytes.Buffer)
		cmd.SetOutput(buf)
		cmd.SetArgs([]string{
			"generate",
			fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config1.yml")),
			fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config2.yml")),
			fmt.Sprintf("--out=%s", filepath.Join(tmpdir, "bundle.tar.gz")),
			"--skip-default=true",
			"--timeout=10",
		})
		err := cmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		contents := GetFileFromBundle("custom/file1")
		Expect(strings.TrimSpace(contents)).To(Equal("File 1"))

		contents = GetFileFromBundle("custom/file2")
		Expect(strings.TrimSpace(contents)).To(Equal("File 2"))
	})

})
