package condition

import "testing"

func TestStringCompare_Eval(t *testing.T) {
	type args struct {
		ref  interface{}
		data map[string]interface{}
	}
	tests := []struct {
		name          string
		stringCompare StringCompare
		args          args
		want          bool
		wantErr       bool
	}{
		{
			name: "eq",
			stringCompare: StringCompare{
				Compare: Compare{
					Eq: "test",
				},
			},
			args: args{
				ref: "test",
			},
			want: true,
		},
		{
			name: "not eq",
			stringCompare: StringCompare{
				Compare: Compare{
					Eq: "test",
				},
			},
			args: args{
				ref: "not test",
			},
			want: false,
		},
		{
			name: "gte",
			stringCompare: StringCompare{
				Compare: Compare{
					Gte: "aaa",
				},
			},
			args: args{
				ref: "bbb",
			},
			want: true,
		},
		{
			name: "not gte",
			stringCompare: StringCompare{
				Compare: Compare{
					Gte: "bbb",
				},
			},
			args: args{
				ref: "aaa",
			},
			want: false,
		},
		{
			name: "in",
			stringCompare: StringCompare{
				Compare: Compare{
					In: []interface{}{
						"bbb", "aaa", "ccc",
					},
				},
			},
			args: args{
				ref: "aaa",
			},
			want: true,
		},
		{
			name: "not in",
			stringCompare: StringCompare{
				Compare: Compare{
					In: []interface{}{
						"bbb", "ccc",
					},
				},
			},
			args: args{
				ref: "aaa",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.stringCompare.Eval(tt.args.ref, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringCompare.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StringCompare.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
