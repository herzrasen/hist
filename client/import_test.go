package client

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestScanHistoryFile(t *testing.T) {
	t.Run("simple lines", func(t *testing.T) {
		content := `my-simple-command-1
my-simple-command-2`
		scanner := bufio.NewScanner(strings.NewReader(content))
		scanner.Split(splitHistFile)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		assert.Len(t, lines, 2)
		assert.Contains(t, lines, "my-simple-command-1")
		assert.Contains(t, lines, "my-simple-command-2")
	})

	t.Run("multiline string", func(t *testing.T) {
		content := `my-simple-command-1\
with-additional-stuff\
and-even-more-additional-stuff
my-simple-command-2`
		scanner := bufio.NewScanner(strings.NewReader(content))
		scanner.Split(splitHistFile)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		assert.Len(t, lines, 2)
		assert.Contains(t, lines, "my-simple-command-1\\\nwith-additional-stuff\\\nand-even-more-additional-stuff")
		assert.Contains(t, lines, "my-simple-command-2")
	})
}
