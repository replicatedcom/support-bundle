package graphql

type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
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
	Spec     string `json:"spec,omitempty"`
	Hydrated string `json:"hydrated,omitempty"`
}

type SupportBundleUploadResponse struct {
	Data   *SupportBundleRequestResult `json:"data,omitempty"`
	Errors []Error                     `json:"errors,omitempty"`
}

type SupportBundleRequestResult struct {
	UploadSupportBundle `json:"uploadSupportBundle,omitempty"`
}

type UploadSupportBundle struct {
	UploadURI             string `json:"uploadUri,omitempty"`
	UploadedSupportBundle `json:"supportBundle,omitempty"`
}

type UploadedSupportBundle struct {
	ID string `json:"id,omitempty"`
}
