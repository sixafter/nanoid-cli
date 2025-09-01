// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package cmd

import (
	"github.com/sixafter/nanoid-cli/cmd/generate"
	"github.com/sixafter/nanoid-cli/cmd/version"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nanoid",
	Short: "A simple, fast, and concurrent CLI for generating secure, URL-friendly unique string IDs",
	Long:  `NanoID CLI is a simple, fast, and concurrent command-line tool for generating secure, URL-friendly unique string IDs using the NanoID Go implementation.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// Execute runs the RootCmd and returns any errors encountered
// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	RootCmd.AddCommand(generate.NewGenerateCommand())
	RootCmd.AddCommand(version.NewVersionCommand())
	return RootCmd.Execute()
}
