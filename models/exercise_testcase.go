package models

import "github.com/google/uuid"

type ExcerciseTestCaseOld struct {
	TestCaseID      int     `gorm:"primary_key;column:testcase_id"`
	ExcerciseID     int     `gorm:"column:excercise_id"`
	IsReady         string  `gorm:"varchar(3);column:is_ready;default:'yes'"`
	TestCaseContent string  `gorm:"varchar(1024);column:testcase_content"`
	Active          *string `gorm:"varchar(3);column:active;type:enum('yes','no');default:'yes'"`
	ShowToStudent   *string `gorm:"varchar(3);column:show_to_student;type:enum('yes','no');default:'yes'"`
	TestCaseNote    *string `gorm:"varchar(1024);column:testcase_note"`
	TestCaseOutput  *string `gorm:"type:mediumtext;column:testcase_output"`
	TestCaseError   *string `gorm:"type:varchar(4096);column:testcase_error"`
}

func (ExcerciseTestCaseOld) TableName() string {
	return "exercise_testcases"
}

type ExerciseTestcase struct {
	TestcaseID      *uuid.UUID `gorm:"column:testcase_id;type:varchar(36);primaryKey"`
	ExerciseID      uuid.UUID  `gorm:"column:exercise_id;type:varchar(36);not null"`
	IsReady         string     `gorm:"column:is_ready;type:varchar(3);not null;default:'yes'"`
	TestcaseContent string     `gorm:"column:testcase_content;type:varchar(1024);not null"`
	IsActive        *bool      `gorm:"column:is_active;type:tinyint(1);default:1"`
	IsShowStudent   *bool      `gorm:"column:is_show_student;type:tinyint(1);default:1"`
	TestcaseNote    *string    `gorm:"column:testcase_note;type:varchar(1024)"`
	TestcaseOutput  *string    `gorm:"column:testcase_output;type:mediumtext"`
	TestcaseError   *string    `gorm:"column:testcase_error;type:varchar(4096)"`
}

func (ExerciseTestcase) TableName() string {
	return "exercise_testcases"
}
