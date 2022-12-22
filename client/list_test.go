package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestClient_List(t *testing.T) {
	t.Run("succeed with empty ListOptions", func(t *testing.T) {
		c, err := CreateTestClient(t)
		require.NoError(t, err)
		defer os.Remove(c.Path)
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		time.Sleep(100 * time.Millisecond)
		if err = c.Update("test-command-2"); err != nil {
			t.FailNow()
		}
		time.Sleep(100 * time.Millisecond)
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		records, err := c.List(ListOptions{})
		require.NoError(t, err)
		assert.Len(t, records, 2)
		r0 := records[0]
		assert.Equal(t, "test-command-2", r0.Command)
		assert.Equal(t, uint64(1), r0.Count)
		r1 := records[1]
		assert.Equal(t, "test-command-1", r1.Command)
		assert.Equal(t, uint64(2), r1.Count)
	})

	t.Run("succeed with selection by ids", func(t *testing.T) {
		c, err := CreateTestClient(t)
		require.NoError(t, err)
		defer os.Remove(c.Path)
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		if err = c.Update("test-command-2"); err != nil {
			t.FailNow()
		}
		if err = c.Update("test-command-3"); err != nil {
			t.FailNow()
		}
		// select the records to get the ids
		records, err := c.List(ListOptions{})
		require.NoError(t, err)
		assert.Len(t, records, 3)
		ids := []int64{records[1].Id, records[2].Id}
		recordsByIds, err := c.List(ListOptions{Ids: ids})
		assert.Len(t, recordsByIds, 2)
	})
}
