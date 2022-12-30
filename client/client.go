package client

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
)

type Client struct {
	Path string
	Db   *sqlx.DB
}

func NewSqliteClient(path string) (*Client, error) {
	regex := func(pattern string, s string) (bool, error) {
		return regexp.MatchString(pattern, s)
	}
	sql.Register("sqlite3_regex", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("regexp", regex, true)
		}})
	db, err := sqlx.Open("sqlite3_regex", path)
	if err != nil {
		return nil, fmt.Errorf("client.NewClient: open: %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hist (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        command TEXT UNIQUE,
        last_update TIMESTAMP,
        count INTEGER DEFAULT 1
    )`)
	if err != nil {
		return nil, fmt.Errorf("client.NewClient: create table: %w", err)
	}
	return &Client{Path: path, Db: db}, nil
}
