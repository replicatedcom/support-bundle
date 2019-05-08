package cli

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

func bundleFromStdin() (string, error) {
	f, err := ioutil.TempFile("", "support-bundle")
	if err != nil {
		return "", errors.Wrap(err, "create temp file")
	}
	if _, err := io.Copy(f, os.Stdin); err != nil {
		return "", errors.Wrap(err, "read from stdin")
	}
	return f.Name(), nil
}
