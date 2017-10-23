package planners

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseUptime(t *testing.T) {
	uptimeSeconds, err := parseUptime([]byte("33524.72 66785.42"))
	require.NoError(t, err)

	assert.Equal(t, float64(33524.72), uptimeSeconds.(uptime).TotalSeconds)

	uptimeSeconds, err = parseUptime([]byte("33524.72"))
	require.NotNil(t, err)
	assert.Nil(t, uptimeSeconds)

	uptimeSeconds, err = parseUptime([]byte("0 0"))
	require.NoError(t, err)
	assert.Equal(t, float64(0), uptimeSeconds.(uptime).IdleSeconds)
}
