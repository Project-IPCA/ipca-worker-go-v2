package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"gorm.io/gorm"
)

type ExerciseSubmissionRepositoryQ interface {
	UpdateSubmission(submission_data *models.UpdateSubmissionInfo) error
}

type ExerciseSubmissionRepository struct {
	DB *gorm.DB
}

func NewExerciseSubmissionRePository(db *gorm.DB) *ExerciseSubmissionRepository {
	return &ExerciseSubmissionRepository{
		DB: db,
	}
}

func (repo *ExerciseSubmissionRepository) UpdateSubmission(submissionData *models.UpdateSubmissionInfo) error {
	tx := repo.DB.Begin()
	if tx.Error != nil {
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "failed to begin transaction", tx.Error.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	update_submission := models.ExerciseSubmission{
		Status:       submissionData.Status,
		Marking:      submissionData.Marking,
		Result:       submissionData.Result,
		ErrorMessage: submissionData.ErrorMessage,
	}

	if err := tx.Model(&models.ExerciseSubmission{}).
		Where("submission_id = ?", submissionData.SubmissionID).
		Updates(update_submission).Error; err != nil {
		tx.Rollback()
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "failed to update submission", err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "failed to commit transaction", err.Error())
	}

	return nil
}
