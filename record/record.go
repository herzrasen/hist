package record

import (
	"github.com/fatih/color"
	"strings"
	"time"
)

type Record struct {
	Id         int64
	Command    string
	LastUpdate time.Time
	Count      uint64
}

type FormatOptions struct {
	NoLastUpdate bool
	NoCount      bool
	WithId       bool
}

func (r *Record) Format(options FormatOptions) string {
	buf := strings.Builder{}
	if !options.NoLastUpdate {
		buf.WriteString(color.GreenString("%s\t", r.LastUpdate.Format(time.RFC1123)))
	}
	if !options.NoCount {
		buf.WriteString(color.BlueString("%d\t", r.Count))
	}
	if options.WithId {
		buf.WriteString(color.YellowString("%d\t", r.Id))
	}
	buf.WriteString(r.Command)
	return buf.String()
}
