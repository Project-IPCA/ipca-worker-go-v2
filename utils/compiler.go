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
	sandboxRunnerPath := filepath.Join(sandboxPath, "box", "runner.py")
	err = MoveFile(runnerPath, sandboxRunnerPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy runner.py to sandbox: %v", err)
	}
	sourceCodeFileName := fmt.Sprintf("%s.py", uuid.New().String())
	sandboxSourceCodePath := filepath.Join(sandboxPath, "box", sourceCodeFileName)
	err = os.WriteFile(sandboxSourceCodePath, []byte(strings.TrimSpace(sourceCode)), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
	}
	testcaseFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxTestcasePath := filepath.Join(sandboxPath, "box", testcaseFileName)
	err = os.WriteFile(sandboxTestcasePath, []byte(strings.TrimSpace(testcase.TestCaseContent)), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
	}

	command := fmt.Sprintf("/usr/bin/python3.12 runner.py %s %s", sourceCodeFileName, testcaseFileName)
	fmt.Println(command)
	return ExecuteCommandWithIsolate(sandboxPath, command)
}

func RunCScript(testcase models.TestCase, sourceCode string) (string, error) {
	sandboxPath, err := InitializeIsolate()
	if err != nil {
		return "", fmt.Errorf("failed to initialize isolate: %v", err)
	}
	defer CleanupIsolate()

	runnerContent, err := os.ReadFile("./c_file/runner.c")
	if err != nil {
		fmt.Printf("Error reading runner.c: %v\n", err)
		return "", fmt.Errorf("failed to initialize isolate: %v", err)
	}

	updatedSourceCode := strings.Replace(string(sourceCode), "int main", "int user_main", 1)

	combinedCode := string(runnerContent) + "\n" + updatedSourceCode

	combinedCodeFileName := fmt.Sprintf("%s.c", uuid.New().String())
	sandboxcombinedCodePath := filepath.Join(sandboxPath, "box", combinedCodeFileName)
	err = os.WriteFile(sandboxcombinedCodePath, []byte(strings.TrimSpace(combinedCode)), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create file in: %v", err)
	}

	buildCommand := fmt.Sprintf("gcc %s -o program", combinedCodeFileName)
	_, err = ExecuteCommand(filepath.Join(sandboxPath, "box"), buildCommand)
	if err != nil {
		os.Remove(combinedCodeFileName)
		return "", err
	}
	defer os.Remove(combinedCodeFileName)

	testcaseFileName := fmt.Sprintf("%s.txt", uuid.New().String())
	sandboxTestcasePath := filepath.Join(sandboxPath, "box", testcaseFileName)
	err = os.WriteFile(sandboxTestcasePath, []byte(strings.TrimSpace(testcase.TestCaseContent)), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create testcase file in sandbox: %v", err)
	}

	runCommand := fmt.Sprintf("./program %s", testcaseFileName)
	return ExecuteCommandWithIsolate(sandboxPath, runCommand)
}
