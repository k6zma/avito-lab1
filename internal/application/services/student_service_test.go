package services_test

import (
	"fmt"
	"testing"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/services"
	infrarepo "github.com/k6zma/avito-lab1/internal/infrastructure/repositories"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	serviceTestPrefix = "StudentService"
)

type registerCase struct {
	name    string
	payload dtos.StudentCreateDTO
	ok      bool
}

func TestStudentService_Register_And_GetByID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Register_And_GetByID] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf(
			"[%s][Register_And_GetByID] error while creating repository: %v",
			serviceTestPrefix,
			err,
		)
	}

	svc := services.NewStudentService(repo)

	tests := []registerCase{
		{
			name: "ok register with grades and fetch",
			payload: dtos.StudentCreateDTO{
				Name:    "Mikhail",
				Surname: "Gunin",
				Age:     19,
				Grades:  []int{90, 60},
			},
			ok: true,
		},
		{
			name: "invalid register because name and surname not capitalized",
			payload: dtos.StudentCreateDTO{
				Name:    "mikhail",
				Surname: "gunin",
				Age:     19,
			},
			ok: false,
		},
	}

	for i, tc := range tests {
		t.Run(
			fmt.Sprintf("[%s]-register-%s-№%d", serviceTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				resp, err := svc.Register(tc.payload)
				gotOK := err == nil

				if gotOK != tc.ok {
					t.Fatalf(
						"[%s][Register] got ok=%v, want ok=%v (err=%v)",
						serviceTestPrefix,
						gotOK,
						tc.ok,
						err,
					)
				}

				if !tc.ok {
					return
				}

				if resp.Name != tc.payload.Name || resp.Surname != tc.payload.Surname ||
					resp.Age != tc.payload.Age {
					t.Fatalf(
						"[%s][Register] mismatch response fields: got{Name:%q,Surname:%q,Age:%d} want{Name:%q,Surname:%q,Age:%d}",
						serviceTestPrefix,
						resp.Name,
						resp.Surname,
						resp.Age,
						tc.payload.Name,
						tc.payload.Surname,
						tc.payload.Age,
					)
				}

				if len(tc.payload.Grades) > 0 {
					if resp.AvgGrade == nil {
						t.Fatalf(
							"[%s][Register] expected AvgGrade not nil for non-empty grades",
							serviceTestPrefix,
						)
					}

					want := float64(150) / 2.0

					if *resp.AvgGrade != want {
						t.Fatalf(
							"[%s][Register] avg mismatch: got=%v want=%v",
							serviceTestPrefix,
							*resp.AvgGrade,
							want,
						)
					}
				}

				getResp, err := svc.GetByID(dtos.GetByIDDTO{ID: resp.ID})
				if err != nil {
					t.Fatalf("[%s][GetByID] unexpected error: %v", serviceTestPrefix, err)
				}

				if getResp.ID != resp.ID || getResp.Name != resp.Name ||
					getResp.Surname != resp.Surname {
					t.Fatalf(
						"[%s][GetByID] mismatch - got{ID:%q,Name:%q,Surname:%q} want{ID:%q,Name:%q,Surname:%q}",
						serviceTestPrefix,
						getResp.ID,
						getResp.Name,
						getResp.Surname,
						resp.ID,
						resp.Name,
						resp.Surname,
					)
				}
			},
		)
	}
}

