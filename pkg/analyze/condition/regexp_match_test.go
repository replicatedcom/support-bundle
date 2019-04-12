package condition

import "testing"

func TestRegexpMatch_Eval(t *testing.T) {
	type args struct {
		ref  interface{}
		data map[string]interface{}
	}
	tests := []struct {
		name        string
		regexpMatch RegexpMatch
		args        args
		want        bool
		wantErr     bool
	}{
		{
			name: "match",
			regexpMatch: RegexpMatch{
				Regexp: "es",
			},
			args: args{
				ref: "test",
			},
			want: true,
		},
		{
			name: "not match",
			regexpMatch: RegexpMatch{
				Regexp: "se",
			},
			args: args{
				ref: "test",
			},
			want: false,
		},
		{
			name: "error",
			regexpMatch: RegexpMatch{
				Regexp: "(",
			},
			args: args{
				ref: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.regexpMatch.Eval(tt.args.ref, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexpMatch.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RegexpMatch.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
