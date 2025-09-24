package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(purple))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle
	noStyle      = lipgloss.NewStyle()
)

type createModel struct {
	focusIndex int
	inputs     []textinput.Model
}

func newCreateModel() createModel {
	m := createModel{inputs: make([]textinput.Model, 4)}

	var t textinput.Model

	for i := range m.inputs {
		t = textinput.New()

		t.Cursor.Style = cursorStyle
		t.CharLimit = 48

		switch i {
		case 0:
			t.Placeholder = "Name (Capitalized)"

			t.Focus()

			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Surname (Capitalized)"
		case 2:
			t.Placeholder = "Age (number)"
		case 3:
			t.Placeholder = "Grades CSV (e.g. 70,85,90)"
			t.CharLimit = 128
		}

		m.inputs[i] = t
	}

	return m
}

func (m createModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m createModel) Update(msg tea.Msg) (createModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return createCancelMsg{} }
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, func() tea.Msg {
					return createSubmittedMsg{
						Name:      strings.TrimSpace(m.inputs[0].Value()),
						Surname:   strings.TrimSpace(m.inputs[1].Value()),
						Age:       strings.TrimSpace(m.inputs[2].Value()),
						GradesCSV: strings.TrimSpace(m.inputs[3].Value()),
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

func (m createModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m createModel) View() string {
	var b strings.Builder

	b.WriteString("Create student\n\n")

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
