package mappers_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/mappers"
	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	mapperTestPrefix = "StudentMappers"
)

type createCase struct {
	name string
	in   dtos.StudentCreateDTO
	ok   bool
}

type updateCase struct {
	name string
	in   dtos.StudentUpdateDTO
	ok   bool
}

func TestMapStudentCreateDTOToDomain_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][CreateDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	tests := []createCase{
		{
			name: "ok with grades",
			in: dtos.StudentCreateDTO{
				Name:    "Mikhail",
				Surname: "Gunin",
				Age:     19,
				Grades:  []int{90, 60},
			},
			ok: true,
		},
		{
			name: "ok without grades",
			in: dtos.StudentCreateDTO{
				Name:    "Alexander",
				Surname: "Gunin",
				Age:     20,
			},
			ok: true,
		},
		{
			name: "invalid name (not capitalized)",
			in: dtos.StudentCreateDTO{
				Name:    "mikhail",
				Surname: "Gunin",
				Age:     19,
			},
			ok: false,
		},
		{
			name: "invalid grades (out of range)",
			in: dtos.StudentCreateDTO{
				Name:    "Mikhail",
				Surname: "Gunin",
				Age:     19,
				Grades:  []int{120},
			},
			ok: false,
		},
	}

	for i, tc := range tests {
		t.Run(
			fmt.Sprintf("[%s]-create-%s-№%d", mapperTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				got, err := mappers.MapStudentCreateDTOToDomain(tc.in)
				gotOK := err == nil

				if gotOK != tc.ok {
					t.Fatalf("[%s][MapStudentCreateDTOToDomain] got ok=%v, want ok=%v (err=%v)",
						mapperTestPrefix, gotOK, tc.ok, err)
				}

				if !tc.ok {
					return
				}

				if got.Name != tc.in.Name || got.Surname != tc.in.Surname || got.Age != tc.in.Age {
					t.Fatalf(
						"[%s][MapStudentCreateDTOToDomain] mismatch fields: got{Name:%q,Surname:%q,Age:%d} want{Name:%q,Surname:%q,Age:%d}",
						mapperTestPrefix,
						got.Name,
						got.Surname,
						got.Age,
						tc.in.Name,
						tc.in.Surname,
						tc.in.Age,
					)
				}

				if fmt.Sprint(got.Grades) != fmt.Sprint(tc.in.Grades) {
					t.Fatalf("[%s][MapStudentCreateDTOToDomain] grades mismatch: got=%v want=%v",
						mapperTestPrefix, got.Grades, tc.in.Grades)
				}
			},
		)
	}
}

func TestMapStudentUpdateDTOToDomain_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][UpdateDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	validID := uuid.NewString()

	tests := []updateCase{
		{
			name: "ok with grades",
			in: dtos.StudentUpdateDTO{
				ID:      validID,
				Name:    "Alexander",
				Surname: "Gunin",
				Age:     21,
				Grades:  []int{70, 85},
			},
			ok: true,
		},
		{
			name: "invalid id (not uuid4)",
			in: dtos.StudentUpdateDTO{
				ID:      "bad-uuid",
				Name:    "Alexander",
				Surname: "Gunin",
				Age:     21,
			},
			ok: false,
		},
		{
			name: "invalid name (not capitalized)",
			in: dtos.StudentUpdateDTO{
				ID:      validID,
				Name:    "alexander",
				Surname: "Gunin",
				Age:     21,
			},
			ok: false,
		},
	}

	for i, tc := range tests {
		t.Run(
			fmt.Sprintf("[%s]-update-%s-№%d", mapperTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				got, err := mappers.MapStudentUpdateDTOToDomain(tc.in)
				gotOK := err == nil

				if gotOK != tc.ok {
					t.Fatalf("[%s][MapStudentUpdateDTOToDomain] got ok=%v, want ok=%v (err=%v)",
						mapperTestPrefix, gotOK, tc.ok, err)
				}

				if !tc.ok {
					return
				}

				if got.Name != tc.in.Name || got.Surname != tc.in.Surname || got.Age != tc.in.Age {
					t.Fatalf(
						"[%s][MapStudentUpdateDTOToDomain] mismatch fields: got{Name:%q,Surname:%q,Age:%d} want{Name:%q,Surname:%q,Age:%d}",
						mapperTestPrefix,
						got.Name,
						got.Surname,
						got.Age,
						tc.in.Name,
						tc.in.Surname,
						tc.in.Age,
					)
				}

				if fmt.Sprint(got.Grades) != fmt.Sprint(tc.in.Grades) {
					t.Fatalf("[%s][MapStudentUpdateDTOToDomain] grades mismatch: got=%v want=%v",
						mapperTestPrefix, got.Grades, tc.in.Grades)
				}

				if got.ID.String() != tc.in.ID {
					t.Fatalf("[%s][MapStudentUpdateDTOToDomain] id mismatch: got=%s want=%s",
						mapperTestPrefix, got.ID, tc.in.ID)
				}
			},
		)
	}
}

