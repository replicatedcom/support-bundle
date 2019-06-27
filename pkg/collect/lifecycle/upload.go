package lifecycle

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
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

	var outputWriter io.Writer = os.Stdout
	if l.Quiet {
		outputWriter = ioutil.Discard
	}

	if l.DenyUploadPrompt {
		if task.Options.Prompt.DeclineMessage != "" {
			return false, runTemplate(outputWriter, "decline", task.Options.Prompt.DeclineMessage+"\n", tplOpts)
		}
		return false, nil
	}

	proceed, err := task.askForConfirmation(l.ConfirmUploadPrompt, tplOpts)
	if err != nil {
		return false, errors.New("confirm upload")
	}

	if !proceed {
		if task.Options.Prompt.DeclineMessage != "" {
			return false, runTemplate(outputWriter, "decline", task.Options.Prompt.DeclineMessage+"\n", tplOpts)
		}
		return false, nil

	}

	if task.Options.Prompt != nil && task.Options.Prompt.AcceptMessage != "" {
		err = runTemplate(outputWriter, "accept", task.Options.Prompt.AcceptMessage+"\n", tplOpts)
		if err != nil {
			return false, errors.Wrap(err, "run accept template")
		}
	}

	if l.UploadCustomerID == "" && l.UploadChannelID == "" && l.UploadWatchID == "" {
		return false, errors.New("upload with no watch id, channel id, or customer id")
	}

	var bundleID string
	var url *url.URL

	if l.UploadChannelID != "" {
		channelBundleID, channelURL, err := l.GraphQLClient.GetSupportBundleChannelUploadURI(l.UploadChannelID, l.FileInfo.Size(), l.Notes)
		if err != nil {
			return false, errors.Wrap(err, "get presigned URL for channel upload")
		}

		bundleID = channelBundleID
		url = channelURL
	} else if l.UploadWatchID != "" {
		watchBundleID, watchURL, err := l.GraphQLClient.GetSupportBundleWatchUploadURI(l.UploadWatchID, l.FileInfo.Size())
		if err != nil {
			return false, errors.Wrap(err, "get presigned URL for watch upload")
		}

		bundleID = watchBundleID
		url = watchURL
	} else if l.UploadCustomerID != "" {
		customerBundleID, customerURL, err := l.GraphQLClient.GetSupportBundleCustomerUploadURI(l.UploadCustomerID, l.FileInfo.Size(), l.Notes)
		if err != nil {
			return false, errors.Wrap(err, "get presigned URL for customer upload")
		}

		bundleID = customerBundleID
		url = customerURL
	}

	fmt.Printf("url = %s\n", url)
	err = putObject(l.FileInfo, l.RealGeneratedBundlePath, url)
	if err != nil {
		return false, errors.Wrap(err, "uploading to presigned URL")
	}

	if err = l.GraphQLClient.MarkSupportBundleUploaded(bundleID); err != nil {
		return false, errors.Wrap(err, "mark support bundle uploaded")
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

func (task *UploadTask) askForConfirmation(autoconfirmed bool, opts *templateOpts) (bool, error) {
	if autoconfirmed || task.Options.Prompt == nil {
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
		fmt.Fprintf(os.Stdout, "%s %s: ", b.String(), def)

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
