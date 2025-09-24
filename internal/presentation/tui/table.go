package tui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table table.Model
}

func newTableModel(t table.Model) tableModel {
	return tableModel{
		table: styleTable(t),
	}
}

func (m tableModel) Init() tea.Cmd {
	return nil
}

func (m tableModel) Update(msg tea.Msg) (tableModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "ctrl+c":
			return m, func() tea.Msg {
				return tableBackMsg{}
			}
		case "enter":
			row := m.table.SelectedRow()
			if len(row) > 1 {
				id := row[1]

				return m, func() tea.Msg { return tableShowMsg{ID: id} }
			}

			return m, nil
		}
	}

	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m tableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n" +
		helpStyle.Render("↑/↓ to move | esc/q to back | enter to show")
}
