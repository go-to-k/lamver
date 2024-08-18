package io

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
)

type UI struct {
	Choices   []string
	Header    string
	Cursor    int
	Selected  map[int]struct{}
	Keyword   string
	isEntered bool
}

var _ tea.Model = (*UI)(nil)

func NewUI(choices []string, header string) *UI {
	return &UI{
		Choices:  choices,
		Header:   header,
		Selected: make(map[int]struct{}),
	}
}

func (u *UI) Init() tea.Cmd {
	return nil
}

func (u *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//nolint:gocritic
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		// Quit the selection
		case "enter":
			u.isEntered = true
			return u, tea.Quit

		// Quit the selection and clear ALL selected items
		case "ctrl+c":
			u.Selected = make(map[int]struct{})
			return u, tea.Quit

		case "up":
			if u.Cursor >= 0 {
				u.Cursor--
			}
			if u.Cursor == -1 {
				u.Cursor = len(u.Choices) - 1
			}

		case "down":
			if u.Cursor <= len(u.Choices)-1 {
				u.Cursor++
			}
			if u.Cursor == len(u.Choices) {
				u.Cursor = 0
			}

		// select or deselect an item
		case " ":
			_, ok := u.Selected[u.Cursor]
			if ok {
				delete(u.Selected, u.Cursor)
			} else {
				u.Selected[u.Cursor] = struct{}{}
			}

		// select all items
		case "right":
			for i := range u.Choices {
				u.Selected[i] = struct{}{}
			}

		// clear all selected items
		case "left":
			u.Selected = make(map[int]struct{})

		case "backspace":
			if len(u.Keyword) > 0 {
				u.Keyword = u.Keyword[:len(u.Keyword)-1]
			}

		default:
			u.Keyword += msg.String()

		}
	}

	return u, nil
}

func (u *UI) View() string {
	s := color.CyanString("? ") + color.New(color.Bold).Sprint(u.Header)

	if u.isEntered {
		return s
	}

	s += u.Keyword

	s += color.CyanString(" [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]\n")

	for i, choice := range u.Choices {
		if u.Keyword != "" && !strings.Contains(choice, u.Keyword) {
			continue
		}

		cursor := " " // no cursor
		if u.Cursor == i {
			cursor = color.CyanString(">") // cursor!
		}

		checked := "[ ]" // not selected
		if _, ok := u.Selected[i]; ok {
			checked = color.GreenString("[x]") // selected!
		}

		s += fmt.Sprintf("%s %s %s\n", cursor, checked, choice)
	}

	return s
}