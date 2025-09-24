package tui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const purple = "170"

var helpStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	PaddingLeft(2).
	PaddingTop(1)

var errorStyle = helpStyle.Copy().
	Foreground(lipgloss.Color("196"))

func styleTable(t table.Model) table.Model {
	s := table.DefaultStyles()

	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color(purple)).
		Bold(false)

	t.SetStyles(s)

	return t
}
