package ginkgo

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	dockernetworktypes "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/replicatedcom/support-bundle/cmd"
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

func LogErrors() {
	contents := ReadFileFromBundle(
		path.Join("bundle.tar.gz"),
		"/error.json",
	)
	fmt.Fprintf(GinkgoWriter, "Errors: %s\n", contents)
}

func WriteFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0666)
	Expect(err).NotTo(HaveOccurred())
}

func WriteBundleConfig(config string) {
	WriteFile("config.yml", config)
}

func GenerateBundle() {
	err := cmd.Generate(
		path.Join(tmpdir, "config.yml"),
		"",
		path.Join(tmpdir, "bundle.tar.gz"),
		true,
		60,
	)

	Expect(err).To(BeNil())
}

func GetFileFromBundle(pathInBundle string) string {
	return ReadFileFromBundle(
		path.Join(tmpdir, "bundle.tar.gz"),
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

		filePath := strings.TrimLeft(header.Name, "0123456789")
		jww.DEBUG.Printf("reading tar entry %s looking for %s", filePath, targetFile)

		if filePath == targetFile && header.Typeflag == tar.TypeReg {
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
func MakeDockerContainer() string {
	client, err := client.NewEnvClient()
	Expect(err).NotTo(HaveOccurred())

	containerSettings := dockercontainertypes.Config{Image: "ubuntu:latest", Cmd: []string{"sleep", "infinity"}}
	hostSettings := dockercontainertypes.HostConfig{}
	networkSettings := dockernetworktypes.NetworkingConfig{}

	container, err := client.ContainerCreate(context.Background(), &containerSettings, &hostSettings, &networkSettings, "")
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
