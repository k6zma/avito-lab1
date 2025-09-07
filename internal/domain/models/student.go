package models

import (
	"fmt"

	"github.com/k6zma/avito-lab1/pkg/validators"
)

type Student struct {
	Name    string `validate:"required,capitalized"`
	Surname string `validate:"required,capitalized"`
	Age     int    `validate:"gte=0,lte=150"`
	Grades  []int  `validate:"omitempty,dive,gte=0,lte=100"`
}

// =============================================
// Realization of custom setters with validation
// =============================================

func (s *Student) SetName(name string) error {
	if err := validators.Validate.Var(name, "required,capitalized"); err != nil {
		return fmt.Errorf("error while validating student name in student name setter: %w", err)
	}

	s.Name = name

	return nil
}

func (s *Student) SetSurname(surname string) error {
	if err := validators.Validate.Var(surname, "required,capitalized"); err != nil {
		return fmt.Errorf(
			"error while validating student surname in student surname setter: %w",
			err,
		)
	}

	s.Surname = surname

	return nil
}

func (s *Student) SetAge(age int) error {
	if err := validators.Validate.Var(age, "gte=0,lte=150"); err != nil {
		return fmt.Errorf("error while validating student age in student age setter: %w", err)
	}

	s.Age = age

	return nil
}

func (s *Student) SetGrades(grades []int) error {
	if err := validators.Validate.Var(grades, "dive,gte=0,lte=100"); err != nil {
		return fmt.Errorf("error while validating grades in student grades setter: %w", err)
	}

	s.Grades = grades

	return nil
}

func (s *Student) AppendGrade(grade int) error {
	if err := validators.Validate.Var(grade, "gte=0,lte=100"); err != nil {
		return fmt.Errorf(
			"error while validating appended grade in student append grade method: %w",
			err,
		)
	}

	s.Grades = append(s.Grades, grade)

	return nil
}

// ================================================
// Realization of Builder pattern for Student model
// ================================================

type StudentBuilder interface {
	SetName(name string) StudentBuilder
	SetSurname(surname string) StudentBuilder
	SetAge(age int) StudentBuilder
	SetGrades(grades []int) StudentBuilder
	Build() (*Student, error)
}

type studentBuilder struct {
	name    string
	surname string
	age     int
	grades  []int
}

func NewStudentBuilder() StudentBuilder {
	return &studentBuilder{}
}

func (s *studentBuilder) SetName(name string) StudentBuilder {
	s.name = name

	return s
}

func (s *studentBuilder) SetSurname(surname string) StudentBuilder {
	s.surname = surname

	return s
}

func (s *studentBuilder) SetAge(age int) StudentBuilder {
	s.age = age

	return s
}

func (s *studentBuilder) SetGrades(grades []int) StudentBuilder {
	s.grades = grades

	return s
}

func (s *studentBuilder) Build() (*Student, error) {
	student := &Student{
		Name:    s.name,
		Surname: s.surname,
		Age:     s.age,
		Grades:  s.grades,
	}

	if err := validators.Validate.Struct(student); err != nil {
		return nil, fmt.Errorf(
			"error while validating student domain model in build method: %w",
			err,
		)
	}

	return student, nil
}
