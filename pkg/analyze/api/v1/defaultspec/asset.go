// Code generated by go-bindata.
// sources:
// pkg/analyze/api/v1/defaultspec/assets/configManagement.yml
// pkg/analyze/api/v1/defaultspec/assets/cpuCores.yml
// pkg/analyze/api/v1/defaultspec/assets/datacenter.yml
// pkg/analyze/api/v1/defaultspec/assets/dockerContainers.yml
// pkg/analyze/api/v1/defaultspec/assets/dockerDevicemapperLoopback.yml
// pkg/analyze/api/v1/defaultspec/assets/dockerIccFalse.yml
// pkg/analyze/api/v1/defaultspec/assets/dockerLoggingDriver.yml
// pkg/analyze/api/v1/defaultspec/assets/dockerVersion.yml
// pkg/analyze/api/v1/defaultspec/assets/loadavg.yml
// pkg/analyze/api/v1/defaultspec/assets/memoryUsage.yml
// pkg/analyze/api/v1/defaultspec/assets/os.yml
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

var _assetsConfigmanagementYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x52\xc1\x6a\x1c\x31\x0c\xbd\xe7\x2b\x44\x2e\x81\x96\xec\x26\x57\x1f\x1b\xe8\xa5\xa4\x84\x10\x7a\x2d\x5a\x5b\xe3\x15\x78\x64\x63\x79\x96\x4c\xbf\xbe\xcc\x7a\x98\x78\x92\xc2\xee\xa1\x7b\xc9\x65\x60\xac\xa7\xa7\xf7\xa4\x87\x82\x61\xfc\x43\xe6\x0a\xe0\x70\x3f\x7d\x01\x6e\x41\xb0\x27\x03\x29\x47\x4b\xaa\xa4\x1b\x1b\xa5\x63\xdf\xa3\xa0\xa7\x9e\xa4\x6c\xd2\x90\x12\x95\x23\x1a\x20\xe0\x8e\x82\x9a\xf9\x0f\x80\x6d\x94\x1f\x34\x1a\xa8\xa8\xdf\x21\xfa\xb8\x14\x1d\x29\x67\x72\x4f\x51\xb9\x70\x14\x03\xd7\xf7\x77\xd7\x73\x35\x93\x67\x2d\x94\x7f\x61\x66\xdc\x05\x6a\x38\x17\x4d\xba\x3c\x01\x74\x1c\xe8\x11\x8b\xdd\x9b\xe6\x11\x20\x61\xd9\x3f\x93\xa7\xd7\xa4\xeb\xc2\x44\xe3\xa8\xc3\x21\x94\xad\x8d\x7d\x8f\xe2\x74\x9b\x74\xab\xc5\xc5\xa1\xac\xa0\xf9\xd8\xff\x80\xa9\x0c\x99\xde\xb3\xd4\xa2\x81\x9b\xcd\xd7\x6a\x11\xd0\x93\x14\xd8\x7c\xb9\x79\x87\x64\x71\xf4\x6a\xe0\x6e\x7e\xb6\x51\x5c\xb5\x7d\xb5\xe6\xfa\xe0\x62\x19\xb1\x1a\xf0\x46\x7f\x98\x57\xf4\x4c\x5d\xb3\x94\x9e\x54\xd1\xb7\x7b\x5b\x26\xbe\xe4\x61\xe5\x23\x65\xee\x31\x8f\x06\x1e\x8f\x57\x75\xf0\x6d\x84\xa7\xf6\xaa\xf5\x58\x05\x39\x18\x78\xd9\xb3\x82\x52\x3e\x50\x06\x56\xc8\x83\x08\x8b\x07\x9c\x3b\xaa\xba\xa6\x4f\xe9\x40\x99\xcb\x68\x80\xa5\x8b\x1f\xc5\x7c\xc7\xa0\xff\x56\xf3\x33\x9e\xad\x42\x62\x69\x94\xa4\x53\x4a\x1c\xed\x06\x7f\x6e\xbe\xed\x9e\xba\x53\xe9\x9e\x30\x9f\x3a\xdb\x93\xc1\x5b\x1b\x78\x5a\xc8\x05\x92\xdd\xd0\x5f\x38\xd7\x0f\x6f\xd7\x84\x73\x52\x3d\xe1\xa1\x2a\xfb\x5f\xa1\x3e\x4f\xc2\x3a\xd2\xa7\x64\xd4\x44\xff\x0d\x00\x00\xff\xff\x04\x35\x14\x02\xc1\x05\x00\x00")

func assetsConfigmanagementYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsConfigmanagementYml,
		"assets/configManagement.yml",
	)
}

