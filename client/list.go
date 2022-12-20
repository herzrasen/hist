package client

import (
	"fmt"
	"github.com/herzrasen/hist/record"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (c *Client) List(options ListOptions) ([]record.Record, error) {
	rows, err := c.db.Query(buildListQuery(options.Limit))
	if err != nil {
		return nil, fmt.Errorf("client.Client.List: query: %w", err)
	}
	defer rows.Close()
	var records []record.Record
	for rows.Next() {
		var r record.Record
		err := rows.Scan(&r.Id, &r.Command, &r.LastUpdate, &r.Count)
		if err != nil {
			log.WithError(err).Warn("Unable to scan row")
		} else {
			records = append(records, r)
		}
	}
	return records, nil
}

func buildListQuery(limit int) string {
	query := strings.Builder{}
	query.WriteString(`SELECT id, command, last_update, count FROM hist ORDER BY last_update DESC`)
	if limit > 0 {
		query.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}
	return query.String()
}
