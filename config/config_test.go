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
