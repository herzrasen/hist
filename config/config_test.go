package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("load succeeds", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-config-*")
		require.NoError(t, err)
		_, err = f.WriteString(`---
excludes:
    - go run
    - ll
`)
		require.NoError(t, err)
		c, err := Load(f.Name())
		require.NoError(t, err)
		assert.Len(t, c.Excludes, 2)
	})
}
