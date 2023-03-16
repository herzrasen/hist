package args

import (
	"fmt"
	"time"
)

// the following vars will be set by goreleaser on build time
var (
	Version = "unknown"
	Commit  = "unknown"
	Date    = "unknown"
)

type RecordCmd struct {
	Command string `arg:"positional"`
}

type SearchCmd struct {
	Input   string `arg:"positional"`
	Verbose bool   `arg:"--verbose"`
}

type ListCmd struct {
	Pattern      string `arg:"--pattern"`
	ByCount      bool   `arg:"--by-count"`
	Reverse      bool   `arg:"--reverse"`
	NoCount      bool   `arg:"--no-count"`
	NoLastUpdate bool   `arg:"--no-last-update"`
	WithId       bool   `arg:"--with-id"`
	Limit        int    `arg:"-l,--limit" default:"-1"`
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
	Delete *DeleteCmd `arg:"subcommand:delete" help:"Delete commands from history"`
	Get    *GetCmd    `arg:"subcommand:get" help:"Get a command by it's index'"`
	Import *ImportCmd `arg:"subcommand:import" help:"Import commands from a legacy history file"`
	List   *ListCmd   `arg:"subcommand:list" help:"List commands"`
	Record *RecordCmd `arg:"subcommand:record" help:"Record a new command"`
	Search *SearchCmd `arg:"subcommand:search" help:"Start the interactive fuzzy selection mode"`
	Stats  *StatsCmd  `arg:"subcommand:stats" help:"Show some statistics"`
	Tidy   *TidyCmd   `arg:"subcommand:tidy" help:"Apply exclude patterns to clean up the hist database"`
	Config string     `arg:"--config" default:"~/.config/hist/config.yml"`
}

func (a *Args) Version() string {
	return fmt.Sprintf("hist %s\nCommit: %s\nDate: %s", Version, Commit, Date)
}
