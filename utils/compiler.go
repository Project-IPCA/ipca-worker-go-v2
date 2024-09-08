package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Project-IPCA/ipca-worker-go-v2/models"
	"github.com/google/uuid"
)

func RunPythonScript(testcase models.TestCase, sourceCode string) (string, error) {
    sandboxPath, err := InitializeIsolate()
    if err != nil {
        return "", fmt.Errorf("failed to initialize isolate: %v", err)
    }
    defer CleanupIsolate()

    runnerPath := "./python_file/runner.py"
    sandboxRunnerPath := filepath.Join(sandboxPath,"box", "runner.py")
    err = MoveFile(runnerPath, sandboxRunnerPath)
    if err != nil {
        return "", fmt.Errorf("failed to copy runner.py to sandbox: %v", err)
    }
    sourceCodeFileName := fmt.Sprintf("%s.py", uuid.New().String())
    sandboxSourceCodePath := filepath.Join(sandboxPath,"box", sourceCodeFileName)
    err =os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
    if err != nil {
        return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
    }
    testcaseFileName := fmt.Sprintf("%s.txt", uuid.New().String())
    sandboxTestcasePath := filepath.Join(sandboxPath,"box", testcaseFileName)
    err =os.WriteFile(sandboxTestcasePath, []byte(strings.TrimSpace(testcase.TestCaseContent)), 0644)
    if err != nil {
        return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
    }

    command := fmt.Sprintf("/usr/bin/python3.12 runner.py %s %s", sourceCodeFileName, testcaseFileName)
	fmt.Println(command)
    return ExecuteCommandWithIsolate(sandboxPath, command)
}

func RunPythonScriptWithoutTestcase(sourceCode string) (string, error) {
        sandboxPath, err := InitializeIsolate()
        if err != nil {
            return "", fmt.Errorf("failed to initialize isolate: %v", err)
        }
        defer CleanupIsolate()
    
        runnerPath := "./python_file/runner.py"
        sandboxRunnerPath := filepath.Join(sandboxPath,"box", "runner.py")
        err = MoveFile(runnerPath, sandboxRunnerPath)
        if err != nil {
            return "", fmt.Errorf("failed to copy runner.py to sandbox: %v", err)
        }
        sourceCodeFileName := fmt.Sprintf("%s.py", uuid.New().String())
        sandboxSourceCodePath := filepath.Join(sandboxPath,"box", sourceCodeFileName)
        err =os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
        if err != nil {
            return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
        }
	command := fmt.Sprintf("/usr/bin/python3.12 runner.py %s", sourceCodeFileName)
	return ExecuteCommandWithIsolate(sandboxPath,command)
}