func assetsConfigmanagementYml() (*asset, error) {
	bytes, err := assetsConfigmanagementYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/configManagement.yml", size: 1473, mode: os.FileMode(420), modTime: time.Unix(1556057947, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsCpucoresYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x6c\x90\xcd\x8a\xe3\x30\x10\x84\xef\x79\x8a\xc6\x97\x9c\xe2\x90\xab\xae\x61\xf7\xb2\xb0\x84\x65\xb3\x97\x65\x18\x64\xa9\x1c\xc4\xe8\x8f\x96\xe5\xc1\x13\xf2\xee\x83\x32\x8e\x11\x13\xfb\x60\x70\x95\xbb\xeb\xab\x96\x5e\xda\xe9\x03\x62\x43\x34\x1e\xca\x9b\x68\x47\x5e\x3a\x08\x52\x31\xab\xc0\x48\x6d\x4e\xf2\x82\xbb\x45\x64\x65\x07\x9b\xc4\xfc\x45\x64\x54\xf0\xbf\x30\x09\x0a\xe9\x55\xc5\xbc\xe8\x1a\xc9\x30\xf4\x29\x24\x33\x98\xe0\x05\x35\x87\x66\x36\x19\x17\x93\x06\xf0\x3f\xc9\x46\x76\x16\xd5\xb6\x47\xb4\xcf\x2e\x72\x50\x8b\x4e\x05\xe6\x58\x60\x04\x5d\x6f\xb3\x1c\x19\x2a\x78\xfd\xb5\x7f\xf9\xd5\x87\x41\x54\x73\x70\x71\x98\xaa\xa1\xf2\x8c\x73\xf0\x1f\xf4\xdf\xa3\x56\x16\x62\x94\x56\x50\x33\x70\xc6\xa3\x81\x43\x2a\x27\xa9\xc0\x6b\x96\x9f\xd2\x26\xd4\x0c\x91\x8d\x93\x3c\x09\x3a\xfb\x37\x1f\xde\x7d\x65\x69\x0c\xd2\x58\x41\xc7\x90\xad\x2e\xec\x45\x01\x3b\xe3\x51\xc8\x3a\x30\x85\x9e\x8e\xa7\x33\xdd\xdb\xb7\xf4\xc3\xa7\xcc\xa0\x29\x64\xa6\x14\xa1\x12\x19\xaf\x6c\xd6\x20\x49\x2a\x38\x27\xbd\xa6\x90\x5a\x86\xd4\xbb\xde\x58\xfc\x6f\xf6\xa5\xde\x5e\xc5\x6c\x7c\x1f\x9a\x97\xb6\x4a\x4f\x18\xc1\xa6\xdc\x47\xa3\xcb\x97\xc5\x59\xaa\xfc\xe5\xbc\xde\x64\x7b\xbd\x32\xa2\xa5\x76\x3e\xdf\xed\xb6\x5d\x69\xf5\xfb\xb9\x01\x99\x44\x4f\xb3\xab\x48\x05\x77\xf3\x19\x00\x00\xff\xff\x1c\x65\xff\x2c\x9f\x02\x00\x00")

func assetsCpucoresYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsCpucoresYml,
		"assets/cpuCores.yml",
	)
}

func assetsCpucoresYml() (*asset, error) {
	bytes, err := assetsCpucoresYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/cpuCores.yml", size: 671, mode: os.FileMode(420), modTime: time.Unix(1556057957, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDatacenterYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x95\x41\x4f\x1b\x3d\x10\x86\xef\xfc\x8a\x51\x2e\x7c\x9f\x0a\x41\x9c\x2a\x59\xea\x81\x2e\xb4\x42\x6d\x25\x94\x40\x7b\xa0\x14\x4d\xec\xc9\x66\x5a\xaf\xbd\xb5\xbd\x09\x21\xca\x7f\xaf\x36\x09\x89\x37\x6c\x68\x97\xd2\xdb\x5e\x22\xc5\xf3\xce\xeb\x19\x7b\x9e\x35\x1a\xd4\xd3\x7b\x12\x7b\x00\xe3\xe3\xf2\x17\xe0\x10\x0c\x66\x24\x40\x61\x40\x49\x26\x90\xeb\xe2\xc4\x2f\x42\x00\x1a\x07\xa4\xbd\x58\xfd\x03\x60\x69\xcd\x07\x9a\xc6\xea\xdb\x8d\x1a\x40\x91\x67\x47\xea\xc2\x7a\x0e\x6c\x8d\x80\xce\xeb\xce\x2a\xe8\x28\x65\x1f\xc8\x7d\x46\xc7\x38\xd0\x14\xb9\x3e\x94\x90\x17\x03\xcd\xf2\x44\x29\x47\x7e\xe3\x09\x30\x64\x4d\x9f\x30\xc8\x91\x88\x16\x01\x72\x0c\xa3\x1e\xa5\x74\x97\xfb\x6a\xa0\x74\x3c\x92\xd6\x04\x64\x43\xce\x1f\x39\xca\x35\x4b\x0c\xa4\x0e\x6d\x4e\x0e\x83\x75\x5f\xbb\xdf\xbd\x35\x8f\xb3\xd8\xf8\x9c\x64\x68\x92\x52\x23\x3d\x5c\xd9\xd4\xa4\xb8\x45\xc1\x09\xe6\xa1\x70\xb4\x5d\xf6\x32\x28\x60\xbf\x73\x71\xf5\xf6\xe3\x79\x72\x7b\x72\x7a\xda\x3b\xeb\xf7\xdf\xfc\x77\xfd\xad\x73\xf3\xea\xff\xce\xfe\x56\x02\x1b\x45\x77\x02\x8e\x1f\x1f\xa5\xb3\x63\x56\xe4\x22\xfd\x64\xc4\x72\x94\x68\x5b\xa8\xea\xb6\xe3\xd5\x85\xf4\x68\x58\x7f\x05\xb9\x23\x69\x8d\x5a\xde\xe8\x3a\xd5\xd8\x10\xfb\x50\x96\x87\xa9\x80\xd9\x7c\xaf\x81\x73\x8d\xad\x0f\x8e\x4d\x9a\xd8\x2c\xc7\xea\xf9\xd0\x4f\x01\xf1\xa8\x55\xcd\xab\xed\x66\xe4\x3d\xa6\xf1\x88\xad\x77\xba\x74\x45\xc5\x36\x77\x9c\xa1\x9b\x0a\x38\xf9\xd2\x8f\x96\x15\x05\x64\x2d\xe0\x72\xc4\x1e\xd8\x03\x9a\x52\x00\x6c\x7c\x40\x23\x29\x52\x7a\x1a\x93\xe3\xb2\x77\x36\x43\xbb\x0b\xa9\x74\x9d\xf3\x27\x48\xa5\xd1\x0e\x2d\x52\x2d\x52\x0d\x9c\x9b\x22\x15\x8f\xda\x4b\x23\xf5\x3e\x39\x7b\x0a\xa9\x32\xfe\x7c\xa2\xf0\xbe\x70\x4d\x98\x8a\xf5\x2d\x55\x2d\x55\xff\xf6\xa1\xaa\x0c\xdb\x8b\x3f\x55\x15\xf7\xfa\xc7\xaa\x94\x3c\x1f\xae\xc2\xfc\x30\x76\x62\x1a\xe0\x55\xcd\x68\x01\x6b\x01\xfb\x5b\xc0\xb6\x6c\x77\xf2\x56\xf6\x2e\xe0\x1a\x27\xfe\x20\x95\x74\xb0\x40\xef\x66\xe7\xde\xbf\xa3\x2f\x6e\xf1\x1d\x6a\x5f\x4f\xe0\xd5\xd6\xb4\x6f\x18\x4c\x6c\xa1\x55\x59\x7b\xb9\x42\x2e\x63\x43\xab\x86\xe1\xfc\x02\x70\xd9\x74\xb7\x16\x48\x45\x83\x22\x6d\xf0\x19\x68\x56\xc4\x86\x55\xd0\x56\x62\xe9\x0c\x76\xf8\x50\x11\xcc\x66\xe5\xf4\x42\xb7\x72\x39\xf3\xf9\x53\x95\xfe\x0a\x00\x00\xff\xff\x96\x50\x56\x86\x4d\x0e\x00\x00")

func assetsDatacenterYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDatacenterYml,
		"assets/datacenter.yml",
	)
}

func assetsDatacenterYml() (*asset, error) {
	bytes, err := assetsDatacenterYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/datacenter.yml", size: 3661, mode: os.FileMode(420), modTime: time.Unix(1556057975, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockercontainersYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x52\x4d\xaf\xd3\x30\x10\xbc\xf7\x57\xac\x7c\x09\x9f\x41\x15\x37\x5f\x38\x54\xe2\x82\x90\x50\x85\xb8\x50\xa8\xb6\xc9\x26\x35\x38\xeb\x68\xed\x84\x86\xaa\xff\xfd\x29\xcd\xc7\x73\xfb\xf2\xde\xa5\x52\x67\x67\x36\x33\xe3\x45\x46\xdb\xfd\x27\xbd\x02\x68\xd7\xfd\x2f\xc0\x7b\x60\xac\x48\x43\xee\xb2\xbf\x24\x69\xe6\x38\xa0\x61\x12\x9f\x66\xae\xe1\x70\xe5\x00\x58\x3c\x90\xf5\x7a\xfc\x07\x60\x32\xc7\x5f\xa8\x9b\x64\xfb\x59\x36\x33\x72\xf2\x46\x28\xff\xe6\xbc\x09\xc6\xb1\x06\xf5\x51\x8d\x43\xa1\xd2\xf8\x40\xf2\x03\xc5\xe0\xc1\x52\xb4\x77\x72\x23\x0d\xb3\xe1\x72\x33\xbb\x99\x19\x00\x85\xb1\xf4\x15\x43\x76\xd4\x11\x08\x50\x63\x38\x6e\xa9\xa4\x53\xed\x6f\x07\xfd\xd6\x0f\xa3\x4f\xc3\x85\xdb\xa5\x7f\xbc\xe3\x1b\x8a\x5c\x75\x1b\xac\x43\x23\x74\xaf\x1e\x86\x1a\x12\xf5\x68\x67\x3b\xf8\x53\x1a\xde\xbc\xfa\xf9\x5b\xbd\xdb\xf1\xaf\xb7\xaf\x93\x3b\xa1\xe1\x9c\x4e\x1a\xd6\x23\x5c\x0b\x65\x8e\xf3\xa1\x8e\x99\xca\x2e\xc4\x1f\xa4\xaa\x0e\x9d\x86\xf3\x25\xc2\xda\xb1\xa7\x2d\x15\xcf\x37\xb3\xb0\x9a\x5a\xb4\x1a\x54\x90\x86\xa6\xea\x2b\xf2\x1e\xcb\xb8\xf1\xd8\xd5\x67\xb4\xfe\x26\x7e\x2d\xa6\x42\xe9\x34\xa8\x4f\x2a\x82\x73\x0a\x68\xac\x86\x8d\x6b\x6c\xde\x27\xe8\x11\x92\xca\x30\x01\x37\xd5\x81\x04\x5c\x31\x39\x1d\x4f\x04\xb2\xa5\xb7\xf4\xd4\x92\x98\x3e\xf1\x3f\x94\xf8\x4d\xee\x0f\x0e\xe2\xa3\x63\xb7\x1f\x96\xce\xe3\x39\xc2\x77\x69\x96\x13\x24\xe7\xb3\x50\x6d\x21\x7d\x52\xe0\xe5\x92\x2c\x64\x7b\x89\x1f\x65\x99\x52\x26\x8b\xa1\xfa\x73\x5b\x3d\x04\x00\x00\xff\xff\xd9\xab\x97\xb3\x75\x03\x00\x00")

func assetsDockercontainersYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockercontainersYml,
		"assets/dockerContainers.yml",
	)
}

