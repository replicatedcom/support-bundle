package graphql

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// deprecated
const customerSpecQuery = `
query {
  supportBundleSpec {
    spec
    hydrated
  }
}
`

const channelSpecQuery = `
query channelCollectors($channelId: String) {
  channelCollectors(channelId: $channelId) {
    spec
    hydrated
  }
}
`

const watchSpecQuery = `
query watchCollectors($watchId: String!) {
  watchCollectors(watchId: $watchId) {
    spec
    hydrated
  }
}
`

const startCustomerUploadMutation = `
mutation GetPresignedURI($size: Int, $notes: String) {
  uploadSupportBundle(size: $size, notes: $notes) {
    uploadUri,
    supportBundle {
      id
    }
  }
}
`

const startChannelUploadMutation = `
mutation GetChannelPresignedURI($channelId: String!, $size: Int, $notes: String) {
  uploadChannelSupportBundle(channelId: $channelId, size: $size, notes: $notes) {
    uploadUri,
    supportBundle {
      id
    }
  }
}
`

const startWatchUploadMutation = `
mutation GetWatchPresignedURI($watchId: String!, $size: Int) {
  uploadSupportBundle(watchId: $watchId, size: $size) {
    uploadUri,
    supportBundle {
      id
    }
  }
}
`

const markSupportBundleUploadedMutation = `
mutation FinishUploadSupportBundle($id: String!) {
  markSupportBundleUploaded(id: $id) {
    id
  }
}
`

type SupportBundleSpec struct {
	CustomerID string
}

type Client struct {
	endpoint string
	client   *http.Client
}

func NewClient(endpoint string, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	return &Client{
		endpoint: endpoint,
		client:   client,
	}
}

// GetCustomerSpec is deprecated
func (c *Client) GetCustomerSpec(id string) ([]byte, error) {
	resp, err := c.executeGraphQLQuery(id, Request{
		Query: customerSpecQuery,
	})

	if err != nil {
		return nil, errors.Wrap(err, "executing graphql request")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	specBody := SupportBundleResponse{}

	if err := decoder.Decode(&specBody); err != nil {
		return nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(specBody.Errors) > 0 {
		return nil, Errors{Errors: specBody.Errors}
	}

	return []byte(specBody.Data.Hydrated), nil
}

// GetChannelSpec will query the endpoint set in the client with a Replicated ChannelID
// and will return a fully rendered spec, containing collect and lifecycle keys
func (c *Client) GetChannelSpec(channelID string) ([]byte, error) {
	resp, err := c.executeGraphQLQuery("", Request{
		Query: channelSpecQuery,
		Variables: map[string]interface{}{
			"channelId": channelID,
		},
	})
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.Wrap(err, "executing graphql request")
	}

	decoder := json.NewDecoder(resp.Body)
	specBody := ChannelCollectorsResponse{}

	if err := decoder.Decode(&specBody); err != nil {
		return nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(specBody.Errors) > 0 {
		return nil, Errors{Errors: specBody.Errors}
	}

	return []byte(specBody.Data.Hydrated), nil
}

// GetWatchSpec will query the endpoint set in the client with a Replicated WatchID
// and will return a fully rendered spec, containing collect and lifecycle keys
func (c *Client) GetWatchSpec(watchID string) ([]byte, error) {
	resp, err := c.executeGraphQLQuery("", Request{
		Query: watchSpecQuery,
		Variables: map[string]interface{}{
			"watchId": watchID,
		},
	})

	if err != nil {
		return nil, errors.Wrap(err, "executing graphql request")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	specBody := WatchCollectorsResponse{}

	if err := decoder.Decode(&specBody); err != nil {
		return nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(specBody.Errors) > 0 {
		return nil, Errors{Errors: specBody.Errors}
	}

	return []byte(specBody.Data.Hydrated), nil
}

// GetSupportBundleCustomerUploadURI queries the Endpoint in Client to retrieve a URI that can be used to upload
// a support bundle for a specific customer
// This function is deprecated
func (c *Client) GetSupportBundleCustomerUploadURI(customerID string, size int64, notes string) (string, *url.URL, error) {
	resp, err := c.executeGraphQLQuery(customerID, Request{
		Query: startCustomerUploadMutation,
		Variables: map[string]interface{}{
			"size":  size,
			"notes": notes,
		},
	})

	if err != nil {
		return "", nil, errors.Wrap(err, "executing graphql request")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	uploadBody := SupportBundleUploadResponse{}

	if err := decoder.Decode(&uploadBody); err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(uploadBody.Errors) > 0 {
		return "", nil, fmt.Errorf("%v", uploadBody.Errors)
	}

	uri, err := url.Parse(uploadBody.Data.UploadURI)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing upload URI")
	}

	return uploadBody.Data.ID, uri, nil
}

// GetSupportBundleChannelUploadURI queries the Endpoint in Client to retrieve a URI that can be used to upload
// a support bundle for a specific channel
func (c *Client) GetSupportBundleChannelUploadURI(channelID string, size int64, notes string) (string, *url.URL, error) {
	resp, err := c.executeGraphQLQuery("", Request{
		Query: startChannelUploadMutation,
		Variables: map[string]interface{}{
			"channelId": channelID,
			"size":      size,
			"notes":     notes,
		},
	})

	if err != nil {
		return "", nil, errors.Wrap(err, "executing graphql request")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	uploadBody := SupportBundleChannelUploadResponse{}

	if err := decoder.Decode(&uploadBody); err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(uploadBody.Errors) > 0 {
		return "", nil, fmt.Errorf("%v", uploadBody.Errors)
	}

	uri, err := url.Parse(uploadBody.Data.UploadURI)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing upload URI")
	}

	return uploadBody.Data.ID, uri, nil
}

// GetSupportBundleWatchUploadURI queries the Endpoint in Client to retrieve a URI that can be used to upload
// a support bundle for a specific watch
func (c *Client) GetSupportBundleWatchUploadURI(watchID string, size int64) (string, *url.URL, error) {
	resp, err := c.executeGraphQLQuery("", Request{
		Query: startWatchUploadMutation,
		Variables: map[string]interface{}{
			"watchId": watchID,
			"size":    size,
		},
	})

	if err != nil {
		return "", nil, errors.Wrap(err, "executing graphql request")
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	uploadBody := SupportBundleWatchUploadResponse{}

	if err := decoder.Decode(&uploadBody); err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if len(uploadBody.Errors) > 0 {
		return "", nil, fmt.Errorf("%v", uploadBody.Errors)
	}

	uri, err := url.Parse(uploadBody.Data.UploadURI)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing upload URI")
	}

	return uploadBody.Data.ID, uri, nil
}

// MarkSupportBundleUploaded mutates the Endpoint in Client to mark a support bundle as uploaded
func (c *Client) MarkSupportBundleUploaded(bundleID string) error {
	_, err := c.executeGraphQLQuery("", Request{
		Query: markSupportBundleUploadedMutation,
		Variables: map[string]interface{}{
			"id": bundleID,
		},
	})

	return err
}

func (c *Client) executeGraphQLQuery(auth string, gr Request) (*http.Response, error) {
	body, err := json.Marshal(gr)

	if err != nil {
		return nil, errors.Wrap(err, "marshalling graphql request")
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest("POST", c.endpoint, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	req.Header.Set("Content-Type", "application/json")

	if auth != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth+":"))))
	}

	return c.client.Do(req)
}
