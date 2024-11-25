package client

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Stats(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		baseTime := time.Now()
		rows := sqlmock.NewRows([]string{"id", "command", "last_update", "count"}).
			AddRow(1, "command-1", baseTime.Add(-2*time.Minute), 1).
			AddRow(2, "command-2", baseTime.Add(-1*time.Minute), 2).
			AddRow(3, "command-3", baseTime.Add(-3*time.Minute), 10).
			AddRow(3, "command-4", baseTime.Add(-2*time.Minute), 10)
		mock.ExpectQuery(`SELECT h.id, h.command, h.last_update, h.count, null 
			FROM hist h 
			ORDER BY h.last_update, h.count DESC`).
			WillReturnRows(rows)
		c := Client{
			Db: sqlx.NewDb(db, "sqlite3"),
		}
		s, err := c.Stats()
		require.NoError(t, err)
		assert.Equal(t, uint64(4), s.NumCommands)
		assert.Len(t, s.MostUsedCommands, 2)
		assert.Contains(t, s.MostUsedCommands, "command-3")
		assert.Contains(t, s.MostUsedCommands, "command-4")
		assert.Equal(t, baseTime.Add(-3*time.Minute), s.OldestCommandTime)
		assert.Equal(t, baseTime.Add(-1*time.Minute), s.NewestCommandTime)
		assert.Len(t, s.OldestCommands, 1)
		assert.Contains(t, s.OldestCommands, "command-3")
		assert.Len(t, s.NewestCommands, 1)
		assert.Contains(t, s.NewestCommands, "command-2")
	})

	t.Run("list fails", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery(`SELECT id, command, last_update, count 
			FROM hist 
			ORDER BY last_update, count DESC`).
			WillReturnError(errors.New("some error"))
		c := Client{
			Db: sqlx.NewDb(db, "sqlite3"),
		}
		_, err := c.Stats()
		require.Error(t, err)
	})
}
