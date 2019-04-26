package analyzer

import (
	"io"
	"reflect"
	"testing"

	"github.com/go-kit/kit/log"
	v1 "github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	bundlereader "github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader"
	"github.com/replicatedcom/support-bundle/pkg/collect/bundle/reader/testfixtures"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzer_distillBundle(t *testing.T) {
	type args struct {
		bundleReader    bundlereader.BundleReader
		resultsToWriter map[collecttypes.Result]io.Writer
	}
	tests := []struct {
		name       string
		bundlePath string
		variables  []v1.Variable
		want       map[string][]interface{}
		wantErr    bool
	}{
		{
			name:       "basic",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "majorVersion",
					},
					FileMatch: &variable.FileMatch{
						PathRegexps: []string{`kubernetes/version/server_version\.json`},
						Distiller: variable.Distiller{
							RegexpCapture: &distiller.RegexpCapture{
								Regexp: `"major": "([^"]+)"`,
								Index:  1,
							},
						},
					},
				},
				{
					Meta: meta.Meta{
						Name: "minorVersion",
					},
					CollectRef: &variable.CollectRef{
						Ref: meta.Ref{
							Selector: meta.Selector{"analyze": "kubernetes-version"},
						},
						Distiller: variable.Distiller{
							RegexpCapture: &distiller.RegexpCapture{
								Regexp: `"minor": "([^"]+)"`,
								Index:  1,
							},
						},
					},
				},
			},
			want: map[string][]interface{}{
				"majorVersion": []interface{}{"1"},
				"minorVersion": []interface{}{"8+"},
			},
		},
		{
			name:       "no name",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					FileMatch: &variable.FileMatch{
						PathRegexps: []string{`kubernetes/version/server_version\.json`},
						Distiller: variable.Distiller{
							RegexpCapture: &distiller.RegexpCapture{
								Regexp: `"gitCommit": "([^"]+)"`,
								Index:  1,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name:       "no variable",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "gitCommit",
					},
				},
			},
			wantErr: true,
		},
		{
			name:       "no match",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "gitCommit",
					},
					FileMatch: &variable.FileMatch{
						PathRegexps: []string{`etc/os-release/stdout`},
						Distiller: variable.Distiller{
							Identity: &distiller.Identity{},
						},
					},
				},
			},
			want: map[string][]interface{}{},
		},
		{
			name:       "no distiller",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "serverVersionJson",
					},
					FileMatch: &variable.FileMatch{
						PathRegexps: []string{`kubernetes/version/server_version\.json`},
					},
				},
			},
			want: map[string][]interface{}{
				"serverVersionJson": []interface{}{"{\n  \"major\": \"1\",\n  \"minor\": \"8+\",\n  \"gitVersion\": \"v1.8.10-gke.0\",\n  \"gitCommit\": \"16ebd0de8e0ab2d1ef86d5b16ab1899b624a77cd\",\n  \"gitTreeState\": \"clean\",\n  \"buildDate\": \"2018-03-20T20:21:01Z\",\n  \"goVersion\": \"go1.8.3b4\",\n  \"compiler\": \"gc\",\n  \"platform\": \"linux/amd64\"\n}\n"},
			},
		},
		{
			name:       "distill error",
			bundlePath: "bundle",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "gitCommit",
					},
					FileMatch: &variable.FileMatch{
						PathRegexps: []string{`kubernetes/version/server_version\.json`},
						Distiller: variable.Distiller{
							RegexpCapture: &distiller.RegexpCapture{
								Regexp: `(`,
								Index:  1,
							},
						},
					},
				},
			},
			wantErr: true,
			want:    map[string][]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			f, err := fs.Create("./bundle.tgz")
			require.NoError(t, err)
			err = testfixtures.WriteBundle(f, tt.bundlePath)
			require.NoError(t, err)
			bundleReader, err := bundlereader.NewBundle(fs, "./bundle.tgz", "")
			require.NoError(t, err)

			a := &Analyzer{
				Logger: log.NewNopLogger(),
				Fs:     fs,
			}
			got, err := a.distillBundle(tt.variables, bundleReader)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyzer.distillBundle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAnalyzer_extractValues(t *testing.T) {
	evalUpper := variable.Eval("{{repl upper .os}}")

	tests := []struct {
		name                     string
		variables                []v1.Variable
		variableNamesToDistilled map[string][]interface{}
		want                     map[string]interface{}
		wantErr                  bool
	}{
		{
			name: "basic",
			variables: []v1.Variable{
				{
					Meta: meta.Meta{
						Name: "os",
					},
					Os: &variable.Os{},
				},
				{
					Meta: meta.Meta{
						Name: "osUpper",
					},
					Eval: &evalUpper,
				},
			},
			variableNamesToDistilled: map[string][]interface{}{
				"os": []interface{}{"centos"},
			},
			want: map[string]interface{}{
				"os":      "centos",
				"osUpper": "CENTOS",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Analyzer{
				Logger: log.NewNopLogger(),
			}
			got, err := a.extractValues(tt.variables, tt.variableNamesToDistilled)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyzer.extractValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Analyzer.extractValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
