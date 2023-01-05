package client

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"strconv"
	"strings"
	"time"
)

type HistoryEntry struct {
	Date    time.Time
	Command string
}

func (c *Client) Import(reader io.Reader) error {
	entries := parseHistFile(reader)
	for _, entry := range entries {
		err := c.RecordWithTime(entry.Command, entry.Date)
		if err != nil {
			log.WithError(err).Warn("Unable to import command")
		}
	}
	return nil
}

func parseHistFile(reader io.Reader) []HistoryEntry {
	scanner := bufio.NewScanner(reader)
	scanner.Split(splitHistFile)
	var entries []HistoryEntry
	for scanner.Scan() {
		entry, err := parseHistoryEntry(scanner.Text())
		if err != nil {
			log.WithError(err).Warn("Unable to parse history entry")
		} else {
			entries = append(entries, entry)
		}
	}
	return entries
}

func splitHistFile(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	indexOfNewLine := indexOfNewlineWithoutLeadingBackslash(data, 0)
	if indexOfNewLine >= 0 {
		return indexOfNewLine + 1, data[0:indexOfNewLine], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	// request more data
	return 0, nil, nil
}

func indexOfNewlineWithoutLeadingBackslash(data []byte, accOffset int) int {
	i := bytes.IndexByte(data, '\n')
	switch {
	case i == 0:
		return accOffset
	case i > 0 && data[i-1] == '\\':
		return indexOfNewlineWithoutLeadingBackslash(data[i+1:], accOffset+i+1)
	case i > 0:
		return accOffset + i
	default:
		return -1
	}
}

func parseHistoryEntry(line string) (HistoryEntry, error) {
	// not the extended format
	if !strings.HasPrefix(line, ":") {
		return HistoryEntry{
			Date:    time.Now(),
			Command: line,
		}, nil
	}
	// extended format
	tokens := strings.SplitN(line, ";", 2)
	if len(tokens) != 2 {
		return HistoryEntry{}, fmt.Errorf("unable to parse %s to ExtendedHistoryEntry", line)
	}
	command := tokens[1]
	t := strings.TrimPrefix(tokens[0], ":")
	tokens = strings.SplitN(t, ":", 2)
	if len(tokens) != 2 {
		return HistoryEntry{}, fmt.Errorf("illegal time format in line %s", line)
	}
	unixSecString := strings.TrimSpace(tokens[0])
	unixSec, err := strconv.ParseInt(unixSecString, 10, 64)
	if err != nil {
		return HistoryEntry{}, fmt.Errorf("unable to parse unix time from %s", unixSecString)
	}
	commandTime := time.Unix(unixSec, 0)
	return HistoryEntry{
		Date:    commandTime,
		Command: command,
	}, nil
}
