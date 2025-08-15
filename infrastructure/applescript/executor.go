package applescript

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/madstone-tech/maestro/domain/music"
)

// ExecutorConfig holds configuration for the AppleScript executor.
type ExecutorConfig struct {
	// ExecPath is the path to the maestro-exec binary
	ExecPath string

	// DefaultTimeout is the default timeout for script execution
	DefaultTimeout time.Duration

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// RetryDelay is the delay between retry attempts
	RetryDelay time.Duration
}

// DefaultExecutorConfig returns a default configuration for the executor.
func DefaultExecutorConfig() *ExecutorConfig {
	return &ExecutorConfig{
		ExecPath:       "maestro-exec", // Will look in PATH or can be overridden
		DefaultTimeout: 10 * time.Second,
		MaxRetries:     3,
		RetryDelay:     500 * time.Millisecond,
	}
}

// Executor provides a wrapper around maestro-exec for executing AppleScript commands.
// It handles timeouts, retries, and error processing.
type Executor struct {
	config *ExecutorConfig
}

// NewExecutor creates a new AppleScript executor with the provided configuration.
func NewExecutor(config *ExecutorConfig) *Executor {
	if config == nil {
		config = DefaultExecutorConfig()
	}

	return &Executor{
		config: config,
	}
}

// ExecuteResult contains the result of AppleScript execution.
type ExecuteResult struct {
	// Output is the stdout from the script execution
	Output string

	// Error is any error that occurred during execution
	Error error

	// Duration is how long the execution took
	Duration time.Duration

	// RetryCount is how many retry attempts were made
	RetryCount int
}

// Execute runs an AppleScript string with the default timeout and retry logic.
func (e *Executor) Execute(ctx context.Context, script string) *ExecuteResult {
	return e.ExecuteWithTimeout(ctx, script, e.config.DefaultTimeout)
}

// ExecuteWithTimeout runs an AppleScript string with a specific timeout and retry logic.
func (e *Executor) ExecuteWithTimeout(ctx context.Context, script string, timeout time.Duration) *ExecuteResult {
	script = strings.TrimSpace(script)
	if script == "" {
		return &ExecuteResult{
			Error: music.NewDomainError(music.ErrInvalidOperation, "AppleScript cannot be empty"),
		}
	}

	var lastErr error
	startTime := time.Now()

	for attempt := 0; attempt <= e.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retrying
			select {
			case <-ctx.Done():
				return &ExecuteResult{
					Error:      music.NewDomainErrorWithCause(music.ErrTimeout, "context cancelled during retry", ctx.Err()),
					Duration:   time.Since(startTime),
					RetryCount: attempt,
				}
			case <-time.After(e.config.RetryDelay):
				// Continue to retry
			}
		}

		result := e.executeOnce(ctx, script, timeout)
		result.RetryCount = attempt

		// If successful, return immediately
		if result.Error == nil {
			result.Duration = time.Since(startTime)
			return result
		}

		lastErr = result.Error

		// Check if this is a permanent error that shouldn't be retried
		if music.IsPermanent(lastErr) {
			result.Duration = time.Since(startTime)
			return result
		}

		// Check if context is cancelled
		if ctx.Err() != nil {
			result.Error = music.NewDomainErrorWithCause(music.ErrTimeout, "context cancelled", ctx.Err())
			result.Duration = time.Since(startTime)
			return result
		}
	}

	// All retries exhausted
	return &ExecuteResult{
		Error:      music.NewDomainErrorWithCause(music.ErrOperationFailed, fmt.Sprintf("AppleScript execution failed after %d attempts", e.config.MaxRetries+1), lastErr),
		Duration:   time.Since(startTime),
		RetryCount: e.config.MaxRetries,
	}
}

