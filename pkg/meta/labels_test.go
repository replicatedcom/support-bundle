package meta

import "testing"

func TestSelectorsMatch(t *testing.T) {
	type args struct {
		selector Selector
		labels   Labels
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match",
			args: args{
				selector: Selector{
					"a": "",
					"d": "1",
				},
				labels: Labels{
					"a": "",
					"b": "C",
					"d": "1",
				},
			},
			want: true,
		},
		{
			name: "no match",
			args: args{
				selector: Selector{
					"a": "",
					"b": "D",
					"d": "1",
				},
				labels: Labels{
					"a": "",
					"b": "C",
					"d": "1",
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SelectorMatches(tt.args.selector, tt.args.labels); got != tt.want {
				t.Errorf("SelectorMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
