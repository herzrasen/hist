package handler

import (
	"fmt"
	"github.com/herzrasen/hist/args"
	"github.com/herzrasen/hist/client"
	"github.com/herzrasen/hist/config"
	"github.com/herzrasen/hist/record"
	"github.com/herzrasen/hist/search"
)

type HistClient interface {
	List(options client.ListOptions) ([]record.Record, error)
	Get(index int64) (string, error)
	Record(command string) error
	Delete(options client.DeleteOptions) error
}

type Handler struct {
	Client HistClient
	Config *config.Config
}

func (h *Handler) Handle(a args.Args) error {
	switch {
	case a.Record != nil:
		command := a.Record.Command
		if h.Config.IsExcluded(command) {
			return nil
		}
		err := h.Client.Record(command)
		if err != nil {
			return fmt.Errorf("unable to record command: %w", err)
		}
	case a.Get != nil:
		command, err := h.Client.Get(a.Get.Index)
		if err != nil {
			return fmt.Errorf("unable to get command by index: %w", err)
		}
		fmt.Printf("%s", command)
	case a.Search != nil:
		searcher := search.NewSearcher(h.Client)
		err := searcher.Show()
		if err != nil {
			return fmt.Errorf("unable to show search dialog: %w", err)
		}
	case a.List != nil:
		options := client.ListOptions{
			Reverse:      a.List.Reverse,
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
	case a.Delete != nil:
		deleteCmd := a.Delete
		options := client.DeleteOptions{
			Ids:     deleteCmd.Ids,
			Pattern: deleteCmd.Pattern,
		}
		err := h.Client.Delete(options)
		if err != nil {
			return fmt.Errorf("unable to delete entry: %w", err)
		}
	}
	return nil
}
