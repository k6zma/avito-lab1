package main

import (
	"os"
	"testing"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/services"
	"github.com/k6zma/avito-lab1/internal/infrastructure/ciphers"
	"github.com/k6zma/avito-lab1/internal/infrastructure/persisters"
	infrastructureRepos "github.com/k6zma/avito-lab1/internal/infrastructure/repositories"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	testCipherKey = "abcdefghijklmnopqrstuvwxyz123456"
	testFilePath  = "test_students.json"

	aliceTestName    = "Alice"
	aliceTestSurname = "Cooper"

	bobTestName    = "Bob"
	bobTestSurname = "Dylan"
)

func newStudentService(t *testing.T, jsonDataPath string) services.StudentServiceContract {
	t.Helper()

	if err := validators.InitValidators(); err != nil {
		t.Fatalf("Failed to init validators: %v", err)
	}

	cipher, err := ciphers.NewAESGCM(testCipherKey)
	if err != nil {
		t.Fatalf("Failed to init AES-GCM cipher: %v", err)
	}

	persister := persisters.NewJSONStudentPersister(jsonDataPath, cipher)

	repository, err := infrastructureRepos.NewStudentStorageWithPersister(persister)
	if err != nil {
		t.Fatalf("Failed to init repository: %v", err)
	}

	return services.NewStudentService(repository)
}

func TestStudentStorage(t *testing.T) {
	_ = os.Remove(testFilePath)
	t.Cleanup(func() {
		_ = os.Remove(testFilePath)
	})

	svc := newStudentService(t, testFilePath)

	var aliceID string

	t.Run("AddStudent", func(t *testing.T) {
		created, err := svc.Register(dtos.StudentCreateDTO{
			Name:    aliceTestName,
			Surname: aliceTestSurname,
			Age:     20,
			Grades:  []int{90, 85, 95},
		})
		if err != nil {
			t.Errorf("Failed to add student: %v", err)
		}

		aliceID = created.ID

		_, err = svc.Register(dtos.StudentCreateDTO{
			Name:    aliceTestName,
			Surname: aliceTestSurname,
			Age:     20,
			Grades:  []int{90, 85, 95},
		})
		if err != nil {
			t.Fatalf("Unexpected error when adding same full name again: %v", err)
		}
	})

	t.Run("GetStudent", func(t *testing.T) {
		got, err := svc.GetByFullName(dtos.GetByFullNameDTO{
			Name:    aliceTestName,
			Surname: aliceTestSurname,
		})
		if err != nil {
			t.Errorf("Failed to get student: %v", err)
		}

		if got.Name != aliceTestName || got.Surname != aliceTestSurname || got.Age != 20 ||
			len(got.Grades) != 3 {
			t.Error("Student data doesn't match")
		}

		if _, err := svc.GetByFullName(dtos.GetByFullNameDTO{
			Name:    bobTestName,
			Surname: bobTestSurname,
		}); err == nil {
			t.Error("Expected error when getting non-existent student")
		}
	})

	t.Run("UpdateStudent", func(t *testing.T) {
		upd, err := svc.Update(dtos.StudentUpdateDTO{
			ID:      aliceID,
			Name:    aliceTestName,
			Surname: aliceTestSurname,
			Age:     21,
			Grades:  []int{95, 90, 100},
		})
		if err != nil {
			t.Errorf("Failed to update student: %v", err)
		}

		gotUpd, err := svc.GetByID(dtos.GetByIDDTO{ID: upd.ID})
		if err != nil {
			t.Errorf("Failed to get student: %v", err)
		}

		if gotUpd.Age != 21 || len(gotUpd.Grades) != 3 || gotUpd.Grades[0] != 95 {
			t.Error("Student data wasn't updated correctly")
		}

		if _, err := svc.Update(dtos.StudentUpdateDTO{
			ID:      "00000000-0000-0000-0000-000000000000",
			Name:    bobTestName,
			Surname: bobTestSurname,
			Age:     21,
			Grades:  []int{95, 90, 100},
		}); err == nil {
			t.Error("Expected error when updating non-existent student")
		}
	})

	t.Run("CalculateAverageGrade", func(t *testing.T) {
		avg, err := svc.AVGByID(dtos.GetByIDDTO{ID: aliceID})
		if err != nil {
			t.Errorf("Failed to calculate average grade: %v", err)
		}

		expected := (95.0 + 90.0 + 100.0) / 3.0
		if avg.AVG != expected {
			t.Errorf("Expected average %.2f, got %.2f", expected, avg.AVG)
		}

		noGrades, err := svc.Register(dtos.StudentCreateDTO{
			Name:    bobTestName,
			Surname: bobTestSurname,
			Age:     22,
			Grades:  []int{},
		})
		if err != nil {
			t.Errorf("Failed to add student: %v", err)
		}

		avg2, err := svc.AVGByID(dtos.GetByIDDTO{ID: noGrades.ID})
		if err != nil {
			t.Errorf("Failed to calculate average grade: %v", err)
		}

		if avg2.AVG != 0 {
			t.Errorf("Expected average 0 for student without grades, got %.2f", avg.AVG)
		}
	})

	t.Run("SaveAndLoad", func(t *testing.T) {
		svc2 := newStudentService(t, testFilePath)

		got, err := svc2.GetByFullName(dtos.GetByFullNameDTO{
			Name:    aliceTestName,
			Surname: aliceTestSurname,
		})
		if err != nil {
			t.Fatalf("Failed to get student after reload: %v", err)
		}

		if got.Name != aliceTestName || got.Surname != aliceTestSurname || got.Age != 21 ||
			len(got.Grades) != 3 {
			t.Error("Loaded student data doesn't match")
		}
	})
}
