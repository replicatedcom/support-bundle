package systemutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadFile(t *testing.T) {

	commandStrings := []string{"readFile_test.go"}

	datas, result, err := ReadFile(context.Background(), commandStrings)
	require.NoError(t, err)
	require.Equal(t, 3, len(datas), "Expected 3 data structs to be returned")
	require.NoError(t, result.HumanError)
	require.NoError(t, result.RawError)
	require.NoError(t, result.JSONError)
}
