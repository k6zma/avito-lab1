package persisters_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/internal/infrastructure/ciphers"
	"github.com/k6zma/avito-lab1/internal/infrastructure/persisters"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

const (
	persisterTestPrefix = "StudentPersister"

	testKey = "12345678901234567890123456789012"
)

func TestPersister_Load_NoFile(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Load_NoFile] failed to init validators: %v", persisterTestPrefix, err)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "students.json")

	persister := persisters.NewJSONStudentPersister(path, cipher)

	got, err := persister.Load(ctx)
	if err != nil {
		t.Fatalf(
			"[%s][Load_NoFile] unexpected error while loading from non-existing file: %v",
			persisterTestPrefix,
			err,
		)
	}

	if got != nil {
		t.Fatalf(
			"[%s][Load_NoFile] want nil slice on missing file, got=%v",
			persisterTestPrefix,
			got,
		)
	}
}

func TestPersister_SaveAndLoad_RoundTrip(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][SaveAndLoad] failed to init validators: %v", persisterTestPrefix, err)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "students.json")

	persister := persisters.NewJSONStudentPersister(path, cipher)

	first, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{90, 95}).
		Build()
	if err != nil {
		t.Fatalf(
			"[%s][SaveAndLoad] failed to build first student model: %v",
			persisterTestPrefix,
			err,
		)
	}

	second, err := models.NewStudentBuilder().
		SetName("Alexander").
		SetSurname("Gunin").
		SetAge(19).
		SetGrades([]int{80}).
		Build()
	if err != nil {
		t.Fatalf(
			"[%s][SaveAndLoad] failed to build second student model: %v",
			persisterTestPrefix,
			err,
		)
	}

	if err := persister.Save(ctx, []*models.Student{first, second}); err != nil {
		t.Fatalf("[%s][SaveAndLoad] save: %v", persisterTestPrefix, err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("[%s][SaveAndLoad] failed to stat snapshot file: %v", persisterTestPrefix, err)
	}

	if info.Size() == 0 {
		t.Fatalf("[%s][SaveAndLoad] something failed, snapshot file is empty", persisterTestPrefix)
	}

	loaded, err := persister.Load(ctx)
	if err != nil {
		t.Fatalf("[%s][SaveAndLoad] failed load student data: %v", persisterTestPrefix, err)
	}

	if len(loaded) != 2 {
		t.Fatalf(
			"[%s][SaveAndLoad] length mismatch: got=%d want=%d",
			persisterTestPrefix,
			len(loaded),
			2,
		)
	}

	want := map[uuid.UUID]*models.Student{
		first.ID:  first,
		second.ID: second,
	}

	for i, st := range loaded {
		ws, ok := want[st.ID]
		if !ok {
			t.Fatalf(
				"[%s][SaveAndLoad] unexpected student id at idx=%d: %s",
				persisterTestPrefix,
				i,
				st.ID,
			)
		}

		if st.Name != ws.Name || st.Surname != ws.Surname || st.Age != ws.Age {
			t.Fatalf(
				"[%s][SaveAndLoad] mismatch for id=%s: got{Name:%q,Surname:%q,Age:%d} want{Name:%q,Surname:%q,Age:%d}",
				persisterTestPrefix,
				st.ID,
				st.Name,
				st.Surname,
				st.Age,
				ws.Name,
				ws.Surname,
				ws.Age,
			)
		}

		if fmt.Sprint(st.Grades) != fmt.Sprint(ws.Grades) {
			t.Fatalf("[%s][SaveAndLoad] grades mismatch for id=%s: got=%v want=%v",
				persisterTestPrefix, st.ID, st.Grades, ws.Grades)
		}
	}
}

func TestPersister_Save_CreatesDirectories(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf(
			"[%s][Save_CreatesDirectories] failed to init validators: %v",
			persisterTestPrefix,
			err,
		)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "nested", "deep", "students.json")

	persister := persisters.NewJSONStudentPersister(path, cipher)

	student, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf(
			"[%s][Save_CreatesDirectories] failed to build student model: %v",
			persisterTestPrefix,
			err,
		)
	}

	if err := persister.Save(ctx, []*models.Student{student}); err != nil {
		t.Fatalf(
			"[%s][Save_CreatesDirectories] failed to save student data: %v",
			persisterTestPrefix,
			err,
		)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf(
			"[%s][Save_CreatesDirectories] expected snapshot file at %s, stat error: %v",
			persisterTestPrefix,
			path,
			err,
		)
	}
}

