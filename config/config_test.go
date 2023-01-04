package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("load succeeds", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-config-*")
		require.NoError(t, err)
		_, err = f.WriteString(`---
patterns:
  excludes:
    - go run
    - ll
`)
		require.NoError(t, err)
		c, err := Load(f.Name())
		require.NoError(t, err)
		assert.Len(t, c.Patterns.Excludes, 2)
	})

	t.Run("load fails (invalid format)", func(t *testing.T) {
		f, err := os.CreateTemp("", "hist-config-*")
		require.NoError(t, err)
		_, err = f.WriteString("<wrong>config</wrong>")
		require.NoError(t, err)
		_, err = Load(f.Name())
		require.Error(t, err)
	})

	t.Run("load fails (path)", func(t *testing.T) {
		// Create temp file just to delete it, to ensure that it's a valid
		// file that does not exist anymore
		f, err := os.CreateTemp("", "hist-config-*")
		require.NoError(t, err)
		err = os.Remove(f.Name())
		require.NoError(t, err)
		_, err = Load(f.Name())
		require.Error(t, err)
	})
}

func TestConfig_resolvePath(t *testing.T) {
	t.Run("resolve with leading ~", func(t *testing.T) {
		home, err := os.UserHomeDir()
		require.NoError(t, err)
		p, err := resolvePath("~/.config/hist/config.hist")
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(p, home))
	})
}

func TestConfig_IsExcluded(t *testing.T) {
	t.Run("is excluded", func(t *testing.T) {
		c := Config{
			Patterns: Patterns{Excludes: []string{
				"^foo.*",
			}},
		}
		assert.True(t, c.IsExcluded("foo bar"))
	})

	t.Run("is excluded II", func(t *testing.T) {
		c := Config{
			Patterns: Patterns{Excludes: []string{
				"foo bar",
			}},
		}
		assert.True(t, c.IsExcluded("this foo bar is excluded"))
	})

	t.Run("not excluded", func(t *testing.T) {
		c := Config{
			Patterns: Patterns{Excludes: []string{
				"^foo.*",
			}},
		}
		assert.False(t, c.IsExcluded("bar baz"))
	})
}
