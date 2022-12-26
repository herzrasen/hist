package handler

import (
	"fmt"
	"github.com/herzrasen/hist/args"
	"github.com/herzrasen/hist/client"
	"github.com/herzrasen/hist/config"
	"github.com/herzrasen/hist/record"
	log "github.com/sirupsen/logrus"
)

type HistClient interface {
	List(options client.ListOptions) ([]record.Record, error)
	Update(command string) error
	Delete(options client.DeleteOptions) error
}

type Handler struct {
	Client HistClient
	Config *config.Config
}

func (h *Handler) Handle(a args.Args) error {
	switch {
	case a.List != nil:
		options := client.ListOptions{
			NoCount:      a.List.NoCount,
			NoLastUpdate: a.List.NoLastUpdate,
			WithId:       a.List.WithId,
			Limit:        a.List.Limit,
		}
		records, err := h.Client.List(options)
		if err != nil {
			return fmt.Errorf("unable to list records: %w", err)
		}
		fmt.Printf("%s", options.ToString(records))
	case a.Record != nil:
		command := a.Record.Command
		if h.Config.IsExcluded(command) {
			log.WithField("command", command).Info("Skipping because command is excluded")
			return nil
		}
		err := h.Client.Update(command)
		if err != nil {
			return fmt.Errorf("unable to record command: %w", err)
		}
	case a.Delete != nil:
		deleteCmd := a.Delete
		options := client.DeleteOptions{
			Ids:    deleteCmd.Ids,
			Filter: deleteCmd.Filter,
		}
		err := h.Client.Delete(options)
		if err != nil {
			return fmt.Errorf("unable to delete entry: %w", err)
		}
	}
	return nil
}
