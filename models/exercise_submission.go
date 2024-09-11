package models

import (
	"time"

	"github.com/google/uuid"
)

type ExerciseSubmission struct {
	SubmissionID        uuid.UUID    `gorm:"type:varchar(36);primary_key;column:submission_id"`
	StuID               uuid.UUID    `gorm:"type:varchar(36);not null;column:stu_id"`
	ExerciseID          uuid.UUID    `gorm:"type:varchar(36);not null;column:exercise_id"`
	Status              string    `gorm:"type:enum('ACCEPTED','PENDING');not null;default:'PENDING';column:status"`
	SourcecodeFilename  string    `gorm:"type:varchar(40);not null;column:sourcecode_filename"`
	Marking             int       `gorm:"type:int;not null;default:0;column:marking"`
	TimeSubmit          time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP;column:time_submit"`
	IsInfLoop           *bool     `gorm:"type:tinyint(1);column:is_inf_loop"`
	Output              *string   `gorm:"type:text;column:output"`
	Result              *string   `gorm:"type:json;column:result"`
	ErrorMessage        *string   `gorm:"type:mediumtext;column:error_message"`
}

func (ExerciseSubmission) TableName() string{
	return "exercise_submissions"
}

type UpdateSubmissionInfo struct {
	SubmissionID      uuid.UUID      
	Status            string    
	Marking           int      
	Result            *string
	ErrorMessage      *string
}