// Code generated by go-bindata.
// sources:
// pkg/analyze/api/v1/requirements/kubernetes_total_memory.yml
// pkg/analyze/api/v1/requirements/kubernetes_version.yml
// DO NOT EDIT!

package v1

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _requirementsKubernetes_total_memoryYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x52\x4d\x6b\xdb\x50\x10\xbc\xe7\x57\x0c\xbe\xa4\x39\x44\x3f\xc0\x87\x1e\x02\x0d\x34\x25\xd0\xd6\xb9\x95\x52\xd6\xd2\x3a\x59\x78\x1f\xf6\xee\x3e\x17\x21\xeb\xbf\x17\x59\x8e\x64\x97\x16\x4a\x89\x4e\x7a\xfb\x66\x76\x66\x67\x9f\xd2\xcf\xe5\x15\x50\xe7\xd4\x88\x4b\x4e\x36\x9c\x80\x5b\xf0\x9e\xc2\xf8\x3f\x7c\x79\xcb\x4a\x9e\x75\x09\x4a\xcd\x54\x35\x27\xe7\xc8\xc9\x6d\x46\x0e\xdc\xeb\xae\x53\xde\x06\x54\x1f\xe2\xd6\xdb\xbe\xbf\x9e\x6e\xf7\xa4\x42\xeb\xc0\xbf\x11\x8e\xb8\xe5\xc4\x4b\xd9\x51\xfd\xf8\xca\x9b\x99\x6a\xbc\x67\x95\x01\xc4\xaa\x59\x4f\xd5\xc8\x66\xf4\xcc\x4b\xbc\xbf\x9d\xfa\xdd\x93\x04\x6e\xe0\x19\x0d\x3b\xd7\x8e\x4f\x65\xcd\x9a\xd8\xd9\xe0\xd9\x29\x20\x72\xcc\xda\xbe\xe1\x9c\xc1\xf1\xee\x33\xa9\xf1\xc7\xe4\xa8\x1e\x8f\xed\x6f\x2e\x4a\x92\x6e\xfe\x21\x86\x91\x39\xe7\xf0\xb0\x1b\x63\xc0\xa2\x12\xe7\x68\xdf\xbe\x57\x83\x95\x62\x55\x4d\x5b\xaa\xc5\xdb\x6a\x9c\x65\x81\x03\x1e\x76\x8f\x63\x1a\xf6\x94\x57\xae\x92\x9e\x57\x41\x6a\xc6\x01\xab\x12\xbf\x14\x4a\x2e\x2e\x6c\xe7\x36\x8e\x9a\x92\x66\xc1\xa3\xe3\x13\xb6\xc5\xa2\xeb\xba\xc1\x7a\xdf\xf7\x8b\xff\xd9\xc4\xd3\x59\xda\xd8\x64\x45\x9b\x8b\x9e\xaf\xa3\x0e\xc5\x9c\x15\x27\xf5\x3f\x64\x78\xc0\x7d\xd6\x48\xfe\xea\xa9\xef\xef\x20\x86\xc0\x66\xf0\x17\x4a\xf0\x17\x46\x94\x24\xb1\xc4\x49\x57\x79\x57\x44\xb9\x79\x95\xce\x1b\xcc\x93\xdc\x5d\x4d\x56\xff\xf2\x30\x10\x8b\x39\xd6\x0c\x72\x04\x26\xf3\x0b\x36\x7e\x05\x00\x00\xff\xff\x5a\x74\xd1\x6f\x32\x03\x00\x00")

func requirementsKubernetes_total_memoryYmlBytes() ([]byte, error) {
	return bindataRead(
		_requirementsKubernetes_total_memoryYml,
		"requirements/kubernetes_total_memory.yml",
	)
}

func requirementsKubernetes_total_memoryYml() (*asset, error) {
	bytes, err := requirementsKubernetes_total_memoryYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "requirements/kubernetes_total_memory.yml", size: 818, mode: os.FileMode(420), modTime: time.Unix(1532810999, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _requirementsKubernetes_versionYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x91\x31\x6f\xe3\x30\x0c\x85\xf7\xfc\x8a\x07\x2f\xb9\x1b\x62\xdc\xec\xe1\x96\xa2\x1d\x52\x64\x69\x82\x02\x9d\x0a\x25\x7e\x71\x04\x58\x92\x43\xd1\x2e\x0a\xd7\xff\xbd\x70\x9d\x2a\x40\x9a\xa1\x28\xaa\x89\xa4\xf8\x9e\xc8\x4f\x62\x5e\x8a\x19\xb0\x0b\xbe\xb4\x6a\x83\x8f\x63\x06\x2c\xc0\xce\xd4\x53\x3c\x9e\xd0\x50\x8c\x06\x29\x60\x7c\x99\xaa\x51\x8d\xd2\xd1\x6b\x3c\x77\x8e\xda\x79\xdf\x0b\x9b\x1a\xf9\xad\x6b\xf4\x75\x18\xe6\xe9\xb6\x33\x62\xcd\xb6\xe6\x85\xe0\xa3\xaf\x48\x3a\x1f\x14\xf9\xf3\x03\xf7\x67\x69\x64\x47\xb1\x63\x13\x45\x82\x9c\xaa\x8e\x31\x9a\x8a\x05\xfe\x2f\x92\xdf\x9d\xb1\x35\x4b\x68\x40\x49\xe5\x4e\x71\xdf\x6e\x29\x9e\xca\x88\x8e\x12\x6d\xf0\xbf\xb8\xe2\x9a\xae\xa3\xdc\x04\xd7\x18\x21\xf2\x29\x45\xd6\xf7\xfd\x29\x5e\x59\x3f\x0c\x43\x86\x37\x54\xc4\xbf\x6f\xb0\x98\x64\x67\x18\x1b\xb1\x0e\x7f\x96\xc7\x89\x08\xb2\xbc\xb2\xfa\x38\xed\x31\xba\x2e\x8f\xab\x09\x42\xdc\x84\xb5\x8a\xf5\xd5\x5f\x64\x5d\xf6\x13\x72\x4f\xa1\x95\x2b\xb4\xf0\xf9\x9b\xd3\x64\xc3\x00\x1b\x51\x33\x46\xe8\xc1\x78\xe8\x81\x70\xd6\x5b\xd7\xba\xe4\x24\x3c\xb6\x56\x58\x26\x8b\xb0\xc7\x25\x92\x59\x1a\xe1\xca\x93\xae\x8d\x8a\x2d\x61\x14\x35\x4d\xd4\xaf\xea\xf7\x00\x00\x00\xff\xff\x56\xf2\x63\x09\xb9\x02\x00\x00")

func requirementsKubernetes_versionYmlBytes() ([]byte, error) {
	return bindataRead(
		_requirementsKubernetes_versionYml,
		"requirements/kubernetes_version.yml",
	)
}

func requirementsKubernetes_versionYml() (*asset, error) {
	bytes, err := requirementsKubernetes_versionYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "requirements/kubernetes_version.yml", size: 697, mode: os.FileMode(420), modTime: time.Unix(1532811002, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"requirements/kubernetes_total_memory.yml": requirementsKubernetes_total_memoryYml,
	"requirements/kubernetes_version.yml":      requirementsKubernetes_versionYml,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"requirements": &bintree{nil, map[string]*bintree{
		"kubernetes_total_memory.yml": &bintree{requirementsKubernetes_total_memoryYml, map[string]*bintree{}},
		"kubernetes_version.yml":      &bintree{requirementsKubernetes_versionYml, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
