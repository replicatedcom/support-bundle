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
	"github.com/replicatedcom/support-bundle/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerate runs stub tasks to ensure results are being parsed and packed properly
func TestGenerate(t *testing.T) {

	singleResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing single"}},
				Path: "/testSingle.txt",
			},
		},
	}
	mixedResults := taskStub{
		elapse: time.Nanosecond,
		results: []*types.Result{
			{
				Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results pass result"}},
				Path: "/testPass.txt",
			},
			{
				Spec:  types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results fail result"}},
				Error: errors.New("This was destined to fail"),
			},
			{
				Spec:  types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results other fail result"}},
				Path:  "/testFail.txt",
				Error: errors.New("This was also meant to fail"),
			},
		},
	}

	tasks := []types.Task{singleResults, mixedResults}

	got, _ := ioutil.TempFile("", "generate-test-bundle")
	defer os.Remove(got.Name())

	_, _, err := Generate(tasks, time.Duration(time.Second*2), got.Name())
	require.NoError(t, err)

	testDir, err := ioutil.TempDir("", "generate-test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	//decompress to temp dir
	uncompressedDir := filepath.Join(testDir, "dir")
	extractor := extractor.NewTgz()
	extractor.Extract(got.Name(), uncompressedDir)

	//verify what we got
	files, err := ioutil.ReadDir(filepath.Join(testDir, "dir"))
	require.NoError(t, err)

	require.Equal(t, 2, len(files))

	//get index.json and error.json
	indexReader, err := os.Open(filepath.Join(uncompressedDir, "index.json"))
	require.NoError(t, err)
	errorReader, err := os.Open(filepath.Join(uncompressedDir, "error.json"))
	require.NoError(t, err)

	//read into byte arrays
	indexBytes, err := ioutil.ReadAll(indexReader)
	require.NoError(t, err)
	errorBytes, err := ioutil.ReadAll(errorReader)
	require.NoError(t, err)

	type testResult struct {
		Description string     `json:"description"`
		Path        string     `json:"path"`
		Spec        types.Spec `json:"spec"`
		Error       string     `json:"error,omitempty"`
	}

	var indexAll []testResult
	var errorAll []testResult

	err = json.Unmarshal(indexBytes, &indexAll)
	require.NoError(t, err)
	err = json.Unmarshal(errorBytes, &errorAll)
	require.NoError(t, err)

	// everything that includes a path
	assert.Contains(t, indexAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing single"}}, Path: "/testSingle.txt"})
	assert.Contains(t, indexAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results pass result"}}, Path: "/testPass.txt"})
	assert.NotContains(t, indexAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results fail result"}}, Error: "This was destined to fail"})
	assert.NotContains(t, indexAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results other fail result"}}, Path: "/testFail.txt", Error: "This was also meant to fail"})

	// everything that includes an error
	assert.NotContains(t, errorAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing single"}}, Path: "/testSingle.txt"})
	assert.NotContains(t, errorAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results pass result"}}, Path: "/testPass.txt"})
	assert.Contains(t, errorAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results fail result"}}, Error: "This was destined to fail"})
	assert.Contains(t, errorAll, testResult{Spec: types.Spec{SpecShared: types.SpecShared{Description: "Testing mixed results other fail result"}}, Path: "/testFail.txt", Error: "This was also meant to fail"})
}
