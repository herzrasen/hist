package client

import (
	"fmt"
	"time"
)

const insertStmt = `INSERT INTO hist (command, last_update) VALUES (?, ?) 
	ON CONFLICT(command) DO UPDATE SET count=count+1, last_update=excluded.last_update`

func (c *Client) RecordWithTime(command string, t time.Time) error {
	if c.Config.IsExcluded(command) {
		return nil
	}
	_, err := c.Db.Exec(insertStmt, command, t)
	if err != nil {
		return fmt.Errorf("hist.Client.Record: exec: %w", err)
	}
	return nil
}

func (c *Client) Record(command string) error {
	return c.RecordWithTime(command, time.Now())
}