func TestPersister_Load_EmptyFile(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Load_EmptyFile] failed to init validators: %v", persisterTestPrefix, err)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "students.json")

	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatalf("[%s][Load_EmptyFile] write empty file: %v", persisterTestPrefix, err)
	}

	persister := persisters.NewJSONStudentPersister(path, cipher)

	got, err := persister.Load(ctx)
	if err != nil {
		t.Fatalf("[%s][Load_EmptyFile] failed to load students data: %v", persisterTestPrefix, err)
	}

	if got != nil {
		t.Fatalf(
			"[%s][Load_EmptyFile] want nil slice for empty file, got=%v",
			persisterTestPrefix,
			got,
		)
	}
}

func TestPersister_Load_InvalidJSON(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf("[%s][Load_InvalidJSON] failed to init validators: %v", persisterTestPrefix, err)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "students.json")

	if err := os.WriteFile(path, []byte("{{ 100% bad/invalid json =) }}"), 0o644); err != nil {
		t.Fatalf("[%s][Load_InvalidJSON] write invalid json: %v", persisterTestPrefix, err)
	}

	persister := persisters.NewJSONStudentPersister(path, cipher)

	if _, err := persister.Load(ctx); err == nil {
		t.Fatalf("[%s][Load_InvalidJSON] expected unmarshal error, got nil", persisterTestPrefix)
	}
}

func TestPersister_Save_OverwriteSnapshot(t *testing.T) {
	if err := validators.InitValidators(); err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to init validators: %v",
			persisterTestPrefix,
			err,
		)
	}

	cipher, err := ciphers.NewAESGCM(testKey)
	if err != nil {
		t.Fatalf("[%s] failed to init cipher: %v", persisterTestPrefix, err)
	}

	ctx := context.Background()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "students.json")

	persister := persisters.NewJSONStudentPersister(path, cipher)

	first, err := models.NewStudentBuilder().
		SetName("Mikhail").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to build first student model: %v",
			persisterTestPrefix,
			err,
		)
	}

	if err := persister.Save(ctx, []*models.Student{first}); err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to save first model: %v",
			persisterTestPrefix,
			err,
		)
	}

	second, err := models.NewStudentBuilder().
		SetName("Alexander").
		SetSurname("Gunin").
		SetAge(19).
		Build()
	if err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to build second student model: %v",
			persisterTestPrefix,
			err,
		)
	}

	if err := persister.Save(ctx, []*models.Student{first, second}); err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to save second student: %v",
			persisterTestPrefix,
			err,
		)
	}

	loaded, err := persister.Load(ctx)
	if err != nil {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] failed to load second student: %v",
			persisterTestPrefix,
			err,
		)
	}

	if len(loaded) != 2 {
		t.Fatalf(
			"[%s][Save_OverwriteSnapshot] want 2 students after overwrite, got=%d",
			persisterTestPrefix,
			len(loaded),
		)
	}
}

func TestPersister_Save_NilCipher(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	p := persisters.NewJSONStudentPersister(filepath.Join(tmp, "s.json"), nil)

	if err := p.Save(ctx, nil); !errors.Is(err, persisters.ErrInvalidCipher) {
		t.Fatalf("want ErrInvalidCipher, got %v", err)
	}
}

func TestPersister_Load_NilCipher(t *testing.T) {
	ctx := context.Background()
	tmp := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmp, "s.json"), []byte("non-empty"), 0o644)

	p := persisters.NewJSONStudentPersister(filepath.Join(tmp, "s.json"), nil)
	if _, err := p.Load(ctx); !errors.Is(err, persisters.ErrInvalidCipher) {
		t.Fatalf("want ErrInvalidCipher, got %v", err)
	}
}
