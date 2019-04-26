package reader

import "strings"

func PathsMatch(a, b string) bool {
	return strings.TrimLeft(a, "/") == strings.TrimLeft(b, "/")
}

func PathTrimPrefix(path, prefix string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	return strings.TrimLeft(strings.TrimPrefix(path, prefix), "/"), true
}
