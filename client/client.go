package client

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	Path string
	db   *sql.DB
}

func NewClient(path string) (*Client, error) {
	db, err := sql.Open("sqlite3", path)
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
	return &Client{Path: path, db: db}, nil
}
