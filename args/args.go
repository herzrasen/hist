package args

import (
	"time"
)

type RecordCmd struct {
	Command string `arg:"positional"`
}

type SearchCmd struct {
}

type ListCmd struct {
	ByCount      bool `arg:"--by-count"`
	Reverse      bool `arg:"--reverse"`
	NoCount      bool `arg:"--no-count"`
	NoLastUpdate bool `arg:"--no-last-update"`
	WithId       bool `arg:"--with-id"`
	Limit        int  `arg:"-l,--limit" default:"-1"`
}

type GetCmd struct {
	Index int64 `arg:"--index"`
}

type DeleteCmd struct {
	Ids           []int64    `arg:"-i,--id"`
	UpdatedBefore *time.Time `arg:"-u,--updated-before"`
	Pattern       string     `arg:"-p,--pattern" help:"Delete all records matching the pattern"`
}

type ImportCmd struct {
	Path string `arg:"positional"`
}

type TidyCmd struct {
}

type StatsCmd struct {
}

type Args struct {
	Delete *DeleteCmd `arg:"subcommand:delete"`
	Get    *GetCmd    `arg:"subcommand:get"`
	Import *ImportCmd `arg:"subcommand:import"`
	List   *ListCmd   `arg:"subcommand:list"`
	Record *RecordCmd `arg:"subcommand:record"`
	Search *SearchCmd `arg:"subcommand:search"`
	Stats  *StatsCmd  `arg:"subcommand:stats"`
	Tidy   *TidyCmd   `arg:"subcommand:tidy"`
	Config string     `arg:"--config" default:"~/.config/hist/config.yml"`
}