func TestStudentService_Update_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Update] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][Update] error while creating repository: %v", serviceTestPrefix, err)
	}

	svc := services.NewStudentService(repo)

	created, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     19,
		Grades:  []int{70},
	})
	if err != nil {
		t.Fatalf(
			"[%s][Register] unexpected error while seeding student: %v",
			serviceTestPrefix,
			err,
		)
	}

	upd, err := svc.Update(dtos.StudentUpdateDTO{
		ID:      created.ID,
		Name:    "Alexander",
		Surname: "Gunin",
		Age:     20,
		Grades:  []int{70, 85},
	})
	if err != nil {
		t.Fatalf(
			"[%s][Update(valid)] unexpected error while updating student: %v",
			serviceTestPrefix,
			err,
		)
	}

	if upd.Age != 20 || len(upd.Grades) != 2 {
		t.Fatalf(
			"[%s][Update(valid)] payload mismatch: got{Age:%d,Grades:%v} want{Age:%d,Grades:%v}",
			serviceTestPrefix, upd.Age, upd.Grades, 20, []int{70, 85},
		)
	}

	if _, err := svc.Update(dtos.StudentUpdateDTO{
		ID:      created.ID,
		Name:    "alexander",
		Surname: "gunin",
		Age:     20,
	}); err == nil {
		t.Fatalf(
			"[%s][Update(invalid)] expected validation error for non-capitalized name and surname, got nil",
			serviceTestPrefix,
		)
	}
}

func TestStudentService_DeleteByID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][DeleteByID] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][DeleteByID] error while creating repository: %v", serviceTestPrefix, err)
	}

	svc := services.NewStudentService(repo)

	created, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     19,
	})
	if err != nil {
		t.Fatalf(
			"[%s][Register] unexpected error while seeding student: %v",
			serviceTestPrefix,
			err,
		)
	}

	if err := svc.DeleteByID(dtos.GetByIDDTO{ID: created.ID}); err != nil {
		t.Fatalf(
			"[%s][DeleteByID] unexpected error while deleting student: %v",
			serviceTestPrefix,
			err,
		)
	}

	if _, err := svc.GetByID(dtos.GetByIDDTO{ID: created.ID}); err == nil {
		t.Fatalf(
			"[%s][GetByID(after delete)] expected error for deleted ID=%s, got nil",
			serviceTestPrefix,
			created.ID,
		)
	}
}

func TestStudentService_GetByFullName(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][GetByFullName] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][GetByFullName] error while creating repository: %v", serviceTestPrefix, err)
	}
	svc := services.NewStudentService(repo)

	created, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     19,
	})
	if err != nil {
		t.Fatalf(
			"[%s][Register] unexpected error while seeding student: %v",
			serviceTestPrefix,
			err,
		)
	}

	got, err := svc.GetByFullName(dtos.GetByFullNameDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
	})
	if err != nil {
		t.Fatalf(
			"[%s][GetByFullName] unexpected error while getting by full name: %v",
			serviceTestPrefix,
			err,
		)
	}

	if got.ID != created.ID {
		t.Fatalf(
			"[%s][GetByFullName] ID mismatch: got=%s want=%s",
			serviceTestPrefix,
			got.ID,
			created.ID,
		)
	}
}

func TestStudentService_List_WithAndWithoutGrades(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][List] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][List] error while creating repository: %v", serviceTestPrefix, err)
	}

	svc := services.NewStudentService(repo)

	_, err = svc.Register(dtos.StudentCreateDTO{
		Name:    "Eleven",
		Surname: "Doctor",
		Age:     100,
		Grades:  []int{100},
	})
	if err != nil {
		t.Fatalf("[%s][Register(a)] unexpected error: %v", serviceTestPrefix, err)
	}

	_, err = svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     19,
		Grades:  []int{80},
	})
	if err != nil {
		t.Fatalf("[%s][Register(b)] unexpected error: %v", serviceTestPrefix, err)
	}

	listNo, err := svc.List(false)
	if err != nil {
		t.Fatalf("[%s][List(false)] unexpected error: %v", serviceTestPrefix, err)
	}

	if len(listNo) != 2 {
		t.Fatalf(
			"[%s][List(false)] length mismatch: got=%d want=%d",
			serviceTestPrefix,
			len(listNo),
			2,
		)
	}

	for i, it := range listNo {
		if len(it.Grades) != 0 {
			t.Fatalf(
				"[%s][List(false)] grades should be omitted at idx=%d, got=%v",
				serviceTestPrefix,
				i,
				it.Grades,
			)
		}
	}

	listYes, err := svc.List(true)
	if err != nil {
		t.Fatalf("[%s][List(true)] unexpected error: %v", serviceTestPrefix, err)
	}

	if len(listYes) != 2 {
		t.Fatalf(
			"[%s][List(true)] length mismatch: got=%d want=%d",
			serviceTestPrefix,
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
			"[%s][List(true)] expected grades for all items, count with grades=%d",
			serviceTestPrefix,
			hasGrades,
		)
	}
}

