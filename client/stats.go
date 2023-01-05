package client

import (
	"fmt"
	"github.com/herzrasen/hist/record"
	"github.com/herzrasen/hist/stats"
	"time"
)

func (c *Client) Stats() (*stats.Stats, error) {
	records, err := c.List(ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("client:Stats: list: %w", err)
	}
	highestCount, mostUsedCommands := findMostUsed(records)
	newestCommandTime := findNewestCommandTime(records)
	newestCommands := findCommandsWithLastUpdateTime(records, newestCommandTime)
	oldestCommandTime := findOldestCommandTime(records)
	oldestCommands := findCommandsWithLastUpdateTime(records, oldestCommandTime)
	return &stats.Stats{
		NumCommands:       uint64(len(records)),
		HighestCount:      highestCount,
		MostUsedCommands:  mostUsedCommands,
		NewestCommandTime: newestCommandTime,
		NewestCommands:    newestCommands,
		OldestCommandTime: oldestCommandTime,
		OldestCommands:    oldestCommands,
	}, nil
}

func findMostUsed(records []record.Record) (uint64, []string) {
	var mostUsed []string
	var highestCount uint64
	for _, r := range records {
		if r.Count > highestCount {
			mostUsed = mostUsed[:0]
			highestCount = r.Count
		}
		if r.Count == highestCount {
			mostUsed = append(mostUsed, r.Command)
		}
	}
	return highestCount, mostUsed
}

func findNewestCommandTime(records []record.Record) time.Time {
	var newest time.Time
	for _, r := range records {
		if r.LastUpdate.After(newest) {
			newest = r.LastUpdate
		}
	}
	return newest
}

func findOldestCommandTime(records []record.Record) time.Time {
	oldest := time.Now()
	for _, r := range records {
		if r.LastUpdate.Before(oldest) {
			oldest = r.LastUpdate
		}
	}
	return oldest
}

func findCommandsWithLastUpdateTime(records []record.Record, lastUpdate time.Time) []string {
	var commands []string
	for _, r := range records {
		if r.LastUpdate == lastUpdate {
			commands = append(commands, r.Command)
		}
	}
	return commands
}
