package client

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	stmtDeleteByIds = `DELETE FROM hist WHERE id IN (?)`
)

type DeleteOptions struct {
	Ids           []int64
	UpdatedBefore *time.Time
	Pattern       string
}

func (c *Client) Delete(options DeleteOptions) error {
	switch {
	case len(options.Ids) > 0:
		return c.deleteByIds(options)
	case options.Pattern != "":
		return c.deleteByPattern(options)
	case options.UpdatedBefore != nil:
		return c.deleteUpdatedBefore(options.UpdatedBefore)
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
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("client:Delete:deleteByPattern: get rows affected: %w", err)
	}
	fmt.Printf("Deleted %d entries\n", rowsAffected)
	return nil
}

func buildDeleteByPatternStatement(pattern string) string {
	return fmt.Sprintf("DELETE FROM hist WHERE regexp('%s', command) = true", pattern)
}

func (c *Client) deleteUpdatedBefore(t *time.Time) error {
	if t == nil {
		return fmt.Errorf("client:Delete:deleteUpdatedBedore: no updatedBefore time provided")
	}
	res, err := c.Db.Exec(`DELETE FROM hist WHERE last_update < ?`, &t)
	if err != nil {
		return fmt.Errorf("client:Delete:deleteUpdatedBefore: exec: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("client:Delete:deleteByPattern: get rows affected: %w", err)
	}
	fmt.Printf("Deleted %d entries\n", rowsAffected)
	return nil
}
