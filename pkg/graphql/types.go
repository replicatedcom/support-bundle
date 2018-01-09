package graphql

import "net/url"

type Client struct {
	GQLServer *url.URL
}

type Request struct {
	Query         string            `json:"query"`
	Variables     map[string]string `json:"variables"`
	OperationName string            `json:"operationName"`
}

type Error struct {
	Locations []map[string]interface{} `json:"locations"`
	Message   string                   `json:"message"`
}

type SupportBundleResponse struct {
	Data   *SupportBundleResult `json:"data,omitempty"`
	Errors []Error              `json:"errors,omitempty"`
}

type SupportBundleResult struct {
	SupportBundle `json:"supportBundleSpec"`
}

type SupportBundle struct {
	Spec     string `json:"spec"`
	Hydrated string `json:"hydrated"`
}
