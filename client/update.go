package client

import (
	"fmt"
	"time"
)

const insertStmt = `INSERT INTO hist (command, last_update) VALUES (?, ?) 
	ON CONFLICT(command) DO UPDATE SET count=count+1, last_update=excluded.last_update`

func (c *Client) Update(command string) error {
	_, err := c.Db.Exec(insertStmt, command, time.Now())
	if err != nil {
		return fmt.Errorf("hist.Client.Update: exec: %w", err)
	}
	return nil
}
