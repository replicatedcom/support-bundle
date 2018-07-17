package analyzer

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"

	multierror "github.com/hashicorp/go-multierror"
	collecttypes "github.com/replicatedcom/support-bundle/pkg/collect/types"
	"github.com/spf13/afero"
)

type BundleReader struct {
	Fs afero.Fs
}

func (r *BundleReader) GetResultsFromIndex(archivePath, index string) ([]collecttypes.Result, error) {
	reader, err := r.FileReaderFromArchive(archivePath, index)
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

func (r *BundleReader) FileReaderFromArchive(archivePath, targetFile string) (io.ReadCloser, error) {
	var closeFns []func() error
	defer func() {
		for _, closeFn := range closeFns {
			closeFn()
		}
	}()

	file, err := r.Fs.Open(archivePath)
	if err != nil {
		return nil, err
	}
	closeFns = append(closeFns, file.Close)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	closeFns = append(closeFns, gzr.Close)

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		} else if err != nil {
			return nil, err
		}
		if header == nil {
			continue
		}

		if header.Name == targetFile && header.Typeflag == tar.TypeReg {
			af := &archiveFile{
				tr:       tr,
				closeFns: closeFns,
			}
			closeFns = nil
			return af, nil
		}
	}
}
