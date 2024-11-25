package client

import "fmt"

func (c *Client) Tag(id int64, tags []string, remove bool) error {
	if !remove {
		stmt, err := c.Db.Prepare(`INSERT OR IGNORE INTO tag (tag, hist_id) VALUES (?, ?)`)
		if err != nil {
			return fmt.Errorf("unable to prepare statement: %w", err)
		}
		defer func() {
			_ = stmt.Close()
		}()
		for _, tag := range tags {
			_, err = stmt.Exec(tag, id)
			if err != nil {
				return fmt.Errorf("unable to insert into '%s' tags: %w", tag, err)
			}
		}
	} else {
		stmt, err := c.Db.Prepare(`DELETE FROM tag WHERE tag = ? AND hist_id = ?`)
		if err != nil {
			return fmt.Errorf("unable to prepare statement: %w", err)
		}
		defer func() {
			_ = stmt.Close()
		}()
		for _, tag := range tags {
			_, err = stmt.Exec(tag, id)
			if err != nil {
				return fmt.Errorf("unable to remove tag '%s' from %d: %w", tag, id, err)
			}
		}
	}
	return nil
}
