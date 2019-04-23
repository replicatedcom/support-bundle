package analyze

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

// Make sure we can get an instance of analyze
func TestDI(t *testing.T) {
	req := require.New(t)

	v := viper.New()
	v.Set("headless", true)

	container, err := buildInjector(v)
	req.NoError(err)

	err = container.Invoke(func(s *Analyze) error {
		// don't do anything with it, just make sure we can get one
		return nil
	})

	req.NoError(err)
}