func assetsDockercontainersYml() (*asset, error) {
	bytes, err := assetsDockercontainersYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/dockerContainers.yml", size: 885, mode: os.FileMode(420), modTime: time.Unix(1556059014, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerdevicemapperloopbackYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xc4\x92\x4b\x6f\xd5\x30\x10\x85\xf7\xfd\x15\xa3\x6c\xca\x43\x5c\xa8\xd8\x79\xdb\x88\x0d\x4f\x55\x15\x1b\x5e\x9a\xc4\x93\x64\xa8\x63\x9b\xb1\x6f\x6e\x03\xe2\xbf\xa3\xa4\x21\xb2\xd3\x8b\x74\xc5\xa6\x9b\x48\x99\xc7\xf1\x77\x7c\x8c\x16\xcd\xf8\x93\xd4\x19\xc0\x70\x31\x7d\x01\x9e\x81\xc5\x9e\x14\x68\x57\xdf\x90\xec\x34\x0d\x5c\x53\x8f\xde\x93\xec\x38\xbc\x71\xce\x57\x58\xdf\xcc\xa3\x00\x06\x2b\x32\x41\x2d\x7f\x00\x5c\x3b\xfb\x9a\xc6\xbf\xdb\xdf\x0e\x1d\x1a\x5a\xbb\x9a\x02\x0b\xe9\x0f\x2e\x70\x64\x67\x15\x14\x2f\x8b\xa5\x29\xd4\x72\x88\x24\x1f\x51\x18\x2b\x43\x89\xe6\x0a\x24\x3c\x90\xac\x65\x80\x86\x0d\xbd\xc5\x58\x77\x2a\x29\x02\x78\x8c\xdd\x15\xb5\x74\xeb\x43\xde\x98\xa4\x9e\x2f\x60\x6c\x1b\xf7\x79\xf7\x3d\x38\x9b\x8d\xc8\xbc\x77\x89\x3e\xee\x85\xb6\xdb\x77\x4d\x05\xe7\x45\x39\xa3\x14\x0a\x9e\x14\x8f\x3e\x7d\x2d\xbe\x3c\x7d\x5c\x9c\x6f\x86\xd9\x6a\xba\x55\x70\x71\xcf\x86\xc9\x6f\xf0\x41\x8d\x94\x18\x71\x06\x9a\x11\xfe\x61\xe1\xc5\x52\xf6\x42\xb5\xb3\xfa\x2e\xba\x75\x34\x44\x61\xdb\x5e\xba\xde\x63\x7e\x10\xfd\x50\x90\x3e\x9e\xb5\x35\x2c\x11\x5f\x51\xb3\x09\xf5\x88\xbc\x75\x31\x13\xed\x7d\x1c\x15\xfc\xfa\x9d\xd4\x32\xbd\xcd\xed\xf6\x14\x02\xb6\xe9\x6b\x5a\xcf\x78\x85\x26\x64\xc0\x5e\xb8\x47\x19\x15\x94\x09\xf5\x04\x00\x6c\x8f\xa5\xa6\x29\x22\x1b\x05\xe5\x1c\x44\xe6\x75\xb1\x05\xef\xde\x5f\xa7\xcb\xd3\xe1\x0d\xb7\x89\x46\xa0\x81\x84\x27\x4b\x9a\xaa\x7d\x7b\x9f\xf2\x5a\xf6\x27\x40\x1e\x38\x76\xff\x89\x78\x22\xde\x01\xc5\x9e\xfd\x09\x00\x00\xff\xff\xdd\xaa\xad\x5e\x2e\x04\x00\x00")

func assetsDockerdevicemapperloopbackYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerdevicemapperloopbackYml,
		"assets/dockerDevicemapperLoopback.yml",
	)
}

func assetsDockerdevicemapperloopbackYml() (*asset, error) {
	bytes, err := assetsDockerdevicemapperloopbackYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/dockerDevicemapperLoopback.yml", size: 1070, mode: os.FileMode(420), modTime: time.Unix(1556057992, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockericcfalseYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x53\x4d\x6f\x13\x31\x10\xbd\xe7\x57\x8c\xf6\xd2\x16\x68\x50\xc5\xcd\xd7\x20\xa4\x0a\x81\x50\x41\xbd\xf0\x51\x39\xf6\x4b\x32\xe0\x1d\x9b\xb1\x37\x6d\x40\xfc\x77\xb4\x9b\x34\x75\xd2\x95\x2a\xf5\xb2\xd2\xce\x9b\xf7\xfc\xde\x8c\x6d\xc5\x86\xcd\x1f\x98\x09\xd1\xfa\xa2\xff\x12\x9d\x93\xd8\x16\x86\x7c\x74\xbf\xa0\x53\x6f\xd1\x46\x99\xb2\x73\x03\x4a\x14\xec\x1c\x21\x9b\xdd\x1f\x11\xbb\x28\xef\xb1\x31\x24\xf1\x66\xcb\xd9\x43\x1e\x99\x15\xfe\x53\xcc\x5c\x38\x8a\xa1\xe6\x4d\xb3\x03\x15\x4b\xce\x05\x7a\x6d\x95\xed\x3c\xa0\x12\x3c\x34\xb0\x2f\x13\x2d\x38\xe0\x83\x2d\x6e\x65\xaa\x22\x51\xb2\x65\x75\x85\x25\xee\x52\x3e\x04\x7a\xa9\xd7\x5b\x95\x1b\x96\x45\xfc\x36\xfd\x99\xa3\x1c\xb4\xe8\xc0\x9b\xd9\x54\x3a\xc5\x31\x7b\x0b\x1a\x3a\x69\x3e\x43\xd7\xd0\x6b\x68\xe6\x28\x8d\xa1\x17\xcd\xe9\xd7\x1f\xcd\xf7\x97\x67\xcd\xc9\x11\x87\xc5\xe3\xce\xd0\xc5\xa3\x34\x0f\x03\x7c\x7e\x94\x61\x15\xcf\x4f\xc1\xce\xf5\xde\x4f\x17\x36\x64\x9c\x3d\xe1\x3c\x29\x5c\x14\xbf\x5d\xdc\xbe\x55\x62\xa9\xf5\xd1\xa6\xb2\x31\xf4\xf7\x5f\x55\x5b\xef\x36\x7a\x85\xc5\xd1\x0e\x47\xf4\x72\x51\x96\xe5\x2c\xb6\xc9\x1e\x3a\xc7\x6f\x43\xcd\x60\xb4\x99\x8c\x2a\x3f\xcc\xb3\x45\xce\x76\x59\xdf\xa0\xfd\x41\xef\x7a\x81\x5a\x36\x29\xb7\x56\x37\x86\xde\x0e\xc6\xe8\x72\x36\x23\x48\x2f\xea\xab\x2e\x8f\x62\x39\x18\xba\x94\x02\x3d\x77\x51\x8a\x65\x81\x92\x8b\x6d\xdb\x09\x3b\xdb\x6b\x13\xe7\x11\x6a\xc6\x1a\xca\xfd\x50\x3c\xe6\xdd\xb2\x42\x8e\xdf\x0d\xd5\x6f\x67\x77\x4b\x6f\x57\x36\xe0\x71\x8c\x2f\xda\x8d\xa7\xf8\x18\xfb\x04\x23\xce\xef\xe3\x3d\x15\xc0\x73\x1e\x12\xbc\x22\x9b\x52\x0f\xa7\x28\x90\x92\xe9\x96\x43\xa0\x39\xa8\x1b\x12\x52\x89\x15\x17\xa3\x79\xa1\x1a\x75\xf2\x3f\x00\x00\xff\xff\x46\xe1\x46\xff\x51\x04\x00\x00")

func assetsDockericcfalseYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockericcfalseYml,
		"assets/dockerIccFalse.yml",
	)
}

func assetsDockericcfalseYml() (*asset, error) {
	bytes, err := assetsDockericcfalseYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/dockerIccFalse.yml", size: 1105, mode: os.FileMode(420), modTime: time.Unix(1556057738, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerloggingdriverYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x8c\x52\x4b\x6f\xd4\x40\x0c\xbe\xf7\x57\x58\xb9\x94\x47\xbb\x55\xc5\x6d\xae\x45\x3d\xf0\x12\xaa\x10\x97\x52\x2a\x6f\xc6\x3b\x75\x99\xd8\x91\x67\x36\x6c\x40\xfc\x77\x94\xdd\x10\x26\x61\x11\x5c\x12\xc5\x9f\xed\xef\x11\xa3\x60\xec\xbf\x91\x3b\x01\xe8\x2e\x87\x27\xc0\x39\x08\x36\xe4\xc0\x6b\xfd\x85\x6c\x15\x35\x04\x96\xf0\xd2\xb8\x23\x5b\x71\x7a\x95\x54\xae\x39\xd2\xbe\x17\x20\xe2\x9a\x62\x72\xe3\x17\x00\xd7\x2a\xaf\xa9\x77\x20\x7a\x1f\x35\xa4\x09\xf0\x94\xd8\xc8\xbf\xd7\xc4\x99\x55\x1c\x54\x2f\xaa\x11\x34\x0a\x9c\x32\xd9\x47\x34\xc6\x75\xa4\x62\xdd\x24\x66\x4f\x3f\x95\x01\x36\x1c\xe9\x2d\xe6\xfa\xc1\x15\x45\x80\x16\xf3\xc3\x0d\x05\xda\xb5\x69\x0e\x0c\xab\x2e\x0e\x96\xee\x59\x36\xfa\x69\xf5\x98\x54\x66\x2d\xb6\x9f\xbb\xc2\x36\x6f\x8d\x96\xd3\x07\xd0\xc1\x69\xf5\xa6\x0c\xa4\x72\xf0\xac\x7a\x72\xfb\xb9\xba\x7b\xfe\xb4\x3a\x5d\xcc\xb0\x78\xda\x39\xb8\x1c\xcb\xad\x51\xad\xe2\x0f\xf6\xa7\x56\xd1\x5c\x72\x51\xd3\xe6\xde\xc1\xf7\x1f\x45\xad\x1b\x73\xb9\xa1\xcd\x22\x89\x7f\xef\x4b\xd9\x58\xc2\x95\x36\x2d\x2e\x4d\xb1\x38\xb8\x1d\x52\x38\x1f\xc2\x3c\x7b\xd4\xad\x09\x46\x7f\xf7\x5f\xcc\x0d\xa5\x84\xa1\xfc\x55\x93\x96\x6b\x8c\x69\x46\xd5\x1a\x37\x68\xbd\x83\x89\xec\xe2\x17\x59\xd1\xe5\x29\x23\x47\x07\x63\xbe\x23\x1d\x70\xfa\x3d\x06\x6a\x70\x64\x32\x51\x47\xc6\x43\x6c\x9e\xd6\xdb\x50\x20\xcb\xeb\x84\xf2\x42\x67\xe7\x39\xa9\xff\x60\xdb\xe3\xe2\xdf\xa9\x14\x4a\xa2\x86\x3f\x8f\xf2\xef\x16\x44\xf3\x71\x1b\x67\x50\xa3\x0c\x68\xad\x31\x52\x3d\xbc\x25\x23\x0b\xd9\x5c\x5f\xe9\xf2\x2b\x9a\x9c\xfc\x0c\x00\x00\xff\xff\x9b\x8a\x14\x27\xba\x03\x00\x00")

func assetsDockerloggingdriverYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerloggingdriverYml,
		"assets/dockerLoggingDriver.yml",
	)
}

