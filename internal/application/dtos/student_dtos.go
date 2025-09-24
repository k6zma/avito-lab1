package dtos

type StudentCreateDTO struct {
	Name    string `json:"name"    validate:"required,capitalized"`
	Surname string `json:"surname" validate:"required,capitalized"`
	Age     int    `json:"age"     validate:"gte=0,lte=150"`
	Grades  []int  `json:"grades"  validate:"omitempty,dive,gte=0,lte=100"`
}

type StudentUpdateDTO struct {
	ID      string `json:"id"      validate:"required,uuid4"`
	Name    string `json:"name"    validate:"required,capitalized"`
	Surname string `json:"surname" validate:"required,capitalized"`
	Age     int    `json:"age"     validate:"gte=0,lte=150"`
	Grades  []int  `json:"grades"  validate:"omitempty,dive,gte=0,lte=100"`
}

type AddGradesDTO struct {
	ID     string `json:"id"     validate:"required,uuid4"`
	Grades []int  `json:"grades" validate:"required,min=1,dive,gte=0,lte=100"`
}

type GetByFullNameDTO struct {
	Name    string `json:"name"    validate:"required,capitalized"`
	Surname string `json:"surname" validate:"required,capitalized"`
}

type GetByIDDTO struct {
	ID string `json:"id" validate:"required,uuid"`
}

type DefaultStudentResponseDTO struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Age      int      `json:"age"`
	Grades   []int    `json:"grades,omitempty"`
	AvgGrade *float64 `json:"avg_grade,omitempty"`
}

type StudentListItemDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
	Grades  []int  `json:"grades,omitempty"`
}

type AVGResponseDTO struct {
	ID  string  `json:"id"`
	AVG float64 `json:"avg"`
}
