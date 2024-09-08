package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
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


func (excerciseTestCaseRepository ExcerciseTestCaseRepository)UpdateTestCase(test_case *models.ExcerciseTestCaseOld,excercise_id int ) error{
	update_test_case := models.ExcerciseTestCaseOld{
		TestCaseID: test_case.TestCaseID,
		ExcerciseID: excercise_id,
		TestCaseContent: test_case.TestCaseContent,
		IsReady: "yes",
		Active: test_case.Active,
		ShowToStudent: test_case.ShowToStudent,
		TestCaseNote: test_case.TestCaseNote,
		TestCaseOutput: test_case.TestCaseOutput,
		TestCaseError: test_case.TestCaseError,
	}

	if err := excerciseTestCaseRepository.DB.Model(&models.ExcerciseTestCaseOld{}).Where("testcase_id = ? AND excercise_id = ?", test_case.TestCaseID,excercise_id).Updates(update_test_case).Error; err != nil {
		excerciseTestCaseRepository.DB.Rollback()
		return utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to update submission", err.Error())
	}

	return nil
}