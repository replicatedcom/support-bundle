package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLoadAverage(t *testing.T) {
	loadAvgValues, err := parseLoadAvg([]byte("0.02 0.01 0.00 4/229 5"))
	require.NoError(t, err)
	assert.Equal(t, loadAvgValues.minuteOne, float64(0.02))

	loadAvgValues, err = parseLoadAvg([]byte("0"))
	require.NotNil(t, err)
	assert.Nil(t, loadAvgValues)
}