func TestMapAddGradesDTOToArgs_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][AddGradesDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	id := uuid.New()

	gotID, gotGrades, err := mappers.MapAddGradesDTOToArgs(dtos.AddGradesDTO{
		ID:     id.String(),
		Grades: []int{10, 20},
	})
	if err != nil {
		t.Fatalf("[%s][MapAddGradesDTOToArgs(ok)] unexpected error: %v", mapperTestPrefix, err)
	}

	if gotID != id {
		t.Fatalf(
			"[%s][MapAddGradesDTOToArgs(ok)] id mismatch: got=%s want=%s",
			mapperTestPrefix,
			gotID,
			id,
		)
	}
	if fmt.Sprint(gotGrades) != fmt.Sprint([]int{10, 20}) {
		t.Fatalf("[%s][MapAddGradesDTOToArgs(ok)] grades mismatch: got=%v want=%v",
			mapperTestPrefix, gotGrades, []int{10, 20})
	}

	if _, _, err := mappers.MapAddGradesDTOToArgs(dtos.AddGradesDTO{
		ID:     id.String(),
		Grades: []int{150},
	}); err == nil {
		t.Fatalf(
			"[%s][MapAddGradesDTOToArgs(invalid grade)] expected validation error for 150, got nil",
			mapperTestPrefix,
		)
	}

	if _, _, err := mappers.MapAddGradesDTOToArgs(dtos.AddGradesDTO{
		ID:     "not-uuid",
		Grades: []int{10},
	}); err == nil {
		t.Fatalf("[%s][MapAddGradesDTOToArgs(invalid uuid)] expected error for bad uuid, got nil",
			mapperTestPrefix)
	}
}

func TestMapGetByFullNameDTOToArgs_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][GetByFullNameDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	name, surname, err := mappers.MapGetByFullNameDTOToArgs(dtos.GetByFullNameDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
	})
	if err != nil {
		t.Fatalf("[%s][MapGetByFullNameDTOToArgs(ok)] unexpected error: %v", mapperTestPrefix, err)
	}

	if name != "Mikhail" || surname != "Gunin" {
		t.Fatalf("[%s][MapGetByFullNameDTOToArgs(ok)] mismatch: got={%q,%q} want={%q,%q}",
			mapperTestPrefix, name, surname, "Mikhail", "Gunin")
	}

	if _, _, err := mappers.MapGetByFullNameDTOToArgs(dtos.GetByFullNameDTO{
		Name:    "mikhail",
		Surname: "gunin",
	}); err == nil {
		t.Fatalf(
			"[%s][MapGetByFullNameDTOToArgs(invalid)] expected validation error for non-capitalized, got nil",
			mapperTestPrefix,
		)
	}
}

func TestMapGetByIDDTOToUUID_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][GetByIDDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	id := uuid.New()

	got, err := mappers.MapGetByIDDTOToUUID(dtos.GetByIDDTO{ID: id.String()})
	if err != nil {
		t.Fatalf("[%s][MapGetByIDDTOToUUID(ok)] unexpected error: %v", mapperTestPrefix, err)
	}

	if got != id {
		t.Fatalf(
			"[%s][MapGetByIDDTOToUUID(ok)] id mismatch: got=%s want=%s",
			mapperTestPrefix,
			got,
			id,
		)
	}

	if _, err := mappers.MapGetByIDDTOToUUID(dtos.GetByIDDTO{ID: "bad-uuid"}); err == nil {
		t.Fatalf(
			"[%s][MapGetByIDDTOToUUID(invalid)] expected error for bad uuid, got nil",
			mapperTestPrefix,
		)
	}
}

