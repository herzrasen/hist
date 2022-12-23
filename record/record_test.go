package record

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRecord_Format(t *testing.T) {
	lastUpdate, err := time.Parse(time.RFC3339, "2022-12-21T09:53:12+01:00")
	require.NoError(t, err)
	r := Record{
		Id:         100,
		Command:    "some-command",
		LastUpdate: lastUpdate,
		Count:      42,
	}

	t.Run("default", func(t *testing.T) {
		formatted := r.Format(FormatOptions{})
		expected := fmt.Sprintf("%s\t%d\t%s", r.LastUpdate.Format(time.RFC1123), r.Count, r.Command)
		assert.Equal(t, expected, formatted)
	})

	t.Run("with id", func(t *testing.T) {
		formatted := r.Format(FormatOptions{
			WithId: true,
		})
		expected := fmt.Sprintf("%s\t%d\t%d\t%s", r.LastUpdate.Format(time.RFC1123), r.Count, r.Id, r.Command)
		assert.Equal(t, expected, formatted)
	})

	t.Run("no count", func(t *testing.T) {
		formatted := r.Format(FormatOptions{
			NoCount: true,
		})
		expected := fmt.Sprintf("%s\t%s", r.LastUpdate.Format(time.RFC1123), r.Command)
		assert.Equal(t, expected, formatted)
	})

	t.Run("no lastUpdate", func(t *testing.T) {
		formatted := r.Format(FormatOptions{
			NoLastUpdate: true,
		})
		expected := fmt.Sprintf("%d\t%s", r.Count, r.Command)
		assert.Equal(t, expected, formatted)
	})
}
