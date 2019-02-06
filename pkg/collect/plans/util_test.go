package plans

import (
	"bytes"
	"math/rand"
	"regexp"
	"testing"

	"github.com/replicatedcom/support-bundle/pkg/collect/types"

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

func Test_filterLongStream(t *testing.T) {
	tests := []struct {
		name    string
		regex   string
		replace string
		n       int
		charset string
	}{
		{
			name:    "1,000",
			regex:   "abc",
			replace: "xyz",
			n:       1000,
			charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			name:    "1,000,000",
			regex:   "abc",
			replace: "xyz",
			n:       1000 * 1000,
			charset: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
		{
			name:    "10,000,000",
			regex:   "abc",
			replace: "xyz",
			n:       1000 * 1000 * 10,
			charset: "abcdefghijklmnopqrstuvwxyz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			scrubber, err := RawScrubber(&types.Scrub{Regex: tt.regex, Replace: tt.replace})
			req.NoError(err)

			in := bytes.NewBuffer(randStringBytes(tt.n, tt.charset))
			out := bytes.NewBuffer(nil)
			out.Grow(tt.n)

			err = filterStreams(in, out, scrubber)
			req.NoError(err)

			outBytes := out.Bytes()
			inRegex := regexp.MustCompile(tt.regex)
			req.Falsef(inRegex.Match(outBytes), "input regex %q must not match the output string %q\n", tt.regex, outBytes)
		})
	}

}

func randStringBytes(n int, letterBytes string) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
