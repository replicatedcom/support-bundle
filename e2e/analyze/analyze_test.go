package analyze_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/kylelemons/godebug/pretty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyze"
	"github.com/replicatedcom/support-bundle/pkg/analyze/analyzer"
	"github.com/replicatedcom/support-bundle/pkg/analyze/api"
	"github.com/replicatedcom/support-bundle/pkg/analyze/collector"
	"github.com/replicatedcom/support-bundle/pkg/analyze/render"
	"github.com/replicatedcom/support-bundle/pkg/analyze/resolver"
	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

type TestMetadata struct {
	ExpectErrAnalysisFailed bool `yaml:"expect_err_analysis_failed"`
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
		var testMetadata TestMetadata

		BeforeEach(func() {
			// read the test metadata
			metadataBytes, err := ioutil.ReadFile(path.Join(testPath, "metadata.yml"))
			Expect(err).NotTo(HaveOccurred())
			err = yaml.Unmarshal(metadataBytes, &testMetadata)
			Expect(err).NotTo(HaveOccurred())
		})

		Context(fmt.Sprintf("When the spec in %q is run", file.Name()), func() {
			testSpecPath := path.Join(testPath, "spec.yml")
			testExpectedPath := path.Join(testPath, "expected.yml")
			testBundlePath := path.Join(testPath, "bundle")

			It("Should output files matching those expected", func() {
				spec, err := ioutil.ReadFile(testSpecPath)
				Expect(err).NotTo(HaveOccurred())

				expected, err := ioutil.ReadFile(testExpectedPath)
				Expect(err).NotTo(HaveOccurred())

				logger := log.NewNopLogger()
				fs := afero.NewMemMapFs()

				c := collector.NewMock(fs, testBundlePath)
				specUnmarshalled, err := api.DeserializeSpec(spec)
				Expect(err).NotTo(HaveOccurred())
				c.On("CollectBundle", specUnmarshalled.Collect)

				a := &analyze.Analyze{
					Logger: logger,

					Resolver:  resolver.New(logger, fs),
					Collector: c,
					Analyzer:  analyzer.New(logger, fs),

					Specs: []string{string(spec)},
				}

				ctx := context.Background()

				results, err := a.Execute(ctx)
				if testMetadata.ExpectErrAnalysisFailed {
					Expect(err).To(Equal(analyze.ErrAnalysisFailed))
				} else {
					Expect(err).NotTo(HaveOccurred())
				}

				c.AssertExpectations(GinkgoT())

				r := render.New(nil, "yaml")
				var actual bytes.Buffer
				err = r.RenderResults(ctx, &actual, results)
				Expect(err).NotTo(HaveOccurred())

				var outExpected, outActual []interface{}
				err = yaml.Unmarshal(expected, &outExpected)
				Expect(err).NotTo(HaveOccurred())
				err = yaml.Unmarshal(actual.Bytes(), &outActual)
				Expect(err).NotTo(HaveOccurred())

				Expect(outActual).To(Equal(outExpected), pretty.Compare(outActual, outExpected))

			}, 60)

		})
	}
})
