package mappers

import (
	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/domain/models"
)

func MapStudentDomainToDefaultResponseDTO(
	student *models.Student,
	withAVG bool,
) dtos.DefaultStudentResponseDTO {
	if student == nil {
		return dtos.DefaultStudentResponseDTO{}
	}

	res := dtos.DefaultStudentResponseDTO{
		ID:      student.ID.String(),
		Name:    student.Name,
		Surname: student.Surname,
		Age:     student.Age,
		Grades:  append([]int(nil), student.Grades...),
	}

	if withAVG && len(student.Grades) > 0 {
		sum := 0

		for _, g := range student.Grades {
			sum += g
		}

		avg := float64(sum) / float64(len(student.Grades))
		res.AvgGrade = &avg
	}

	return res
}

func MapStudentsDomainToListDTO(
	list []*models.Student,
	includeGrades bool,
) []dtos.StudentListItemDTO {
	out := make([]dtos.StudentListItemDTO, 0, len(list))
	for _, student := range list {
		if student == nil {
			continue
		}

		item := dtos.StudentListItemDTO{
			ID:      student.ID.String(),
			Name:    student.Name,
			Surname: student.Surname,
			Age:     student.Age,
		}

		if includeGrades && len(student.Grades) > 0 {
			item.Grades = append([]int(nil), student.Grades...)
		}

		out = append(out, item)
	}

	return out
}
