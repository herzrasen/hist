package client

import (
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herzrasen/hist/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestClient_Update(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec(`INSERT INTO hist (command, last_update) 
			VALUES (?, ?) 
			ON CONFLICT(command) 
			    DO UPDATE SET count=count+1, last_update=excluded.last_update`).
			WithArgs("ls -alF", AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		c := Client{
			Db:     sqlx.NewDb(db, "sqlite3"),
			Config: &config.Config{},
		}
		err := c.Record("ls -alF")
		require.NoError(t, err)
	})
	t.Run("insert returns error", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec(`INSERT INTO hist (command, last_update) 
			VALUES (?, ?) 
			ON CONFLICT(command) 
			    DO UPDATE SET count=count+1, last_update=excluded.last_update`).
			WillReturnError(errors.New("some error"))
		c := Client{
			Db:     sqlx.NewDb(db, "sqlite3"),
			Config: &config.Config{},
		}
		err := c.Record("ls -alF")
		require.Error(t, err)
	})
}
