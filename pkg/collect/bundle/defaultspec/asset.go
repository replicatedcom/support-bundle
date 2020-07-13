// Code generated by go-bindata.
// sources:
// pkg/collect/bundle/defaultspec/assets/core.yml
// pkg/collect/bundle/defaultspec/assets/docker.yml
// pkg/collect/bundle/defaultspec/assets/kubernetes.yml
// pkg/collect/bundle/defaultspec/assets/replicated.yml
// DO NOT EDIT!

package defaultspec

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

var _assetsCoreYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x96\x4b\x6f\xdb\x30\x10\x84\xef\xfe\x15\x0b\x14\xe8\x29\x8a\x0e\x05\x7a\xd0\x39\xe8\xb1\x97\x1c\x8b\xc2\xa0\xc5\x95\x4c\x84\xe2\xb2\xdc\x65\xdc\x07\xfa\xdf\x0b\xea\xe5\x26\xf1\x83\xb6\x05\xb4\x97\x20\xb6\x66\xbe\x19\x51\x4b\x5a\xec\xb1\xe6\x6a\xb5\x02\x78\x07\x1c\xbd\xa7\x20\xc5\x26\x3a\x6d\x71\x05\x50\x40\x4b\x81\xa2\x18\x87\x5c\xc1\xaf\xdf\x2b\x00\x00\x8d\x5c\x07\xe3\xc5\x90\xab\xe0\x71\xb0\xc0\x60\xd9\xcb\x41\xc7\xce\xf7\x72\x8a\xe2\xa3\xac\xb5\x09\x15\x68\x6c\x54\xb4\x52\xbe\x09\x4a\xd4\x06\x43\x05\x12\xe2\x10\x6c\xa9\xcd\x8b\x4c\xc2\x6b\x83\x86\xbb\xfe\xc1\x82\x5d\x1f\x4a\x7c\x1f\xa2\x2b\x6a\xea\x3a\xe5\x74\xd5\xeb\x01\x9c\xea\xb0\x02\xad\x04\xc7\x2f\x0e\x25\x8d\x1e\x2e\x47\xdd\x19\x58\x33\x7e\x54\x21\xdd\xe6\x97\x42\xd9\xaf\x59\xf0\xe6\x1a\xb4\xc9\x64\xaf\x8d\x23\x8d\x7c\x3e\xa2\x43\x6e\xb3\x90\xa3\xf0\x34\xae\x09\x88\xaf\x3a\x77\x59\x8d\x47\xe3\x69\xba\xf1\xaf\xd8\x74\x07\x4a\xeb\x70\x07\xbc\xa5\x5d\x56\x8e\xf1\xeb\xe4\x58\x27\xc3\x55\x79\xd6\xb8\xa7\xcb\xf2\x92\xe3\xfa\xbc\xb4\x09\xf1\xb2\xc0\xde\x92\x97\xe8\xf9\x65\x62\xa3\xe2\xf7\xdd\x2e\x2f\xc9\xcf\xd3\x65\x49\x69\xf5\xdc\x56\x39\xae\x51\x3b\x59\xb7\xc4\xd2\x57\xc9\xf1\x4e\xe2\xc9\x1c\xbd\x98\x4c\xeb\x20\x9d\xd7\x03\x95\x2e\x1a\x63\x67\x6f\xfa\xdf\x2b\xd9\x56\x50\xa2\xd4\x65\xc3\xa2\x36\x27\xb0\x28\x75\x2e\xea\xaf\xce\x0b\xd1\x78\x19\x14\x71\x11\xd0\xa2\xe2\x85\xaa\x0d\x67\xef\xb2\xcc\x1a\x9d\x2c\xdd\x33\xcd\xc1\x4f\x72\x0b\xd0\x7c\xa0\xba\xac\x7d\x34\xae\xa1\x13\xb4\x24\xcb\xc6\x75\xd8\x2d\x8a\xa3\xe8\x4e\x0e\xcc\x45\xb4\x79\x13\x2d\x42\xdb\x9f\x04\x8b\xe0\x9e\x31\xb0\x21\xb7\x18\xae\x63\x51\x72\x8e\xd6\xbf\x75\x38\x94\x1d\x85\x27\xe3\x5a\x78\x0f\x8d\x09\xb8\x53\xd6\x5e\xb0\x6d\x6a\x72\x8d\x69\x4b\xe3\x45\x6d\x2c\x4e\x4f\xeb\xc5\x8b\xd2\x74\x0d\x66\x39\x30\x8a\x18\xd7\x9e\x39\x0d\xf6\xfc\xeb\x0b\x15\x33\xe0\x7f\xe9\xf5\xf1\x9f\xaf\x54\xe4\x50\x5a\xb3\x29\x3f\x8d\x8f\xfb\xe1\x50\x97\x69\x16\x74\x4e\x85\x91\x98\xfd\x13\x35\xb1\x6f\x0e\x4e\xa7\x5d\x3f\xc7\xda\x71\x6e\x7a\x40\x26\xfb\x7c\x9f\xd6\xeb\x50\xfe\xc3\xe7\x47\x78\x2b\xb9\xf1\xe0\xd6\x8e\x3b\xc5\xdf\x8e\x86\x1e\xb8\x7e\x6b\xe2\xb6\xf6\x1f\xd2\x5f\x6b\xd0\xc9\xf1\xe0\x24\x83\x43\xb2\x63\xa3\xd6\x3b\x56\x7f\x02\x00\x00\xff\xff\xa9\x83\xf5\xdb\xb0\x0d\x00\x00")

func assetsCoreYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsCoreYml,
		"assets/core.yml",
	)
}

func assetsCoreYml() (*asset, error) {
	bytes, err := assetsCoreYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/core.yml", size: 3504, mode: os.FileMode(420), modTime: time.Unix(1587562930, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x92\x51\x8e\xb4\x20\x10\x84\xdf\xe7\x14\x5c\x40\xc8\x9f\xfc\x4f\x5c\x66\x43\xa0\x75\x98\x45\x9a\x74\x37\x66\xbd\xfd\xc6\x11\x67\xcd\x64\xe3\xaa\x2f\x12\x53\xf5\x55\x49\x37\x17\xf0\x6c\x6f\x37\xa5\x3a\x15\xd0\x7f\x02\xe9\x98\x7b\xb4\x37\xf5\x7c\xb0\x4a\xa9\xf2\x11\x22\x59\x15\xa0\x77\x35\x89\x59\x65\x7b\xc7\x04\xc4\x11\xf3\x35\x53\x1c\xdd\x00\x5d\xe2\x6b\x2e\x8f\x59\x5c\xcc\x40\x3b\xa7\x4b\xc9\x2a\xa1\x0a\xe7\x48\xc8\x9a\xc0\x85\xae\x8f\x09\x36\xc4\x72\x2e\x4e\xee\x56\x19\x10\xdf\xd4\x26\x38\x18\x31\xeb\x07\x63\x6e\xba\x00\xec\x29\x16\x59\xfe\xb6\x55\x52\xab\x4a\xed\x54\xbf\x15\xf8\xc1\x9e\x2e\xf1\x5e\xfd\x10\xbd\x9e\xcf\xb2\x79\x66\x8f\xb9\x8f\xc3\x39\xfa\x4b\x7e\x81\x2f\x30\x86\xf6\x6e\x21\x9a\x81\xa6\xe8\x41\x07\x73\x17\x29\x5d\x21\xfc\x9a\x97\x89\xf6\x7f\xc7\x1f\xd1\xce\x94\x9a\x1c\x99\x84\x83\xa9\x85\xc5\xd1\x76\xa9\x3a\xe1\x70\x90\xfd\x66\x7a\xe6\x3c\xb0\x52\x76\x29\x2c\xd6\xd7\x06\xd6\x1c\x65\xdb\x87\xf6\x89\x63\xf6\x60\x55\xf7\xef\xbf\x0a\x6e\xe6\x83\x94\x8d\xb8\x8d\xe2\x3b\x00\x00\xff\xff\x00\x95\x34\x8a\x93\x03\x00\x00")

func assetsDockerYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerYml,
		"assets/docker.yml",
	)
}

