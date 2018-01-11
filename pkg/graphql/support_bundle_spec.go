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

const customerSpecQuery = `
query {
  supportBundleSpec {
	spec
	hydrated
  }
}
`

const startUploadMutation = `
mutation GetPresignedURI($size: Int) {
  uploadSupportBundle(size: $size) {
    uploadUri,
    supportBundle {
      id
    }
  }
}
`

const finishUploadMutation = `
mutation SetBundleStatus($id: String!, $status: String!) {
  updateSupportBundle(id: $id, status: $status) {
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

func (c *Client) GetCustomerSpec(id string) ([]byte, error) {
	resp, err := c.executeGraphQLQuery(id, Request{
		Query: customerSpecQuery,
	})

	if err != nil {
		return nil, errors.Wrap(err, "executing graphql request")
	}
	decoder := json.NewDecoder(resp.Body)
	specBody := SupportBundleResponse{}

	if err := decoder.Decode(&specBody); err != nil {
		return nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if specBody.Errors != nil && len(specBody.Errors) > 0 {
		return nil, fmt.Errorf("%v", specBody.Errors)
	}

	return []byte(specBody.Data.Hydrated), nil
}

func (c *Client) GetSupportBundleUploadURI(id string, size int64) (string, *url.URL, error) {
	resp, err := c.executeGraphQLQuery(id, Request{
		Query: startUploadMutation,
		Variables: map[string]interface{}{
			"size": size,
		},
	})

	if err != nil {
		return "", nil, errors.Wrap(err, "executing graphql request")
	}

	decoder := json.NewDecoder(resp.Body)
	uploadBody := SupportBundleUploadResponse{}

	if err := decoder.Decode(&uploadBody); err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if uploadBody.Errors != nil && len(uploadBody.Errors) > 0 {
		return "", nil, fmt.Errorf("%v", uploadBody.Errors)
	}

	uri, err := url.Parse(uploadBody.Data.UploadURI)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing upload URI")
	}

	return uploadBody.Data.ID, uri, nil
}

func (c *Client) UpdateSupportBundleStatus(customerID string, bundleID, status string) error {
	_, err := c.executeGraphQLQuery(customerID, Request{
		Query: finishUploadMutation,
		Variables: map[string]interface{}{
			"id":     bundleID,
			"status": status,
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
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth+":"))))

	return c.client.Do(req)
}
