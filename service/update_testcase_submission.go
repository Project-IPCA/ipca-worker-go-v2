package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/redis_client"
	"github.com/Project-IPCA/ipca-worker-go-v2/repositories"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func AddAndUpdateTestCase(channel *amqp.Channel, db_pool *gorm.DB, msg amqp.Delivery, msgBody models.ReciveMessage, redisClient *redis.Client){
	publisher := redis_client.NewRedisAction(redisClient)
	err := compileCodeTestcase(db_pool,msgBody)
	if(err!=nil){
		channel.Nack(msg.DeliveryTag,false,false)
		return
	}
	err = publisher.PublishMessage(fmt.Sprintf("submission-result:%s", msgBody.JobID),"done")
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		return
	}
	fmt.Println("FINISHED RUNNING")
	channel.Ack(msg.DeliveryTag, false)
}

func compileCodeTestcase (db_pool *gorm.DB, msgBody models.ReciveMessage) error{
	exerciseTestcaseRepo := repositories.NewExcerciseTestCaseRePository(db_pool)
	if(len(msgBody.TestCaseList)>0){
		for i:=0; i<len(msgBody.TestCaseList); i++ {
			testCaseInt,err := strconv.Atoi(msgBody.TestCaseList[i].TestCaseID)
			if(err!=nil){
				utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
			}
			excerciseInt,err := strconv.Atoi(*msgBody.ExcerciseID)
			if(err!=nil){
				utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
			}
			result, err := utils.RunPythonScript(msgBody.TestCaseList[i], msgBody.SourceCode)
			if err != nil {
				appErr, ok := err.(*utils.AppError)
				if(ok){
					fmt.Println("Error running Python script:", appErr)
					errorMessage := string(appErr.Err.Error())
					testCaseData := models.ExcerciseTestCaseOld{
						TestCaseID: testCaseInt,
						ExcerciseID: excerciseInt,
						TestCaseContent: msgBody.TestCaseList[i].TestCaseContent,
						IsReady: msgBody.TestCaseList[i].IsReady,
						Active: &msgBody.TestCaseList[i].Active,
						ShowToStudent: &msgBody.TestCaseList[i].ShowToStudent,
						TestCaseNote: &msgBody.TestCaseList[i].TestCaseNote,
						TestCaseOutput: &appErr.Stdout,
						TestCaseError: &errorMessage,
					}
					err = exerciseTestcaseRepo.UpdateTestCase(&testCaseData,excerciseInt)
					if(err != nil){
						utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
					}
					for j:=i+1; j<len(msgBody.TestCaseList); j++{
						errorText := "Error occurred before this testcase"
						testCaseData := models.ExcerciseTestCaseOld{
							TestCaseID: testCaseInt,
							ExcerciseID: excerciseInt,
							TestCaseContent: msgBody.TestCaseList[j].TestCaseContent,
							IsReady: msgBody.TestCaseList[j].IsReady,
							Active: &msgBody.TestCaseList[j].Active,
							ShowToStudent: &msgBody.TestCaseList[j].ShowToStudent,
							TestCaseNote: &msgBody.TestCaseList[j].TestCaseNote,
							TestCaseOutput: &appErr.Stdout,
							TestCaseError: &errorText,
						}
						err = exerciseTestcaseRepo.UpdateTestCase(&testCaseData,excerciseInt)
						if(err != nil){
							utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
						}
					}
					return utils.NewAppError(appErr.Name,appErr.Error(), appErr.Stdout)
				}
			}else{
				output := strings.TrimSpace(result)
				testCaseData := models.ExcerciseTestCaseOld{
					TestCaseID: testCaseInt,
					ExcerciseID: excerciseInt,
					TestCaseContent: msgBody.TestCaseList[i].TestCaseContent,
					IsReady: msgBody.TestCaseList[i].IsReady,
					Active: &msgBody.TestCaseList[i].Active,
					ShowToStudent: &msgBody.TestCaseList[i].ShowToStudent,
					TestCaseNote: &msgBody.TestCaseList[i].TestCaseNote,
					TestCaseOutput: &output,
					TestCaseError: nil,
				}

				err = exerciseTestcaseRepo.UpdateTestCase(&testCaseData,excerciseInt)
				if(err != nil){
					utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
				}
			}
		}
	}
	return nil
}