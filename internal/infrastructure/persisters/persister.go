package persisters

import (
	"github.com/k6zma/avito-lab1/internal/domain/models"
)

type StudentPersister interface {
	Save(students []*models.Student) error
	Load() ([]*models.Student, error)
}
