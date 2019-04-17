package reader

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/spf13/afero"
)

type MultiBundleReader interface {
	GetBundles() map[string]BundleReader
}

type MultiBundle struct {
	Fs      afero.Fs
	Path    string
	bundles map[string]BundleReader
}

func NewMultiBundle(fs afero.Fs, path string) (MultiBundleReader, error) {
	b := &MultiBundle{
		Fs:   fs,
		Path: path,
	}
	if err := b.discoverBundles(); err != nil {
		return b, errors.Wrap(err, "discover bundles")
	}
	return b, nil
}

func (b *MultiBundle) GetBundles() map[string]BundleReader {
	return b.bundles
}

func (b *MultiBundle) discoverBundles() error {
	b.bundles = map[string]BundleReader{}

	if _, err := b.Fs.Stat(b.Path); os.IsNotExist(err) {
		return err
	}

	file, err := b.Fs.Open(b.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		if header == nil || header.Typeflag != tar.TypeReg {
			continue
		}

		parts := strings.Split(header.Name, "/")
		if parts[len(parts)-1] != "index.json" {
			continue
		}

		var results []collecttypes.Result
		if err := json.NewDecoder(tr).Decode(&results); err != nil || len(results) == 0 && results[0].Path == "" {
			continue
		}

		// It seems safe to assume this is the root of a bundle
		prefix := filepath.Join(parts[:len(parts)-1]...)
		bundle, err := NewBundle(b.Fs, b.Path, prefix)
		if err != nil {
			return errors.Wrapf(err, "new bundle at prefix %s", prefix)
		}
		b.bundles[prefix] = bundle
	}
}
