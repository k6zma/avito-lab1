package tui

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/services"
)

type rootModel struct {
	svc        services.StudentServiceContract
	mode       mode
	prevMode   mode
	currentAct string
	menu       menuModel
	tbl        tableModel
	form       createModel
	grades     addGradesModel
	idInput    idInputModel
	detail     detailModel
	status     string
}

func Run(svc services.StudentServiceContract) error {
	m := newRootModel(svc)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI app: %w", err)
	}

	return nil
}

func newRootModel(svc services.StudentServiceContract) rootModel {
	return rootModel{
		svc:      svc,
		mode:     modeMenu,
		prevMode: modeMenu,
		menu: newMenuModel([]string{
			"Add student",
			"List students",
			"Show student (by ID)",
			"Average by ID",
			"Add grades",
			"Delete student",
			"Quit",
		}),
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case menuChoiceMsg:
		switch string(msg) {
		case "Add student":
			m.mode = modeCreate
			m.form = newCreateModel()

			return m, m.form.Init()

		case "List students":
			list, err := m.svc.List(true)
			if err != nil {
				m.status = fmt.Sprintf("list error: %v", err)

				return m, nil
			}

			m.tbl = newTableModel(studentsToTable(list))
			m.mode = modeTable

			return m, nil

		case "Show student (by ID)":
			m.mode = modeIDInput
			m.currentAct = actionShow
			m.idInput = newIDInputModel("Student ID in UUID)")

			return m, m.idInput.Init()

		case "Average by ID":
			m.mode = modeIDInput
			m.currentAct = actionAVG
			m.idInput = newIDInputModel("Student ID in UUID")

			return m, m.idInput.Init()

		case "Add grades":
			m.mode = modeAddGrades
			m.grades = newAddGradesModel()

			return m, m.grades.Init()

		case "Delete student":
			m.mode = modeIDInput
			m.currentAct = actionDel
			m.idInput = newIDInputModel("Student ID in UUID")

			return m, m.idInput.Init()

		case "Quit":
			return m, tea.Quit
		}

	case tableBackMsg:
		if m.mode == modeDetail {
			m.mode = m.prevMode
		} else {
			m.mode = modeMenu
		}

		return m, nil

	case tableShowMsg:
		id := strings.TrimSpace(msg.ID)

		r, err := m.svc.GetByID(dtos.GetByIDDTO{ID: id})
		if err != nil {
			m.status = fmt.Sprintf("fetch failed: %v", err)
			m.mode = modeMenu

			return m, nil
		}

		m.detail = newDetailModel(studentLines(r))
		m.prevMode = modeTable
		m.mode = modeDetail

		return m, nil

	case createCancelMsg:
		m.mode = modeMenu

		return m, nil

	case createSubmittedMsg:
		name := strings.TrimSpace(msg.Name)
		surname := strings.TrimSpace(msg.Surname)
		ageStr := strings.TrimSpace(msg.Age)
		gradesStr := strings.TrimSpace(msg.GradesCSV)

		age := 0
		if ageStr != "" {
			v, err := strconv.Atoi(ageStr)
			if err != nil {
				m.status = fmt.Sprintf("invalid age: %v", err)

				return m, nil
			}

			age = v
		}

		var grades []int

		if gradesStr != "" {
			for _, p := range strings.Split(gradesStr, ",") {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}

				v, err := strconv.Atoi(p)
				if err != nil {
					m.status = fmt.Sprintf("invalid grade %q: %v", p, err)

					return m, nil
				}

				grades = append(grades, v)
			}
		}

		resp, err := m.svc.Register(dtos.StudentCreateDTO{
			Name: name, Surname: surname, Age: age, Grades: grades,
		})
		if err != nil {
			m.status = fmt.Sprintf("register failed: %v", err)
			m.mode = modeMenu

			return m, nil
		}

		m.detail = newDetailModel(studentLines(resp))
		m.prevMode = modeMenu
		m.mode = modeDetail

		return m, nil

	case addGradesCancelMsg:
		m.mode = modeMenu

		return m, nil

	case addGradesSubmittedMsg:
		id := strings.TrimSpace(msg.ID)

		var grades []int

		for _, p := range strings.Split(strings.TrimSpace(msg.Grades), ",") {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			v, err := strconv.Atoi(p)
			if err != nil {
				m.status = fmt.Sprintf("invalid grade %q: %v", p, err)
				m.mode = modeMenu

				return m, nil
			}

			grades = append(grades, v)
		}

		resp, err := m.svc.AddGrades(dtos.AddGradesDTO{ID: id, Grades: grades})
		if err != nil {
			m.status = fmt.Sprintf("add grades failed: %v", err)
			m.mode = modeMenu

			return m, nil
		}

		m.detail = newDetailModel(append([]string{"Grades added successfully"}, studentLines(resp)...))
		m.prevMode = modeMenu
		m.mode = modeDetail

		return m, nil

	case idCancelMsg:
		m.mode = modeMenu

		return m, nil

	case idSubmittedMsg:
		id := strings.TrimSpace(string(msg))

		switch m.currentAct {
		case actionAVG:
			r, err := m.svc.AVGByID(dtos.GetByIDDTO{ID: id})
			if err != nil {
				m.status = fmt.Sprintf("avg error: %v", err)
				m.mode = modeMenu
				return m, nil
			}

			m.detail = newDetailModel([]string{
				fmt.Sprintf("ID: %s", r.ID),
				fmt.Sprintf("AVG: %.2f", r.AVG),
			})

			m.prevMode = modeMenu
			m.mode = modeDetail

			return m, nil

		case actionDel:
			if err := m.svc.DeleteByID(dtos.GetByIDDTO{ID: id}); err != nil {
				m.status = fmt.Sprintf("delete failed: %v", err)
				m.mode = modeMenu

				return m, nil
			}

			m.status = ""
			m.mode = modeMenu

			return m, nil

		case actionShow:
			r, err := m.svc.GetByID(dtos.GetByIDDTO{ID: id})
			if err != nil {
				m.status = fmt.Sprintf("fetch failed: %v", err)
				m.mode = modeMenu
				return m, nil
			}

			m.detail = newDetailModel(studentLines(r))
			m.prevMode = modeMenu
			m.mode = modeDetail

			return m, nil
		}
	}

	switch m.mode {
	case modeMenu:
		var cmd tea.Cmd

		m.menu, cmd = m.menu.Update(msg)

		return m, cmd
	case modeTable:
		var cmd tea.Cmd

		m.tbl, cmd = m.tbl.Update(msg)

		return m, cmd
	case modeCreate:
		var cmd tea.Cmd

		m.form, cmd = m.form.Update(msg)

		return m, cmd
	case modeAddGrades:
		var cmd tea.Cmd

		m.grades, cmd = m.grades.Update(msg)

		return m, cmd
	case modeIDInput:
		var cmd tea.Cmd

		m.idInput, cmd = m.idInput.Update(msg)

		return m, cmd
	case modeDetail:
		var cmd tea.Cmd

		m.detail, cmd = m.detail.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m rootModel) View() string {
	switch m.mode {
	case modeMenu:
		out := "\n" + m.menu.View()
		if m.status != "" {
			out += "\n" + renderStatus(m.status)
		}

		return out
	case modeTable:
		return m.tbl.View()
	case modeCreate:
		return m.form.View()
	case modeAddGrades:
		return m.grades.View()
	case modeIDInput:
		return m.idInput.View()
	case modeDetail:
		return m.detail.View()
	default:
		return ""
	}
}
