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
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
)

// Generate a new support bundle and write the results as an archive at pathname
func Generate(tasks []types.Task, timeout time.Duration, pathname string) (os.FileInfo, string, error) {
	var isURL bool
	var isStdout bool

	callbackURL, err := url.Parse(pathname)
	if err == nil && (callbackURL.Scheme == "http" || callbackURL.Scheme == "https") {
		isURL = true
		pathname = "/tmp/bundle.tar.gz"
	} else if pathname == "-" {
		isStdout = true
		pathname = "/tmp/bundle.tar.gz"
	}

	collectDir, err := ioutil.TempDir(filepath.Dir(pathname), "")
	if err != nil {
		return nil, pathname, errors.Wrap(err, "creating a temporary directory to store results failed")
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
		if r.Error != nil {
			resultsWithError = append(resultsWithError, r)
		} else {
			resultsWithOutput = append(resultsWithOutput, r)
		}
	}

	//write index and error json files
	indexJSON, err := json.MarshalIndent(resultsWithOutput, "", "  ")
	if err != nil {
		return nil, pathname, errors.Wrap(err, "marshalled index.json")
	}
	ioutil.WriteFile(filepath.Join(collectDir, "index.json"), indexJSON, 0666)

	errorJSON, err := json.MarshalIndent(resultsWithError, "", "  ")
	if err != nil {
		return nil, "", errors.Wrap(err, "marshalled error.json")
	}
	ioutil.WriteFile(filepath.Join(collectDir, "error.json"), errorJSON, 0666)

	comp := compressor.NewTgz()
	comp.SetTarConfig(compressor.Tar{TruncateLongFiles: true})
	// trailing slash keeps the parent directory from being included in archive

	if err := comp.Compress(collectDir+"/", pathname); err != nil {
		return nil, pathname, errors.Wrap(err, "compressed results directory")
	}

	file, err := os.Open(pathname)
	if err != nil {
		return nil, pathname, errors.Wrap(err, "open compressed file")
	}
	defer file.Close()

	if isURL {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		fw, err := w.CreateFormFile("file", pathname)
		if err != nil {
			return nil, pathname, errors.Wrap(err, "multipart form")
		}

		if _, err = io.Copy(fw, file); err != nil {
			return nil, pathname, errors.Wrap(err, "copy buffer")
		}

		w.Close()

		req, err := http.NewRequest("POST", callbackURL.String(), &b)
		if err != nil {
			return nil, pathname, errors.Wrap(err, "post to callback url")
		}

		req.Header.Set("Content-Type", w.FormDataContentType())

		client := &http.Client{}
		_, err = client.Do(req)
		if err != nil {
			return nil, pathname, errors.Wrap(err, "post request request")
		}

	} else if isStdout {
		_, err := io.Copy(os.Stdout, file)
		if err != nil {
			return nil, pathname, errors.Wrap(err, "copy to stdout")
		}
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, pathname, errors.Wrap(err, "file stats")
	}

	return fileInfo, pathname, nil
}
