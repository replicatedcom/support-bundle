package v1

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable"
	"github.com/replicatedcom/support-bundle/pkg/analyze/variable/distiller"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	collectbundle "github.com/replicatedcom/support-bundle/pkg/test-mocks/collect/bundle"
)

func TestVariable_Register(t *testing.T) {
	tests := []struct {
		name            string
		variable        Variable
		registerExpects func(*collectbundle.MockBundleReader)
		want            map[string]interface{}
		wantErr         bool
	}{
		{
			name: "os ubuntu",
			variable: Variable{
				Os: &variable.Os{},
				Meta: meta.Meta{
					Name: "os",
				},
			},
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return([]collecttypes.Result{
						{
							Path: "/default/etc/os-release",
							Spec: collecttypes.Spec{
								CoreReadFile: &collecttypes.CoreReadFileOptions{
									Filepath: "/etc/os-release",
								},
							},
							Size: 1,
						},
					})

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(`NAME="Ubuntu"
VERSION="16.04.5 LTS (Xenial Xerus)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 16.04.5 LTS"
VERSION_ID="16.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
VERSION_CODENAME=xenial
UBUNTU_CODENAME=xenial`)), nil)
			},
			want: map[string]interface{}{"os": "ubuntu"},
		},
		{
			name: "os rhel 6",
			variable: Variable{
				Os: &variable.Os{},
				Meta: meta.Meta{
					Name: "os",
				},
			},
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					GetIndex().
					Return([]collecttypes.Result{
						{
							Path: "/default/etc/system-release",
							Spec: collecttypes.Spec{
								CoreReadFile: &collecttypes.CoreReadFileOptions{
									Filepath: "/etc/system-release",
								},
							},
							Size: 1,
						},
					})

				bundleReader.
					EXPECT().
					Open("/default/etc/system-release").
					Return(ioutil.NopCloser(strings.NewReader(`Red Hat Enterprise Linux Server release 6.10 (Santiago)`)), nil)
			},
			want: map[string]interface{}{"os": "rhel"},
		},
		{
			name: "collectRef",
			variable: Variable{
				CollectRef: &variable.CollectRef{
					Ref: meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					},
					Distiller: variable.Distiller{
						RegexpCapture: &distiller.RegexpCapture{
							Regexp: `(?m)^VERSION_ID="([^"]+)"`,
							Index:  1,
						},
					},
				},
				Meta: meta.Meta{
					Name: "osVersion",
				},
			},
			registerExpects: func(bundleReader *collectbundle.MockBundleReader) {
				bundleReader.
					EXPECT().
					ResultsFromRef(meta.Ref{
						Selector: meta.Selector{"analyze": "/etc/os-release"},
					}).
					Return([]collecttypes.Result{
						{
							Path: "/default/etc/os-release",
							Spec: collecttypes.Spec{
								CoreReadFile: &collecttypes.CoreReadFileOptions{
									Filepath: "/etc/os-release",
								},
							},
							Size: 1,
						},
					})

				bundleReader.
					EXPECT().
					Open("/default/etc/os-release").
					Return(ioutil.NopCloser(strings.NewReader(`NAME="Ubuntu"
VERSION="16.04.5 LTS (Xenial Xerus)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 16.04.5 LTS"
VERSION_ID="16.04"
HOME_URL="http://www.ubuntu.com/"
SUPPORT_URL="http://help.ubuntu.com/"
BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
VERSION_CODENAME=xenial
UBUNTU_CODENAME=xenial`)), nil)
			},
			want: map[string]interface{}{"osVersion": "16.04"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := gomock.NewController(t)
			bundleReader := collectbundle.NewMockBundleReader(mc)
			defer mc.Finish()

			if tt.registerExpects != nil {
				tt.registerExpects(bundleReader)
			}

			got, err := tt.variable.Register(bundleReader, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Variable.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variable.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}
