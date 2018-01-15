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

type templateOpts struct {
	BundlePath string
}

func UploadTask(task *types.LifecycleTask) Task {
	return func(l *Lifecycle) (bool, error) {
		tplOpts := &templateOpts{
			BundlePath: l.GenerateBundlePath,
		}
		proceed, err := askForConfirmation(l.SkipPrompts, task, tplOpts)
		if err != nil {
			return false, errors.New("confirm upload")
		}

		if !proceed {
			if task.Upload.Prompt.DeclineMessage != "" {
				return false, runTemplate(os.Stdout, "decline", task.Upload.Prompt.DeclineMessage+"\n", tplOpts)
			}
			return false, nil
		}

		if task.Upload.Prompt.AcceptMessage != "" {
			err = runTemplate(os.Stdout, "accept", task.Upload.Prompt.AcceptMessage+"\n", tplOpts)
		}

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

func askForConfirmation(skip bool, task *types.LifecycleTask, opts *templateOpts) (bool, error) {
	if skip || task.Upload.Prompt == nil {
		return true, nil
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		def := "[y/N]"
		if task.Upload.Prompt.Default {
			def = "[Y/n]"
		}
		var b bytes.Buffer
		if err := runTemplate(&b, "message", task.Upload.Prompt.Message, opts); err != nil {
			return false, errors.Wrap(err, "template message")
		}
		fmt.Printf("%s %s: ", b.String(), def)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, errors.Wrap(err, "prompt user")
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" {
			return task.Upload.Prompt.Default, nil
		}

		if response == "y" {
			return true, nil
		}

		if response == "n" {
			return false, nil
		}
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
