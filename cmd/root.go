// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nanoid",
	Short: "A simple, fast, and concurrent CLI for generating secure, URL-friendly unique string IDs",
	Long:  `NanoID CLI is a simple, fast, and concurrent command-line tool for generating secure, URL-friendly unique string IDs using the NanoID Go implementation.`,
	// Uncomment the following line if your bare application has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing nanoid: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Here you can define persistent flags and configuration settings if needed.
	// Example:
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nanoid-cli.yaml)")
}
