package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type idInputModel struct {
	textInput textinput.Model
	label     string
}

func newIDInputModel(label string) idInputModel {
	ti := textinput.New()

	ti.Placeholder = label

	ti.Focus()

	ti.CharLimit = 156
	ti.Width = 40

	return idInputModel{
		textInput: ti,
		label:     label,
	}
}

func (m idInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m idInputModel) Update(msg tea.Msg) (idInputModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg {
				return idCancelMsg{}
			}
		case "enter":
			id := m.textInput.Value()

			return m, func() tea.Msg {
				return idSubmittedMsg(id)
			}
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m idInputModel) View() string {
	return fmt.Sprintf(
		"%s:\n\n%s\n\n%s",
		m.label,
		m.textInput.View(),
		helpStyle.Render("type the ID | enter to submit | esc to back"),
	)
}
