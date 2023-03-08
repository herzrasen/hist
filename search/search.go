package search

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/herzrasen/hist/client"
	"github.com/herzrasen/hist/record"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
	"sort"
	"strings"
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
		SetWrapAround(true)
	input := tview.NewInputField()
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(list, 0, 20, false).
		AddItem(input, 0, 2, true)
	app := tview.NewApplication().
		SetRoot(flex, true).
		EnableMouse(true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			numItems := list.GetItemCount() - 1
			currentItem := list.GetCurrentItem()
			oldItemText, _ := list.GetItemText(currentItem)
			list.SetItemText(currentItem, strings.TrimPrefix(oldItemText, "> "), "")
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
				itemText, _ := list.GetItemText(nextItem)
				list.SetItemText(nextItem, "> "+itemText, "")
				return nil
			case tcell.KeyDown:
				nextItem := currentItem + 1
				if currentItem >= numItems {
					nextItem = numItems
				}
				list.SetCurrentItem(nextItem)
				itemText, _ := list.GetItemText(nextItem)
				list.SetItemText(nextItem, "> "+itemText, "")
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

func (s *Searcher) Show(input string) error {
	records, err := s.ListClient.List(client.ListOptions{
		Reverse: true,
	})
	if err != nil {
		return fmt.Errorf("search:Searcher:Show: list: %w", err)
	}
	var commands []string
	for _, rec := range records {
		commands = append(commands, rec.Command)
	}
	for _, command := range commands {
		s.List.AddItem(command, "", 0, nil)
	}
	s.Input.SetChangedFunc(func(text string) {
		ranks := fuzzy.RankFind(text, commands)
		sort.Sort(ranks)
		s.List.Clear()
		for _, rankedCommand := range ranks {
			s.List.AddItem(rankedCommand.Target, "", 0, nil)
		}
	})
	currentItem, _ := s.List.GetItemText(0)
	s.List.SetCurrentItem(0)
	s.List.SetItemText(0, "> "+currentItem, "")
	s.Input.SetText(input)
	if err := s.App.Run(); err != nil {
		return fmt.Errorf("search:Searcher:Show: run: %w", err)
	}
	return nil
}
