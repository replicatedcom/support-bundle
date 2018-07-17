package templates

import (
	"reflect"
	"testing"
)

func TestJSONQuery(t *testing.T) {
	type args struct {
		obj string
		exp string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "1",
			args: args{
				obj: `{"gitVersion": "v1.10.3"}`,
				exp: ".gitVersion",
			},
			want: []interface{}{"v1.10.3"},
		},
		{
			name: "1",
			args: args{
				obj: `{"items":[{"status":{"capacity":{"memory":"4042004Ki"}}},{"status":{"capacity":{"memory":"5042004Ki"}}}]}`,
				exp: ".items[].status.capacity.memory",
			},
			want: []interface{}{"4042004Ki", "5042004Ki"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JSONQuery(tt.args.obj, tt.args.exp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
