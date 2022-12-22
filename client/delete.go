package client

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	stmtDeleteByIds    = `DELETE FROM hist WHERE id IN (?)`
	stmtSelectByPrefix = `SELECT id, command, last_update, count FROM hist WHERE command LIKE ?`
	stmtDeleteByPrefix = `DELETE FROM hist WHERE command LIKE ?`
)

type DeleteOptions struct {
	Ids    []int64
	Prefix string
}

func (c *Client) Delete(options DeleteOptions) error {
	if len(options.Ids) > 0 {
		return c.deleteByIds(options)
	} else if options.Prefix != "" {
		return c.deleteByPrefix(options)
	}
	return nil
}

func (c *Client) deleteByIds(options DeleteOptions) error {
	query, args, err := sqlx.In(stmtDeleteByIds, options.Ids)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: in: %w", err)
	}
	_, err = c.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: exec: %w", err)
	}
	return nil
}

func (c *Client) deleteByPrefix(options DeleteOptions) error {
	prefix := fmt.Sprintf("%s%%", options.Prefix)
	res, err := c.db.Exec(stmtDeleteByPrefix, prefix)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: exec prefix: %w", err)
	}
	x, err := res.RowsAffected()
	fmt.Printf("Deleted %d entries\n", x)
	return nil
}
