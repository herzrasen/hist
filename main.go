package main

import (
	"fmt"
	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/herzrasen/hist/client"
	log "github.com/sirupsen/logrus"
)

type RecordCmd struct {
	Command string `arg:"positional"`
}

type ListCmd struct {
	NoCount      bool `arg:"--no-count"`
	NoLastUpdate bool `arg:"--no-last-update"`
	WithId       bool `arg:"--with-id"`
	Limit        int  `arg:"-l,--limit" default:"-1"`
}

var args struct {
	Record *RecordCmd `arg:"subcommand:record"`
	List   *ListCmd   `arg:"subcommand:list"`
	Config string     `arg:"--config" default:"~/.config/hist/config.yml"`
}

func main() {
	arg.MustParse(&args)
	dbPath, err := xdg.DataFile("hist/hist.db")
	if err != nil {
		log.WithError(err).Fatal("Unable to get data dir")
	}
	c, err := client.NewClient(dbPath)
	if err != nil {
		log.WithError(err).Fatal("Unable to create client client")
	}

	switch {
	case args.List != nil:
		options := client.ListOptions{
			NoCount:      args.List.NoCount,
			NoLastUpdate: args.List.NoLastUpdate,
			WithId:       args.List.WithId,
			Limit:        args.List.Limit,
		}
		records, err := c.List(options)
		if err != nil {
			log.WithError(err).Fatal("Unable to list records")
		}
		fmt.Printf("%s", options.ToString(records))
	case args.Record != nil:
		command := args.Record.Command
		err = c.Update(command)
		if err != nil {
			log.WithError(err).WithField("command", command).Error("Unable to record command")
		}
	}
}
