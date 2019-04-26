package condition

import (
	"testing"

	"github.com/pkg/errors"
)

func TestEmpty_Eval(t *testing.T) {
	type args struct {
		ref  interface{}
		data map[string]interface{}
		err  error
	}
	tests := []struct {
		name    string
		empty   Empty
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:  "empty",
			empty: Empty{},
			args: args{
				ref: "",
			},
			want: true,
		},
		{
			name:  "not empty",
			empty: Empty{},
			args: args{
				ref: "test",
			},
			want: false,
		},
		{
			name:  "ref error",
			empty: Empty{},
			args: args{
				err: errors.New("ERROR"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.empty.Eval(tt.args.ref, tt.args.data, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Empty.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Empty.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
