package graphql

import "os"

var (
	endpoint = "https://pg.staging.replicated.com/graphql"
)

func init() {
	ep := os.Getenv("REPLICATED_API_ENDPOINT")
	if ep != "" {
		endpoint = ep
	}
}
