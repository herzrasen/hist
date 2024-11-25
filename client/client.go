package client

import (
	"database/sql/driver"
	"fmt"
	"github.com/herzrasen/hist/config"
	"github.com/jmoiron/sqlx"
	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
	"regexp"
)

type Client struct {
	Db     *sqlx.DB
	Config *config.Config
	path   string
}

func NewSqliteClient(path string, cfg *config.Config) (*Client, error) {
	err := sqlite.RegisterScalarFunction("regexp", 2, func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
		pattern := args[0].(string)
		s := args[1].(string)
		return regexp.MatchString(pattern, s)
	})
	if err != nil {
		return nil, fmt.Errorf("error registering scalar function %w", err)
	}
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("client:NewSqliteClient: open: %w", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("unable to turn foreign keys on: %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS hist (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        command TEXT UNIQUE,
        last_update TIMESTAMP,
        count INTEGER DEFAULT 1
    )`)
	if err != nil {
		return nil, fmt.Errorf("client: NewSqliteClient: create table 'hist': %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tag (
        id INTEGER PRIMARY KEY AUTOINCREMENT, 
    	tag TEXT NOT NULL,
    	hist_id INTEGER,
    	FOREIGN KEY (hist_id) REFERENCES hist(id) ON DELETE CASCADE,
    	UNIQUE (hist_id, tag)
	)`)
	if err != nil {
		return nil, fmt.Errorf("client: NewSqliteClient: create table 'tag': %w", err)
	}
	return &Client{
		Db:     db,
		Config: cfg,
		path:   path,
	}, nil
}
