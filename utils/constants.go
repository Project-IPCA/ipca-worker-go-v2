package utils

type ExerciseStatusStruct struct {
	Accepted    string
	WrongAnswer string
	Pending     string
	Rejected    string
	Error       string
}

var ExerciseStatus = ExerciseStatusStruct{
	Accepted:    "ACCEPTED",
	WrongAnswer: "WRONG_ANSWER",
	Pending:     "PENDING",
	Rejected:    "REJECTED",
	Error:       "ERROR",
}