package getter

import (
	"io/ioutil"
	"os"

	getter "github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
)

type Getter struct {
	dstDir string
}

func NewWithDefaults() (*Getter, error) {
	return New(os.Getenv("TMP_DIR"))
}

func New(dst string) (*Getter, error) {
	dstDir, err := ioutil.TempDir(dst, "analyze_getter")
	if err != nil {
		return nil, errors.Wrap(err, "create tmp dir")
	}
	return &Getter{
		dstDir: dstDir,
	}, nil
}

func (g *Getter) Get(src string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "get cwd")
	}
	dst, err := g.tmpFile()
	if err != nil {
		return "", errors.Wrap(err, "tmp file")
	}
	if err := getter.GetFile(dst, src, func(c *getter.Client) error {
		c.Pwd = pwd                                        // allow for a relative path
		c.Decompressors = map[string]getter.Decompressor{} // do not decompress
		return nil
	}); err != nil {
		return "", errors.Wrap(err, "get file")
	}
	return dst, nil
}

func (g *Getter) DstDir() string {
	return g.dstDir
}

// tmpFile generates a temporary file name
func (g *Getter) tmpFile() (string, error) {
	file, err := ioutil.TempFile(g.DstDir(), "bundle_*.tgz")
	if err != nil {
		return "", errors.Wrap(err, "create tmp file")
	}
	dst := file.Name()
	if err := os.Remove(dst); err != nil {
		return "", errors.Wrap(err, "remove tmp file")
	}
	return dst, nil
}
