package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LabExercise struct {
	ExerciseID             uuid.UUID        `gorm:"type:varchar(36);primaryKey;column:exercise_id"`
	ChapterID              *uuid.UUID       `gorm:"type:varchar(36);column:chapter_id"`
	Level                  *string          `gorm:"type:enum('0','1','2','3','4','5','6')"`
	Name                   *string          `gorm:"type:varchar(1024)"`
	Content                *string          `gorm:"type:mediumtext"`
	Testcase               string           `gorm:"type:enum('NO_INPUT','YES','UNDEFINED');not null;default:'NO_INPUT'"`
	Sourcecode             *string          `gorm:"type:varchar(50)"`
	FullMark               int              `gorm:"type:int;not null;default:10"`
	AddedDate              time.Time        `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
	LastUpdate             *time.Time       `gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
	UserDefinedConstraints *json.RawMessage `gorm:"type:json"`
	SuggestedConstraints   *json.RawMessage `gorm:"type:json"`
	AddedBy                *string          `gorm:"type:varchar(40)"`
	CreatedBy              *uuid.UUID       `gorm:"type:varchar(36)"`
}

func (LabExercise) TableName() string {
	return "lab_exercises"
}
