package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	exitSuccess = 0
	exitError   = 1
	exitTimeout = 2
	timeout     = 10 * time.Second
)

func main() {
	// Read AppleScript from stdin
	scanner := bufio.NewScanner(os.Stdin)
	var script strings.Builder

	for scanner.Scan() {
		script.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		os.Exit(exitError)
	}

	scriptContent := strings.TrimSpace(script.String())
	if scriptContent == "" {
		fmt.Fprintf(os.Stderr, "Error: empty AppleScript input\n")
		os.Exit(exitError)
	}

	// Execute AppleScript with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "osascript", "-e", scriptContent)
	output, err := cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintf(os.Stderr, "Error: script execution timed out after %v\n", timeout)
		os.Exit(exitTimeout)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing AppleScript: %v\n", err)
		os.Exit(exitError)
	}

	// Write result to stdout
	fmt.Print(string(output))
	os.Exit(exitSuccess)
}
