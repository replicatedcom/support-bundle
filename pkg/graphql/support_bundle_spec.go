package graphql

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const query = `
query {
  supportBundleSpec {
	spec
	hydrated
  }
}
`

var defaultRequest = Request{
	Query: query,
}

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
	body, err := json.Marshal(defaultRequest)

	if err != nil {
		return nil, errors.Wrap(err, "marshalling graphql request")
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest("POST", endpoint, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(id+":"))))

	resp, err := c.client.Do(req)
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
