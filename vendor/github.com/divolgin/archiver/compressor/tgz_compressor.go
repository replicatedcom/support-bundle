package compressor

import (
	"compress/gzip"
	"os"
)

type Compressor interface {
	SetTarConfig(config Tar)
	Compress(src string, dst string) error
}

func NewTgz() Compressor {
	return &tgzCompressor{}
}

type tgzCompressor struct {
	tar Tar
}

func (compressor *tgzCompressor) SetTarConfig(config Tar) {
	compressor.tar = config
}

func (compressor *tgzCompressor) Compress(src string, dest string) error {
	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	return compressor.tar.writeTar(src, gw)
}
