package client

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/herzrasen/hist/record"
	"sort"
	"strings"
	"time"
)

type ListOptions struct {
	NoCount      bool
	NoLastUpdate bool
	WithId       bool
	Limit        int
}

func (l *ListOptions) ToString(records []record.Record) string {
	sort.Slice(records, func(i, j int) bool {
		left := records[i]
		right := records[j]
		return left.LastUpdate.Before(right.LastUpdate)
	})
	buf := strings.Builder{}
	for _, r := range records {
		if !l.NoLastUpdate {
			buf.WriteString(color.GreenString("%s\t", r.LastUpdate.Format(time.RFC1123)))
		}
		if !l.NoCount {
			buf.WriteString(color.BlueString("%d\t", r.Count))
		}
		if l.WithId {
			buf.WriteString(color.YellowString("%d\t", r.Id))
		}
		buf.WriteString(fmt.Sprintf("%s\n", r.Command))
	}
	return buf.String()
}