func assetsDockerYml() (*asset, error) {
	bytes, err := assetsDockerYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/docker.yml", size: 915, mode: os.FileMode(420), modTime: time.Unix(1548452059, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsKubernetesYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x55\x41\x6e\xdc\x30\x0c\xbc\xef\x2b\x88\xf4\x6c\x3b\x6d\x0f\x45\x05\xb4\xa7\xbe\xa0\x1f\x58\x28\x12\x37\x56\x2c\x8b\xae\x48\xa5\xc8\xef\x0b\xc9\x5a\xc7\x4e\x0e\x71\xd0\xed\x61\xab\xa3\x4c\xce\x0c\x67\x28\x98\x27\x34\xac\x0e\x07\x80\x0f\x30\xa4\x3b\xf4\x28\x07\x80\x06\x1e\x28\xc5\xa0\xbd\x6d\x3d\xdd\xb3\x3a\x40\x39\x29\x38\x51\xab\xb2\x7c\x22\x3e\x62\x64\x54\x20\x31\x61\xbd\xa3\x24\x53\x92\xa3\x75\x51\x81\xc5\x93\x4e\x5e\xba\xdc\x15\x03\x0a\x72\x57\x01\xba\x0c\x3d\x33\x07\x3d\x22\x4f\xda\x60\x01\x6f\xf8\x89\x05\xc7\xa2\xc3\x92\x19\x30\x6e\x54\x58\x64\x13\xdd\x24\x8e\x82\x82\xef\x4d\xbd\x05\xf8\x51\x4a\x21\x97\xc2\x89\x22\x68\xef\x41\x7a\x04\x43\x41\xb4\x0b\x18\x19\xbc\xce\xcc\x16\x7e\x3b\xe9\x97\x3e\x47\xed\xb3\xb8\x76\x22\xdb\x2e\x72\xbe\x6d\xe5\xe4\x23\x6e\x44\x4a\x72\x64\x34\x14\x2c\x2b\xf8\x78\x7b\x5b\x3f\x2d\x44\xc7\xac\xe1\x48\x45\xe2\x22\x7b\x6e\x65\xd1\xe3\xc4\x0a\x4e\xda\x33\xbe\xee\x73\x2c\xaf\xfb\xb4\xf7\x1b\x7b\x01\x4e\xce\x0b\xc6\x55\x09\xcc\xa3\xad\x2f\xb2\x7b\xef\x1b\x6d\x47\x6c\xb5\x61\x8e\x6e\x95\xcf\x32\x42\xe3\x42\x5e\x29\xd9\x19\x56\xad\xfe\x87\x79\x5d\xa1\xb9\xd5\x94\xf9\x69\x44\xa2\xe1\xbd\x2f\xe1\x27\xd1\xf0\x3c\x38\xd4\xb0\x5e\xaa\x58\xd3\xbb\x20\x98\xdf\x7b\x09\xb6\xab\x9c\xff\xd5\xba\xaf\x66\x7a\x23\x8a\x5c\xf9\xf7\x0b\xfe\x22\x83\x73\xa6\x57\x64\xc1\x66\x0d\x51\x8c\x2d\x7e\x10\xb7\xbd\xc8\xd4\x44\xfc\x95\x90\x17\x1b\x52\xf4\x0a\xf2\x07\xd5\x75\x9e\x8c\xf6\x3d\xb1\xa8\x4f\x9f\xbf\x7c\xed\x46\x94\xe8\xcc\x79\x03\x47\x94\x9e\xac\x82\xfb\xe5\x0f\xf2\x86\x96\xda\xde\x15\x05\x45\x4b\x2d\x2a\x72\x56\xc3\x46\x64\x4a\xd1\x60\x93\x3d\x3d\xeb\x1a\x5c\xb0\x19\x76\xf2\xf4\x34\x62\x90\xb3\x8c\xc5\x15\x05\x37\x37\x3b\x4d\xa9\xf8\xdc\x6d\xe1\xf6\x8a\xc0\xc7\x4b\xf1\x2f\x48\x7b\xa9\x1f\xe8\xee\x22\xc4\x15\x67\x2f\xed\x44\xf6\x22\xb4\x05\xe7\x4f\x00\x00\x00\xff\xff\x0c\x49\x2e\x65\xaa\x08\x00\x00")

func assetsKubernetesYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsKubernetesYml,
		"assets/kubernetes.yml",
	)
}

func assetsKubernetesYml() (*asset, error) {
	bytes, err := assetsKubernetesYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/kubernetes.yml", size: 2218, mode: os.FileMode(420), modTime: time.Unix(1548452059, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsReplicatedYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x92\x31\x73\x83\x30\x0c\x85\x77\x7e\x85\xee\x32\x53\x76\xcf\x1d\x7b\x97\xa5\x3b\xe7\xda\x4a\xaa\xc3\xd8\x9c\x24\xd2\xbf\xdf\x03\x4c\xea\x84\x0e\x29\xd9\xca\xc6\x93\x25\x7d\xef\xd9\x32\xa0\x13\x53\x55\x00\x07\x88\x56\xe9\x82\x15\x40\x0d\x3e\xb9\x0e\xf9\xc5\xa5\xa8\x96\x22\x72\x4d\x71\x3a\xa9\xa6\x82\xf9\xbb\x16\xda\x40\xa2\x6d\x1a\x94\x52\x94\xb5\x0a\x60\x43\x30\xa0\x3c\xe2\x55\x39\x51\x50\xe4\xe2\x08\x40\xb4\x3d\x96\xff\xd3\x62\xc6\x21\x90\xb3\x8a\x7e\x53\x50\xb6\x6e\x23\x0f\x8c\x7d\x47\x7a\xa7\x8a\x8e\x9e\x52\x16\xd3\xa8\xc3\xa8\xad\x27\x36\xe0\xf1\x64\xc7\xa0\xcd\xcf\x9a\x66\x71\xdd\x64\x83\xa5\xfb\x90\xce\xf2\x7f\x0d\x4f\xee\x96\x7b\x97\x2f\xcb\x7d\x69\x5c\xd4\xba\xae\x16\xe4\x0b\x39\xac\xcb\x18\x26\x07\x32\x58\x87\x66\x0b\xee\x51\x1c\xd3\x1c\x8c\x81\xb7\x74\x16\x38\x71\xea\x21\x8f\x11\xa0\x08\xfa\x89\x45\x1f\xcc\x7b\x1e\x83\x9e\x19\x6f\x84\x82\x6e\x07\xfb\x4d\xb6\x0f\x93\x2f\x5d\xcf\x70\xd7\xeb\x90\x66\x03\xfe\x0b\x72\xd6\xcc\x7a\xeb\xed\xa3\xa1\xdf\x05\x9d\xdb\xd7\x79\x3b\xd1\xf3\x94\xbf\x91\x8b\x5a\x15\xbf\x17\x7c\xe9\x7e\x92\x7b\x19\x92\x1f\xfc\x01\xba\xf1\x03\x39\xa2\xa2\xcc\xaf\xff\xfd\xf8\x7a\xac\xbe\x03\x00\x00\xff\xff\x9e\x86\x37\xa4\x07\x05\x00\x00")

func assetsReplicatedYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsReplicatedYml,
		"assets/replicated.yml",
	)
}

func assetsReplicatedYml() (*asset, error) {
	bytes, err := assetsReplicatedYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/replicated.yml", size: 1287, mode: os.FileMode(420), modTime: time.Unix(1569526063, 0)}
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
	"assets/core.yml":       assetsCoreYml,
	"assets/docker.yml":     assetsDockerYml,
	"assets/kubernetes.yml": assetsKubernetesYml,
	"assets/replicated.yml": assetsReplicatedYml,
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
	"assets": &bintree{nil, map[string]*bintree{
		"core.yml":       &bintree{assetsCoreYml, map[string]*bintree{}},
		"docker.yml":     &bintree{assetsDockerYml, map[string]*bintree{}},
		"kubernetes.yml": &bintree{assetsKubernetesYml, map[string]*bintree{}},
		"replicated.yml": &bintree{assetsReplicatedYml, map[string]*bintree{}},
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
