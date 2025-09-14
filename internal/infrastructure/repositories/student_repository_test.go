package repositories

import (
	"context"
	"fmt"
	"testing"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	repoImplTestPrefix = "StudentRepositoryImpl"
)

type smokeCase struct {
	name      string
	build     func() *models.Student
	wantErr   bool
	wantFound bool
}

func TestRepository_Create_And_GetByID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Create_And_GetByID] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	st, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{90, 95}).
		Build()
	if err != nil {
		t.Fatalf("[%s][Create_And_GetByID] error while creating Student domain model: %v", repoImplTestPrefix, err)
	}

	id, err := repo.Create(ctx, st)
	if err != nil {
		t.Fatalf("[%s][Create] unexpected error while creating student from storage: %v", repoImplTestPrefix, err)
	}

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("[%s][GetByID-2] unexpected error while getting student from storage: %v", repoImplTestPrefix, err)
	}

	if got.ID != id || got.Name != st.Name || got.Surname != st.Surname || got.Age != st.Age {
		t.Fatalf("[%s][GetByID] mismatch - got: ID=%s Name=%q Surname=%q Age=%d\n  want: ID=%s Name=%q Surname=%q Age=%d",
			repoImplTestPrefix, got.ID, got.Name, got.Surname, got.Age, id, st.Name, st.Surname, st.Age)
	}

	got.Name = "Change name to check copy"

	back, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("[%s][GetByID-2] unexpected error while getting student by id: %v", repoImplTestPrefix, err)
	}

	if back.Name == "Change name to check copy" {
		t.Fatalf("[%s][GetByID] shallow copy leak - returned value changed repository state", repoImplTestPrefix)
	}
}

func TestRepository_Create_DuplicateID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Create_DuplicateID] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	st1, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf("[%s][Create_DuplicateID] error while building first Student domain model: %v", repoImplTestPrefix, err)
	}

	st2, err := models.NewStudentBuilder().
		SetName("Alexander").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{100}).
		Build()
	if err != nil {
		t.Fatalf("[%s][Create_DuplicateID] error while building second Student domain model: %v", repoImplTestPrefix, err)
	}

	st2.ID = st1.ID

	if _, err := repo.Create(ctx, st1); err != nil {
		t.Fatalf("[%s][Create(first)] unexpected error while creating frist student in storage: %v", repoImplTestPrefix, err)
	}

	if _, err := repo.Create(ctx, st2); err == nil {
		t.Fatalf("[%s][Create(duplicate)] expected error while creating second student on duplicate ID=%s, got nil", repoImplTestPrefix, st1.ID)
	}
}

func TestRepository_Create_InvalidStudent(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Create_Invalid] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	s, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf("[%s][Create_Invalid] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	s.Name = "mikhail"

	if _, err := repo.Create(ctx, s); err == nil {
		t.Fatalf("[%s][Create_Invalid] expected validation error for Name=%q, got nil", repoImplTestPrefix, s.Name)
	}
}

