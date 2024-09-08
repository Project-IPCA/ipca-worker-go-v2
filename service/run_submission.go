package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/redis_client"
	"github.com/Project-IPCA/ipca-worker-go-v2/repositories"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"gorm.io/gorm"
)

type SubmissionResult struct{
	SubmissionID string
    Status   string
    Marking string  
	Result  string
    ErrorMessage string
}

type TestCaseResult struct{
	TestCaseNo 	int `json:"testcase_no"`
	IsPassed 	bool `json:"is_passed"`
	ShowToStudent bool `json:"show_to_student"`
	Expected string `json:"expected"`
    Actual string `json:"actual"`
}

func RunSubmission(channel *amqp.Channel, db_pool *gorm.DB, msg amqp.Delivery, msgBody models.ReciveMessage, redisClient *redis.Client){
	publisher := redis_client.NewRedisAction(redisClient)
	activityLogRepo := repositories.NewActivityLogRePository(db_pool)
	excerciseSubmissionRepo := repositories.NewExcerciseSubmissionRePository(db_pool)
	var publishLog *models.ActivityLogOld

	tempLog,err := compileCode(db_pool,msgBody);
	publishLog = tempLog
	if(err!=nil){
		appErr, ok := err.(*utils.AppError)
		if(ok && (appErr.Name == utils.ERROR_NAME.DATABASE_ERROR || appErr.Name == utils.ERROR_NAME.FUNCTION_ERROR)){
			channel.Nack(msg.DeliveryTag,false,false)
		}else{
			submission_int,err := strconv.Atoi(*msgBody.SubmissionID)
			newAction := msgBody.LogData.Actoin

			if(err!=nil){
				fmt.Println("fail to convert to int")
				return 
			}

			var output interface{}

			if appErr.Stdout != "" {
				output = appErr.Stdout
			} else {
				output = []string{}
			}

			resultJson, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling testcaseResult:", err)
				return 
			}
			errorMessage := string(appErr.Err.Error())

			marking := 0

			rawResult := json.RawMessage(resultJson)

			submission := models.UpdateSubmissionInfo{
				SubmissionID: submission_int,
				Status: "error",
				Marking: marking,
				Result: &rawResult,
				ErrorMessage: &errorMessage,
			}

			saveLog := &msgBody.LogData
			newAction.Status = "error"
			newAction.Marking = &marking

			err = excerciseSubmissionRepo.UpdateSubmission(&submission)
		if err != nil {
			channel.Nack(msg.DeliveryTag,false,false)
			fmt.Println("Error updating submission:", err)
			return 
		}

		tempLog, err := activityLogRepo.AddSubmissionLog(saveLog)
		publishLog = tempLog
		if err != nil {
			channel.Nack(msg.DeliveryTag,false,false)
			fmt.Println("Error adding submission log:", err)
			return 
		}
		}
	}

	err = publisher.PublishMessage(fmt.Sprintf("submission-result:%s", msgBody.JobID),"done")
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		return
	}
	if publishLog != nil {
		logJSON, _ := json.Marshal(publishLog)
		err = publisher.PublishMessage(fmt.Sprintf("logs:%s", msgBody.LogData.GroupID), string(logJSON))
		if err != nil {
			fmt.Println("Error publishing log to Redis:", err)
		}
	}

	fmt.Println("FINISHED RUNNING")
	channel.Ack(msg.DeliveryTag, false)
}

func compileCode (db_pool *gorm.DB, msgBody models.ReciveMessage) (*models.ActivityLogOld,error){
	submission_int,err := strconv.Atoi(*msgBody.SubmissionID)
	activityLogRepo := repositories.NewActivityLogRePository(db_pool)
	excerciseSubmissionRepo := repositories.NewExcerciseSubmissionRePository(db_pool)

	if(err!=nil){
		fmt.Println("fail to convert to int")
		return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"failed to convert", err.Error())
	}

	testcaseResult := []TestCaseResult{}
	newAction := msgBody.LogData.Actoin
	insertedLog := models.ActivityLogOld{}

	if(len(msgBody.TestCaseList)>0){
		for i, testcase := range msgBody.TestCaseList {
			result, err := utils.RunPythonScript(testcase, msgBody.SourceCode)
			if err != nil {
				appErr, ok := err.(*utils.AppError)
				if(ok){
					fmt.Println("Error running Python script:", appErr)
					return nil,utils.NewAppError(appErr.Name,appErr.Error(), appErr.Stdout)
				}
			}
			passed := strings.TrimSpace(result) == strings.TrimSpace(testcase.TestCaseOutput)
			fmt.Printf("Testcase %d: %v\n", i+1, passed)
	
			testcaseResult = append(testcaseResult,TestCaseResult{
				TestCaseNo: i+1,
				IsPassed: passed,
				ShowToStudent: testcase.ShowToStudent == "yes",
				Expected: strings.TrimSpace(testcase.TestCaseOutput),
				Actual: strings.TrimSpace(result),
			})
		}

		isPassedAllTestcase := true
		for _, testcase := range testcaseResult {
		if !testcase.IsPassed {
			isPassedAllTestcase = false
			break
		}
		}

		studentMarking := 0
		if isPassedAllTestcase {
			studentMarking = 2
		}

		jsonData, err := json.Marshal(testcaseResult)
		if err != nil {
			fmt.Println("Error marshalling testcaseResult:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}

		rawMessage := json.RawMessage(jsonData)
		status := "wrong_answer";
		if(studentMarking == 2){
			status = "accepted";
		}

		submission := models.UpdateSubmissionInfo{
			SubmissionID: submission_int,
			Status: status,
			Marking: studentMarking,
			Result: &rawMessage,
			ErrorMessage: nil,
		}

		newAction.Status = status
		newAction.Marking = &studentMarking
		saveLog := &msgBody.LogData
		saveLog.Actoin = newAction

		err = excerciseSubmissionRepo.UpdateSubmission(&submission)
		if err != nil {
			fmt.Println("Error updating submission:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}

		tempLog, err := activityLogRepo.AddSubmissionLog(saveLog)
		if err != nil {
			fmt.Println("Error adding submission log:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}
		insertedLog = *tempLog
	}else{
		result, err := utils.RunPythonScriptWithoutTestcase(msgBody.SourceCode)
		if err != nil {
			fmt.Println("Error running Python script:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error running Python script", err.Error())
		}
		fmt.Println("Output : ",strings.TrimSpace(result))

		jsonData, err := json.Marshal(strings.TrimSpace(result))
		if err != nil {
			fmt.Println("Error marshalling testcaseResult:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}

		rawMessage := json.RawMessage(jsonData)

		studentMarking :=2

		submission := models.UpdateSubmissionInfo{
			SubmissionID: submission_int,
			Status: "accepted",
			Marking: studentMarking,
			Result: &rawMessage,
			ErrorMessage: nil,
		}

		newAction.Status = "accepted"
		newAction.Marking = &studentMarking
		saveLog := &msgBody.LogData
		saveLog.Actoin = newAction

		err = excerciseSubmissionRepo.UpdateSubmission(&submission)
		if err != nil {
			fmt.Println("Error updating submission:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}

		tempLog, err := activityLogRepo.AddSubmissionLog(saveLog)
		if err != nil {
			fmt.Println("Error adding submission log:", err)
			return nil,utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR,"Error marshalling testcaseResult", err.Error())
		}
		insertedLog = *tempLog
	}
	return &insertedLog,nil
}