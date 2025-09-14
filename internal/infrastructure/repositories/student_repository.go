package repositories

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/domain/models"
	"github.com/k6zma/avito-lab1/internal/domain/repositories"
	"github.com/k6zma/avito-lab1/pkg/validators"
)

type StudentStorage struct {
	students map[uuid.UUID]*models.Student
	mu       sync.RWMutex
}

func NewStudentStorage() *StudentStorage {
	return &StudentStorage{
		students: make(map[uuid.UUID]*models.Student),
	}
}

func (s *StudentStorage) Create(_ context.Context, student *models.Student) (uuid.UUID, error) {
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

	return cp.ID, nil
}

func (s *StudentStorage) Update(_ context.Context, student *models.Student) error {
	cp := student.Clone()

	if err := validators.Validate.Struct(cp); err != nil {
		return fmt.Errorf("input student is invalid: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.students[cp.ID]; !ok {
		return repositories.ErrStudentNotFound
	}

	s.students[cp.ID] = cp

	return nil
}

func (s *StudentStorage) DeleteByID(_ context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return repositories.ErrInvalidStudentID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.students[id]; !ok {
		return repositories.ErrStudentNotFound
	}

	delete(s.students, id)

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
	ctx context.Context,
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

func (s *StudentStorage) AddGrades(_ context.Context, id uuid.UUID, grades ...int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	student, ok := s.students[id]
	if !ok {
		return repositories.ErrStudentNotFound
	}

	if err := student.AddGrades(grades...); err != nil {
		return fmt.Errorf("error while adding grades: %w", err)
	}

	return nil
}
