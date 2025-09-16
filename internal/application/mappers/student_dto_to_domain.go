package mappers

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

func MapStudentCreateDTOToDomain(d dtos.StudentCreateDTO) (*models.Student, error) {
	if err := validators.Validate.Struct(d); err != nil {
		return nil, fmt.Errorf("failed to validate student create dto: %w", err)
	}

	student, err := models.NewStudentBuilder().
		SetName(d.Name).
		SetSurname(d.Surname).
		SetAge(d.Age).
		SetGrades(d.Grades).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build domain student from create dto: %w", err)
	}

	return student, nil
}

func MapStudentUpdateDTOToDomain(d dtos.StudentUpdateDTO) (*models.Student, error) {
	if err := validators.Validate.Struct(d); err != nil {
		return nil, fmt.Errorf("failed to validate student update dto: %w", err)
	}

	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse id from string to uuid: %w", err)
	}

	student, err := models.NewStudentBuilder().
		SetName(d.Name).
		SetSurname(d.Surname).
		SetAge(d.Age).
		SetGrades(d.Grades).
		Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build domain student from update dto: %w", err)
	}

	student.ID = id

	return student, nil
}

func MapAddGradesDTOToArgs(d dtos.AddGradesDTO) (uuid.UUID, []int, error) {
	if err := validators.Validate.Struct(d); err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to validate add-grades dto: %w", err)
	}

	id, err := uuid.Parse(d.ID)
	if err != nil {
		return uuid.Nil, nil, fmt.Errorf("failed to parse id from string to uuid: %w", err)
	}

	return id, append([]int(nil), d.Grades...), nil
}

func MapGetByFullNameDTOToArgs(d dtos.GetByFullNameDTO) (string, string, error) {
	if err := validators.Validate.Struct(d); err != nil {
		return "", "", fmt.Errorf("failed to validate get-by-fullname dto: %w", err)
	}

	return d.Name, d.Surname, nil
}

func MapGetByIDDTOToUUID(d dtos.GetByIDDTO) (uuid.UUID, error) {
	if err := validators.Validate.Struct(d); err != nil {
		return uuid.Nil, fmt.Errorf("failed to validate get-by-id dto: %w", err)
	}

	id, err := uuid.Parse(d.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse id from string to uuid: %w", err)
	}

	return id, nil
}
