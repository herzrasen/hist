package client

import (
	"fmt"
)

func (c *Client) Get(index int64) (string, error) {
	rows, err := c.Db.Query(`SELECT command 
			FROM hist 
			ORDER BY last_update DESC
			LIMIT 1
			OFFSET ?`, index)
	if err != nil {
		return "", fmt.Errorf("client:Client:Get: query: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var command string
		err := rows.Scan(&command)
		if err != nil {
			return "", fmt.Errorf("client:Client:Get: scan row: %w", err)
		}
		return command, nil
	}
	return "", fmt.Errorf("unable to get record with index: %d", index)
}
