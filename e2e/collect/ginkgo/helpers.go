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
	"os/exec"
	"path/filepath"
	"regexp"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainertypes "github.com/docker/docker/api/types/container"
	dockernetworktypes "github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmd "github.com/replicatedcom/support-bundle/cmd/support-bundle/commands"
	"github.com/replicatedcom/support-bundle/pkg/collect/cli"
	"github.com/replicatedcom/support-bundle/pkg/collect/types"
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

func SetupLogger() {
	// jww.SetLogOutput(GinkgoWriter)
	jww.SetStdoutThreshold(jww.LevelTrace)
}

func EnterNewTempDir() {
	cwd, err = os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	tmpdir, err = ioutil.TempDir("", "support-bundle")
	Expect(err).NotTo(HaveOccurred())
	fmt.Fprintln(GinkgoWriter, "Entering temp dir", tmpdir)
	err = os.Chdir(GetTempDir())
	Expect(err).NotTo(HaveOccurred())
}

func GetTempDir() string {
	return tmpdir
}

func CleanupDir() {
	err = os.Chdir(cwd)
	Expect(err).NotTo(HaveOccurred())
	err = os.RemoveAll(GetTempDir())
	Expect(err).NotTo(HaveOccurred())
}

func LogResultsFromBundle() {
	src := filepath.Join(GetTempDir(), "bundle.tar.gz")
	if _, err := os.Stat(src); err != nil {
		fmt.Fprintf(GinkgoWriter, "bundle.tar.gz not found for log results: %v\n", err)
		return
	}
	contents, err := ReadFileFromBundle(src, "index.json")
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to read index.json: %v\n", err)
	} else {
		jww.DEBUG.Printf("Index: %s", contents)
	}
	contents, err = ReadFileFromBundle(src, "error.json")
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to read error.json: %v\n", err)
	} else {
		jww.DEBUG.Printf("Errors: %s", contents)
	}
}

func LogDockerInfo() {
	commands := [][]string{
		{"docker", "version"},
		{"docker", "info"},
		{"docker", "ps", "-a"},
		{"docker", "images", "--format", "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"},
	}
	for _, cmdArgs := range commands {
		out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).CombinedOutput()
		label := cmdArgs[0] + " " + cmdArgs[1]
		if err != nil {
			fmt.Fprintf(GinkgoWriter, "%s failed: %v\n%s\n", label, err, string(out))
		} else {
			fmt.Fprintf(GinkgoWriter, "%s output:\n%s\n", label, string(out))
		}
	}
}

func PreserveBundleArtifact() {
	src := filepath.Join(GetTempDir(), "bundle.tar.gz")
	if _, err := os.Stat(src); err != nil {
		fmt.Fprintf(GinkgoWriter, "bundle.tar.gz stat error: %v\n", err)
		files, _ := filepath.Glob(filepath.Join(GetTempDir(), "*"))
		fmt.Fprintf(GinkgoWriter, "Files in temp dir %s: %v\n", GetTempDir(), files)
		return
	}
	artifactsDir := "/tmp/e2e-artifacts"
	_ = os.MkdirAll(artifactsDir, 0755)
	dst := filepath.Join(artifactsDir, fmt.Sprintf("bundle-%s.tar.gz", filepath.Base(GetTempDir())))
	in, err := os.Open(src)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to open bundle for preservation: %v\n", err)
		return
	}
	defer CloseLogErr(in)
	out, err := os.Create(dst)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to create artifact file: %v\n", err)
		return
	}
	defer CloseLogErr(out)
	_, err = io.Copy(out, in)
	if err != nil {
		fmt.Fprintf(GinkgoWriter, "Failed to copy bundle artifact: %v\n", err)
		return
	}
	fmt.Fprintf(GinkgoWriter, "Preserved bundle artifact to %s\n", dst)
}

func WriteFile(path string, contents string) {
	err := ioutil.WriteFile(path, []byte(contents), 0666)
	Expect(err).NotTo(HaveOccurred())
}

func WriteBundleConfig(config string) {
	WriteFile("config.yml", config)
}

func GenerateBundle(args ...string) {
	cmd := cmd.NewSupportBundleCommand(cli.NewCli())
	buf := new(bytes.Buffer)
	cmd.SetOutput(buf)
	cmd.SetArgs(append([]string{
		"generate",
		fmt.Sprintf("--spec-file=%s", filepath.Join(GetTempDir(), "config.yml")),
		fmt.Sprintf("--out=%s", filepath.Join(GetTempDir(), "bundle.tar.gz")),
		"--timeout=10",
		"--skip-default",
	}, args...))
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
		filepath.Join(GetTempDir(), "bundle.tar.gz"),
		index,
	)
	Expect(err).NotTo(HaveOccurred())
	err = json.Unmarshal([]byte(contents), &results)
	Expect(err).NotTo(HaveOccurred())
	return
}

func ExpectBundleErrorToHaveOccurred(path, reStr string) {
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
		filepath.Join(GetTempDir(), "bundle.tar.gz"),
		pathInBundle,
	)
	Expect(err).NotTo(HaveOccurred())
	return contents
}

func ExpectFileNotInBundle(pathInBundle string) {
	_, err := ReadFileFromBundle(
		filepath.Join(GetTempDir(), "bundle.tar.gz"),
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
		Image:  "ubuntu:16.04",
		Cmd:    cmd,
		Labels: labels,
		Env: []string{
			"ENVNORMAL=normal",
			"ENVSCRUBBED=secret",
			"ENVSCRUBBEDANOTHER=anothersecret",
			"ENVNORMALTWO=normaltwo",
		},
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
