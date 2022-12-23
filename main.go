package main

import (
	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/herzrasen/hist/args"
	"github.com/herzrasen/hist/client"
	"github.com/herzrasen/hist/handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	var a args.Args
	arg.MustParse(&a)
	dbPath, err := xdg.DataFile("hist/hist.db")
	if err != nil {
		log.WithError(err).Fatal("Unable to get data dir")
	}
	c, err := client.NewSqliteClient(dbPath)
	if err != nil {
		log.WithError(err).Fatal("Unable to create client client")
	}
	h := handler.Handler{Client: c}
	err = h.Handle(a)
	if err != nil {
		log.WithError(err).Fatal("Error executing command")
	}
}
