package bundle

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/divolgin/archiver/extractor"
	"github.com/replicatedcom/support-bundle/metrics"
	"github.com/replicatedcom/support-bundle/systemutil"
	"github.com/stretchr/testify/require"
)

// TestGenerate runs all the local data collection tools (read file, run command, hostname, loadavg, uptime)
// some tasks are not fully tested on windows (run command, loadavg, uptime)
func TestGenerate(t *testing.T) {

	successfulFile := "./generate_test.go"
	unsuccessfulFile := "/path/does/not/exist.xyz"

	var tasks = []Task{
		Task{
			Description: "Get File",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{successfulFile},
		},

		Task{
			Description: "Get nonexistent file",
			ExecFunc:    systemutil.ReadFile,
			Timeout:     time.Duration(time.Second * 15),
			Args:        []string{unsuccessfulFile},
		},

		Task{
			Description: "System hostname",
			ExecFunc:    metrics.Hostname,
			Timeout:     time.Duration(time.Second * 1),
		},

		Task{
			Description: "System loadavg",
			ExecFunc:    metrics.LoadAvg,
			Timeout:     time.Duration(time.Second * 1),
		},

		Task{
			Description: "System uptime",
			ExecFunc:    metrics.Uptime,
			Timeout:     time.Duration(time.Second * 1),
		},
	}

	if !(runtime.GOOS == "windows") {
		// command tasks are not tested on windows
		tasks = append(tasks,
			Task{
				Description: "Run command",
				ExecFunc:    systemutil.RunCommand,
				Timeout:     time.Duration(time.Second * 15),
				Args:        []string{"ls", "-a"},
			},
			Task{
				Description: "Run long Command",
				ExecFunc:    systemutil.RunCommand,
				Timeout:     time.Duration(time.Second * 1),
				Args:        []string{"sleep", "1m"},
			},
		)
	}

	got, err := Generate(tasks)
	require.NoError(t, err)

	testDir, _ := ioutil.TempDir("", "generate-test")
	defer os.RemoveAll(testDir)

	//decompress to temp dir
	extractor := extractor.NewTgz()
	extractor.ExtractFromReader(got, filepath.Join(testDir, "dir"))

	//verify what we got
	files, err := ioutil.ReadDir(filepath.Join(testDir, "dir"))
	require.NoError(t, err)

	require.Equal(t, 1, len(files))
	require.True(t, files[0].IsDir())

	uncompressedDir := files[0].Name()

	//get index.json and error.json
	indexReader, err := os.Open(filepath.Join(testDir, "dir", uncompressedDir, "index.json"))
	require.NoError(t, err)
	errorReader, err := os.Open(filepath.Join(testDir, "dir", uncompressedDir, "error.json"))
	require.NoError(t, err)

	//read into byte arrays
	indexBytes, err := ioutil.ReadAll(indexReader)
	require.NoError(t, err)
	errorBytes, err := ioutil.ReadAll(errorReader)
	require.NoError(t, err)

	var indexAll []resultInfo
	var errorAll []errorInfo

	err = json.Unmarshal(indexBytes, &indexAll)
	require.NoError(t, err)
	err = json.Unmarshal(errorBytes, &errorAll)
	require.NoError(t, err)

	// check for presence of what should be there
	// directory for successful file, error for nonexsistent file
	// directory for successful command, error for timeout command
	// directories for hostname, uptime, loadavg

	// jww.FEEDBACK.Print(len(indexAll))
	// jww.FEEDBACK.Print(len(errorAll))

	// search for successful file copy & unsuccessful file copy
	fileCopyPath := ""
	foundUnsuccessfulCopy := false
	for _, resultInfo := range indexAll {
		if resultInfo.Task == "readFile" && resultInfo.Args[0] == successfulFile {
			require.Equal(t, 1, len(resultInfo.Paths))
			fileCopyPath = resultInfo.Paths[0]
		}
		if resultInfo.Task == "readFile" && resultInfo.Args[0] == unsuccessfulFile {
			require.Equal(t, 0, len(resultInfo.Paths))
			foundUnsuccessfulCopy = true
		}
	}
	require.True(t, foundUnsuccessfulCopy, "A results index was not found for an unsuccessful copy")
	require.NotEqual(t, "", fileCopyPath, "No path was found for the successful file copy")

	fileCopyReader, err := os.Open(filepath.Join(testDir, "dir", uncompressedDir, fileCopyPath))
	require.NoError(t, err)
	fileCopyBytes, err := ioutil.ReadAll(fileCopyReader)
	require.NoError(t, err)

	// ensure the file was actually copied by checking for this magic string
	// the file we're reading is this test's source file, so by definition it
	require.True(t, strings.Contains(string(fileCopyBytes), "GlRIh6YfVnnJBo4TY3Q3"))

	foundFailedCommand := false
	if !(runtime.GOOS == "windows") {
		// search for successful command and timed out command
		// these commands aren't tested on windows platforms
		commandPath := ""
		for _, resultInfo := range indexAll {
			if resultInfo.Task == "runCommand" && resultInfo.Args[0] == "ls" {
				require.Equal(t, 1, len(resultInfo.Paths))
				commandPath = resultInfo.Paths[0]
			}
			if resultInfo.Task == "runCommand" && resultInfo.Args[0] == "sleep" {
				require.Equal(t, 0, len(resultInfo.Paths))
				foundFailedCommand = true
			}
		}
		require.True(t, foundFailedCommand, "A results index was not found for a timed out command")
		require.NotEqual(t, "", commandPath, "No path was found for the successful command run")
	}

	// look in the errors json and ensure entries are present for the failed copy and timed out command

	foundUnsuccessfulCopy = false
	foundFailedCommand = false
	for _, errorInfo := range errorAll {
		if errorInfo.Task == "readFile" && errorInfo.Args[0] == unsuccessfulFile {
			foundUnsuccessfulCopy = true
			require.NotEmpty(t, errorInfo.Errors)
		}
		if errorInfo.Task == "runCommand" && errorInfo.Args[0] == "sleep" {
			foundFailedCommand = true
			require.NotEmpty(t, errorInfo.Errors)
		}
	}

	require.True(t, foundUnsuccessfulCopy, "An error entry was not found for a failed file copy")

	if !(runtime.GOOS == "windows") {
		// command tasks not tested on windows
		require.True(t, foundFailedCommand, "An error entry was not found for a timed out command")
	}

	// look for uptime, loadavg, hostname task successes
	uptimeFound, loadavgFound, hostnameFound := false, false, false
	for _, resultInfo := range indexAll {
		if resultInfo.Task == "uptime" {
			if runtime.GOOS == "windows" {
				// uptime does not run properly on windows (the file doesn't exist)
				require.Empty(t, resultInfo.Paths)
			} else {
				require.NotEmpty(t, resultInfo.Paths)
			}
			uptimeFound = true
		}
		if resultInfo.Task == "loadavg" {
			if runtime.GOOS == "windows" {
				// loadavg does not run properly on windows (the file doesn't exist)
				require.Empty(t, resultInfo.Paths)
			} else {
				require.NotEmpty(t, resultInfo.Paths)
			}
			loadavgFound = true
		}
		if resultInfo.Task == "hostname" {
			require.NotEmpty(t, resultInfo.Paths)
			hostnameFound = true
		}
	}

	require.True(t, uptimeFound)
	require.True(t, loadavgFound)
	require.True(t, hostnameFound)
}
