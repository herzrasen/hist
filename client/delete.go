package client

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	stmtDeleteByIds = `DELETE FROM hist WHERE id IN (?)`
)

type DeleteOptions struct {
	Ids     []int64
	Pattern string
}

func (c *Client) Delete(options DeleteOptions) error {
	if len(options.Ids) > 0 {
		return c.deleteByIds(options)
	} else if options.Pattern != "" {
		return c.deleteByPattern(options)
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

func (c *Client) deleteByPattern(options DeleteOptions) error {
	stmt := buildDeleteByPatternStatement(options.Pattern)
	res, err := c.Db.Exec(stmt)
	if err != nil {
		return fmt.Errorf("hist.Client.Delete: exec prefix: %w", err)
	}
	x, err := res.RowsAffected()
	fmt.Printf("Deleted %d entries\n", x)
	return nil
}

func buildDeleteByPatternStatement(pattern string) string {
	return fmt.Sprintf("DELETE FROM hist WHERE regexp('%s', command) = true", pattern)
}
