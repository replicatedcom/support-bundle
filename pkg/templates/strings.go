package templates

import (
	"strings"
)

func init() {
	RegisterFunc("Trim", Trim)
	RegisterFunc("StringLength", StringLength)
}

// Trim returns a slice of the string with all leading and trailing code points
// contained in the optional cutset removed (default unicode.IsSpace).
func Trim(s string, args ...string) string {
	if len(args) == 0 {
		return strings.TrimSpace(s)
	}
	return strings.Trim(s, args[0])
}

func StringLength(i string) int {
	return len(i)
}
