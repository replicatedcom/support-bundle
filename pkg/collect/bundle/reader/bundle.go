package reader

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/replicatedcom/support-bundle/pkg/meta"
	"github.com/spf13/afero"
)

type BundleReader interface {
	ResultsFromRef(ref meta.Ref) []collecttypes.Result
	ReaderFromRef(ref meta.Ref) (io.ReadCloser, error)
	Open(name string) (io.ReadCloser, error)
	GetIndex() []collecttypes.Result
	GetErrorIndex() []collecttypes.Result
}

type Bundle struct {
	Fs     afero.Fs
	Path   string
	Prefix string
	index  []collecttypes.Result
	errors []collecttypes.Result
}

func NewBundle(fs afero.Fs, path, prefix string) (*Bundle, error) {
	b := &Bundle{
		Fs:     fs,
		Path:   path,
		Prefix: prefix,
	}
	var err error
	b.index, err = b.getResultsFromIndex("index.json")
	if err != nil {
		return b, errors.Wrap(err, "get results from index.json")
	}
	b.errors, err = b.getResultsFromIndex("error.json")
	if err != nil {
		return b, errors.Wrap(err, "get results from error.json")
	}
	return b, nil
}

func (b *Bundle) ResultsFromRef(ref meta.Ref) (results []collecttypes.Result) {
	for _, result := range b.GetIndex() {
		if meta.RefMatches(ref, result.Spec.Shared().Meta) {
			results = append(results, result)
		}
	}
	return
}

func (b *Bundle) ReaderFromRef(ref meta.Ref) (io.ReadCloser, error) {
	results := b.ResultsFromRef(ref)
	for _, result := range results {
		// TODO: sometimes we have stdout and stderr, how do we choose one?
		if result.Size > 0 {
			return b.Open(strings.TrimLeft(result.Path, "/"))
		}
	}
	// TODO: what should we return here?
	return nil, nil
}

func (b *Bundle) Open(name string) (io.ReadCloser, error) {
	var closeFns []func() error
	defer func() {
		for _, closeFn := range closeFns {
			closeFn()
		}
	}()

	if _, err := b.Fs.Stat(b.Path); os.IsNotExist(err) {
		return nil, err
	}

	file, err := b.Fs.Open(b.Path)
	if err != nil {
		return nil, err
	}
	closeFns = append(closeFns, file.Close)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	closeFns = append(closeFns, gzr.Close)

	if b.Prefix != "" {
		name = filepath.Join(b.Prefix, name)
	}

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil, &os.PathError{
				Op:   "stat",
				Path: name,
				Err:  os.ErrNotExist,
			}
		} else if err != nil {
			return nil, err
		}
		if header == nil || header.Typeflag != tar.TypeReg {
			continue
		}

		if header.Name == name {
			af := &archiveFile{
				tr:       tr,
				closeFns: closeFns,
			}
			closeFns = nil
			return af, nil
		}
	}
}

func (b *Bundle) GetIndex() []collecttypes.Result {
	return b.index
}

func (b *Bundle) GetErrorIndex() []collecttypes.Result {
	return b.errors
}

func (b *Bundle) getResultsFromIndex(indexPath string) ([]collecttypes.Result, error) {
	reader, err := b.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	var results []collecttypes.Result
	return results, json.NewDecoder(reader).Decode(&results)
}

type archiveFile struct {
	tr       *tar.Reader
	closeFns []func() error
}

func (f *archiveFile) Close() (err error) {
	for _, closeFn := range f.closeFns {
		errI := closeFn()
		if errI != nil {
			err = multierror.Append(err, errI)
		}
	}
	return
}

func (f *archiveFile) Read(p []byte) (int, error) {
	return f.tr.Read(p)
}
