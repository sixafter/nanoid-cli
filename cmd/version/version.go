// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package version

import (
	"strings"

	"github.com/blang/semver/v4"
)

// Prefix is the prefix of the git tag for a version
const Prefix = "v"

// version is a private field and should be set when compiling with --ldflags="-X github.com/sixafter/nanoid-cli/cmd/version.version=vX.Y.Z"
var version = "v0.0.0-unset"

// gitCommitID is a private field and should be set when compiling with --ldflags="-X github.com/sixafter/nanoid-cli/cmd/version.gitCommitID=<commit-id>"
var gitCommitID = ""

// GetVersion returns the current minikube version
func GetVersion() string {
	return version
}

// GetGitCommitID returns the git commit id from which it is being built
func GetGitCommitID() string {
	return gitCommitID
}

// GetSemverVersion returns the current semantic version (semver)
func GetSemverVersion() (semver.Version, error) {
	return semver.Make(strings.TrimPrefix(GetVersion(), Prefix))
}
