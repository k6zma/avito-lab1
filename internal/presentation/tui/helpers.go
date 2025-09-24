package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
)

func studentsToTable(list []dtos.StudentListItemDTO) table.Model {
	cols := []table.Column{
		{Title: "#", Width: 3},
		{Title: "ID", Width: 36},
		{Title: "Name", Width: 14},
		{Title: "Surname", Width: 16},
		{Title: "Age", Width: 4},
		{Title: "Grades", Width: 40},
	}

	var rows []table.Row

	for i, s := range list {
		grades := ""

		if len(s.Grades) > 0 {
			ss := make([]string, len(s.Grades))

			for i, g := range s.Grades {
				ss[i] = strconv.Itoa(g)
			}

			grades = strings.Join(ss, ",")
		}

		rows = append(rows, table.Row{
			strconv.Itoa(i + 1),
			s.ID,
			s.Name,
			s.Surname,
			strconv.Itoa(s.Age),
			grades,
		})
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	return styleTable(t)
}

func studentLines(s dtos.DefaultStudentResponseDTO) []string {
	lines := []string{
		fmt.Sprintf("ID: %s", s.ID),
		fmt.Sprintf("Name: %s", s.Name),
		fmt.Sprintf("Surname: %s", s.Surname),
		fmt.Sprintf("Age: %d", s.Age),
	}

	if len(s.Grades) > 0 {
		var ss []string

		for _, g := range s.Grades {
			ss = append(ss, strconv.Itoa(g))
		}

		lines = append(lines, "Grades: "+strings.Join(ss, ", "))
	}

	if s.AvgGrade != nil {
		lines = append(lines, fmt.Sprintf("AVG: %.2f", *s.AvgGrade))
	}

	return lines
}

func renderStatus(msg string) string {
	l := strings.ToLower(msg)
	if strings.Contains(l, "error") || strings.Contains(l, "failed") ||
		strings.Contains(l, "invalid") {
		return errorStyle.Render(msg)
	}

	return helpStyle.Render(msg)
}
