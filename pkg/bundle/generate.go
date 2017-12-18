package bundle

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/divolgin/archiver/compressor"
	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

// Generate a new support bundle and write the results as an archive at pathname
func Generate(tasks []types.Task, timeout time.Duration, pathname string) error {
	var isURL bool

	callbackURL, err := url.Parse(pathname)
	if err == nil && (callbackURL.Scheme == "http" || callbackURL.Scheme == "https") {
		isURL = true
		pathname = "/tmp/bundle.tar.gz"
	}

	collectDir, err := ioutil.TempDir(filepath.Dir(pathname), "")
	if err != nil {
		return errors.Wrap(err, "Creating a temporary directory to store results failed")
	}
	defer os.RemoveAll(collectDir)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	results := Exec(ctx, collectDir, tasks)

	// any result that wrote a file, whether it has an error or not
	var resultsWithOutput []*types.Result
	// any result with an error, whether or not it wrote a file
	var resultsWithError []*types.Result

	for _, r := range results {
		if r.Path != "" {
			resultsWithOutput = append(resultsWithOutput, r)
		}
		if r.Error != nil {
			resultsWithError = append(resultsWithError, r)
		}
	}

	//write index and error json files
	indexJSON, err := json.MarshalIndent(resultsWithOutput, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Marshalling index.json failed")
	}
	ioutil.WriteFile(filepath.Join(collectDir, "index.json"), indexJSON, 0666)

	errorJSON, err := json.MarshalIndent(resultsWithError, "", "  ")
	if err != nil {
		return errors.Wrap(err, "Marshalling error.json failed")
	}
	ioutil.WriteFile(filepath.Join(collectDir, "error.json"), errorJSON, 0666)

	comp := compressor.NewTgz()
	comp.SetTarConfig(compressor.Tar{TruncateLongFiles: true})
	// trailing slash keeps the parent directory from being included in archive

	if err := comp.Compress(collectDir+"/", pathname); err != nil {
		return errors.Wrap(err, "Compressing results directory failed")
	}

	if isURL {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		file, err := os.Open(pathname)
		if err != nil {
			return errors.Wrap(err, "finding the file that was just compressed")
		}
		defer file.Close()
		fw, err := w.CreateFormFile("file", pathname)
		if err != nil {
			return errors.Wrap(err, "creating multipart form")
		}

		if _, err = io.Copy(fw, file); err != nil {
			return errors.Wrap(err, "copying buffer")
		}

		w.Close()

		req, err := http.NewRequest("POST", callbackURL.String(), &b)
		if err != nil {
			return errors.Wrap(err, "making request")
		}

		req.Header.Set("Content-Type", w.FormDataContentType())

		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			return errors.Wrap(err, "completing request")
		}
	}

	return nil
}
