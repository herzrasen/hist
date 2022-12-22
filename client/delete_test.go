package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestClient_Delete(t *testing.T) {
	t.Run("delete by prefix", func(t *testing.T) {
		c, err := CreateTestClient(t)
		require.NoError(t, err)
		defer os.Remove(c.Path)
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		if err = c.Update("test-command-2"); err != nil {
			t.FailNow()
		}
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		if err = c.Update("other-command"); err != nil {
			t.FailNow()
		}
		err = c.Delete(DeleteOptions{Prefix: "test-command"})
		require.NoError(t, err)
		records, err := c.List(ListOptions{})
		require.NoError(t, err)
		assert.Len(t, records, 1)
	})

	t.Run("delete by ids", func(t *testing.T) {
		c, err := CreateTestClient(t)
		require.NoError(t, err)
		defer os.Remove(c.Path)
		if err = c.Update("test-command-1"); err != nil {
			t.FailNow()
		}
		if err = c.Update("test-command-2"); err != nil {
			t.FailNow()
		}
		if err = c.Update("other-command"); err != nil {
			t.FailNow()
		}
		// get the entries to get some ids
		records, err := c.List(ListOptions{})
		require.NoError(t, err)
		require.Len(t, records, 3)
		err = c.Delete(DeleteOptions{Ids: []int64{
			records[0].Id,
			records[1].Id,
		}})
		require.NoError(t, err)
		recordsAfterDelete, err := c.List(ListOptions{})
		require.NoError(t, err)
		assert.Len(t, recordsAfterDelete, 1)
		assert.Equal(t, records[2], recordsAfterDelete[0])
	})
}
