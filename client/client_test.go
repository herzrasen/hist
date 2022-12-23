package client

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := CreateTestClient(t)
	defer os.Remove(c.Path)
	rows, err := c.Db.Query("SELECT name FROM sqlite_schema")
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

func CreateTestClient(t *testing.T) *Client {
	t.Helper()
	dbPath, err := os.CreateTemp("../", "test-*")
	require.NoError(t, err)
	c, err := NewSqliteClient(dbPath.Name())
	require.NoError(t, err)
	return c
}
