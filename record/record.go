package record

import (
	"time"
)

type Record struct {
	Id         int
	Command    string
	LastUpdate time.Time
	Count      uint64
}
