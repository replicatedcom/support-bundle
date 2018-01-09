package graphql

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const query = `
query {
  supportBundleSpec {
    spec
  }
}
`

var defaultRequest = Request{
	Query: query,
}

type SupportBundleSpec struct {
	CustomerID string
}

func (s *SupportBundleSpec) Get() ([]byte, error) {
	if endpoint == "" {
		panic("endpoint not defined, this should have been defined at build, define the environment variable REPLICATED_API_ENDPOINT to override")
	}

	if s.CustomerID == "" {
		return nil, errors.New("customer id must be defined to fetch a support bundle spec")
	}

	body, err := json.Marshal(defaultRequest)

	if err != nil {
		return nil, errors.Wrap(err, "marshalling graphql request")
	}

	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequest("GET", endpoint, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s:", base64.StdEncoding.EncodeToString([]byte(s.CustomerID))))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "executing graphql request")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	specBody := SupportBundleResponse{}

	if err := json.Unmarshal(respBody, &specBody); err != nil {
		return nil, errors.Wrap(err, "unmarshalling graphql response")
	}

	if specBody.Errors != nil && len(specBody.Errors) > 0 {
		return nil, errors.New(fmt.Sprintf("%v", specBody.Errors))
	}

	return []byte(specBody.Data.Spec), nil
}
