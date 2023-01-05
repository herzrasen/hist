package stats

import (
	"encoding/json"
	"fmt"
	"time"
)

type Stats struct {
	NumCommands       uint64    `json:"numCommands"`
	HighestCount      uint64    `json:"highestCount"`
	MostUsedCommands  []string  `json:"mostUsedCommands"`
	OldestCommandTime time.Time `json:"oldestCommandTime"`
	OldestCommands    []string  `json:"oldestCommands"`
	NewestCommandTime time.Time `json:"newestCommandTime"`
	NewestCommands    []string  `json:"newestCommands"`
}

func (s *Stats) ToString() string {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s\n", string(data))
}
