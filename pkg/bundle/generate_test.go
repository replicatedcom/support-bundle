package bundle

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/divolgin/archiver/extractor"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/replicatedcom/support-bundle/pkg/types"
)

// TestGenerate runs stub tasks to ensure results are being parsed and packed properly
func TestGenerate(t *testing.T) {

	singleResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Description: "Testing single",
				Path:        "/testSingle.txt",
			},
		},
	}
	mixedResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Description: "Testing mixed results pass result",
				Path:        "/testPass.txt",
			},
			{
				Description: "Testing mixed results fail result",
				Error:       errors.New("This was destined to fail"),
			},
			{
				Description: "Testing mixed results other fail result",
				Path:        "/testFail.txt",
				Error:       errors.New("This was also meant to fail"),
			},
		},
	}

	tasks := []types.Task{singleResults, mixedResults}

	got, _ := ioutil.TempFile("", "generate-test-bundle")
	defer os.Remove(got.Name())

	err := Generate(tasks, time.Duration(time.Second*2), got.Name())
	require.NoError(t, err)

	testDir, err := ioutil.TempDir("", "generate-test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	//decompress to temp dir
	extractor := extractor.NewTgz()
	extractor.Extract(got.Name(), filepath.Join(testDir, "dir"))

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

	type testResult struct {
		Description string `json:"description"`
		Path        string `json:"path"`
		Error       string `json:"error,omitempty"`
	}

	var indexAll []testResult
	var errorAll []testResult

	err = json.Unmarshal(indexBytes, &indexAll)
	require.NoError(t, err)
	err = json.Unmarshal(errorBytes, &errorAll)
	require.NoError(t, err)

	// everything that includes a path
	require.Contains(t, indexAll, testResult{Description: "Testing single", Path: "/testSingle.txt"})
	require.Contains(t, indexAll, testResult{Description: "Testing mixed results pass result", Path: "/testPass.txt"})
	require.NotContains(t, indexAll, testResult{Description: "Testing mixed results fail result", Error: "This was destined to fail"})
	require.Contains(t, indexAll, testResult{Description: "Testing mixed results other fail result", Path: "/testFail.txt", Error: "This was also meant to fail"})

	// everything that includes an error
	require.NotContains(t, errorAll, testResult{Description: "Testing single", Path: "/testSingle.txt"})
	require.NotContains(t, errorAll, testResult{Description: "Testing mixed results pass result", Path: "/testPass.txt"})
	require.Contains(t, errorAll, testResult{Description: "Testing mixed results fail result", Error: "This was destined to fail"})
	require.Contains(t, errorAll, testResult{Description: "Testing mixed results other fail result", Path: "/testFail.txt", Error: "This was also meant to fail"})
}
