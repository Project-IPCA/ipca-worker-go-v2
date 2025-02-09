package utils

import (
	"errors"
	"fmt"
)

var ERROR_NAME = struct {
	RUNTIME_ERROR         string
	TIMEOUT_ERROR         string
	OUTPUT_LIMIT_EXCEEDED string
	MEMORY_LIMIT_EXCEEDED string
	DATABASE_ERROR        string
	RABBITMQ_ERROR        string
	SERVER_ERROR          string
	NO_OUTPUT_PRODUCED    string
	PUBLISH_ERROR         string
	FUNCTION_ERROR        string
}{
	RUNTIME_ERROR:         "RuntimeError",
	TIMEOUT_ERROR:         "TimeoutError",
	OUTPUT_LIMIT_EXCEEDED: "OutputLimitExceeded",
	MEMORY_LIMIT_EXCEEDED: "MemoryLimitExceeded",
	DATABASE_ERROR:        "DatabaseError",
	RABBITMQ_ERROR:        "RabbitMQError",
	SERVER_ERROR:          "ServerError",
	NO_OUTPUT_PRODUCED:    "NoOutputProduced",
	PUBLISH_ERROR:         "PublishError",
	FUNCTION_ERROR:        "FunctionError",
}

type AppError struct {
	Name   string
	Stdout string
	Err    error
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s\nstdout: %s", e.Name, e.Err.Error(), e.Stdout)
}

func NewAppError(name, message, stdout string) *AppError {
	return &AppError{
		Name:   name,
		Stdout: stdout,
		Err:    errors.New(message),
	}
}
