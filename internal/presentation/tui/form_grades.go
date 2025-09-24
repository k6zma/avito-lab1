package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type addGradesModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func newAddGradesModel() addGradesModel {
	m := addGradesModel{inputs: make([]textinput.Model, 2)}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()

		t.Cursor.Style = cursorStyle
		t.CharLimit = 128

		switch i {
		case 0:
			t.Placeholder = "Student ID UUID"

			t.Focus()

			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Grades separated example: 80,90"
		}

		m.inputs[i] = t
	}

	return m
}

func (m addGradesModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addGradesModel) Update(msg tea.Msg) (addGradesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return addGradesCancelMsg{} }
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				id := strings.TrimSpace(m.inputs[0].Value())
				grades := strings.TrimSpace(m.inputs[1].Value())

				return m, func() tea.Msg {
					return addGradesSubmittedMsg{
						ID:     id,
						Grades: grades,
					}
				}
			}

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()

					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
				} else {
					m.inputs[i].Blur()

					m.inputs[i].PromptStyle = noStyle
					m.inputs[i].TextStyle = noStyle
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m addGradesModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m addGradesModel) View() string {
	var b strings.Builder

	b.WriteString("Add grades\n\n")

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := blurredStyle.Render("[ Submit ]")
	if m.focusIndex == len(m.inputs) {
		button = focusedStyle.Render("[ Submit ]")
	}

	b.WriteString("\n\n" + button + "\n")
	b.WriteString(helpStyle.Render("tab/shift+tab to move | enter to submit | esc to back"))

	return b.String()
}
