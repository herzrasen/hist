package client

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Delete(t *testing.T) {
	t.Run("delete by pattern", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE regexp('^test-command.*', command) = true").
			WillReturnResult(sqlmock.NewResult(0, 2))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Pattern: "^test-command.*"})
		require.NoError(t, err)
	})

	t.Run("exec returns err with pattern", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE regexp('test', command) = true").
			WillReturnError(errors.New("some error"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Pattern: "test"})
		require.Error(t, err)
	})

	t.Run("exec returns error for rowsAffected (pattern)", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec("DELETE FROM hist WHERE regexp('foo', command) = true").
			WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{Pattern: "foo"})
		require.Error(t, err)
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

	t.Run("delete records updated before  date", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		date := time.Now().Add(-24 * time.Hour)
		mock.ExpectExec("DELETE FROM hist WHERE last_update < ?").
			WithArgs(date).
			WillReturnResult(sqlmock.NewResult(0, 2))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{UpdatedBefore: &date})
		require.NoError(t, err)
	})

	t.Run("exec returns error for rowsAffected (updatedBefore)", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		date := time.Now()
		mock.ExpectExec("DELETE FROM hist WHERE last_update < ?").
			WithArgs(date).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("some error")))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{UpdatedBefore: &date})
		require.Error(t, err)
	})

	t.Run("exec returns err when deleting records updated before date", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		date := time.Now().Add(-24 * time.Hour)
		mock.ExpectExec("DELETE FROM hist WHERE last_update < ?").
			WithArgs(date).
			WillReturnError(errors.New("some error"))
		c := Client{Db: sqlx.NewDb(db, "sqlite3")}
		err := c.Delete(DeleteOptions{UpdatedBefore: &date})
		require.Error(t, err)
	})

	t.Run("return an error when empty deleteOptions are passed", func(t *testing.T) {
		c := Client{}
		err := c.Delete(DeleteOptions{})
		require.Error(t, err)
	})

}
