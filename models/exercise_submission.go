package models

import (
	"encoding/json"
	"time"
)

type ExcerciseSubmission struct {
	SubmissionID      int            `gorm:"column:submission_id;primaryKey;autoIncrement"`
	StuID             int            `gorm:"column:stu_id"`
	ExerciseID        int            `gorm:"column:exercise_id"`
	Status            string         `gorm:"column:status;type:enum('accepted','wrong_answer','pending','rejected','error');default:'pending'"`
	SourcecodeFilename string         `gorm:"column:sourcecode_filename;type:varchar(40)"`
	Marking           int            `gorm:"column:marking"`
	TimeSubmit        *time.Time      `gorm:"column:time_submit;default:CURRENT_TIMESTAMP"`
	InfLoop           *string         `gorm:"column:inf_loop;type:enum('Yes','No')"`
	Output            *string         `gorm:"column:output;type:varchar(16384)"`
	Result            *json.RawMessage      `gorm:"column:result;type:json"`
	ErrorMessage      *string         `gorm:"column:error_message;type:mediumtext"`
}

func (ExcerciseSubmission) TableName() string {
	return "exercise_submission"
}

type UpdateSubmissionInfo struct {
	SubmissionID      int         
	Status            string    
	Marking           int      
	Result            *json.RawMessage
	ErrorMessage      *string
}