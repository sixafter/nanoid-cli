// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package version

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
)

const logo string = "  _   _                     ___ ____  \n | \\ | | __ _ _ __   ___   |_ _|  _ \\ \n |  \\| |/ _` | '_ \\ / _ \\   | || | | |\n | |\\  | (_| | | | | (_) |  | || |_| |\n |_| \\_|\\__,_|_| |_|\\___/  |___|____/ \n                                      \n"

// NewVersionCommand creates and returns the version command
func NewVersionCommand() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Display the version of NanoID CLI",
		Long:  `Display the current version of NanoID CLI.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Use a buffered writer for efficient writing
			writer := bufio.NewWriter(cmd.OutOrStdout())

			// Write the logo to the output
			_, err := fmt.Fprint(writer, logo)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing logo: %v\n", err)
				return
			}

			_, err = fmt.Fprintf(writer, "version: %s\n", Version())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing version: %v\n", err)
				return
			}

			_, err = fmt.Fprintf(writer, "commit:  %s\n", GitCommitID())
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error writing commit: %v\n", err)
				return
			}

			defer func(writer *bufio.Writer) {
				err = writer.Flush()
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Error flushing writer: %v\n", err)
				}
			}(writer)
		},
	}

	return versionCmd
}

// Prefix is the prefix of the git tag for a version
const Prefix = "v"

// version is a private field and should be set when compiling with --ldflags="-X github.com/sixafter/nanoid-cli/cmd/version.version=vX.Y.Z"
var version = "v0.0.0-unset"

// gitCommitID is a private field and should be set when compiling with --ldflags="-X github.com/sixafter/nanoid-cli/cmd/version.gitCommitID=<commit-id>"
var gitCommitID = ""

// Version returns the current minikube version
func Version() string {
	return version
}

// GitCommitID returns the git commit id from which it is being built
func GitCommitID() string {
	return gitCommitID
}

// SemverVersion returns the current semantic version (semver)
func SemverVersion() (semver.Version, error) {
	return semver.Make(strings.TrimPrefix(Version(), Prefix))
}
