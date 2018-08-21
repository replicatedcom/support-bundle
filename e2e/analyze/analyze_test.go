package analyze_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/mholt/archiver"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/cli"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type TestMetadata struct {
	ExpectErr bool `yaml:"expect_err"`
}

func TestCore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "integration")
}

var _ = Describe("integration", func() {
	integrationDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(integrationDir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		testPath := path.Join(integrationDir, file.Name())
		testSpecPath := path.Join(testPath, "spec.yml")
		testExpectedPath := path.Join(testPath, "expected.yml")
		testBundlePath := path.Join(testPath, "bundle")
		testBundleDestPath := path.Join(testPath, "bundle.tgz")
		var testMetadata TestMetadata

		BeforeEach(func() {
			// read the test metadata
			metadataBytes, err := ioutil.ReadFile(path.Join(testPath, "metadata.yml"))
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(metadataBytes, &testMetadata)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			// remove the temporary bundle
			_ = os.RemoveAll(testBundleDestPath)
		})

		Context(fmt.Sprintf("When the spec in %q is run", file.Name()), func() {

			It("Should output the expected results", func() {

				expected, err := ioutil.ReadFile(testExpectedPath)
				Expect(err).NotTo(HaveOccurred())
				var outExpected []interface{}
				err = yaml.Unmarshal(expected, &outExpected)
				Expect(err).NotTo(HaveOccurred())

				_, err = makeBundle(afero.NewOsFs(), testBundlePath, testBundleDestPath)
				Expect(err).NotTo(HaveOccurred())

				cmd := cli.RootCmd()
				buf := new(bytes.Buffer)
				cmd.SetOutput(buf)
				cmd.SetArgs([]string{
					"run",
					fmt.Sprintf("--collect-bundle-path=%s", testBundleDestPath),
					fmt.Sprintf("--spec-file=%s", testSpecPath),
					"--output=yaml",
					"--log-level=off",
				})

				err = cmd.Execute()
				if testMetadata.ExpectErr {
					Expect(err).To(Equal(analyze.ErrSeverityThreshold))
				} else {
					Expect(err).NotTo(HaveOccurred())
				}

				var outActual []interface{}
				err = yaml.Unmarshal(buf.Bytes(), &outActual)
				Expect(err).NotTo(HaveOccurred())

				Expect(outActual).To(Equal(outExpected), pretty.Compare(outActual, outExpected))

			}, 60)

		})
	}
})

func makeBundle(fs afero.Fs, src, dest string) (os.FileInfo, error) {
	f, err := fs.Create(dest)
	if err != nil {
		return nil, errors.Wrapf(err, "create file %s", dest)
	}
	err = func() error {
		defer f.Close()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := os.Chdir(src); err != nil {
			return err
		}
		defer os.Chdir(cwd)

		var filePaths []string
		files, err := ioutil.ReadDir(src)
		if err != nil {
			return err
		}
		for _, info := range files {
			filePaths = append(filePaths, info.Name())
		}

		return archiver.TarGz.Write(f, filePaths)
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "create archive from %s", src)
	}
	return fs.Stat(dest)
}
