package client

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	stmtDeleteByIds    = `DELETE FROM hist WHERE id IN (?)`
	stmtDeleteByFilter = `DELETE FROM hist WHERE command LIKE ?`
)

type DeleteOptions struct {
	Ids    []int64
	Filter string
}

func (c *Client) Delete(options DeleteOptions) error {
	if len(options.Ids) > 0 {
		return c.deleteByIds(options)
	} else if options.Filter != "" {
		return c.deleteByFilter(options)
	}
	return nil
}

func (c *Client) deleteByIds(options DeleteOptions) error {
	query, args, err := sqlx.In(stmtDeleteByIds, options.Ids)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: in: %w", err)
	}
	_, err = c.Db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: exec: %w", err)
	}
	return nil
}

func (c *Client) deleteByFilter(options DeleteOptions) error {
	prefix := fmt.Sprintf("%s%%", options.Filter)
	res, err := c.Db.Exec(stmtDeleteByFilter, prefix)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: exec prefix: %w", err)
	}
	x, err := res.RowsAffected()
	fmt.Printf("Deleted %d entries\n", x)
	return nil
}
