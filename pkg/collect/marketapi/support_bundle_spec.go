package marketapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

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

// GetChannelSpec will query the endpoint set in the client with a Replicated ChannelID
// and will return a fully rendered spec, containing collect and lifecycle keys
func (c *Client) GetChannelSpec(channelID string) ([]byte, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/v1/support/channel/%s/spec", c.endpoint, channelID))
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Error getting channel spec: StatusCode=%v, %s", resp.StatusCode, body)
	}

	collectorSpec := ChannelCollectorsResponse{}
	if err := json.Unmarshal(body, &collectorSpec); err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}

	jsonSpec, err := yaml.YAMLToJSON([]byte(collectorSpec.Spec))
	if err != nil {
		return nil, errors.Wrap(err, "converting yaml to json")
	}

	return jsonSpec, nil
}

// GetSupportBundleChannelUploadURI queries the Endpoint in Client to retrieve a URI that can be used to upload
// a support bundle for a specific channel
func (c *Client) GetSupportBundleChannelUploadURI(channelID string, size int64, notes string) (string, *url.URL, error) {
	resp, err := c.client.Post(fmt.Sprintf("%s/v1/support/upload/url", c.endpoint), "application/json", nil)
	if err != nil {
		return "", nil, errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, errors.Wrap(err, "reading response body")
	}

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.Errorf("Error getting channel spec: StatusCode=%v, %s", resp.StatusCode, body)
	}

	uploadResponse := UploadURLResponse{}
	if err := json.Unmarshal(body, &uploadResponse); err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling response")
	}

	uri, err := url.Parse(uploadResponse.UploadURL)
	if err != nil {
		return "", nil, errors.Wrap(err, "parsing upload URI")
	}

	return uploadResponse.BundleID, uri, nil
}

// MarkChannelSupportBundleUploaded mutates the Endpoint in Client to mark a support bundle as uploaded
func (c *Client) MarkChannelSupportBundleUploaded(bundleID string, channelID string) error {
	payload := struct {
		ChannelID string `json:"channel_id"`
	}{
		ChannelID: channelID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return errors.Wrap(err, "marshalling payload")
	}

	resp, err := c.client.Post(fmt.Sprintf("%s/v1/support/bundle/%s/uploaded", c.endpoint, bundleID), "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Error getting channel spec: StatusCode=%v, %s", resp.StatusCode, body)
	}

	return nil
}
