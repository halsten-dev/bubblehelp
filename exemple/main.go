package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/halsten-dev/bubblehelp"
	"log"
	"strings"
)

const (
	StartContext bubblehelp.KeymapContext = "start"
	NextContext  bubblehelp.KeymapContext = "next"
)

var (
	ScreenWidth  = 0
	ScreenHeight = 0

	EscKey = key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	)
	SKey = key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "s key"),
	)
	HKey = key.NewBinding(
		key.WithKeys("h"),
		key.WithHelp("h", "help"),
	)
	EnterKey = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	)
)

func main() {
	bubblehelp.Init()

	startKeymap := bubblehelp.NewKeymap(2)

	startKeymap.NewKeyBinding(SKey, true)
	startKeymap.SetHelpDesc(SKey, "switch context")
	startKeymap.NewKeyBinding(EscKey, false)
	startKeymap.SetHelpDesc(EscKey, "quit")
	startKeymap.NewKeyBinding(HKey, true)

	bubblehelp.RegisterContext(StartContext, startKeymap)

	nextKeymap := bubblehelp.NewKeymap(2)

	nextKeymap.NewKeyBinding(EscKey, true)
	nextKeymap.NewKeyBinding(EnterKey, true)
	nextKeymap.SetHelpDesc(EnterKey, "hide/show esc key")
	nextKeymap.NewKeyBinding(HKey, true)

	bubblehelp.RegisterContext(NextContext, nextKeymap)

	bubblehelp.SwitchContext(StartContext)

	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SKey):
			if bubblehelp.IsKeybindVisible(SKey) {
				bubblehelp.SwitchContext(NextContext)
			}

			return m, nil

		case key.Matches(msg, EnterKey):
			if bubblehelp.IsKeybindVisible(EnterKey) {
				bubblehelp.SetKeybindVisible(EscKey, !bubblehelp.IsKeybindVisible(EscKey))
			}

			return m, nil

		case key.Matches(msg, EscKey):
			if bubblehelp.IsKeybindVisible(EscKey) {
				if bubblehelp.CurrentContext == StartContext {
					return m, tea.Quit
				}

				bubblehelp.SwitchContext(StartContext)
			}

			return m, nil

		case key.Matches(msg, HKey):
			bubblehelp.ShowAll = !bubblehelp.ShowAll

			return m, nil
		}

	case tea.WindowSizeMsg:
		ScreenWidth = msg.Width
		ScreenHeight = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString("bubblehelp demo")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Current context : %s", bubblehelp.CurrentContext))
	b.WriteString("\n\n")
	b.WriteString(bubblehelp.View(120))

	return lipgloss.Place(
		ScreenWidth, ScreenHeight,
		lipgloss.Center, lipgloss.Center,
		b.String(),
	)
}
