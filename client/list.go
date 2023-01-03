package client

import (
	"fmt"
	"github.com/herzrasen/hist/record"
	log "github.com/sirupsen/logrus"
	"strings"
)

type ListOptions struct {
	ByCount      bool
	Reverse      bool
	NoCount      bool
	NoLastUpdate bool
	WithId       bool
	Limit        int
}

func (c *Client) List(options ListOptions) ([]record.Record, error) {
	statement, args, err := buildListQuery(options)
	if err != nil {
		return nil, fmt.Errorf("client.Client.List: build query: %w", err)
	}

	rows, err := c.Db.Query(statement, args...)
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

	if !options.Reverse {
		records = reverse(records)
	}

	return records, nil
}

func buildListQuery(options ListOptions) (string, []interface{}, error) {
	query := strings.Builder{}
	query.WriteString(`SELECT id, command, last_update, count FROM hist`)
	if options.ByCount {
		query.WriteString(" ORDER BY count, last_update DESC")
	} else {
		query.WriteString(" ORDER BY last_update, count DESC")
	}
	if options.Limit > 0 {
		query.WriteString(" LIMIT ?")
		return query.String(), []interface{}{options.Limit}, nil
	}
	return query.String(), nil, nil
}

func reverse(records []record.Record) []record.Record {
	var reversed []record.Record
	for i := len(records) - 1; i >= 0; i-- {
		reversed = append(reversed, records[i])
	}
	return reversed
}

func (l *ListOptions) ToString(records []record.Record) string {
	options := record.FormatOptions{
		NoLastUpdate: l.NoLastUpdate,
		NoCount:      l.NoCount,
		WithId:       l.WithId,
	}
	buf := strings.Builder{}
	for _, r := range records {
		buf.WriteString(fmt.Sprintf("%s\n", r.Format(options)))
	}
	return buf.String()
}
