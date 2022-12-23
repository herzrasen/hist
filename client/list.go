package client

import (
	"fmt"
	"github.com/herzrasen/hist/record"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
)

const (
	stmtSelect      = `SELECT id, command, last_update, count FROM hist ORDER BY last_update DESC`
	stmtSelectByIds = `SELECT id, command, last_update, count FROM hist WHERE id IN (?) ORDER BY last_update DESC`
)

type ListOptions struct {
	NoCount      bool
	NoLastUpdate bool
	WithId       bool
	Limit        int
	Ids          []int64
}

func (c *Client) List(options ListOptions) ([]record.Record, error) {
	var statement string
	var args []interface{}
	var err error

	if len(options.Ids) > 0 {
		statement, args, err = buildByIdsQuery(options)
	} else {
		statement, args, err = buildListQuery(options)
	}

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
	options.sort(records)
	return records, nil
}

func buildListQuery(options ListOptions) (string, []interface{}, error) {
	query := strings.Builder{}
	query.WriteString(stmtSelect)
	if options.Limit > 0 {
		query.WriteString(" LIMIT ?")
		return query.String(), []interface{}{options.Limit}, nil
	}
	return query.String(), nil, nil
}

func buildByIdsQuery(options ListOptions) (string, []interface{}, error) {
	statement := strings.Builder{}
	statement.WriteString(stmtSelectByIds)
	if options.Limit > 0 {
		statement.WriteString(" LIMIT ?")
		return sqlx.In(statement.String(), options.Ids, options.Limit)
	}
	return sqlx.In(statement.String(), options.Ids)
}

func (l *ListOptions) sort(records []record.Record) {
	sort.Slice(records, func(i, j int) bool {
		left := records[i]
		right := records[j]
		return left.LastUpdate.Before(right.LastUpdate)
	})
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
