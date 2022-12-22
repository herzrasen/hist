package client

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := CreateTestClient(t)
	require.NoError(t, err)
	defer os.Remove(c.Path)
	rows, err := c.db.Query("SELECT name FROM sqlite_schema")
	require.NoError(t, err)
	defer rows.Close()
	var tableExists bool
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		require.NoError(t, err)
		if name == "hist" {
			tableExists = true
		}
	}
	if !tableExists {
		t.Fatal("No table named 'hist' found")
	}
}

func CreateTestClient(t *testing.T) (*Client, error) {
	t.Helper()
	dbPath, err := os.CreateTemp("../", "test-*")
	require.NoError(t, err)
	return NewClient(dbPath.Name())
}
