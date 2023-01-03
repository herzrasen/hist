package client

import (
	"github.com/herzrasen/hist/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := CreateTestClient(t)
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
	c, err := NewSqliteClient(dbPath.Name(), &config.Config{})
	require.NoError(t, err)
	return c
}
