package collector

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
)

type MockCollector struct {
	mock.Mock

	Fs         afero.Fs
	BundlePath string
}

func NewMock(fs afero.Fs, bundlePath string) *MockCollector {
	return &MockCollector{
		Fs:         fs,
		BundlePath: bundlePath,
	}
}

func (c *MockCollector) CollectBundle(ctx context.Context, input Input) error {
	c.Called(input.Specs, input.SpecFiles)

	f, err := c.Fs.Create(input.Dest)
	if err != nil {
		return errors.Wrapf(err, "create file %s", input.Dest)
	}
	err = func() error {
		defer f.Close()

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := os.Chdir(c.BundlePath); err != nil {
			return err
		}
		defer os.Chdir(cwd)

		var filePaths []string
		files, err := ioutil.ReadDir(c.BundlePath)
		if err != nil {
			return err
		}
		for _, info := range files {
			filePaths = append(filePaths, info.Name())
		}

		return archiver.TarGz.Write(f, filePaths)
	}()
	return errors.Wrapf(err, "create archive from %s", c.BundlePath)
}