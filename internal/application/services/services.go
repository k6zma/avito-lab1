package services

import (
	"context"
	"fmt"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/mappers"
	"github.com/k6zma/avito-lab1/internal/domain/repositories"
)

type StudentServiceContract interface {
	Register(ctx context.Context, in dtos.StudentCreateDTO) (dtos.DefaultStudentResponseDTO, error)
	Update(ctx context.Context, in dtos.StudentUpdateDTO) (dtos.DefaultStudentResponseDTO, error)
	DeleteByID(ctx context.Context, in dtos.GetByIDDTO) error

	GetByID(ctx context.Context, in dtos.GetByIDDTO) (dtos.DefaultStudentResponseDTO, error)
	GetByFullName(
		ctx context.Context,
		in dtos.GetByFullNameDTO,
	) (dtos.DefaultStudentResponseDTO, error)
	List(ctx context.Context, includeGrades bool) ([]dtos.StudentListItemDTO, error)

	AddGrades(ctx context.Context, in dtos.AddGradesDTO) (dtos.DefaultStudentResponseDTO, error)
	AVGByID(ctx context.Context, in dtos.GetByIDDTO) (dtos.AVGResponseDTO, error)
}

type StudentService struct {
	studentRepo repositories.StudentRepository
}

func NewStudentService(repo repositories.StudentRepository) StudentServiceContract {
	return &StudentService{
		studentRepo: repo,
	}
}

func (s *StudentService) Register(
	ctx context.Context,
	in dtos.StudentCreateDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	student, err := mappers.MapStudentCreateDTOToDomain(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map create dto to domain: %w",
			err,
		)
	}

	id, err := s.studentRepo.Create(ctx, student)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to create student in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after create: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) Update(
	ctx context.Context,
	in dtos.StudentUpdateDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	student, err := mappers.MapStudentUpdateDTOToDomain(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map update dto to domain: %w",
			err,
		)
	}

	if err := s.studentRepo.Update(ctx, student); err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to update student in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(ctx, student.ID)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after update: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) DeleteByID(ctx context.Context, in dtos.GetByIDDTO) error {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return fmt.Errorf("failed to map get-by-id dto to uuid: %w", err)
	}

	if err := s.studentRepo.DeleteByID(ctx, id); err != nil {
		return fmt.Errorf("failed to delete student in repository: %w", err)
	}

	return nil
}

func (s *StudentService) GetByID(
	ctx context.Context,
	in dtos.GetByIDDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map get-by-id dto to uuid: %w",
			err,
		)
	}

	student, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf("failed to get student by id: %w", err)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(student, true), nil
}

func (s *StudentService) GetByFullName(
	ctx context.Context,
	in dtos.GetByFullNameDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	name, surname, err := mappers.MapGetByFullNameDTOToArgs(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map get-by-fullname dto: %w",
			err,
		)
	}

	student, err := s.studentRepo.GetByFullName(ctx, name, surname)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to get student by full name: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(student, true), nil
}

func (s *StudentService) List(
	ctx context.Context,
	includeGrades bool,
) ([]dtos.StudentListItemDTO, error) {
	list, err := s.studentRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	return mappers.MapStudentsDomainToListDTO(list, includeGrades), nil
}

func (s *StudentService) AddGrades(
	ctx context.Context,
	in dtos.AddGradesDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	id, grades, err := mappers.MapAddGradesDTOToArgs(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf("failed to map add-grades dto: %w", err)
	}

	if err := s.studentRepo.AddGrades(ctx, id, grades...); err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to add grades in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after add-grades: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) AVGByID(
	ctx context.Context,
	in dtos.GetByIDDTO,
) (dtos.AVGResponseDTO, error) {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return dtos.AVGResponseDTO{}, fmt.Errorf("failed to map get-by-id dto to uuid: %w", err)
	}

	st, err := s.studentRepo.GetByID(ctx, id)
	if err != nil {
		return dtos.AVGResponseDTO{}, fmt.Errorf("failed to get student by id: %w", err)
	}

	avg := 0.0
	if n := len(st.Grades); n > 0 {
		sum := 0

		for _, g := range st.Grades {
			sum += g
		}

		avg = float64(sum) / float64(n)
	}

	return dtos.AVGResponseDTO{ID: st.ID.String(), AVG: avg}, nil
}
