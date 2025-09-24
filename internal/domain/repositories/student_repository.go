package repositories

import (
	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/domain/models"
)

type StudentRepository interface {
	Create(student *models.Student) (uuid.UUID, error)
	Update(student *models.Student) error
	DeleteByID(id uuid.UUID) error
	GetByID(id uuid.UUID) (*models.Student, error)
	GetByFullName(name, surname string) (*models.Student, error)
	List() ([]*models.Student, error)
	AddGrades(id uuid.UUID, grades ...int) error
}
