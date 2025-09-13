package models_test

import (
	"fmt"
	"testing"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	studentTestPrefix = "StudentDomainModel"
)

type studentSetterTestCase struct {
	name      string
	setFunc   func(s *models.Student) error
	wantError bool
}

type studentBuilderTestCase struct {
	name      string
	nameVal   string
	surname   string
	age       int
	grades    []int
	wantError bool
}

func TestStudent_Setters(t *testing.T) {
	err := validators.InitValidators()
	if err != nil {
		t.Fatalf("error while initializing validators: %v", err)
	}

	s := &models.Student{}

	tests := []studentSetterTestCase{
		{
			name: "valid name",
			setFunc: func(s *models.Student) error {
				return s.SetName("K6zma")
			},
			wantError: false,
		},
		{
			name: "invalid name (empty)",
			setFunc: func(s *models.Student) error {
				return s.SetName("")
			},
			wantError: true,
		},
		{
			name: "invalid name (not capitalized)",
			setFunc: func(s *models.Student) error {
				return s.SetName("k6zma")
			},
			wantError: true,
		},
		{
			name: "valid surname",
			setFunc: func(s *models.Student) error {
				return s.SetSurname("Gunin")
			},
			wantError: false,
		},
		{
			name: "invalid surname (empty)",
			setFunc: func(s *models.Student) error {
				return s.SetSurname("")
			},
			wantError: true,
		},
		{
			name: "invalid surname (not capitalized)",
			setFunc: func(s *models.Student) error {
				return s.SetSurname("gunin")
			},
			wantError: true,
		},
		{
			name: "valid age",
			setFunc: func(s *models.Student) error {
				return s.SetAge(20)
			},
			wantError: false,
		},
		{
			name: "invalid age (negative)",
			setFunc: func(s *models.Student) error {
				return s.SetAge(-1)
			},
			wantError: true,
		},
		{
			name: "valid grades",
			setFunc: func(s *models.Student) error {
				return s.SetGrades([]int{80, 90, 100})
			},
			wantError: false,
		},
		{
			name: "invalid grades (out of range)",
			setFunc: func(s *models.Student) error {
				return s.SetGrades([]int{187, 104, -653})
			},
			wantError: true,
		},
		{
			name: "invalid grades (mix valid with out of range)",
			setFunc: func(s *models.Student) error {
				return s.SetGrades([]int{50, 104, 90, 32})
			},
			wantError: true,
		},
		{
			name: "append valid grade",
			setFunc: func(s *models.Student) error {
				return s.AddGrades(75)
			},
			wantError: false,
		},
		{
			name: "append invalid grade (too high)",
			setFunc: func(s *models.Student) error {
				return s.AddGrades(150)
			},
			wantError: true,
		},
		{
			name: "append invalid grade (too low)",
			setFunc: func(s *models.Student) error {
				return s.AddGrades(-14)
			},
			wantError: true,
		},
	}

	for i, tt := range tests {
		t.Run(
			fmt.Sprintf("[%s]-setter-%s-№%d", studentTestPrefix, tt.name, i+1),
			func(t *testing.T) {
				err := tt.setFunc(s)
				gotError := err != nil

				if gotError != tt.wantError {
					t.Errorf(
						"Test %d (%s): got error = %v, want error = %v (err: %v)",
						i+1, tt.name, gotError, tt.wantError, err,
					)
				}
			},
		)
	}
}

func TestStudentBuilder_Build(t *testing.T) {
	err := validators.InitValidators()
	if err != nil {
		t.Fatalf("error while initializing validators: %v", err)
	}

	testCases := []studentBuilderTestCase{
		{
			name:      "valid student",
			nameVal:   "K6zma",
			surname:   "Gunin",
			age:       20,
			grades:    []int{80, 90, 100},
			wantError: false,
		},
		{
			name:      "invalid name (empty)",
			nameVal:   "",
			surname:   "Gunin",
			age:       20,
			grades:    []int{80, 90, 100},
			wantError: true,
		},
		{
			name:      "invalid name (not capitalized)",
			nameVal:   "k6zma",
			surname:   "Gunin",
			age:       20,
			grades:    []int{80, 90, 100},
			wantError: true,
		},
		{
			name:      "invalid surname (empty)",
			nameVal:   "K6zma",
			surname:   "",
			age:       20,
			grades:    []int{80, 90, 100},
			wantError: true,
		},
		{
			name:      "invalid surname (not capitalized)",
			nameVal:   "K6zma",
			surname:   "gunin",
			age:       20,
			grades:    []int{80, 90, 100},
			wantError: true,
		},
		{
			name:      "invalid age (negative)",
			nameVal:   "K6zma",
			surname:   "Gunin",
			age:       -1,
			grades:    []int{80, 90, 100},
			wantError: true,
		},
		{
			name:      "invalid grades (out of range)",
			nameVal:   "K6zma",
			surname:   "Gunin",
			age:       20,
			grades:    []int{187, 104, -653},
			wantError: true,
		},
		{
			name:      "invalid grades (mix valid with out of range)",
			nameVal:   "K6zma",
			surname:   "Gunin",
			age:       20,
			grades:    []int{50, 104, 90, 32},
			wantError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(
			fmt.Sprintf("[%s]-builder-%s-№%d", studentTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				student := models.NewStudentBuilder().
					SetName(tc.nameVal).
					SetSurname(tc.surname).
					SetAge(tc.age).
					SetGrades(tc.grades)

				_, err := student.Build()
				gotError := err != nil

				if gotError != tc.wantError {
					t.Errorf(
						"Builder Test %d (%s): got error = %v, want error = %v (err: %v)",
						i+1, tc.name, gotError, tc.wantError, err,
					)
				}
			},
		)
	}
}
