package persisters

import (
	"context"

	"github.com/k6zma/avito-lab1/internal/domain/models"
)

type StudentPersister interface {
	Save(ctx context.Context, students []*models.Student) error
	Load(ctx context.Context) ([]*models.Student, error)
}
