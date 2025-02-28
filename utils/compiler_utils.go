package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	TIME_LIMIT      = 3 * time.Second
	MAX_OUTPUT_SIZE = 1024 * 1024
)

func GetStdOutBeforeError(stdout string) string {
	lines := strings.Split(stdout, "\n")
	if len(lines) > 50 {
		first50Lines := strings.Join(lines[:50], "\n")
		return fmt.Sprintf("%s\n... and %d more lines", first50Lines, len(lines)-50)
	}
	return stdout
}

func ProcessStdout(stdout string) (string, error) {
	if len(stdout) > MAX_OUTPUT_SIZE {
		return "", NewAppError(ERROR_NAME.OUTPUT_LIMIT_EXCEEDED, "Output limit exceeded", stdout)
	}

	if stdout == "" {
		return "", NewAppError(ERROR_NAME.NO_OUTPUT_PRODUCED, "No output produced", stdout)
	}

	return stdout, nil
}

func HandleExecError(err error, stderr string, stdout string) error {
	stdoutBeforeError := GetStdOutBeforeError(stdout)
	if exitError, ok := err.(*exec.ExitError); ok {
		switch exitError.ExitCode() {
		case 137:
			return NewAppError(ERROR_NAME.MEMORY_LIMIT_EXCEEDED, "Memory limit exceeded", stdoutBeforeError)
		default:
			if strings.Contains(err.Error(), "MemoryError") {
				return NewAppError(ERROR_NAME.MEMORY_LIMIT_EXCEEDED, "Memory limit exceeded", stdoutBeforeError)
			}
			return NewAppError(ERROR_NAME.RUNTIME_ERROR, stderr, stdoutBeforeError)
		}
	}

	return NewAppError(ERROR_NAME.RUNTIME_ERROR, err.Error(), stdoutBeforeError)
}

func ExecuteCommandWithIsolate(sandboxPath, command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_LIMIT)
	defer cancel()

	isolateCmd := fmt.Sprintf("isolate --run --time=%d --wall-time=%d --extra-time=1 --mem=128000 -- %s --no-net",
		int(TIME_LIMIT.Seconds()),
		int(TIME_LIMIT.Seconds()),
		command)

	cmd := exec.CommandContext(ctx, "bash", "-c", isolateCmd)
	cmd.Dir = filepath.Join(sandboxPath, "box")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	var stdoutBuf, stderrBuf strings.Builder
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		readOutput(stdout, &stdoutBuf)
	}()

	go func() {
		defer wg.Done()
		readOutput(stderr, &stderrBuf)
	}()

	err = cmd.Start()
	if err != nil {
		return "", HandleExecError(err, stderrBuf.String(), stdoutBuf.String())
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		// Attempt to kill the process
		cmd.Process.Kill()
		wg.Wait() // Wait for output goroutines to finish
		return "", NewAppError(ERROR_NAME.TIMEOUT_ERROR, "Time limit exceeded", "")
	case err := <-done:
		wg.Wait() // Wait for output goroutines to finish
		if err != nil {
			return "", HandleExecError(err, stderrBuf.String(), stdoutBuf.String())
		}
		output, err := ProcessStdout(stdoutBuf.String())
		if err != nil {
			return "", err
		}
		return output, nil
	}
}

func ExecuteCommand(path,command string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), TIME_LIMIT)
	defer cancel()

	cmd := exec.CommandContext(ctx , "bash", "-c", command)
	cmd.Dir = filepath.Join(path)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		return "", HandleExecError(err,err.Error(),string(cmdOutput))
	}
	return string(cmdOutput), nil
}

func readOutput(r io.Reader, buf *strings.Builder) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		buf.WriteString(scanner.Text() + "\n")
	}
}

func MoveFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return NewAppError(ERROR_NAME.FUNCTION_ERROR, "fail to open file", err.Error())
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	fmt.Printf("File copied successfully from %s to %s\n", sourcePath, destinationPath)
	return nil
}

func InitializeIsolate() (string, error) {
	cmd := exec.Command("isolate", "--init")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to initialize isolate: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func CleanupIsolate() error {
	fmt.Println("cleanup")
	cmd := exec.Command("isolate", "--cleanup")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to cleanup isolate: %v", err)
	}
	return nil
}
