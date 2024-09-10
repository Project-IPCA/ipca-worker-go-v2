package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExcerciseTestCaseRepositoryQ interface{
	UpdateTestCase(db_pool *gorm.DB, test_case *models.ExcerciseTestCaseOld,excercise_id int ) error
}

type ExcerciseTestCaseRepository struct{
	DB *gorm.DB
}

func NewExcerciseTestCaseRePository(db *gorm.DB) *ExcerciseTestCaseRepository{
	return &ExcerciseTestCaseRepository{
		DB: db,
	}
}


func (excerciseTestCaseRepository ExcerciseTestCaseRepository)UpdateTestCase(test_case *models.ExerciseTestcase,exercise_id uuid.UUID) error{
	update_test_case := models.ExerciseTestcase{
		TestcaseID: test_case.TestcaseID,
		ExerciseID: exercise_id,
		TestcaseContent: test_case.TestcaseContent,
		IsReady: test_case.IsReady,
		IsActive: test_case.IsActive,
		IsShowStudent: test_case.IsShowStudent,
		TestcaseNote: test_case.TestcaseNote,
		TestcaseOutput: test_case.TestcaseOutput,
		TestcaseError: test_case.TestcaseError,
	}

	if err := excerciseTestCaseRepository.DB.Model(&models.ExerciseTestcase{}).Where("testcase_id = ? AND exercise_id = ?", test_case.TestcaseID,exercise_id).Updates(update_test_case).Error; err != nil {
		excerciseTestCaseRepository.DB.Rollback()
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to update submission", err.Error())
	}

	return nil
}