package ginkgo

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd/support-bundle/commands"
	"github.com/replicatedcom/support-bundle/pkg/cli"
)

var _ = Describe("Multiple specs", func() {

	BeforeEach(EnterNewTempDir)
	AfterEach(LogResultsFomBundle)
	AfterEach(CleanupDir)

	It("Successfully reads from multiple specs", func() {
		WriteFile("file1.txt", `File 1`)
		WriteFile("file2.txt", `File 2`)
		WriteFile("file3.txt", `File 3`)
		WriteFile("file4.txt", `File 4`)

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

		spec1 := `
specs:
  - builtin: core.read-file
    raw: /custom/file3
    config:
      file_path: file3.txt`
		spec2 := `
specs:
  - builtin: core.read-file
    raw: /custom/file4
    config:
      file_path: file4.txt`

		cmd := commands.NewSupportBundleCommand(cli.NewCli())
		buf := new(bytes.Buffer)
		cmd.SetOutput(buf)
		args := []string{
			"generate",
			fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config1.yml")),
			fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config2.yml")),
			fmt.Sprintf("--spec=%s", spec1),
			fmt.Sprintf("--spec=%s", spec2),
			fmt.Sprintf("--out=%s", filepath.Join(tmpdir, "bundle.tar.gz")),
			"--skip-default=true",
			"--timeout=10",
		}
		fmt.Println(args)
		cmd.SetArgs(args)
		err := cmd.Execute()
		Expect(err).NotTo(HaveOccurred())

		contents := GetFileFromBundle("custom/file1")
		Expect(strings.TrimSpace(contents)).To(Equal("File 1"))

		contents = GetFileFromBundle("custom/file2")
		Expect(strings.TrimSpace(contents)).To(Equal("File 2"))

		contents = GetFileFromBundle("custom/file3")
		Expect(strings.TrimSpace(contents)).To(Equal("File 3"))

		contents = GetFileFromBundle("custom/file4")
		Expect(strings.TrimSpace(contents)).To(Equal("File 4"))
	})

})
