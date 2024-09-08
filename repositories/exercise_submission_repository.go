package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"gorm.io/gorm"
)

type ExcerciseSubmissionRepositoryQ interface{
	UpdateSubmission(submission_data *models.UpdateSubmissionInfo) error
}

type ExcerciseSubmissionRepository struct{
	DB *gorm.DB
}

func NewExcerciseSubmissionRePository(db *gorm.DB) *ExcerciseSubmissionRepository{
	return &ExcerciseSubmissionRepository{
		DB: db,
	}
}

func (excerciseSubmissionRepository ExcerciseSubmissionRepository)UpdateSubmission(submission_data *models.UpdateSubmissionInfo) error{
	update_submission := models.ExcerciseSubmission{
		Status: submission_data.Status,
		Marking: submission_data.Marking,
		Result: submission_data.Result,
		ErrorMessage: submission_data.ErrorMessage,
	}

	if err := excerciseSubmissionRepository.DB.Model(&models.ExcerciseSubmission{}).Where("submission_id = ?", submission_data.SubmissionID).Updates(update_submission).Error; err != nil {
		excerciseSubmissionRepository.DB.Rollback()
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to update submission", err.Error())
	}

	return nil
}