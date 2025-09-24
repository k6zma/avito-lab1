package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type detailModel struct {
	body string
}

func newDetailModel(lines []string) detailModel {
	return detailModel{
		body: strings.Join(lines, "\n"),
	}
}

func (m detailModel) Init() tea.Cmd {
	return nil
}

func (m detailModel) Update(msg tea.Msg) (detailModel, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg { return tableBackMsg{} }
	}

	return m, nil
}

func (m detailModel) View() string {
	return m.body + "\n\n" + helpStyle.Render("any key/esc to back")
}
