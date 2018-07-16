package supportbundle

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd/support-bundle/commands"
	. "github.com/replicatedcom/support-bundle/e2e/collect/ginkgo"
	"github.com/replicatedcom/support-bundle/pkg/collect/cli"
)

var _ = Describe("supportbundle.multiple-specs", func() {

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

			WriteFile("file1.txt", `File 1`)
			WriteFile("file2.txt", `File 2`)
			WriteFile("file3.txt", `File 3`)
			WriteFile("file4.txt", `File 4`)

			WriteFile("config1.yml", `
specs:
  - os.read-file:
      filepath: file1.txt
    output_dir: /os/read-file/`)

			WriteFile("config2.yml", `
specs:
  - os.read-file:
      filepath: file2.txt
    output_dir: /os/read-file/`)

			spec1 := `
specs:
  - os.read-file:
      filepath: file3.txt
    output_dir: /os/read-file/`

			spec2 := `
specs:
  - os.read-file:
      filepath: file4.txt
    output_dir: /os/read-file/`

			cmd := commands.NewSupportBundleCommand(cli.NewCli())
			buf := new(bytes.Buffer)
			cmd.SetOutput(buf)
			args := []string{
				"generate",
				fmt.Sprintf("--spec-file=%s", filepath.Join(GetTempDir(), "config1.yml")),
				fmt.Sprintf("--spec-file=%s", filepath.Join(GetTempDir(), "config2.yml")),
				fmt.Sprintf("--spec=%s", spec1),
				fmt.Sprintf("--spec=%s", spec2),
				fmt.Sprintf("--out=%s", filepath.Join(GetTempDir(), "bundle.tar.gz")),
				"--skip-default=true",
				"--timeout=10",
			}
			cmd.SetArgs(args)
			err := cmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			contents := GetFileFromBundle("os/read-file/file1.txt")
			Expect(strings.TrimSpace(contents)).To(Equal("File 1"))

			contents = GetFileFromBundle("os/read-file/file2.txt")
			Expect(strings.TrimSpace(contents)).To(Equal("File 2"))

			contents = GetFileFromBundle("os/read-file/file3.txt")
			Expect(strings.TrimSpace(contents)).To(Equal("File 3"))

			contents = GetFileFromBundle("os/read-file/file4.txt")
			Expect(strings.TrimSpace(contents)).To(Equal("File 4"))
		})
	})
})
