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
	"github.com/docker/docker/client"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/pkg/cli"
	"github.com/replicatedcom/support-bundle/pkg/cli/commands"
	"github.com/replicatedcom/support-bundle/pkg/types"
	jww "github.com/spf13/jwalterweatherman"
)

var tmpdir string
var cwd string
var err error

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
	LogResults(filepath.Join(tmpdir, "bundle.tar.gz"))()
}

func LogResults(archivePath string) func() {
	return func() {
		contents := ReadFileFromBundle(
			archivePath,
			"index.json",
		)
		jww.DEBUG.Printf("Index: %s\n", contents)
		contents = ReadFileFromBundle(
			archivePath,
			"error.json",
		)
		jww.DEBUG.Printf("Errors: %s\n", contents)
	}
}

func WriteFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0666)
	Expect(err).NotTo(HaveOccurred())
}

func WriteBundleConfig(config string) {
	WriteFile("config.yml", config)
}

func GenerateBundle() {
	cmd := commands.NewSupportBundleCommand(cli.NewCli())
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs([]string{
		"generate",
		fmt.Sprintf("--spec-file=%s", filepath.Join(tmpdir, "config.yml")),
		fmt.Sprintf("--out=%s", filepath.Join(tmpdir, "bundle.tar.gz")),
		"--skip-default=true",
		"--timeout=10",
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
		if result.Description == "/"+path {
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
	contents := ReadFileFromBundle(
		filepath.Join(tmpdir, "bundle.tar.gz"),
		index,
	)
	err := json.Unmarshal([]byte(contents), &results)
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
	return ReadFileFromBundle(
		filepath.Join(tmpdir, "bundle.tar.gz"),
		pathInBundle,
	)
}

func ReadFile(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	Expect(err).NotTo(HaveOccurred())
	return data
}

func ReadFileFromBundle(archivePath, targetFile string) string {
	file, err := os.Open(archivePath)
	defer CloseLogErr(file)
	Expect(err).NotTo(HaveOccurred())

	gzr, err := gzip.NewReader(file)
	defer CloseLogErr(gzr)
	Expect(err).NotTo(HaveOccurred())

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			Expect(err).NotTo(HaveOccurred(), "Failed to find "+targetFile+" in support bundle.")
		}
		Expect(err).NotTo(HaveOccurred())
		if header == nil {
			continue
		}

		jww.DEBUG.Printf("reading tar entry %s looking for %s", header.Name, targetFile)

		if header.Name == targetFile && header.Typeflag == tar.TypeReg {
			contents, err := ioutil.ReadAll(tr)
			Expect(err).NotTo(HaveOccurred())
			return string(contents)
		}
	}
}

func CloseLogErr(c io.Closer) {
	if err := c.Close(); err != nil {
		jww.ERROR.Print(err)
	}
}

// MakeDockerContainer makes a docker container to be used in tests, returning the container ID.
// name and labels are optional
func MakeDockerContainer(name string, labels map[string]string) string {
	client, err := client.NewEnvClient()
	Expect(err).NotTo(HaveOccurred())

	containerSettings := dockercontainertypes.Config{
		Image:  "ubuntu:latest",
		Cmd:    []string{"sleep", "infinity"},
		Labels: labels,
	}
	hostSettings := dockercontainertypes.HostConfig{}
	networkSettings := dockernetworktypes.NetworkingConfig{}

	container, err := client.ContainerCreate(context.Background(), &containerSettings, &hostSettings, &networkSettings, name)
	Expect(err).NotTo(HaveOccurred())
	Expect(container.Warnings).To(BeEmpty())

	err = client.ContainerStart(context.Background(), container.ID, dockertypes.ContainerStartOptions{})
	Expect(err).NotTo(HaveOccurred())

	return container.ID
}

// RemoveDockerContainer removes a docker container by ID as cleanup.
func RemoveDockerContainer(ID string) {
	client, err := client.NewEnvClient()
	Expect(err).NotTo(HaveOccurred())

	err = client.ContainerRemove(context.Background(), ID, dockertypes.ContainerRemoveOptions{Force: true})
	Expect(err).NotTo(HaveOccurred())
}