func TestMapStudentDomainToDefaultResponseDTO_WithAndWithoutAVG_And_Nil(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][DomainToDTO] failed to init validators: %v", mapperTestPrefix, err)
	}

	st, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{90, 60}).
		Build()
	if err != nil {
		t.Fatalf("[%s][DomainToDTO] failed to build domain student: %v", mapperTestPrefix, err)
	}

	with := mappers.MapStudentDomainToDefaultResponseDTO(st, true)
	if with.ID != st.ID.String() || with.Name != st.Name || with.Surname != st.Surname ||
		with.Age != st.Age {
		t.Fatalf("[%s][DomainToDTO(with avg)] fields mismatch: got{%q,%q,%q,%d} want{%q,%q,%q,%d}",
			mapperTestPrefix, with.ID, with.Name, with.Surname, with.Age,
			st.ID.String(), st.Name, st.Surname, st.Age)
	}

	if with.AvgGrade == nil {
		t.Fatalf("[%s][DomainToDTO(with avg)] expected AvgGrade not nil", mapperTestPrefix)
	}

	want := (90.0 + 60.0) / 2.0
	if *with.AvgGrade != want {
		t.Fatalf(
			"[%s][DomainToDTO(with avg)] avg mismatch: got=%v want=%v",
			mapperTestPrefix,
			*with.AvgGrade,
			want,
		)
	}

	if fmt.Sprint(with.Grades) != fmt.Sprint(st.Grades) {
		t.Fatalf(
			"[%s][DomainToDTO(with avg)] grades mismatch: got=%v want=%v",
			mapperTestPrefix,
			with.Grades,
			st.Grades,
		)
	}

	without := mappers.MapStudentDomainToDefaultResponseDTO(st, false)
	if without.AvgGrade != nil {
		t.Fatalf(
			"[%s][DomainToDTO(without avg)] expected AvgGrade nil, got=%v",
			mapperTestPrefix,
			*without.AvgGrade,
		)
	}

	stNo, err := models.NewStudentBuilder().
		SetName("No").
		SetSurname("Grades").
		SetAge(20).
		Build()
	if err != nil {
		t.Fatalf("[%s][DomainToDTO(empty)] build: %v", mapperTestPrefix, err)
	}

	no := mappers.MapStudentDomainToDefaultResponseDTO(stNo, true)
	if no.AvgGrade != nil {
		t.Fatalf(
			"[%s][DomainToDTO(empty)] expected AvgGrade nil for empty grades",
			mapperTestPrefix,
		)
	}

	zero := mappers.MapStudentDomainToDefaultResponseDTO(nil, true)
	if zero.ID != "" || zero.Name != "" || zero.Surname != "" || zero.Age != 0 ||
		len(zero.Grades) != 0 ||
		zero.AvgGrade != nil {
		t.Fatalf("[%s][DomainToDTO(nil)] expected zero value dto, got=%+v", mapperTestPrefix, zero)
	}
}

func TestMapStudentsDomainToListDTO_IncludeGrades_And_NilItems(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][ListMap] failed to init validators: %v", mapperTestPrefix, err)
	}

	a, err := models.NewStudentBuilder().
		SetName("Eleven").
		SetSurname("Doctor").
		SetAge(100).
		SetGrades([]int{100}).
		Build()
	if err != nil {
		t.Fatalf("[%s][ListMap] failed to build first student: %v", mapperTestPrefix, err)
	}

	b, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{80}).
		Build()
	if err != nil {
		t.Fatalf("[%s][ListMap] failed to build second student: %v", mapperTestPrefix, err)
	}

	listNo := mappers.MapStudentsDomainToListDTO([]*models.Student{a, nil, b}, false)
	if len(listNo) != 2 {
		t.Fatalf(
			"[%s][ListMap(false)] length mismatch: got=%d want=%d",
			mapperTestPrefix,
			len(listNo),
			2,
		)
	}

	for i, it := range listNo {
		if len(it.Grades) != 0 {
			t.Fatalf(
				"[%s][ListMap(false)] grades should be omitted at idx=%d, got=%v",
				mapperTestPrefix,
				i,
				it.Grades,
			)
		}
	}

	listYes := mappers.MapStudentsDomainToListDTO([]*models.Student{a, nil, b}, true)
	if len(listYes) != 2 {
		t.Fatalf(
			"[%s][ListMap(true)] length mismatch: got=%d want=%d",
			mapperTestPrefix,
			len(listYes),
			2,
		)
	}

	hasGrades := 0
	for _, it := range listYes {
		if len(it.Grades) > 0 {
			hasGrades++
		}
	}

	if hasGrades != 2 {
		t.Fatalf(
			"[%s][ListMap(true)] expected grades for all items, count with grades=%d",
			mapperTestPrefix,
			hasGrades,
		)
	}
}
