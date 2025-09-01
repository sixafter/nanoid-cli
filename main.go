// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"fmt"
	"os"

	"github.com/sixafter/nanoid-cli/cmd"
)

// run is the main entry point for the CLI, allowing it to be tested without os.Exit.
func run() error {
	return cmd.Execute()
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
