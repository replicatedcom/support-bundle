package util

import "encoding/json"

func SpewJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
