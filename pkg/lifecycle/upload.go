package lifecycle

import (
	"net/http"
	"net/url"
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

func UploadTask(task *types.LifecycleTask) Task {
	return func(l *Lifecycle) (bool, error) {
		if l.UploadCustomerID == "" {
			return false, errors.New("upload with no customer id")
		}

		bundleID, url, err := l.GraphQLClient.GetSupportBundleUploadURI(l.UploadCustomerID, l.FileInfo.Size())

		if err != nil {
			return false, errors.Wrap(err, "get presigned URL")
		}

		err = putObject(l.FileInfo, url)
		if err != nil {
			return false, errors.Wrap(err, "uploading to presigned URL")
		}

		if err = l.GraphQLClient.UpdateSupportBundleStatus(l.UploadCustomerID, bundleID, "uploaded"); err != nil {
			return false, errors.Wrap(err, "updating bundle status")
		}

		return true, nil
	}
}

func putObject(fi os.FileInfo, url *url.URL) error {
	file, err := os.Open(fi.Name())
	if err != nil {
		return errors.Wrap(err, "opening file for upload")
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", url.String(), file)
	if err != nil {
		return errors.Wrap(err, "making request")
	}
	req.ContentLength = fi.Size()
	req.Header.Set("Content-Type", "application/tar+gzip")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "completing request")
	}

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("Error uploading support bundle, got %s", res.Status)
	}
	return nil
}
