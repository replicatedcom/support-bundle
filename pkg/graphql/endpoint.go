package graphql

import "os"

var endpoint string

func init() {
	ep := os.Getenv("REPLICATED_API_ENDPOINT")
	if ep != "" {
		endpoint = ep
	}
}
