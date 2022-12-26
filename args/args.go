package args

type RecordCmd struct {
	Command string `arg:"positional"`
}

type ListCmd struct {
	NoCount      bool `arg:"--no-count"`
	NoLastUpdate bool `arg:"--no-last-update"`
	WithId       bool `arg:"--with-id"`
	Limit        int  `arg:"-l,--limit" default:"-1"`
}

type DeleteCmd struct {
	Ids      []int64 `arg:"-i,--id"`
	Filter   string  `arg:"-f,--filter"`
	MaxCount int64   `arg:"--max-count" help:"Delete all records with a count of at most max-count"`
}

type Args struct {
	Record *RecordCmd `arg:"subcommand:record"`
	List   *ListCmd   `arg:"subcommand:list"`
	Delete *DeleteCmd `arg:"subcommand:delete"`
	Config string     `arg:"--config" default:"~/.config/hist/config.yml"`
}
