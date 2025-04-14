// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package version

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := Version()
	is.NotNil(v)
}

func TestGitCommitID(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := GitCommitID()
	is.NotNil(v)
}

func TestSemverVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v, err := SemverVersion()
	is.NoError(err)
	is.NotNil(v)
}

func TestVersionCommand_Defaults(t *testing.T) {
	is := assert.New(t)

	// Set default values for version and gitCommitID
	version = "v0.0.0-unset"
	gitCommitID = ""

	cmd := NewVersionCommand()

	// Capture output
	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	// Execute command
	err := cmd.Execute()
	is.NoError(err, "Expected no error on version command with default values")

	// Verify output
	output := strings.TrimSpace(outBuf.String())
	lines := strings.Split(output, "\n")
	is.Contains(lines[6], "version: v0.0.0-unset", "Expected default version")
	is.Contains(lines[7], "commit:", "Expected empty commit message for default")
}

func TestVersionCommand_CustomValues(t *testing.T) {
	//t.Parallel()
	is := assert.New(t)

	// Set custom values for testing
	version = "v1.0.0-test"
	gitCommitID = "abcdef1234567890"

	cmd := NewVersionCommand()

	// Capture output
	var outBuf bytes.Buffer
	cmd.SetOut(&outBuf)

	// Execute command
	err := cmd.Execute()
	is.NoError(err, "Expected no error on version command with custom values")

	// Verify output
	output := strings.TrimSpace(outBuf.String())
	lines := strings.Split(output, "\n")
	is.Contains(lines[6], "version: v1.0.0-test", "Expected custom version")
	is.Contains(lines[7], "commit:  abcdef1234567890", "Expected custom commit ID")
}
