package lifecycle

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/types"
)

type UploadTask struct {
	Options types.UploadOptions
}

type templateOpts struct {
	BundlePath string
}

func (task *UploadTask) Execute(l *Lifecycle) (bool, error) {
	tplOpts := &templateOpts{
		BundlePath: l.GenerateBundlePath,
	}
	proceed, err := task.askForConfirmation(l.SkipPrompts, tplOpts)
	if err != nil {
		return false, errors.New("confirm upload")
	}

	if !proceed {
		if task.Options.Prompt.DeclineMessage != "" {
			return false, runTemplate(os.Stdout, "decline", task.Options.Prompt.DeclineMessage+"\n", tplOpts)
		}
		return false, nil

	}

	if task.Options.Prompt.AcceptMessage != "" {
		err = runTemplate(os.Stdout, "accept", task.Options.Prompt.AcceptMessage+"\n", tplOpts)
		if err != nil {
			return false, errors.Wrap(err, "run accept template")
		}
	}

	if l.UploadCustomerID == "" {
		return false, errors.New("upload with no customer id")
	}

	bundleID, url, err := l.GraphQLClient.GetSupportBundleUploadURI(l.UploadCustomerID, l.FileInfo.Size())

	if err != nil {
		return false, errors.Wrap(err, "get presigned URL")
	}

	err = putObject(l.FileInfo, l.RealGeneratedBundlePath, url)
	if err != nil {
		return false, errors.Wrap(err, "uploading to presigned URL")
	}

	if err = l.GraphQLClient.UpdateSupportBundleStatus(l.UploadCustomerID, bundleID, "uploaded"); err != nil {
		return false, errors.Wrap(err, "updating bundle status")
	}

	return true, nil
}

func runTemplate(w io.Writer, name string, tpl string, opts *templateOpts) error {
	if tpl != "" {
		tmpl, err := template.New(name).Parse(tpl)
		if err != nil {
			return err
		}

		return tmpl.Execute(w, opts)
	}

	return nil
}

func (task *UploadTask) askForConfirmation(skip bool, opts *templateOpts) (bool, error) {
	if skip || task.Options.Prompt == nil {
		return true, nil
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		def := "[y/N]"
		if task.Options.Prompt.Default {
			def = "[Y/n]"
		}
		var b bytes.Buffer
		if err := runTemplate(&b, "message", task.Options.Prompt.Message, opts); err != nil {
			return false, errors.Wrap(err, "template message")
		}
		fmt.Fprintf(os.Stderr, "%s %s: ", b.String(), def)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, errors.Wrap(err, "prompt user")
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" {
			return task.Options.Prompt.Default, nil
		}

		if response == "y" {
			return true, nil
		}

		if response == "n" {
			return false, nil
		}
	}

}

func putObject(fi os.FileInfo, fullpath string, url *url.URL) error {
	file, err := os.Open(fullpath)
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
