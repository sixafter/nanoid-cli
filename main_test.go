// main_test.go
package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sixafter/nanoid-cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestRun_GenerateCommand(t *testing.T) {
	//t.Parallel()
	is := assert.New(t)

	// Set command-line arguments to simulate "generate" command with count 1
	os.Args = []string{"nanoid", "generate", "--count", "1"}

	// Capture output
	var outBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outBuf)
	cmd.RootCmd.SetErr(&outBuf)

	// Execute Run and check for no errors
	err := run()
	is.NoError(err, "Expected no error on run with generate command")

	// Check if output contains one NanoID of default length
	output := strings.TrimSpace(outBuf.String())
	is.Equal(21, len(output), "Expected single NanoID of default length 21")
}

func TestRun_VersionCommand(t *testing.T) {
	//t.Parallel()
	is := assert.New(t)

	// Set command-line arguments to simulate "version" command
	os.Args = []string{"nanoid", "version"}

	// Capture output
	var outBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outBuf)
	cmd.RootCmd.SetErr(&outBuf)

	// Execute Run and check for no errors
	err := run()
	is.NoError(err, "Expected no error on run with version command")

	// Check if output contains "version" and "commit" information
	output := strings.TrimSpace(outBuf.String())
	is.Contains(output, "version:", "Expected version information in output")
	is.Contains(output, "commit:", "Expected commit information in output")
}

func TestRun_InvalidCommand(t *testing.T) {
	//t.Parallel()
	is := assert.New(t)

	// Set command-line arguments to an invalid command
	os.Args = []string{"nanoid", "invalidcmd"}

	// Capture output
	var outBuf bytes.Buffer
	cmd.RootCmd.SetOut(&outBuf)
	cmd.RootCmd.SetErr(&outBuf)

	// Execute Run and check for an error
	err := run()
	is.Error(err, "Expected an error on run with invalid command")

	// Verify error message
	output := outBuf.String()
	is.Contains(output, "unknown command", "Expected unknown command error")
}