func TestStudentService_AddGrades_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][AddGrades] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][AddGrades] error while creating repository: %v", serviceTestPrefix, err)
	}

	svc := services.NewStudentService(repo)

	created, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     19,
		Grades:  []int{60},
	})
	if err != nil {
		t.Fatalf(
			"[%s][Register] unexpected error while seeding student: %v",
			serviceTestPrefix,
			err,
		)
	}

	back, err := svc.AddGrades(dtos.AddGradesDTO{
		ID:     created.ID,
		Grades: []int{80, 90},
	})
	if err != nil {
		t.Fatalf(
			"[%s][AddGrades(valid)] unexpected error while adding grades: %v",
			serviceTestPrefix,
			err,
		)
	}

	if len(back.Grades) != 3 {
		t.Fatalf("[%s][AddGrades(valid)] grades length mismatch, got=%d want=%d (grades=%v)",
			serviceTestPrefix, len(back.Grades), 3, back.Grades)
	}

	if back.AvgGrade == nil {
		t.Fatalf("[%s][AddGrades(valid)] expected AvgGrade not nil", serviceTestPrefix)
	}

	if _, err := svc.AddGrades(dtos.AddGradesDTO{
		ID:     created.ID,
		Grades: []int{150},
	}); err == nil {
		t.Fatalf(
			"[%s][AddGrades(invalid)] expected validation error for grade=150, got nil",
			serviceTestPrefix,
		)
	}
}

func TestStudentService_AVGByID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][AVGByID] failed to init validators: %v", serviceTestPrefix, err)
	}

	repo, err := infrarepo.NewStudentStorageWithPersister(nil)
	if err != nil {
		t.Fatalf("[%s][AVGByID] error while creating repository: %v", serviceTestPrefix, err)
	}

	svc := services.NewStudentService(repo)

	a, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "Mikhail",
		Surname: "Gunin",
		Age:     20,
	})
	if err != nil {
		t.Fatalf("[%s][Register(a)] unexpected error: %v", serviceTestPrefix, err)
	}

	avgA, err := svc.AVGByID(dtos.GetByIDDTO{ID: a.ID})
	if err != nil {
		t.Fatalf("[%s][AVGByID(no grades)] unexpected error: %v", serviceTestPrefix, err)
	}

	if avgA.AVG != 0.0 {
		t.Fatalf(
			"[%s][AVGByID(no grades)] avg mismatch: got=%v want=%v",
			serviceTestPrefix,
			avgA.AVG,
			0.0,
		)
	}

	b, err := svc.Register(dtos.StudentCreateDTO{
		Name:    "With",
		Surname: "Grades",
		Age:     21,
		Grades:  []int{50, 75, 100},
	})
	if err != nil {
		t.Fatalf("[%s][Register(b)] unexpected error: %v", serviceTestPrefix, err)
	}

	avgB, err := svc.AVGByID(dtos.GetByIDDTO{ID: b.ID})
	if err != nil {
		t.Fatalf("[%s][AVGByID(with grades)] unexpected error: %v", serviceTestPrefix, err)
	}

	want := (50.0 + 75.0 + 100.0) / 3.0
	if avgB.AVG != want {
		t.Fatalf(
			"[%s][AVGByID(with grades)] avg mismatch: got=%v want=%v",
			serviceTestPrefix,
			avgB.AVG,
			want,
		)
	}
}