// executeOnce executes the AppleScript once without retries.
func (e *Executor) executeOnce(ctx context.Context, script string, timeout time.Duration) *ExecuteResult {
	startTime := time.Now()

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Create command to run maestro-exec
	cmd := exec.CommandContext(execCtx, e.config.ExecPath)

	// Set up pipes
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Create stdin pipe and write script
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return &ExecuteResult{
			Error:    music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to create stdin pipe", err),
			Duration: time.Since(startTime),
		}
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		stdin.Close()
		return &ExecuteResult{
			Error:    music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to start maestro-exec", err),
			Duration: time.Since(startTime),
		}
	}

	// Write script to stdin and close
	_, writeErr := stdin.Write([]byte(script))
	stdin.Close()

	if writeErr != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return &ExecuteResult{
			Error:    music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to write script to stdin", writeErr),
			Duration: time.Since(startTime),
		}
	}

	// Wait for command to complete
	err = cmd.Wait()

	result := &ExecuteResult{
		Output:   strings.TrimSpace(stdout.String()),
		Duration: time.Since(startTime),
	}

	// Handle different types of errors
	if err != nil {
		stderrOutput := strings.TrimSpace(stderr.String())

		// Check if it was a timeout
		if execCtx.Err() == context.DeadlineExceeded {
			result.Error = music.WrapOperationTimeout("AppleScript execution", music.NewDurationFromTime(timeout), err)
			return result
		}

		// Check if it was a context cancellation
		if ctx.Err() != nil {
			result.Error = music.NewDomainErrorWithCause(music.ErrTimeout, "context cancelled", ctx.Err())
			return result
		}

		// Check exit code for specific error types
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitError.ExitCode() {
			case 1: // General error
				result.Error = music.NewDomainErrorWithCause(music.ErrOperationFailed, fmt.Sprintf("AppleScript execution failed: %s", stderrOutput), err)
			case 2: // Timeout
				result.Error = music.WrapOperationTimeout("AppleScript execution", music.NewDurationFromTime(timeout), err)
			default:
				result.Error = music.NewDomainErrorWithCause(music.ErrOperationFailed, fmt.Sprintf("AppleScript execution failed with exit code %d: %s", exitError.ExitCode(), stderrOutput), err)
			}
		} else {
			result.Error = music.NewDomainErrorWithCause(music.ErrOperationFailed, fmt.Sprintf("failed to execute AppleScript: %s", stderrOutput), err)
		}
	}

	return result
}

// LoadScript loads an AppleScript template from the scripts directory.
func (e *Executor) LoadScript(scriptName string) (string, error) {
	// Get the directory of the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", music.NewDomainError(music.ErrOperationFailed, "unable to determine current file location")
	}

	// Construct path to scripts directory
	scriptsDir := filepath.Join(filepath.Dir(currentFile), "scripts")
	scriptPath := filepath.Join(scriptsDir, scriptName+".scpt")

	// Read the script file
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return "", music.NewDomainErrorWithCause(
			music.ErrOperationFailed,
			fmt.Sprintf("failed to read script file %s", scriptPath),
			err,
		)
	}

	return string(content), nil
}

// ExecuteTemplate loads and executes an AppleScript template with the given parameters.
func (e *Executor) ExecuteTemplate(ctx context.Context, templateName string, params map[string]interface{}) *ExecuteResult {
	script, err := e.LoadScript(templateName)
	if err != nil {
		return &ExecuteResult{
			Error: music.NewDomainErrorWithCause(music.ErrOperationFailed, "failed to load script template", err),
		}
	}

	// Replace template parameters
	for key, value := range params {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		script = strings.ReplaceAll(script, placeholder, replacement)
	}

	return e.Execute(ctx, script)
}

// IsExecutable checks if the maestro-exec binary is available and executable.
func (e *Executor) IsExecutable() error {
	// Try to find the executable
	execPath, err := exec.LookPath(e.config.ExecPath)
	if err != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "maestro-exec not found in PATH", err)
	}

	// Update config with full path
	e.config.ExecPath = execPath

	// Test execution with a simple script
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := e.executeOnce(ctx, "return \"test\"", 2*time.Second)
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "maestro-exec is not working properly", result.Error)
	}

	if result.Output != "test" {
		return music.NewDomainError(music.ErrOperationFailed, "maestro-exec test returned unexpected output")
	}

	return nil
}

// HealthCheck performs a basic health check to ensure the executor is working.
func (e *Executor) HealthCheck(ctx context.Context) error {
	result := e.Execute(ctx, "return \"health\"")
	if result.Error != nil {
		return music.NewDomainErrorWithCause(music.ErrOperationFailed, "executor health check failed", result.Error)
	}

	if result.Output != "health" {
		return music.NewDomainError(music.ErrOperationFailed, "executor health check returned unexpected output")
	}

	return nil
}
