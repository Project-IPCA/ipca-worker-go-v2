package service

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/Project-IPCA/ipca-worker-go-v2/redis_client"
	"github.com/Project-IPCA/ipca-worker-go-v2/repositories"
	"github.com/Project-IPCA/ipca-worker-go-v2/utils"
	"gorm.io/gorm"
)

type SubmissionResult struct {
	SubmissionID string
	Status       string
	Marking      string
	Result       string
	ErrorMessage string
}

type TestCaseResult struct {
	TestCaseNo    int    `json:"testcase_no"`
	IsPassed      bool   `json:"is_passed"`
	ShowToStudent bool   `json:"show_to_student"`
	Expected      string `json:"expected"`
	Actual        string `json:"actual"`
}

func RunSubmission(channel *amqp.Channel, db_pool *gorm.DB, msg amqp.Delivery, msgBody models.ReciveMessage, redisClient *redis.Client) error{
	publisher := redis_client.NewRedisAction(redisClient)
	activityLogRepo := repositories.NewActivityLogRePository(db_pool)
	excerciseSubmissionRepo := repositories.NewExerciseSubmissionRePository(db_pool)
	var publishLog *models.ActivityLog

	tempLog, err := compileCode(db_pool, msgBody)
	publishLog = tempLog
	if err != nil {
		appErr, ok := err.(*utils.AppError)
		if ok && (appErr.Name == utils.ERROR_NAME.DATABASE_ERROR || appErr.Name == utils.ERROR_NAME.FUNCTION_ERROR) {
			channel.Nack(msg.DeliveryTag, false, false)
		} else {
			newAction := &msgBody.LogData.Actoin

			var output interface{}

			if appErr.Stdout != ""{
				output = appErr.Stdout
			} else {
				output = []string{}
			}

			resultJson, err := json.Marshal(output)
			if err != nil {
				fmt.Println("Error marshalling testcaseResult:", err)
				return err
			}
			outputStr := string(resultJson)
			errorMessage := string(appErr.Err.Error())

			marking := 0

			submissionUuid, err := uuid.Parse(*msgBody.SubmissionID)
			if err != nil {
				fmt.Println("fail to convert uuid")
				return err
			}

			submission := models.UpdateSubmissionInfo{
				SubmissionID: submissionUuid,
				Status:       utils.ExerciseStatus.Error,
				Marking:      marking,
				Result:       &outputStr,
				ErrorMessage: &errorMessage,
			}

			saveLog := &msgBody.LogData
			newAction.Status = utils.ExerciseStatus.Error
			newAction.Marking = &marking

			err = excerciseSubmissionRepo.UpdateSubmission(&submission)
			if err != nil {
				channel.Nack(msg.DeliveryTag, false, false)
				fmt.Println("Error updating submission:", err)
				return err
			}

			tempLog, err := activityLogRepo.AddSubmissionLog(saveLog)
			publishLog = tempLog
			if err != nil {
				channel.Nack(msg.DeliveryTag, false, false)
				fmt.Println("Error adding submission log:", err)
				return err
			}
		}
	}

	err = publisher.PublishMessage(fmt.Sprintf("submission-result:%s", msgBody.JobID), "done")
	if err != nil {
		fmt.Println("Error publishing to Redis:", err)
		return err
	}
	if publishLog != nil {
		err = publisher.PublishMessage(fmt.Sprintf("logs:%s", msgBody.LogData.GroupID), publishLog)
		if err != nil {
			fmt.Println("Error publishing log to Redis:", err)
		}
	}

	fmt.Println("FINISHED RUNNING")
	channel.Ack(msg.DeliveryTag, false)
	return nil
}

func compileCode(db_pool *gorm.DB, msgBody models.ReciveMessage) (*models.ActivityLog, error) {
	submissionUuid, err := uuid.Parse(*msgBody.SubmissionID)
	if err != nil {
		fmt.Println("fail to convert uuid")
		return nil, utils.NewAppError(utils.ERROR_NAME.FUNCTION_ERROR, "failed to convert uuid", err.Error())
	}
	activityLogRepo := repositories.NewActivityLogRePository(db_pool)
	excerciseSubmissionRepo := repositories.NewExerciseSubmissionRePository(db_pool)

	testcaseResult := []TestCaseResult{}
	newAction := msgBody.LogData.Actoin

	for i, testcase := range msgBody.TestCaseList {
		var result string
		var err error
		if msgBody.Language == utils.LanguageList.Python {
			result, err = utils.RunPythonScript(testcase, msgBody.SourceCode)
		} else if msgBody.Language == utils.LanguageList.C {
			result, err = utils.RunCScript(testcase, msgBody.SourceCode)
		} else {
			return nil, utils.NewAppError("UnsupportedLanguage", "The provided language is not supported", "")
		}
		if err != nil {
			fmt.Println("come error")
			fmt.Println(err)
			appErr, ok := err.(*utils.AppError)
			if ok {
				fmt.Println("Error running Python script:", appErr)
				return nil, utils.NewAppError(appErr.Name, appErr.Error(), appErr.Stdout)
			}
		}
		passed := result == testcase.TestCaseOutput
		fmt.Printf("Testcase %d: %v\n", i+1, passed)

		testcaseResult = append(testcaseResult, TestCaseResult{
			TestCaseNo:    i + 1,
			IsPassed:      passed,
			ShowToStudent: testcase.ShowToStudent,
			Expected:      testcase.TestCaseOutput,
			Actual:        result,
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
		return nil, utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "Error marshalling testcaseResult", err.Error())
	}

	outputStr := string(jsonData)
	status := utils.ExerciseStatus.WrongAnswer
	if studentMarking == 2 {
		studentAssignItemRepo := repositories.NewStudentAssignChapterItemRepository(db_pool)
		studentAssignItemRepo.UpdateStudentAssignItemMarking(msgBody.StudentId, msgBody.ChapterId, msgBody.ItemId, studentMarking)
		status = utils.ExerciseStatus.Accepted
	}

	submission := models.UpdateSubmissionInfo{
		SubmissionID: submissionUuid,
		Status:       status,
		Marking:      studentMarking,
		Result:       &outputStr,
		ErrorMessage: nil,
	}

	newAction.Status = status
	newAction.Marking = &studentMarking
	saveLog := &msgBody.LogData
	saveLog.Actoin = newAction

	err = excerciseSubmissionRepo.UpdateSubmission(&submission)
	if err != nil {
		fmt.Println("Error updating submission:", err)
		return nil, utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "Error marshalling testcaseResult", err.Error())
	}

	insertedLog, err := activityLogRepo.AddSubmissionLog(saveLog)
	if err != nil {
		fmt.Println("Error adding submission log:", err)
		return nil, utils.NewAppError(utils.ERROR_NAME.DATABASE_ERROR, "Error marshalling testcaseResult", err.Error())
	}

	return insertedLog, nil
}
