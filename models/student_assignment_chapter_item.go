package models

import (
	"time"

	"github.com/google/uuid"
)

type StudentAssignmentChapterItem struct {
	StuID      uuid.UUID    `gorm:"type:varchar(36);primary_key;column:stu_id"`
	ChapterID  uuid.UUID    `gorm:"type:varchar(36);primary_key;column:chapter_id"`
	ItemID     int    `gorm:"type:varchar(36);primary_key;column:item_id"`
	ExerciseID *uuid.UUID    `gorm:"type:varchar(36);primary_key;column:exercise_id"`
	FullMark   int       `gorm:"type:int;not null;default:0;column:full_mark"`
	Marking    int       `gorm:"type:int;not null;default:0;column:marking"`
	AddedDate  time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;column:added_date"`
	TimeStart  *string   `gorm:"type:varchar(8);column:time_start"`
	TimeEnd    *string   `gorm:"type:varchar(8);column:time_end"`
}

func (StudentAssignmentChapterItem) TableName() string{
	return "student_assignment_chapter_items"
}