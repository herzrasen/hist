package search

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/herzrasen/hist/client"
	"github.com/herzrasen/hist/record"
	"github.com/rivo/tview"
)

type ListClient interface {
	List(options client.ListOptions) ([]record.Record, error)
}

type Searcher struct {
	ListClient ListClient
	App        *tview.Application
	List       *tview.List
	Input      *tview.InputField
	Flex       *tview.Flex
}

func NewSearcher(listClient ListClient) *Searcher {
	list := tview.NewList().
		ShowSecondaryText(false).
		SetShortcutStyle(tcell.Style{}).
		SetSelectedStyle(tcell.StyleDefault.
			Foreground(tcell.ColorPaleGreen)).
		SetMainTextColor(tcell.ColorViolet).
		SetWrapAround(true)
	input := tview.NewInputField()
	input.SetFieldStyle(tcell.StyleDefault.
		Italic(true).
		Foreground(tcell.ColorLightBlue))
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 20, false).
		AddItem(input, 0, 2, true)
	app := tview.NewApplication().
		SetRoot(flex, true).
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			numItems := list.GetItemCount() - 1
			currentItem := list.GetCurrentItem()
			switch event.Key() {
			case tcell.KeyEnter:
				selectedIndex := list.GetCurrentItem()
				selected, _ := list.GetItemText(selectedIndex)
				fmt.Printf("%s", selected)
				return event
			case tcell.KeyUp:
				nextItem := 0
				if currentItem > 0 {
					nextItem = currentItem - 1
				}
				list.SetCurrentItem(nextItem)
				return nil
			case tcell.KeyDown:
				nextItem := currentItem + 1
				if currentItem >= numItems {
					nextItem = numItems
				}
				list.SetCurrentItem(nextItem)
				return nil
			}
			return event
		})
	input.SetDoneFunc(func(key tcell.Key) {
		app.Stop()
	})
	return &Searcher{
		ListClient: listClient,
		App:        app,
		List:       list,
		Input:      input,
		Flex:       flex,
	}
}

func (s *Searcher) Show(input string, verbose bool) error {
	recs, err := s.ListClient.List(client.ListOptions{
		Reverse: false,
	})
	if err != nil {
		return fmt.Errorf("list: %w", err)
	}
	records := record.Records(recs)
	s.showRecords(records, verbose)
	s.Input.SetChangedFunc(func(text string) {
		updatedRecords := records.Search(text)
		s.showRecords(updatedRecords, verbose)
	})
	s.Input.SetText(input)
	if err := s.App.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}
	return nil
}

func (s *Searcher) showRecords(records record.Records, verbose bool) {
	s.List.Clear()
	for _, updatedRecord := range records {
		item := updatedRecord.Command
		if verbose {
			item = fmt.Sprintf("(weight: %d, len: %d) %s",
				updatedRecord.Weight,
				len(records),
				updatedRecord.Command,
			)
		}
		s.List.AddItem(item, "", 0, nil)
	}
}
