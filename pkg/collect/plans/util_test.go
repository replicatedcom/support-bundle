package plans

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_filterStreams(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		scrubber func([]byte) []byte
	}{
		{
			name: "no newline",
			input: `no
newline
at
end`,
			expected: `no
newline
at
end`,
			scrubber: func(b []byte) []byte {
				return b
			},
		},
		{
			name: "newline",
			input: `
newline
at
end
`,
			expected: `
newline
at
end
`,
			scrubber: func(b []byte) []byte {
				return b
			},
		},
		{
			name: "scrub",
			input: `
newline
at
end
`,
			expected: `

at
end
`,
			scrubber: func(b []byte) []byte {
				if string(b) == "newline" {
					return nil
				}
				return b
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in, out := bytes.NewBuffer([]byte(tt.input)), bytes.NewBuffer(nil)
			err := filterStreams(in, out, tt.scrubber)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, out.String())
		})
	}
}
