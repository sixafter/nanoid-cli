// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package cmd

import (
	"fmt"
	"github.com/sixafter/nanoid-cli/cmd/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of NanoID CLI",
	Long:  `Display the current version of NanoID CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version: %s\n", version.GetVersion())
		fmt.Printf("commit: %s\n", version.GetGitCommitID())
	},
}

func init() {
	// Add versionCmd as a subcommand to rootCmd
	RootCmd.AddCommand(versionCmd)
}
