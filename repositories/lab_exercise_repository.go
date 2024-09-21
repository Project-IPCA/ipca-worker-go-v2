package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LabExerciseRepository struct {
	DB *gorm.DB
}

func NewLabExerciseRePository(db *gorm.DB) *LabExerciseRepository {
	return &LabExerciseRepository{
		DB: db,
	}
}

func (labExerciseRepo *LabExerciseRepository) GetLabExerciseById(labExercise *models.LabExercise,exerciseId uuid.UUID) {
	labExerciseRepo.DB.Where("exercise_id",exerciseId).Find(labExercise)
}

func (labExerciseRepo *LabExerciseRepository) UpdateExerciseTestcaseEnum(labExercise *models.LabExercise,testcaseEnum string) {
	labExerciseRepo.DB.Model(labExercise).Update("testcase" , testcaseEnum)
}