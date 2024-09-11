package repositories

import (
	"encoding/json"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityLogRepositoryQ interface{
	AddSubmissionLog(db_pool *gorm.DB, log_data *models.LogDataType) (*models.ActivityLogOld, error)
}

type ActivityLogRepository struct{
	DB *gorm.DB
}

func NewActivityLogRePository(db *gorm.DB) *ActivityLogRepository{
	return &ActivityLogRepository{
		DB: db,
	}
}

func (activityLogRepository ActivityLogRepository)AddSubmissionLog(log_data *models.LogDataType) (*models.ActivityLog, error) {
	groupUuid,err := uuid.Parse(log_data.GroupID)
	if(err!=nil){
		return &models.ActivityLog{},utils.NewAppError(utils.ERROR_NAME.FUNCTION_ERROR,"failed to convert group id", err.Error())
	}
	action_str,err := json.Marshal(log_data.Actoin)
	if(err!=nil){
		return &models.ActivityLog{},utils.NewAppError(utils.ERROR_NAME.FUNCTION_ERROR,"failed to deserialize", err.Error())
	}
	logId := utils.NewULID()
	add_Log := models.ActivityLog{
		LogID: logId,
		GroupID: &groupUuid,
		Username: log_data.Username,
		RemoteIP: log_data.RemoteIP,
		Agent: &log_data.Agent,
		PageName: log_data.PageName,
		Action: string(action_str),
	}	

	if err := activityLogRepository.DB.Create(&add_Log).Error; err != nil {
		activityLogRepository.DB.Rollback()
		return &models.ActivityLog{}, utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to insert log", err.Error())
	}
	return &add_Log, nil
}
