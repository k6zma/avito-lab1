package repositories

import (
	"context"
	"fmt"
	"github.com/k6zma/avito-lab1/internal/domain/persisters"
	"sync"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/internal/domain/repositories"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

type StudentStorage struct {
	students  map[uuid.UUID]*models.Student
	persister persisters.StudentPersister
	mu        sync.RWMutex
}

func NewStudentStorageWithPersister(ctx context.Context, p persisters.StudentPersister) (*StudentStorage, error) {
	s := &StudentStorage{
		students:  make(map[uuid.UUID]*models.Student),
		persister: p,
	}

	if p == nil {
		return s, nil
	}

	sts, err := p.Load(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load students snapshot: %w", err)
	}

	for _, st := range sts {
		if st == nil || st.ID == uuid.Nil {
			return nil, repositories.ErrStudentNotFound
		}

		if err := validators.Validate.Struct(st); err != nil {
			return nil, fmt.Errorf("failed to validate student from snapshot: %w", err)
		}

		s.students[st.ID] = st.Clone()
	}

	return s, nil
}

func (s *StudentStorage) Create(ctx context.Context, student *models.Student) (uuid.UUID, error) {
	cp := student.Clone()

	if err := validators.Validate.Struct(cp); err != nil {
		return uuid.Nil, fmt.Errorf("input student is invalid: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.students[cp.ID]; ok {
		return uuid.Nil, repositories.ErrStudentAlreadyExists
	}

	s.students[cp.ID] = cp

	if s.persister != nil {
		students := make([]*models.Student, 0, len(s.students))

		for _, st := range s.students {
			students = append(students, st.Clone())
		}

		if err := s.persister.Save(ctx, students); err != nil {
			delete(s.students, cp.ID)

			return uuid.Nil, fmt.Errorf("persist student data after create failed: %w", err)
		}
	}

	return cp.ID, nil
}

func (s *StudentStorage) Update(ctx context.Context, student *models.Student) error {
	cp := student.Clone()

	if err := validators.Validate.Struct(cp); err != nil {
		return fmt.Errorf("input student is invalid: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.students[cp.ID]
	if !ok {
		return repositories.ErrStudentNotFound
	}

	s.students[cp.ID] = cp

	if s.persister != nil {
		students := make([]*models.Student, 0, len(s.students))

		for _, st := range s.students {
			students = append(students, st.Clone())
		}

		if err := s.persister.Save(ctx, students); err != nil {
			s.students[cp.ID] = prev

			return fmt.Errorf("persist student data after update failed: %w", err)
		}
	}

	return nil
}

func (s *StudentStorage) DeleteByID(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return repositories.ErrInvalidStudentID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	prev, ok := s.students[id]
	if !ok {
		return repositories.ErrStudentNotFound
	}

	delete(s.students, id)

	if s.persister != nil {
		students := make([]*models.Student, 0, len(s.students))

		for _, st := range s.students {
			students = append(students, st.Clone())
		}

		if err := s.persister.Save(ctx, students); err != nil {
			s.students[id] = prev

			return fmt.Errorf("persist student data after delete failed: %w", err)
		}
	}

	return nil
}

func (s *StudentStorage) GetByID(_ context.Context, id uuid.UUID) (*models.Student, error) {
	if id == uuid.Nil {
		return nil, repositories.ErrInvalidStudentID
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	student, ok := s.students[id]
	if !ok {
		return nil, repositories.ErrStudentNotFound
	}

	return student.Clone(), nil
}

func (s *StudentStorage) GetByFullName(
	_ context.Context,
	name string,
	surname string,
) (*models.Student, error) {
	if err := validators.Validate.Var(name, "required,capitalized"); err != nil {
		return nil, fmt.Errorf("invalid student name: %w", err)
	}

	if err := validators.Validate.Var(surname, "required,capitalized"); err != nil {
		return nil, fmt.Errorf("invalid student surname: %w", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, student := range s.students {
		if student.Name == name && student.Surname == surname {
			return student.Clone(), nil
		}
	}

	return nil, repositories.ErrStudentNotFound
}

func (s *StudentStorage) List(_ context.Context) ([]*models.Student, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	students := make([]*models.Student, 0, len(s.students))
	for _, student := range s.students {
		students = append(students, student.Clone())
	}

	return students, nil
}

func (s *StudentStorage) AddGrades(ctx context.Context, id uuid.UUID, grades ...int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.students[id]
	if !ok {
		return repositories.ErrStudentNotFound
	}

	cp := current.Clone()

	if err := cp.AddGrades(grades...); err != nil {
		return fmt.Errorf("error while adding grades for student in storage: %w", err)
	}

	s.students[id] = cp

	if s.persister != nil {
		students := make([]*models.Student, 0, len(s.students))

		for _, st := range s.students {
			students = append(students, st.Clone())
		}

		if err := s.persister.Save(ctx, students); err != nil {
			s.students[id] = current

			return fmt.Errorf("persist student data after add-grades failed: %w", err)
		}
	}

	return nil
}
