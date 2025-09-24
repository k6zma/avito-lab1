package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minMenuHeight = 8
	defaultWidth  = 32
	menuTittle    = "STUDIFY - Choose an action:"
)

var (
	menuContainer     = lipgloss.NewStyle().PaddingLeft(1)
	titleStyle        = lipgloss.NewStyle()
	itemStyle         = lipgloss.NewStyle().PaddingLeft(1)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color(purple))
)

type menuItem string

func (i menuItem) FilterValue() string {
	return ""
}

type itemDelegate struct{}

func (d itemDelegate) Height() int {
	return 1
}

func (d itemDelegate) Spacing() int {
	return 0
}

func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(menuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)
	fn := itemStyle.Render

	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	_, err := io.WriteString(w, fn(str))
	if err != nil {
		return
	}
}

type menuModel struct {
	list     list.Model
	quitting bool
}

func newMenuModel(items []string) menuModel {
	its := make([]list.Item, 0, len(items))
	for _, it := range items {
		its = append(its, menuItem(it))
	}

	h := len(its) + 2
	if h < minMenuHeight {
		h = minMenuHeight
	}

	l := list.New(its, itemDelegate{}, defaultWidth, h)

	l.Title = menuTittle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle

	return menuModel{
		list: l,
	}
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (menuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true

			return m, tea.Quit
		case "enter":
			if it, ok := m.list.SelectedItem().(menuItem); ok {
				return m, func() tea.Msg { return menuChoiceMsg(it) }
			}

			return m, nil
		}
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m menuModel) View() string {
	menuHelp := helpStyle.Copy().PaddingTop(0)
	content := m.list.View() + "\n" +
		menuHelp.Render("↑/↓ to move | enter to select | esc/q to quit")

	return menuContainer.Render(content)
}