func TestRepository_GetByFullName(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][GetByFullName] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	st, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{100}).
		Build()
	if err != nil {
		t.Fatalf("[%s][GetByFullName] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	id, err := repo.Create(ctx, st)
	if err != nil {
		t.Fatalf("[%s][Create] unexpected error while creating Student in storage: %v", repoImplTestPrefix, err)
	}

	got, err := repo.GetByFullName(ctx, "Mikhail", "Gunin")
	if err != nil {
		t.Fatalf("[%s][GetByFullName] unexpected error while getting student by fullname from storage: %v", repoImplTestPrefix, err)
	}

	if got.ID != id {
		t.Fatalf("[%s][GetByFullName] ID mismatch: got=%s want=%s", repoImplTestPrefix, got.ID, id)
	}
}

func TestRepository_Update_Success_And_Validation(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Update] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	orig, err := models.NewStudentBuilder().
		SetName("Mikahil").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{70}).
		Build()
	if err != nil {
		t.Fatalf("[%s][Update] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	id, err := repo.Create(ctx, orig)
	if err != nil {
		t.Fatalf("[%s][Create] unexpected error: %v", repoImplTestPrefix, err)
	}

	upd, err := models.NewStudentBuilder().
		SetName("Alexander").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{70, 85}).
		Build()
	if err != nil {
		t.Fatalf("[%s][Update] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	upd.ID = id
	if err := repo.Update(ctx, upd); err != nil {
		t.Fatalf("[%s][Update(valid)] unexpected error while updating student in storage: %v", repoImplTestPrefix, err)
	}

	back, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("[%s][GetByID(after update)] unexpected error: %v", repoImplTestPrefix, err)
	}

	if back.Age != 19 || len(back.Grades) != 2 {
		t.Fatalf("[%s][Update(valid)] payload mismatch, got:Age=%d Grades=%v; want: Age=%d Grades=%v",
			repoImplTestPrefix, back.Age, back.Grades, 19, []int{70, 85})
	}

	bad, err := models.NewStudentBuilder().
		SetName("Mihail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{70}).
		Build()
	if err != nil {
		t.Fatalf("[%s][Update] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	bad.ID = id
	bad.Name = "mikhail"

	if err := repo.Update(ctx, bad); err == nil {
		t.Fatalf("[%s][Update(invalid)] expected validation error for Name=%q, got nil", repoImplTestPrefix, bad.Name)
	}
}

func TestRepository_AddGrades(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][AddGrades] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	st, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{60}).
		Build()
	if err != nil {
		t.Fatalf("[%s][AddGrades] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	id, err := repo.Create(ctx, st)
	if err != nil {
		t.Fatalf("[%s][Create] unexpected error while creating Student in storage: %v", repoImplTestPrefix, err)
	}

	if err := repo.AddGrades(ctx, id, 80, 90); err != nil {
		t.Fatalf("[%s][AddGrades(valid)] unexpected error while adding grades for student: %v", repoImplTestPrefix, err)
	}

	after, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("[%s][GetByID(after add)] unexpected error while getting Student by ID: %v", repoImplTestPrefix, err)
	}

	if len(after.Grades) != 3 {
		t.Fatalf("[%s][AddGrades(valid)] grades length mismatch, got=%d want=%d (grades=%v)",
			repoImplTestPrefix, len(after.Grades), 3, after.Grades)
	}

	if err := repo.AddGrades(ctx, id, 150); err == nil {
		t.Fatalf("[%s][AddGrades(invalid)] expected validation error for grade=150, got nil", repoImplTestPrefix)
	}
}

func TestRepository_DeleteByID(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][DeleteByID] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	st, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf("[%s][DeleteByID] error while building Student domain model: %v", repoImplTestPrefix, err)
	}

	id, err := repo.Create(ctx, st)
	if err != nil {
		t.Fatalf("[%s][Create] unexpected error while creating student in storage: %v", repoImplTestPrefix, err)
	}

	if err := repo.DeleteByID(ctx, id); err != nil {
		t.Fatalf("[%s][DeleteByID] unexpected error while deleting student from storage by ID: %v", repoImplTestPrefix, err)
	}

	if _, err := repo.GetByID(ctx, id); err == nil {
		t.Fatalf("[%s][GetByID(after delete)] expected not found for ID=%s, got nil error", repoImplTestPrefix, id)
	}
}

func TestRepository_List_ReturnsCopies(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][List_ReturnsCopies] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	a, err := models.NewStudentBuilder().
		SetName("Eleven").
		SetSurname("Doctor").
		SetAge(100).
		SetGrades([]int{100}).
		Build()
	if err != nil {
		t.Fatalf("[%s][List_ReturnsCopies] error while building first Student domain model: %v", repoImplTestPrefix, err)
	}

	id1, err := repo.Create(ctx, a)
	if err != nil {
		t.Fatalf("[%s][Create(a)] unexpected error while creating first student: %v", repoImplTestPrefix, err)
	}

	b, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{80}).
		Build()
	if err != nil {
		t.Fatalf("[%s][List_ReturnsCopies] error while building second Student domain model: %v", repoImplTestPrefix, err)
	}

	if _, err = repo.Create(ctx, b); err != nil {
		t.Fatalf("[%s][Create(b)] unexpected error while creating second student: %v", repoImplTestPrefix, err)
	}

	list, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("[%s][List] unexpected error while getting list of student from storages: %v", repoImplTestPrefix, err)
	}

	if len(list) != 2 {
		t.Fatalf("[%s][List] length mismatch: got=%d want=%d", repoImplTestPrefix, len(list), 2)
	}

	list[0].Name = "Change name to check copy"

	back, err := repo.GetByID(ctx, id1)
	if err != nil {
		t.Fatalf("[%s][GetByID(after list-mutate)] unexpected error while getting student by ID from storage: %v", repoImplTestPrefix, err)
	}

	if back.Name == "Change name to check copy" {
		t.Fatalf("[%s][List] shallow copy leak: mutating list element affected repository state", repoImplTestPrefix)
	}
}

func TestRepository_Smoke_TableDriven(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Smoke] init validators: %v", repoImplTestPrefix, err)
	}

	ctx := context.Background()
	repo := NewStudentStorage()

	tests := []smokeCase{
		{
			name: "ok create and get",
			build: func() *models.Student {
				st, err := models.NewStudentBuilder().
					SetName("Mikhail").
					SetSurname("Gunin").
					SetAge(19).
					SetGrades([]int{75}).
					Build()
				if err != nil {
					t.Fatalf("[%s][Smoke/build-ok] error while building Student domain model: %v", repoImplTestPrefix, err)
				}

				return st
			},
			wantErr:   false,
			wantFound: true,
		},
		{
			name: "invalid create (name not capitalized, mutated after build)",
			build: func() *models.Student {
				st, err := models.NewStudentBuilder().
					SetName("Mikhail").
					SetSurname("Gunin").
					SetAge(19).
					Build()
				if err != nil {
					t.Fatalf("[%s][Smoke/build-ok] error while building Student domain model: %v", repoImplTestPrefix, err)
				}

				st.Name = "mikhail"

				return st
			},
			wantErr:   true,
			wantFound: false,
		},
	}

	for i, tc := range tests {
		t.Run(
			fmt.Sprintf("[%s]-smoke-%s-№%d", repoImplTestPrefix, tc.name, i+1),
			func(t *testing.T) {
				st := tc.build()

				id, err := repo.Create(ctx, st)
				gotErr := err != nil

				if gotErr != tc.wantErr {
					t.Fatalf("[%s][Smoke/Create] got error=%v, want error=%v (err=%v)", repoImplTestPrefix, gotErr, tc.wantErr, err)
				}

				if !tc.wantErr {
					_, err := repo.GetByID(ctx, id)
					present := err == nil

					if present != tc.wantFound {
						t.Fatalf("[%s][Smoke/GetByID] presence mismatch: got=%v want=%v (err=%v)", repoImplTestPrefix, present, tc.wantFound, err)
					}
				}
			},
		)
	}
}
