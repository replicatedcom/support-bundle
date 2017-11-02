package bundle

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
)

type responseData struct {
	Bundle struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
}

func streamingUploadFile(file *os.File, backingWriter *io.PipeWriter, writer *multipart.Writer, data string) {
	defer backingWriter.Close()
	defer writer.Close()

	// add the file to the multipart form
	formFile, err := writer.CreateFormFile("file", file.Name())
	if err != nil {
		jww.ERROR.Fatal(err)
	}
	if _, err := io.Copy(formFile, file); err != nil {
		jww.ERROR.Fatal(err)
	}

	if err := writer.WriteField("data", data); err != nil {
		jww.ERROR.Fatal(err)
	}
}

// Upload uploads a support bundle in such a way that it becomes accessable to replicated support.
func Upload(file *os.File, fname, lname, email, company, description string) (string, error) {

	// bundle the data fields
	bundleData := struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		EMail       string `json:"email_address"`
		Company     string `json:"company"`
		Description string `json:"description"`
	}{
		FirstName:   fname,
		LastName:    lname,
		EMail:       email,
		Company:     company,
		Description: description,
	}
	bundleBytes, err := json.Marshal(bundleData)
	if err != nil {
		return "", err
	}

	read, write := io.Pipe()
	multipartWriter := multipart.NewWriter(write)
	contentType := multipartWriter.FormDataContentType()
	go streamingUploadFile(file, write, multipartWriter, string(bundleBytes))

	uploadURL := "https://support-bundle-secure-upload.replicated.com/v1/upload"
	req, err := http.NewRequest("POST", uploadURL, read)
	if err != nil {
		return "", err
	}

	// content type includes the boundary
	req.Header.Set("Content-Type", contentType)

	// submit request
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		err = errors.Errorf("Bad status when uploading support bundle: %s", res.Status)
		return "", err
	}

	response := json.NewDecoder(res.Body)
	respData := responseData{}
	err = response.Decode(&respData)
	if err != nil {
		return "", err
	}

	return respData.Bundle.ID, nil
}
