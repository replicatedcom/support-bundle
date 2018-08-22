package analyzer

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/api/v1"
)

func TestEvalCondition(t *testing.T) {
	type args struct {
		condition v1.Condition
		vars      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "and true",
			args: args{
				condition: v1.Condition{
					Eval: &v1.EvalCondition{
						Operator: v1.AndOperator,
						Statements: []string{
							"true",
							"true",
						},
					},
				},
				vars: nil,
			},
			want: true,
		},
		{
			name: "and false",
			args: args{
				condition: v1.Condition{
					Eval: &v1.EvalCondition{
						Operator: v1.AndOperator,
						Statements: []string{
							"false",
							"true",
						},
					},
				},
				vars: nil,
			},
			want: false,
		},
		{
			name: "or true",
			args: args{
				condition: v1.Condition{
					Eval: &v1.EvalCondition{
						Operator: v1.OrOperator,
						Statements: []string{
							"true",
							"false",
						},
					},
				},
				vars: nil,
			},
			want: true,
		},
		{
			name: "or false",
			args: args{
				condition: v1.Condition{
					Eval: &v1.EvalCondition{
						Operator: v1.OrOperator,
						Statements: []string{
							"false",
							"false",
						},
					},
				},
				vars: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalCondition(tt.args.condition, tt.args.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EvalCondition() = %v, want %v", got, tt.want)
			}
		})
	}
}
