package repositories

import (
	"context"

	"github.com/google/uuid"

	"github.com/k6zma/avito-lab1/internal/domain/models"
)

type StudentRepository interface {
	Create(ctx context.Context, student *models.Student) (uuid.UUID, error)
	Update(ctx context.Context, student *models.Student) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Student, error)
	GetByFullName(ctx context.Context, name, surname string) (*models.Student, error)
	List(ctx context.Context) ([]*models.Student, error)
	AddGrades(ctx context.Context, id uuid.UUID, grades ...int) error
}
