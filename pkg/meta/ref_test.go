package meta

import "testing"

func TestRefMatches(t *testing.T) {
	type args struct {
		ref  Ref
		meta *Meta
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match",
			args: args{
				ref: Ref{
					Name: "n1",
					Selector: Selector{
						"a": "",
						"d": "1",
					},
				},
				meta: &Meta{
					Name: "n1",
					Labels: Labels{
						"a": "",
						"b": "C",
						"d": "1",
					},
				},
			},
			want: true,
		},
		{
			name: "no match labels",
			args: args{
				ref: Ref{
					Name: "n1",
					Selector: Selector{
						"a": "",
						"b": "D",
						"d": "1",
					},
				},
				meta: &Meta{
					Name: "n1",
					Labels: Labels{
						"a": "",
						"b": "C",
						"d": "1",
					},
				},
			},
			want: false,
		},
		{
			name: "no match name",
			args: args{
				ref: Ref{
					Name: "n1",
					Selector: Selector{
						"a": "",
						"d": "1",
					},
				},
				meta: &Meta{
					Name: "n2",
					Labels: Labels{
						"a": "",
						"b": "C",
						"d": "1",
					},
				},
			},
			want: false,
		},
		{
			name: "nil meta",
			args: args{
				ref: Ref{
					Name: "n1",
					Selector: Selector{
						"a": "",
						"b": "D",
						"d": "1",
					},
				},
				meta: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RefMatches(tt.args.ref, tt.args.meta); got != tt.want {
				t.Errorf("RefMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}
