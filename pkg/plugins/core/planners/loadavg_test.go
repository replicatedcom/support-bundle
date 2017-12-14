package planners

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLoadAvg(t *testing.T) {
	loadAvgValues, err := parseLoadAvg(strings.NewReader("0.02 0.01 0.00 4/229 5"))
	require.NoError(t, err)
	assert.Equal(t, float64(0.02), loadAvgValues.(*LoadAverage).MinuteOne)

	loadAvgValues, err = parseLoadAvg(strings.NewReader("0"))
	require.NotNil(t, err)
	assert.Nil(t, loadAvgValues)
}
