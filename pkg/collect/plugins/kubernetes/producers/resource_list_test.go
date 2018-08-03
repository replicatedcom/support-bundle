package producers

import (
	"reflect"
	"runtime"
	"testing"
)

func Test_findResourceListFnFromAliasAndGv(t *testing.T) {
	type args struct {
		alias string
		gv    string
	}
	tests := []struct {
		name    string
		args    args
		want    resourceListFn
		wantErr bool
	}{
		{
			name: "default",
			args: args{
				alias: "deployments",
				gv:    "",
			},
			want:    resourceListDeploymentsAppsV1,
			wantErr: false,
		},
		{
			name: "alternate gv 1",
			args: args{
				alias: "deployments",
				gv:    "apps/v1beta2",
			},
			want:    resourceListDeploymentsAppsV1beta2,
			wantErr: false,
		},
		{
			name: "alternate gv 2",
			args: args{
				alias: "deployments",
				gv:    "extensions/v1beta1",
			},
			want:    resourceListDeploymentsExtensionsV1beta1,
			wantErr: false,
		},
		{
			name: "alias",
			args: args{
				alias: "deploy",
				gv:    "apps/v1beta2",
			},
			want:    resourceListDeploymentsAppsV1beta2,
			wantErr: false,
		},
		{
			name: "not found kind",
			args: args{
				alias: "notfound",
				gv:    "",
			},
			wantErr: true,
		},
		{
			name: "not found gv",
			args: args{
				alias: "deployments",
				gv:    "apps/vblahbeta1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findResourceListFnFromAliasAndGv(tt.args.alias, tt.args.gv)
			if (err != nil) != tt.wantErr {
				t.Errorf("findResourceListFnFromAliasAndGv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf(
					"findResourceListFnFromAliasAndGv() = %v, want %v",
					runtime.FuncForPC(reflect.ValueOf(got).Pointer()).Name(),
					runtime.FuncForPC(reflect.ValueOf(tt.want).Pointer()).Name())
			}
		})
	}
}
