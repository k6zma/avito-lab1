package repositories

import "errors"

var (
	ErrStudentAlreadyExists   = errors.New("student already exists")
	ErrStudentNotFound        = errors.New("student not found")
	ErrInvalidStudentID       = errors.New("invalid student id")
	ErrInvalidStudentSnapshot = errors.New("invalid student snapshot")
)
