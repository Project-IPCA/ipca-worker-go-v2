package models

import (
	// "time"

	// "github.com/google/uuid"
)

type ReciveMessage struct {
	JobID	string `json:"job_id"`
	JobType string `json:"job_type"`
	LogData LogDataType `json:"log_data"`
	SubmissionID *string `json:"submission_id"`
	ExcerciseID *string `json:"exercise_id"`
	SourceCode string `json:"sourcecode"`
	TestCaseList []TestCase `json:"testcase_list"`
	ChapterId string `json:"chapter_id"`
	ItemId int	`json:"item_id"`
	StudentId string `json:"stu_id"`
}

type TestCase struct{
	TestCaseID string `json:"testcase_id"`
	ExceriseID string `json:"excerise_id"`
	IsReady string `json:"is_ready"`
	TestCaseContent string `json:"testcase_content"`
	Active bool `json:"is_active"`
	ShowToStudent bool `json:"show_to_student"`
	TestCaseNote string `json:"testcase_note"`
	TestCaseOutput string `json:"testcase_output"`
	TestCaseError string `json:"testcase_error"`
}

type LogDataType struct{
	GroupID string `json:"group_id"`
	Username string `json:"username"`
	RemoteIP string `json:"remote_ip"`
	Agent string `json:"agent"`
	PageName string `json:"page_name"`
	Actoin ActionData `json:"action"`
}

type ActionData struct{
	StudentID string `json:"stu_id"`
	JobID string `json:"job_id"`
	Status string `json:"status"`
	SubmissionID string `json:"submission_id"`
	Attempt string `json:"attempt"`
	SourcecodeFilename string `json:"sourcecode_filename"`
	Marking *int `json:"marking"`
}