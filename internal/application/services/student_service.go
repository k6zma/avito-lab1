package services

import (
	"fmt"

	"github.com/k6zma/avito-lab1/internal/application/dtos"
	"github.com/k6zma/avito-lab1/internal/application/mappers"
	"github.com/k6zma/avito-lab1/internal/domain/repositories"
)

type StudentServiceContract interface {
	Register(in dtos.StudentCreateDTO) (dtos.DefaultStudentResponseDTO, error)
	Update(in dtos.StudentUpdateDTO) (dtos.DefaultStudentResponseDTO, error)
	DeleteByID(in dtos.GetByIDDTO) error
	GetByID(in dtos.GetByIDDTO) (dtos.DefaultStudentResponseDTO, error)
	GetByFullName(in dtos.GetByFullNameDTO) (dtos.DefaultStudentResponseDTO, error)
	List(includeGrades bool) ([]dtos.StudentListItemDTO, error)
	AddGrades(in dtos.AddGradesDTO) (dtos.DefaultStudentResponseDTO, error)
	AVGByID(in dtos.GetByIDDTO) (dtos.AVGResponseDTO, error)
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
	in dtos.StudentCreateDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	student, err := mappers.MapStudentCreateDTOToDomain(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map create dto to domain: %w",
			err,
		)
	}

	id, err := s.studentRepo.Create(student)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to create student in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after create: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) Update(
	in dtos.StudentUpdateDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	student, err := mappers.MapStudentUpdateDTOToDomain(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map update dto to domain: %w",
			err,
		)
	}

	if err := s.studentRepo.Update(student); err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to update student in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(student.ID)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after update: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) DeleteByID(in dtos.GetByIDDTO) error {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return fmt.Errorf("failed to map get-by-id dto to uuid: %w", err)
	}

	if err := s.studentRepo.DeleteByID(id); err != nil {
		return fmt.Errorf("failed to delete student in repository: %w", err)
	}

	return nil
}

func (s *StudentService) GetByID(
	in dtos.GetByIDDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map get-by-id dto to uuid: %w",
			err,
		)
	}

	student, err := s.studentRepo.GetByID(id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf("failed to get student by id: %w", err)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(student, true), nil
}

func (s *StudentService) GetByFullName(
	in dtos.GetByFullNameDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	name, surname, err := mappers.MapGetByFullNameDTOToArgs(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to map get-by-fullname dto: %w",
			err,
		)
	}

	student, err := s.studentRepo.GetByFullName(name, surname)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to get student by full name: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(student, true), nil
}

func (s *StudentService) List(
	includeGrades bool,
) ([]dtos.StudentListItemDTO, error) {
	list, err := s.studentRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}

	return mappers.MapStudentsDomainToListDTO(list, includeGrades), nil
}

func (s *StudentService) AddGrades(
	in dtos.AddGradesDTO,
) (dtos.DefaultStudentResponseDTO, error) {
	id, grades, err := mappers.MapAddGradesDTOToArgs(in)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf("failed to map add-grades dto: %w", err)
	}

	if err := s.studentRepo.AddGrades(id, grades...); err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to add grades in repository: %w",
			err,
		)
	}

	back, err := s.studentRepo.GetByID(id)
	if err != nil {
		return dtos.DefaultStudentResponseDTO{}, fmt.Errorf(
			"failed to fetch student after add-grades: %w",
			err,
		)
	}

	return mappers.MapStudentDomainToDefaultResponseDTO(back, true), nil
}

func (s *StudentService) AVGByID(
	in dtos.GetByIDDTO,
) (dtos.AVGResponseDTO, error) {
	id, err := mappers.MapGetByIDDTOToUUID(in)
	if err != nil {
		return dtos.AVGResponseDTO{}, fmt.Errorf("failed to map get-by-id dto to uuid: %w", err)
	}

	st, err := s.studentRepo.GetByID(id)
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
