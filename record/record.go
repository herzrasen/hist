package record

import (
	"time"
)

type Record struct {
	Id         int64
	Command    string
	LastUpdate time.Time
	Count      uint64
}
