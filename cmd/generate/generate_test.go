// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package generate

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCommand_Default(t *testing.T) {
	is := assert.New(t)

	// Set up command
	cmd := NewGenerateCommand()
	cmd.SetArgs([]string{"--count", "2"})

	// Capture output
	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	// Execute command
	err := cmd.Execute()
	is.NoError(err, "Expected no error on generate command with default options")

	// Verify output contains two IDs
	output := strings.TrimSpace(outBuf.String())
	ids := strings.Split(output, "\n")
	is.Len(ids, 2, "Expected two IDs in the output")
	is.Equal(21, len(ids[0]), "Expected ID of default length 21")
	is.Equal(21, len(ids[1]), "Expected ID of default length 21")
}

func TestGenerateCommand_CustomLength(t *testing.T) {
	is := assert.New(t)

	// Set up command
	cmd := NewGenerateCommand()
	cmd.SetArgs([]string{"--id-length", "30", "--count", "1"})

	// Capture output
	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	// Execute command
	err := cmd.Execute()
	is.NoError(err, "Expected no error on generate command with custom length")

	// Verify output contains one ID of length 30
	output := strings.TrimSpace(outBuf.String())
	is.Equal(30, len(output), "Expected ID of custom length 30")
}

func TestGenerateCommand_CustomAlphabet(t *testing.T) {
	is := assert.New(t)

	customAlphabet := "abcdef123456"
	cmd := NewGenerateCommand()
	cmd.SetArgs([]string{"--alphabet", customAlphabet, "--count", "3"})

	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	err := cmd.Execute()
	is.NoError(err, "Expected no error on generate command with custom alphabet")

	// Verify output contains three IDs using custom alphabet
	output := strings.TrimSpace(outBuf.String())
	ids := strings.Split(output, "\n")
	is.Len(ids, 3, "Expected three IDs in the output")
	for _, id := range ids {
		for _, char := range id {
			is.Contains(customAlphabet, string(char), "Expected characters in ID to match custom alphabet")
		}
	}
}

func TestGenerateCommand_Verbose(t *testing.T) {
	is := assert.New(t)

	// Set up command
	cmd := NewGenerateCommand()
	cmd.SetArgs([]string{"--id-length", "30", "--count", "10", "--verbose"})

	// Capture output
	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	// Execute command
	err := cmd.Execute()
	is.NoError(err, "Expected no error on generate command with custom length")

	// Verify output contains one ID of length 30
	output := strings.TrimSpace(outBuf.String())

	// Split on \n to get individual lines
	lines := strings.Split(output, "\n")

	is.Equal(19, len(lines), "Expected output to contain 19 lines (10 IDs of length 30 + 10 verbose messages)")
}

func TestGenerateCommand_ErrorOutput(t *testing.T) {
	is := assert.New(t)

	// Set up the command with invalid arguments
	cmd := NewGenerateCommand()
	cmd.SetArgs([]string{"--id-length", "-1"}) // Invalid id-length to trigger error

	// Buffers for stdout and stderr
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.SetOut(&outBuf) // Capture standard output
	cmd.SetErr(&errBuf) // Capture standard error

	// Execute the command
	err := cmd.Execute()
	is.Error(err, "Expected an error on invalid argument")

	// Verify stderr captured the error message
	stderrOutput := strings.TrimSpace(errBuf.String())
	is.Contains(stderrOutput, "--id-length must be a positive integer", "Expected specific error message in stderr")

	// stdout should be empty since an error occurred
	stdoutOutput := strings.TrimSpace(outBuf.String())
	is.NotEmpty(stdoutOutput, "Expected output showing usage.")
}

func TestGenerateCommand_WriteError(t *testing.T) {
	is := assert.New(t)
	var stdoutBuf, rawStderrBuf bytes.Buffer
	stderr := bufio.NewWriter(&rawStderrBuf)

	cmd := &cobra.Command{}
	cmd.SetOut(bufio.NewWriter(&stdoutBuf))
	cmd.SetErr(stderr)

	errMsg := "test error"
	origErr := errors.New("underlying failure")

	returnedErr := writeError(cmd, errMsg, origErr)

	_ = stderr.Flush()

	expectedOutput := fmt.Sprintf("%s: %v", errMsg, origErr)
	is.Contains(rawStderrBuf.String(), expectedOutput, "stderr should contain the error message")
	is.ErrorContains(returnedErr, errMsg)
	is.ErrorIs(returnedErr, origErr)
}

func TestGenerateCommand_WriteString(t *testing.T) {
	is := assert.New(t)
	var stdoutBuf, rawStderrBuf bytes.Buffer
	stderr := bufio.NewWriter(&rawStderrBuf)

	cmd := &cobra.Command{}
	cmd.SetOut(bufio.NewWriter(&stdoutBuf))
	cmd.SetErr(stderr)

	errMsg := "test error"

	returnedErr := writeString(cmd, errMsg)

	_ = stderr.Flush()

	expectedOutput := fmt.Sprintf("%s", errMsg)
	is.Contains(rawStderrBuf.String(), expectedOutput, "stderr should contain the error message")
	is.ErrorContains(returnedErr, errMsg)
}
