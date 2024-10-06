package service

import (
	"fmt"
	"strings"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/redis_client"

	"github.com/Project-IPCA/ipca-worker-go-v2/repositories"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"github.com/google/uuid"
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
	err = publisher.PublishMessage(fmt.Sprintf("testcase-result:%s", msgBody.JobID),"done")
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		return
	}
	fmt.Println("FINISHED RUNNING")
	channel.Ack(msg.DeliveryTag, false)
}

func compileCodeTestcase (db_pool *gorm.DB, msgBody models.ReciveMessage) error{
	exerciseTestcaseRepo := repositories.NewExcerciseTestCaseRePository(db_pool)
	labExerciseRepo := repositories.NewLabExerciseRePository(db_pool)

	exerciseUuid,err := uuid.Parse(*msgBody.ExcerciseID)
	if(err!=nil){
		return utils.NewAppError(utils.ERROR_NAME.FUNCTION_ERROR,"failed to convert exercise uuid", err.Error())
	}
	var labExercise models.LabExercise
	labExerciseRepo.GetLabExerciseById(&labExercise,exerciseUuid)
	if(len(msgBody.TestCaseList)>0){
		for i:=0; i<len(msgBody.TestCaseList); i++ {
			testcaseUuid,err := uuid.Parse(msgBody.TestCaseList[i].TestCaseID)
			if(err!=nil){
				return utils.NewAppError(utils.ERROR_NAME.FUNCTION_ERROR,"failed to convert testcase uuid", err.Error())
			}
			result, err := utils.RunPythonScript(msgBody.TestCaseList[i], msgBody.SourceCode)
			if err != nil {
				appErr, ok := err.(*utils.AppError)
				if(ok){
					fmt.Println("Error running Python script:", appErr)
					errorMessage := string(appErr.Err.Error())
					testcaseData := models.ExerciseTestcase{
						TestcaseID: &testcaseUuid,
						ExerciseID: exerciseUuid,
						TestcaseContent: msgBody.TestCaseList[i].TestCaseContent,
						IsReady: "no",
						IsActive: &msgBody.TestCaseList[i].Active,
						IsShowStudent: &msgBody.TestCaseList[i].ShowToStudent,
						TestcaseNote: &msgBody.TestCaseList[i].TestCaseNote,
						TestcaseOutput: &appErr.Stdout,
						TestcaseError: &errorMessage,
					}
					err = exerciseTestcaseRepo.UpdateTestCase(&testcaseData,exerciseUuid)
					if(err != nil){
						utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
					}
					for j:=i+1; j<len(msgBody.TestCaseList); j++{
						errorText := "Error occurred before this testcase"
						testCaseData := models.ExerciseTestcase{
							TestcaseID: &testcaseUuid,
							ExerciseID: exerciseUuid,
							TestcaseContent: msgBody.TestCaseList[i].TestCaseContent,
							IsReady: "no",
							IsActive: &msgBody.TestCaseList[i].Active,
							IsShowStudent: &msgBody.TestCaseList[i].ShowToStudent,
							TestcaseNote: &msgBody.TestCaseList[i].TestCaseNote,
							TestcaseOutput: nil,
							TestcaseError: &errorText,
						}
						err = exerciseTestcaseRepo.UpdateTestCase(&testCaseData,exerciseUuid)
						if(err != nil){
							utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to write file", err.Error())
						}
					}
					return utils.NewAppError(appErr.Name,appErr.Error(), appErr.Stdout)
				}
			}else{
				output := strings.TrimSpace(result)
				fmt.Println("output : " + output)
				testCaseData := models.ExerciseTestcase{
					TestcaseID: &testcaseUuid,
					ExerciseID: exerciseUuid,
					TestcaseContent: msgBody.TestCaseList[i].TestCaseContent,
					IsReady: "yes",
					IsActive: &msgBody.TestCaseList[i].Active,
					IsShowStudent: &msgBody.TestCaseList[i].ShowToStudent,
					TestcaseNote: &msgBody.TestCaseList[i].TestCaseNote,
					TestcaseOutput: &output,
					TestcaseError: nil,
				}

				err = exerciseTestcaseRepo.UpdateTestCase(&testCaseData,exerciseUuid)
				if(err != nil){
					utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to update testcase to db", err.Error())
				}
				labExerciseRepo.UpdateExerciseTestcaseEnum(&labExercise,"YES")
			}
		}
	}
	return nil
}