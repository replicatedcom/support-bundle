package templates

import (
	"github.com/stretchr/testify/require"
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
			if got := JSONQuery(tt.args.exp, tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONQueryTemplating(t *testing.T) {
	type args struct {
		obj interface{}
		exp string
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "no message slice parsing",
			args: args{
				obj: `{"gitVersion": "v1.10.3"}`,
				exp: "{{repl . | jq \".gitVersion\"}}",
			},
			want: "[v1.10.3]",
		},
		{
			name: "parse message slice to string",
			args: args{
				obj: `{"gitVersion": "v1.10.3"}`,
				exp: "{{repl . | jq \".gitVersion\" | jqMessagesToString}}",
			},
			want: "v1.10.3",
		},
		{
			name: "usage with parens instead of pipe",
			args: args{
				obj: `{"gitVersion": "v1.10.3"}`,
				exp: "{{repl (jq \".gitVersion\" .) | jqMessagesToString}}",
			},
			want: "v1.10.3",
		},
		{
			name: "number parsing example",
			args: args{
				obj: `{"replicas": 2.00}`,
				exp: "{{repl printf \"%d\" (. | jq \".replicas\" | jqMessagesToString | float64 | int)}}",
			},
			want: "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err :=  Execute(tt.args.exp, tt.args.obj)
			require.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
