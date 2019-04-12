package condition

import "testing"

func TestEvalCondition_Eval(t *testing.T) {
	type args struct {
		ref  interface{}
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		eval    EvalCondition
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "scalar",
			eval: EvalCondition(`true`),
			args: args{
				ref: "test",
			},
			want: true,
		},
		{
			name: "ref",
			eval: EvalCondition(`{{repl eq .Ref "true"}}`),
			args: args{
				ref: "true",
			},
			want: true,
		},
		{
			name: "from vars",
			eval: EvalCondition(`{{repl eq .var 22}}`),
			args: args{
				ref: "true",
				data: map[string]interface{}{
					"var": 22,
				},
			},
			want: true,
		},
		{
			name: "from vars false",
			eval: EvalCondition(`{{repl eq .var 33}}`),
			args: args{
				ref: "true",
				data: map[string]interface{}{
					"var": 22,
				},
			},
			want: false,
		},
		{
			name:    "error",
			eval:    EvalCondition(`{{repl blah}}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.eval.Eval(tt.args.ref, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalCondition.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvalCondition.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
