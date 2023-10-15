package client

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_Get(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery("SELECT command FROM hist ORDER BY last_update DESC LIMIT 1 OFFSET ?").
			WithArgs(101).
			WillReturnRows(sqlmock.NewRows([]string{"command"}).
				AddRow("my-command --help"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		_, err := c.Get(101)
		require.NoError(t, err)
		//assert.Equal(t, "my-command --help", command)
	})

	t.Run("scan returns error", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		rows := sqlmock.NewRows([]string{"command"}).
			AddRow(nil)
		mock.ExpectQuery("SELECT command FROM hist ORDER BY last_update DESC LIMIT 1 OFFSET ?").
			WithArgs(100).
			WillReturnRows(rows)
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		_, err := c.Get(100)
		require.Error(t, err)
	})

	t.Run("index not found", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		rows := sqlmock.NewRows([]string{"command"})
		mock.ExpectQuery("SELECT command FROM hist ORDER BY last_update DESC LIMIT 1 OFFSET ?").
			WithArgs(100).
			WillReturnRows(rows)
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		_, err := c.Get(100)
		require.Error(t, err)
	})

	t.Run("query fails", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectQuery("SELECT command FROM hist ORDER BY last_update DESC LIMIT 1 OFFSET ?").
			WithArgs(101).
			WillReturnError(errors.New("some error"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		_, err := c.Get(101)
		require.Error(t, err)
	})
}
