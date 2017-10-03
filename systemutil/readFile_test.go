package systemutil

import (
	"testing"
	"time"

	"github.com/replicatedcom/support-bundle/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {
	resultsCh := make(chan types.Result)
	dataCh := make(chan types.Data)
	completeCh := make(chan bool, 1)

	go func() {
		for {
			select {
			case <-resultsCh:

			case <-dataCh:
			}
		}
	}()

	commandStrings := []string{"readFile_test.go"}

	err := ReadFile(dataCh, completeCh, resultsCh, time.Second*1, commandStrings)
	require.NoError(t, err)

	complete := <-completeCh
	assert.Equal(t, true, complete)
}
