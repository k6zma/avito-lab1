package tui

type mode int

const (
	modeMenu mode = iota
	modeTable
	modeCreate
	modeAddGrades
	modeIDInput
	modeDetail

	actionAVG  = "avg"
	actionDel  = "del"
	actionShow = "show"
)

type (
	menuChoiceMsg string

	tableBackMsg struct{}

	tableShowMsg struct {
		ID string
	}

	createSubmittedMsg struct {
		Name, Surname, Age, GradesCSV string
	}

	createCancelMsg struct{}

	addGradesSubmittedMsg struct {
		ID, Grades string
	}

	addGradesCancelMsg struct{}

	idSubmittedMsg string

	idCancelMsg struct{}
)
