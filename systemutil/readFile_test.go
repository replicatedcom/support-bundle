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
	require.Equal(t, 1, len(datas), "Expected 1 data struct to be returned")
	require.NoError(t, result.Error)
}
