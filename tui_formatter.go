package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	item struct {
		contents string
		done     bool
	}

	section struct {
		header string
		items  []*item
		count  int
		done   bool
	}
)

type tuiFormatter struct {
	sections []*section
	output   []string
}

func newTuiFormatter() *tuiFormatter {
	return &tuiFormatter{
		sections: []*section{},
	}
}

func (tf *tuiFormatter) FormatSection(header string, count int, items []string) {
	uiItems := []*item{}
	for _, it := range items {
		uiItems = append(uiItems, &item{it, false})
	}
	tf.sections = append(tf.sections, &section{
		header: header,
		items:  uiItems,
		count:  count,
	})
}

func (tf *tuiFormatter) Output() error {
	// start bubble tea app
	p := tea.NewProgram(initialModel(tf.sections))
	_, err := p.Run()
	return err
}

type model struct {
	sections                  []*section
	cursorSection, cursorItem int
}

func initialModel(sections []*section) model {
	return model{
		sections:      sections,
		cursorSection: 0, cursorItem: -1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// TODO: Add ability to export this to a markdown file for my phone
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursorSection <= 0 && m.cursorItem <= -1 {
				return m, nil
			}
			if m.cursorItem > -1 {
				m.cursorItem--
				return m, nil
			}

			m.cursorSection--
			m.cursorItem = len(m.sections[m.cursorSection].items) - 1

		case "down", "j":
			currentItems := m.sections[m.cursorSection].items
			if m.cursorSection >= len(m.sections)-1 && m.cursorItem >= len(currentItems)-1 {
				return m, nil
			}

			if m.cursorItem < len(currentItems)-1 {
				m.cursorItem++
				return m, nil
			}

			m.cursorItem = -1
			m.cursorSection++

		case "enter", " ":
			if m.cursorItem != -1 {
				currSec := m.sections[m.cursorSection]
				item := currSec.items[m.cursorItem]
				item.done = !item.done
				for _, item := range currSec.items {
					if !item.done {
						currSec.done = false
						return m, nil
					}
				}
				currSec.done = true
				return m, nil
			}
			for i, section := range m.sections {
				if i == m.cursorSection {
					section.done = !section.done

					for _, item := range section.items {
						item.done = section.done
					}
					return m, nil
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "Packing list\n\n"

	for i, section := range m.sections {
		s += fmt.Sprintln(
			cursor(func() bool { return m.cursorItem == -1 && m.cursorSection == i }),
			checkbox(section.done),
			section.header,
		)
		for j, item := range section.items {
			s += fmt.Sprintln(
				cursor(func() bool { return m.cursorSection == i && m.cursorItem == j }),
				"-",
				checkbox(item.done),
				item.contents,
			)
		}
		s += fmt.Sprintln()
	}

	s += "\nPress q to quit.\n"

	return s
}

func cursor(showFunc func() bool) string {
	cursor := " "
	if showFunc() {
		cursor = ">"
	}
	return cursor
}
func checkbox(done bool) string {
	cross := " "
	if done {
		cross = "x"
	}

	return fmt.Sprintf("[%s]", cross)
}
