package client

import (
	"bufio"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/herzrasen/hist/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestClient_Import(t *testing.T) {
	t.Run("succeed", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		mock.ExpectExec(`INSERT INTO hist (command, last_update) 
			VALUES (?, ?)
			ON CONFLICT(command) 
			    DO UPDATE SET count=count+1, last_update=excluded.last_update`).
			WithArgs("git push", AnyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec(`INSERT INTO hist (command, last_update) 
			VALUES (?, ?)
			ON CONFLICT(command) 
			    DO UPDATE SET count=count+1, last_update=excluded.last_update`).
			WithArgs("git pull", AnyTime{}).
			WillReturnResult(sqlmock.NewResult(2, 1))
		content := `git push
: 1671208522:0;git pull`
		c := Client{
			Db:     sqlx.NewDb(db, "sqlite3"),
			Config: &config.Config{},
		}
		err := c.Import(strings.NewReader(content))
		require.NoError(t, err)
	})
}

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

func TestClient_parseHistoryEntry(t *testing.T) {
	t.Run("error parsing timestamp", func(t *testing.T) {
		_, err := parseHistoryEntry(": not a date:0; foo")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unable to parse unix time from")
	})

	t.Run("illegal time format", func(t *testing.T) {
		_, err := parseHistoryEntry(": missing second colon; foo")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "illegal time format in line")
	})
}
