package systemutil

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {

	// it's a command that is sure to be installed AND will be the same on both windows and linux
	commandStrings := []string{"go", "help"}

	datas, result, err := RunCommand(context.Background(), commandStrings)
	require.NoError(t, err)
	require.Equal(t, 1, len(datas), "Expected 1 data struct to be returned")
	require.NoError(t, result.HumanError)
	require.NoError(t, result.RawError)
	require.NoError(t, result.JSONError)
}

func TestRunCommandTimeout(t *testing.T) {

	commandStrings := []string{"sleep", "10s"}
	if runtime.GOOS == "windows" {
		t.Skipf("This test is not yet compatible with windows")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	datas, result, err := RunCommand(ctx, commandStrings)
	require.Error(t, err)
	require.Equal(t, 0, len(datas), "Expected no data structs to be returned due to timeout")
	require.Error(t, result.HumanError)
	require.Error(t, result.RawError)
	require.Error(t, result.JSONError)
}