func assetsDockerloggingdriverYml() (*asset, error) {
	bytes, err := assetsDockerloggingdriverYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/dockerLoggingDriver.yml", size: 954, mode: os.FileMode(420), modTime: time.Unix(1556058389, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsDockerversionYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x92\x49\x8f\xdb\x30\x0c\x85\xef\xf3\x2b\x08\x5d\xa6\x0b\x3a\xc5\xf4\xa8\x4b\x0f\x53\xf4\x52\x14\x28\xa6\xc5\x5c\xba\x04\x8c\xfd\x92\xa8\x95\x25\x81\x92\x9d\xb8\x41\xfe\x7b\xe1\x78\xa9\x6c\x18\x93\x43\x00\xf3\x91\xd4\xf7\x48\xb2\x63\xdb\xfe\x85\xbe\x21\x6a\xee\xbb\x7f\xa2\x37\xe4\xb8\x82\xa6\xd2\x17\x7f\x20\x77\x0d\x24\x1a\xef\xae\x12\x91\xe5\x2d\x6c\xd4\xc3\x17\x91\x29\xbc\xfb\x84\x76\xcc\xde\x1c\x0f\x6c\x31\xa9\x25\xa2\x11\x94\x5f\x7c\x34\xc9\x78\xa7\x49\xbd\x53\x83\x28\xd8\x9b\x98\x20\x4f\x2c\x86\xb7\x16\x59\xcf\x11\x60\xfe\x72\xf7\xdb\x19\x8b\xcf\x9c\x8a\x83\xce\x82\x44\x81\xd3\xe1\x11\x7b\x9c\x42\x9c\x0b\x5d\xaf\xb7\x03\x99\x71\x3b\xff\xe3\xee\x77\x9c\x35\xbc\x72\xe0\x14\x1e\x38\xa4\x5a\xb0\xac\xee\x45\x4d\xb7\xea\x2b\xa4\x81\x3c\xf5\x44\x4a\xd3\x2b\xf5\xe2\xfb\x2f\xf5\xf3\xf5\x4b\x75\xbb\xa8\x31\xae\xc4\x49\xd3\xfd\x10\x0e\x82\xc2\xbb\xb2\xf7\x3f\xa5\x3a\x9f\xf2\xb7\x50\x85\xd4\x6a\x3a\x5f\xb2\x58\x33\x0c\xe6\x11\xbb\xe5\x28\x56\x1a\xa2\x61\xab\x49\x25\xa9\x31\x4e\xb8\x42\x8c\xbc\xcf\x07\x9b\xb3\x7c\x64\x1b\x67\x7e\x83\x98\x8a\xa5\xd5\xa4\xde\xab\x2c\x5c\x22\xb1\xb1\x9a\x1e\x7c\x6d\xcb\x8e\xbb\x8b\x40\x2a\xe3\x30\xec\x7c\x65\x4f\x11\x0d\xc4\x74\x96\x8e\x2c\xb9\xb0\x3c\x1f\xca\x4f\xc8\xf9\x4d\xdf\x71\x92\x27\xda\x6f\x52\xaf\xc3\x7e\xe8\x11\xce\x67\x41\xb0\x34\x1e\xeb\xe5\xb2\xe2\x60\x48\x8d\xd7\x4d\x8e\xd0\x64\xe2\xb3\xc5\xff\x9d\x74\xf7\x73\xf3\x2f\x00\x00\xff\xff\x4d\xc0\x38\x98\x2f\x03\x00\x00")

func assetsDockerversionYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsDockerversionYml,
		"assets/dockerVersion.yml",
	)
}

func assetsDockerversionYml() (*asset, error) {
	bytes, err := assetsDockerversionYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/dockerVersion.yml", size: 815, mode: os.FileMode(420), modTime: time.Unix(1556058784, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsLoadavgYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x92\x41\x8f\xd3\x30\x10\x85\xef\xfb\x2b\x9e\x22\xa1\xbd\xec\x46\x32\x52\x24\xe4\x2b\x82\x0b\x97\x15\x62\xb9\xb2\x6e\x3c\x09\x16\xf6\x38\xb2\x9d\xa0\x50\xfa\xdf\x91\xdb\xa6\xb8\x6d\x28\x07\xf0\x21\x52\x66\xbe\x79\x33\xf3\x6c\xc5\xca\xce\x3f\x48\xde\x01\x93\xc8\x5f\xe0\x11\xac\x1c\x49\xf8\x58\x5b\xaf\xb4\x9a\xfa\x7d\x18\xb0\x6a\x43\x36\xca\xe3\x1f\x60\x5a\xcf\x1f\x68\xce\xe4\x97\x73\x12\xd0\x14\x4d\x20\xfd\xe4\xa3\x49\xc6\xb3\x44\xf5\xa6\x3a\x26\x03\xf5\x26\x26\x0a\x9f\x55\x30\x6a\x63\xa9\x50\x5c\x5a\x5f\xaa\x61\x89\x48\x6c\x77\x7f\xa2\xc5\x35\x2e\x6e\xf1\xcd\x35\xdf\xdc\xd4\x5f\x29\x10\x65\x45\xeb\x59\x1f\xb6\x3d\x81\x8a\xb5\x2c\xaa\x1e\xc1\x3e\x95\x81\x7c\xc8\x0d\x69\x3e\x6b\x7c\x38\xd3\xd1\x9f\x8f\xd4\xad\xee\xf8\xcf\x62\xcd\x7f\x9d\x6c\x51\x73\x14\xa3\xea\xcb\x5b\x3d\xf9\xf2\x5e\xd9\x48\x65\x93\x21\x18\xa7\xc2\x2c\xf1\xcc\xdf\xd8\x7f\xe7\x22\xa5\x29\x29\x63\x25\xde\xfa\xd1\xea\x3c\x5c\x8e\x50\x70\x86\x69\x69\x59\xe3\x1d\xc7\x31\x10\x66\x3f\x06\xc4\x81\xda\x08\xc3\xad\x1d\x35\x21\x7d\x25\xb4\xde\x39\xc5\x1a\x2f\xbf\x5f\xf2\x4b\x5d\xb4\x88\x34\x51\x30\x79\x43\x4d\x9b\xb1\xbf\x9e\xf7\x53\x18\xd7\xc7\xbd\xdf\x6e\x03\x0d\x16\x8b\xac\xc0\xcf\x9c\xe4\xd4\xa1\x7a\x55\xbf\xee\xaa\xdd\x0e\x17\x48\xf3\x77\x44\xac\x30\xf7\x6b\x96\x3c\x3d\xef\x2d\x80\x9a\x28\x64\xab\xe1\x27\x0a\x10\x0f\x68\x1e\xf2\x8b\x83\x68\xe0\x0c\x8f\x89\xe2\xea\xb6\x86\x3b\x7f\xf7\x2b\x00\x00\xff\xff\x6c\x96\xba\x82\xf9\x03\x00\x00")

func assetsLoadavgYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsLoadavgYml,
		"assets/loadavg.yml",
	)
}

func assetsLoadavgYml() (*asset, error) {
	bytes, err := assetsLoadavgYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/loadavg.yml", size: 1017, mode: os.FileMode(420), modTime: time.Unix(1556132016, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsMemoryusageYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x53\x41\x6f\xdb\x3c\x0c\xbd\xe7\x57\x10\x06\x8a\x7c\xdf\x21\x6e\x0a\x0c\x05\xe6\xdb\x50\xac\x97\x5d\x8a\xad\xdd\x65\x18\x06\xda\xa2\x33\x61\x12\x65\x50\x92\x0b\x37\xf1\x7f\x1f\x94\xb8\x99\xb2\x1a\x4b\x0f\xbb\xcd\x07\xc3\xe6\x23\xdf\x13\x1f\x29\x64\x34\xc3\x13\x55\x0b\x80\xfe\x2a\xbd\x01\x56\xc0\x68\xa9\x02\x4b\xd6\xc9\x50\x46\x8f\x1b\xda\x03\x00\x06\x6b\x32\xbe\x9a\xfe\x00\x74\xe3\xf8\x03\x0d\x15\x38\xff\xed\x90\x7e\x84\x14\x79\x2d\xa4\xee\x9c\xd7\x41\x3b\xae\xa0\xb8\x2e\x26\x50\x68\xa3\x7d\x20\xf9\x8c\xa2\xb1\x36\x94\x11\x9e\x6a\x3f\x24\xe9\x77\x3d\x6a\x93\xd2\x8e\x49\x30\x0b\x57\xb0\x1d\xff\xc4\x73\xef\x02\x9a\x79\x8e\x3d\x74\xae\xfe\x8e\xa4\x21\x0e\x19\x03\xf5\xa9\x6a\xb9\xdd\x0a\x75\x06\x74\x0b\xc8\x0a\xca\xb9\xb3\x9d\x44\xf7\x6a\xe3\x38\x95\x29\xdd\xdf\x1a\x87\x01\xfe\xf3\xb1\x3e\x7c\xbd\x48\x9e\x27\xfd\x7f\x26\x71\x07\x36\x9a\x03\xcb\xd5\x7a\x7d\x14\x21\x56\xe3\xb8\x9c\x8e\xde\x09\x35\x8e\xd5\x61\x2c\xc7\x76\x90\x55\x95\xf5\xb6\x02\x76\x21\x0f\xec\x1b\xb6\x5d\x18\x4e\x7c\x3a\x3c\xfd\x34\xc8\x8f\xd4\x9e\x1d\xdd\x5f\x21\xce\x67\x39\xd3\x0b\x47\x5b\x93\xdc\x38\xdb\xa1\x50\xae\xb5\x09\x15\xbc\x5d\x2f\xce\xd1\x9f\x8e\xda\x92\x4f\xc1\x6c\x4b\x73\x07\x6f\xd1\xf8\x13\x8d\x4e\xb4\x45\x19\x2a\x78\xe0\x1f\xec\x1e\x39\x83\x14\x05\xd4\xa6\x82\x1b\x17\x8d\x4a\x3e\xa4\x08\x89\xd5\x4c\x93\x3c\xec\x6f\x5b\x09\xef\xd9\x47\x21\x18\x5c\x14\xf0\x1d\x35\x1e\x34\x37\x26\x2a\x02\x84\xc6\x59\x9b\x36\xcd\xf9\x52\x08\xd5\xaa\xd5\x86\xbe\x14\x97\x9d\xb8\xe6\xd2\x92\xd5\xdc\xba\xe2\x6b\x99\xc9\x7a\xea\x49\x74\x32\x58\x51\x1d\x37\x47\xe4\x15\x3d\x3c\x2f\xb7\xb8\xf8\xdb\x6e\x4f\x1e\xc1\x1a\x76\x29\x9d\x43\x0b\xc5\x45\xb9\x6e\x8b\x71\xbc\x58\xce\xf4\xfc\xcc\xe4\x63\xfd\xda\x05\x87\x1d\xb4\x69\x95\xaf\xdf\xc0\x0e\xbe\x47\x8b\xfc\x49\x3f\xd1\x38\xfe\xf2\x8a\x14\xb8\x16\x26\xe6\xb9\xdb\x90\x55\x2d\x67\x1d\x49\x6e\xbd\x34\xe4\x5e\xe2\x3f\xea\xc7\x23\x0a\x2f\x7e\x06\x00\x00\xff\xff\xbe\x14\x02\x9a\x15\x06\x00\x00")

func assetsMemoryusageYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsMemoryusageYml,
		"assets/memoryUsage.yml",
	)
}

func assetsMemoryusageYml() (*asset, error) {
	bytes, err := assetsMemoryusageYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/memoryUsage.yml", size: 1557, mode: os.FileMode(420), modTime: time.Unix(1556135467, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _assetsOsYml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x91\xc1\x4a\xc3\x40\x10\x86\xef\x3e\xc5\xbc\x80\x82\xd7\xbd\x49\x55\x04\x85\x42\xab\x5e\x65\xd2\x4c\xe3\xc0\x66\x36\xee\x6e\x0a\x51\x7c\x77\xa9\x94\x36\xd3\x84\x34\x81\x25\x97\x40\x96\xef\x9f\x7f\x98\x0f\x05\x6d\xf3\x4d\xe6\x0a\x60\x77\xbb\xff\x02\x5c\x83\x60\x49\x06\x5c\xb8\xa9\xb3\x5a\x62\xfd\xff\x0a\x60\x31\x23\x1b\xcc\xe1\x0f\x80\x37\x4e\x9e\xa9\xd9\x83\x1f\x0a\xf4\x54\x70\x88\xe4\xdf\xd1\x33\x66\x96\x5a\x99\xd3\xec\xe3\x13\x80\x0b\x06\x7e\x7e\x0f\x0f\x1b\x27\x39\x47\x76\x72\x0a\x85\xe8\x59\x8a\x85\x2b\x2b\xf4\x64\x5a\x41\xfa\x32\xa0\x9a\x01\x76\x87\xce\x15\x6d\x5b\x2d\x25\x85\x80\x45\x7b\x91\x63\xcd\xab\xaf\xd5\xcc\xca\x73\x89\xbe\x31\xb0\x5c\x03\x07\x78\xd3\xe3\x01\x72\x8a\xc8\xd6\xc0\xb2\x22\x8f\x91\xa5\x80\x75\x13\x22\x95\xbd\x70\xa0\x1d\x79\x8e\x8d\x01\x96\xad\xeb\x96\x3f\xa2\x0d\x43\xed\xe2\xe2\xb4\x0d\x7a\x03\xa7\x2d\x72\xca\xea\xe2\x5c\x32\xda\x8a\x85\x46\x48\x56\xe0\xac\x92\x55\x73\x7a\xc9\x77\x7a\xfc\xf0\x89\x3b\x70\x0a\xc9\x93\x36\xe8\x0d\x5c\x90\xbc\x21\x89\xc7\x4b\x0d\x49\x56\xe0\xac\x92\x55\x73\x7a\xc9\x0b\x92\xb8\x5c\x8f\x3c\x71\x07\x4e\x21\x79\xd2\x06\xbd\x81\x0b\x92\xfd\x27\xd9\x11\x8a\x5b\xd8\xac\x82\x5b\xbd\xe9\xf5\xae\x9e\x1e\x5e\x46\x9e\xf6\x0c\x4d\xa1\x76\x42\x7b\x0f\x7e\x41\x6b\x4e\x19\xa3\x8c\x10\xab\xc0\x59\xd5\xaa\xe6\xf4\x72\xef\xf5\xf8\xe1\x03\x77\xe0\x14\x82\x27\x6d\xd0\x1b\x38\x97\xfc\x17\x00\x00\xff\xff\x63\xdc\x7e\x31\x7a\x09\x00\x00")

func assetsOsYmlBytes() ([]byte, error) {
	return bindataRead(
		_assetsOsYml,
		"assets/os.yml",
	)
}

func assetsOsYml() (*asset, error) {
	bytes, err := assetsOsYmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "assets/os.yml", size: 2426, mode: os.FileMode(420), modTime: time.Unix(1555703226, 0)}
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
	"assets/configManagement.yml": assetsConfigmanagementYml,
	"assets/cpuCores.yml": assetsCpucoresYml,
	"assets/datacenter.yml": assetsDatacenterYml,
	"assets/dockerContainers.yml": assetsDockercontainersYml,
	"assets/dockerDevicemapperLoopback.yml": assetsDockerdevicemapperloopbackYml,
	"assets/dockerIccFalse.yml": assetsDockericcfalseYml,
	"assets/dockerLoggingDriver.yml": assetsDockerloggingdriverYml,
	"assets/dockerVersion.yml": assetsDockerversionYml,
	"assets/loadavg.yml": assetsLoadavgYml,
	"assets/memoryUsage.yml": assetsMemoryusageYml,
	"assets/os.yml": assetsOsYml,
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
		"configManagement.yml": &bintree{assetsConfigmanagementYml, map[string]*bintree{}},
		"cpuCores.yml": &bintree{assetsCpucoresYml, map[string]*bintree{}},
		"datacenter.yml": &bintree{assetsDatacenterYml, map[string]*bintree{}},
		"dockerContainers.yml": &bintree{assetsDockercontainersYml, map[string]*bintree{}},
		"dockerDevicemapperLoopback.yml": &bintree{assetsDockerdevicemapperloopbackYml, map[string]*bintree{}},
		"dockerIccFalse.yml": &bintree{assetsDockericcfalseYml, map[string]*bintree{}},
		"dockerLoggingDriver.yml": &bintree{assetsDockerloggingdriverYml, map[string]*bintree{}},
		"dockerVersion.yml": &bintree{assetsDockerversionYml, map[string]*bintree{}},
		"loadavg.yml": &bintree{assetsLoadavgYml, map[string]*bintree{}},
		"memoryUsage.yml": &bintree{assetsMemoryusageYml, map[string]*bintree{}},
		"os.yml": &bintree{assetsOsYml, map[string]*bintree{}},
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

