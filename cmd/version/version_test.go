// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := GetVersion()
	is.NotNil(v)
}

func TestGetGitCommitID(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v := GetGitCommitID()
	is.NotNil(v)
}

func TestGetSemverVersion(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	v, err := GetSemverVersion()
	is.NoError(err)
	is.NotNil(v)
}
