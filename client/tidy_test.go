package client

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herzrasen/hist/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Tidy(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery(`SELECT id, command, last_update, count 
          FROM hist 
          ORDER BY last_update, count DESC`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "command", "last_update", "count"}).
				AddRow(1, "ls -alF", time.Now().Add(-5*time.Second), 100).
				AddRow(2, "git push", time.Now(), 10))
		mock.ExpectExec("DELETE FROM hist WHERE id IN (?)").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		c := Client{
			Db: sqlx.NewDb(db, "sqlite3"),
			Config: &config.Config{
				Patterns: config.Patterns{
					Excludes: []string{
						"^ls .*",
					}},
			},
		}
		err := c.Tidy()
		require.NoError(t, err)
	})

	t.Run("list returns error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery(`SELECT id, command, last_update, count 
          FROM hist 
          ORDER BY last_update, count DESC`).
			WillReturnError(errors.New("some error"))
		c := Client{
			Db: sqlx.NewDb(db, "sqlite3"),
			Config: &config.Config{
				Patterns: config.Patterns{
					Excludes: []string{
						"^ls .*",
					}},
			},
		}
		err := c.Tidy()
		require.Error(t, err)
	})
}
