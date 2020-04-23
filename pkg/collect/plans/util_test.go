package plans

import (
	"bytes"
	"math/rand"
	"reflect"
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

func TestFileRedactor(t *testing.T) {
	type args struct {
		files    []string
		patterns []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				files: []string{
					"a/b/c.txt",
					"a/b/c.log",
					"a/b/c.zip",
					"a/b/c.tar",
					"/d/e.txt",
					"/d/e.log",
					"/d/e.zip",
					"/d/e.tar",
				},
				patterns: []string{
					"**/*.zip",
					"/d/e.txt",
					"*/e.tar",
				},
			},
			want: []string{
				"a/b/c.txt",
				"a/b/c.log",
				"a/b/c.tar",
				"/d/e.log",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileRedactor, err := FileRedactor(tt.args.patterns)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRedactor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := []string{}
			for _, name := range tt.args.files {
				if !fileRedactor(name) {
					got = append(got, name)
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileRedactor() = %v, want %v", got, tt.want)
			}
		})
	}
}
