package ginkgo

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	dockernetworktypes "github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	. "github.com/onsi/gomega"
	cmd "github.com/replicatedcom/support-bundle/cmd/support-bundle/commands"
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

var tmpdir string
var cwd string
var err error

type ErrFileNotFound struct {
	Filename string
}

func (e ErrFileNotFound) Error() string {
	return fmt.Sprintf("Failed to find %s in support bundle.", e.Filename)
}

func EnterNewTempDir() {
	cwd, err = os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	tmpdir, err = ioutil.TempDir("", "support-bundle")
	Expect(err).NotTo(HaveOccurred())
	err = os.Chdir(tmpdir)
	Expect(err).NotTo(HaveOccurred())
}

func CleanupDir() {
	err = os.Chdir(cwd)
	Expect(err).NotTo(HaveOccurred())
	err = os.RemoveAll(tmpdir)
	Expect(err).NotTo(HaveOccurred())
}

func LogResultsFomBundle() {
	contents := GetFileFromBundle("index.json")
	jww.DEBUG.Printf("Index: %s", contents)
	contents = GetFileFromBundle("error.json")
	jww.DEBUG.Printf("Errors: %s", contents)
}

func WriteFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0666)
	Expect(err).NotTo(HaveOccurred())
}

func WriteBundleConfig(config string) {
	WriteFile("config.yml", config)
}

func GenerateBundle() {
	cmd := cmd.NewSupportBundleCommand(cli.NewCli())
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{
		"generate",
		fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config.yml")),
		fmt.Sprintf("--out=%s", filepath.Join(tmpdir, "bundle.tar.gz")),
		"--timeout=10",
		"--skip-default",
		"--journald",
		"--kubernetes",
		"--retraced",
	})
	err := cmd.Execute()
	Expect(err).NotTo(HaveOccurred())
	// output := buf.String()
}

func GetResultFromBundle(path string) *types.Result {
	results := GetResultsFromBundle()
	for _, result := range results {
		if result.Path == "/"+path {
			return result
		}
	}
	Expect(fmt.Errorf("failed to find result at path %s", path)).NotTo(HaveOccurred())
	return nil
}

func GetResultFromBundleErrors(path string) *types.Result {
	results := GetResultsFromBundleErrors()
	for _, result := range results {
		if result.Path == "/"+path {
			return result
		}
	}
	Expect(fmt.Errorf("failed to find result at path %s", path)).NotTo(HaveOccurred())
	return nil
}

func GetResultsFromBundle() []*types.Result {
	return getResultsFromBundleIndex("index.json")
}

func GetResultsFromBundleErrors() []*types.Result {
	return getResultsFromBundleIndex("error.json")
}

func getResultsFromBundleIndex(index string) (results []*types.Result) {
	contents, err := ReadFileFromBundle(
		filepath.Join(tmpdir, "bundle.tar.gz"),
		index,
	)
	Expect(err).NotTo(HaveOccurred())
	err = json.Unmarshal([]byte(contents), &results)
	Expect(err).NotTo(HaveOccurred())
	return
}

func ExpectBundleErrorToHaveOccured(path, reStr string) {
	result := GetResultFromBundleErrors(path)
	if reStr == "" {
		return
	}
	re, err := regexp.Compile(reStr)
	Expect(err).NotTo(HaveOccurred())
	if !re.MatchString(result.Error.Error()) {
		Expect(fmt.Errorf("error %q for path %s does not match", result.Error, path)).NotTo(HaveOccurred())
	}
}

func GetFileFromBundle(pathInBundle string) string {
	contents, err := ReadFileFromBundle(
		filepath.Join(tmpdir, "bundle.tar.gz"),
		pathInBundle,
	)
	Expect(err).NotTo(HaveOccurred())
	return contents
}

func ExpectFileNotInBundle(pathInBundle string) {
	_, err := ReadFileFromBundle(
		filepath.Join(tmpdir, "bundle.tar.gz"),
		pathInBundle,
	)
	Expect(err).To(HaveOccurred())
	Expect(err).To(BeEquivalentTo(&ErrFileNotFound{pathInBundle}))
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func ReadFileFromBundle(archivePath, targetFile string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer CloseLogErr(file)

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer CloseLogErr(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			return "", &ErrFileNotFound{targetFile}
		}
		Expect(err).NotTo(HaveOccurred())
		if header == nil {
			continue
		}

		jww.DEBUG.Printf("reading tar entry %s looking for %s", header.Name, targetFile)

		if header.Name == targetFile && header.Typeflag == tar.TypeReg {
			contents, err := ioutil.ReadAll(tr)
			Expect(err).NotTo(HaveOccurred())
			return string(contents), nil
		}
	}
}

func CloseLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		jww.ERROR.Printf("Failed to close closer: %v", err)
	}
}

// MakeDockerContainer makes a docker container to be used in tests, returning the container ID.
// name and labels are optional
func MakeDockerContainer(client docker.CommonAPIClient, name string, labels map[string]string, cmd []string) string {
	Expect(err).NotTo(HaveOccurred())

	config := dockercontainertypes.Config{
		Image:  "ubuntu:latest",
		Cmd:    cmd,
		Labels: labels,
	}
	if config.Cmd == nil {
		config.Cmd = []string{"sleep", "infinity"}
	}
	hostConfig := dockercontainertypes.HostConfig{}
	networkConfig := dockernetworktypes.NetworkingConfig{}

	container, err := client.ContainerCreate(context.Background(), &config, &hostConfig, &networkConfig, name)
	Expect(err).NotTo(HaveOccurred())
	Expect(container.Warnings).To(BeEmpty())

	err = client.ContainerStart(context.Background(), container.ID, dockertypes.ContainerStartOptions{})
	Expect(err).NotTo(HaveOccurred())

	return container.ID
}

// RemoveDockerContainer removes a docker container by ID as cleanup.
func RemoveDockerContainer(client docker.CommonAPIClient, containerID string) {
	err = client.ContainerRemove(context.Background(), containerID, dockertypes.ContainerRemoveOptions{Force: true})
	Expect(err).NotTo(HaveOccurred())
}
