package models

type ExcerciseTestCaseOld struct {
	TestCaseID int  `gorm:"primary_key;column:testcase_id"`
	ExcerciseID int  `gorm:"column:excercise_id"`
	IsReady string  `gorm:"varchar(3);column:is_ready;default:'yes'"`
	TestCaseContent string  `gorm:"varchar(1024);column:testcase_content"`
	Active *string  `gorm:"varchar(3);column:active;type:enum('yes','no');default:'yes'"`
	ShowToStudent *string  `gorm:"varchar(3);column:show_to_student;type:enum('yes','no');default:'yes'"`
	TestCaseNote *string  `gorm:"varchar(1024);column:testcase_note"`
	TestCaseOutput *string  `gorm:"type:mediumtext;column:testcase_output"`
	TestCaseError *string  `gorm:"type:varchar(4096);column:testcase_error"`
}

func (ExcerciseTestCaseOld) TableName() string {
	return "excercise_testcase"
}