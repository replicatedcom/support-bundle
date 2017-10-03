package systemutil

import (
	"runtime"
	"testing"
	"time"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool, 1)

	// accept and discard channel data
	go func() {
		for {
			select {
			case <-resultsCh:

			case <-dataCh:
			}
		}
	}()

	// it's a command that is sure to be installed AND will be the same on both windows and linux
	commandStrings := []string{"go", "help"}

	err := RunCommand(dataCh, completeCh, resultsCh, time.Second*1, commandStrings)
	require.NoError(t, err)

	complete := <-completeCh
	assert.Equal(t, true, complete)
}

func TestRunCommandTimeout(t *testing.T) {
	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool, 1)

	// accept and discard channel data
	go func() {
		for {
			select {
			case <-resultsCh:
			case <-dataCh:
			}
		}
	}()

	commandStrings := []string{"sleep", "10s"}
	if runtime.GOOS == "windows" {
		t.Skipf("This test is not yet compatible with windows")
	}

	err := RunCommand(dataCh, completeCh, resultsCh, time.Second*1, commandStrings)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed")

	complete := <-completeCh
	assert.Equal(t, true, complete)
}
