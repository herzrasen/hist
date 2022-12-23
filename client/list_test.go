package client

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herzrasen/hist/record"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_List(t *testing.T) {
	t.Run("succeed with empty ListOptions", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery("SELECT id, command, last_update, count FROM hist ORDER BY last_update DESC").
			WillReturnRows(sqlmock.NewRows([]string{"id", "command", "last_update", "count"}).
				AddRow(1, "test-command-1", time.Now(), 42).
				AddRow(2, "test-command-2", time.Now().Add(-1*time.Minute), 10))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		records, err := c.List(ListOptions{})
		require.NoError(t, err)
		assert.Len(t, records, 2)
		r0 := records[0]
		assert.Equal(t, "test-command-2", r0.Command)
		assert.Equal(t, uint64(10), r0.Count)
		r1 := records[1]
		assert.Equal(t, "test-command-1", r1.Command)
		assert.Equal(t, uint64(42), r1.Count)
	})

	t.Run("succeed with specified ids", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectQuery(`SELECT id, command, last_update, count 
				FROM hist WHERE id IN (?, ?, ?) ORDER BY last_update DESC`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "command", "last_update", "count"}))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		records, err := c.List(ListOptions{
			Ids: []int64{100, 101, 102},
		})
		require.NoError(t, err)
		assert.Len(t, records, 0)
	})
}

func TestListOptions_ToString(t *testing.T) {
	// records are sorted by list and not by ToString
	records := []record.Record{
		{Id: 2, Command: "command-2", LastUpdate: time.Now().Add(-1 * time.Second), Count: 8},
		{Id: 1, Command: "command-1", LastUpdate: time.Now(), Count: 5},
	}
	options := ListOptions{NoLastUpdate: true}
	got := options.ToString(records)
	wanted := "8\tcommand-2\n5\tcommand-1\n"
	assert.Equal(t, wanted, got)
}
