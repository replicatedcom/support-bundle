package v1

import (
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/analyze/condition"
)

func TestCondition_Eval(t *testing.T) {
	osEqRhelEval := condition.EvalCondition(`{{repl eq .os "rhel"}}`)
	errorEval := condition.EvalCondition(`{{repl thisisnotafunc .os "rhel"}}`)

	tests := []struct {
		name      string
		condition Condition
		data      map[string]interface{}
		want      bool
		wantErr   bool
	}{
		{
			name: "or predicate",
			condition: Condition{
				Or: &OrPredicate{
					{
						StringCompare: &condition.StringCompare{
							Compare: condition.Compare{Eq: "centos"},
						},
						Ref: "os",
					},
					{
						StringCompare: &condition.StringCompare{
							Compare: condition.Compare{Eq: "rhel"},
						},
						Ref: "os",
					},
				},
			},
			data: map[string]interface{}{
				"os": "rhel",
			},
			want: true,
		},
		{
			name: "and predicate",
			condition: Condition{
				And: &AndPredicate{
					{
						StringCompare: &condition.StringCompare{
							Compare: condition.Compare{Eq: "centos"},
						},
						Ref: "os",
					},
					{
						StringCompare: &condition.StringCompare{
							Compare: condition.Compare{Eq: "rhel"},
						},
						Ref: "os",
					},
				},
			},
			data: map[string]interface{}{
				"os": "rhel",
			},
			want: false,
		},
		{
			name: "not predicate",
			condition: Condition{
				Not: &NotPredicate{
					Condition: Condition{
						StringCompare: &condition.StringCompare{
							Compare: condition.Compare{Eq: "centos"},
						},
						Ref: "os",
					},
				},
			},
			data: map[string]interface{}{
				"os": "rhel",
			},
			want: true,
		},
		{
			name: "eval true",
			condition: Condition{
				EvalCondition: &osEqRhelEval,
			},
			data: map[string]interface{}{
				"os": "rhel",
			},
			want: true,
		},
		{
			name: "eval false",
			condition: Condition{
				EvalCondition: &osEqRhelEval,
			},
			data: map[string]interface{}{
				"os": "centos",
			},
			want: false,
		},
		{
			name: "eval error",
			condition: Condition{
				EvalCondition: &errorEval,
			},
			data: map[string]interface{}{
				"os": "centos",
			},
			wantErr: true,
		},
		{
			name:      "empty",
			condition: Condition{},
			data: map[string]interface{}{
				"os": "centos",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.condition.Eval(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Condition.Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Condition.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}
