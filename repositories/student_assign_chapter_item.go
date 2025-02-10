package repositories

import (
	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"gorm.io/gorm"
)

type StudentAssignChapterItemRepository struct {
	DB *gorm.DB
}

func NewStudentAssignChapterItemRepository(db *gorm.DB) *StudentAssignChapterItemRepository {
	return &StudentAssignChapterItemRepository{DB: db}
}

func (studentAssignChapterItemRepo *StudentAssignChapterItemRepository) UpdateStudentAssignItemMarking(stuId string, chapterId string, itemId int, marking int) {
	studentAssignChapterItemRepo.DB.Model(&models.StudentAssignmentChapterItem{}).Where("stu_id = ? AND chapter_id = ? AND item_id = ?", stuId, chapterId, itemId).Update("marking", marking)
}
