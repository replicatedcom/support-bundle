package marketapi

type ChannelCollectorsResponse struct {
	ID        string `json:"id"`
	Spec      string `json:"spec,omitempty"`
	Lifecycle string `json:"lifecycle,omitempty"`
}

type UploadURLResponse struct {
	BundleID  string `json:"bundle_id"`
	UploadURL string `json:"upload_url"`
}
