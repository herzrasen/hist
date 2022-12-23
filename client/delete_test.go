package client

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_Delete(t *testing.T) {
	t.Run("delete by prefix", func(t *testing.T) {
		db, mock, _ := sqlmock.New()
		mock.ExpectExec("DELETE FROM hist WHERE command LIKE ?").
			WithArgs("test-command%").
			WillReturnResult(sqlmock.NewResult(0, 2))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Prefix: "test-command"})
		require.NoError(t, err)
	})

	t.Run("delete by ids", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE id IN (?, ?)").
			WithArgs(100, 101).
			WillReturnResult(sqlmock.NewResult(0, 2))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Ids: []int64{100, 101}})
		require.NoError(t, err)
	})

	t.Run("exec returns err with ids", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE id IN (?, ?)").
			WillReturnError(errors.New("some error"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Ids: []int64{100, 101}})
		require.Error(t, err)
	})

	t.Run("exec returns err with prefix", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE command LIKE ?").
			WillReturnError(errors.New("some error"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Prefix: "test"})
		require.Error(t, err)
	})
}
