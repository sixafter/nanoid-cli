// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package generate

import (
	"bufio"
	"fmt"

	"github.com/sixafter/nanoid"
	"github.com/spf13/cobra"
)

var (
	idLength int
	alphabet string
	count    int
	verbose  bool
)

// NewGenerateCommand creates and returns the generate command
func NewGenerateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate one or more Nano IDs",
		Long: `Generate one or more Nano IDs with customizable length and alphabet.

If --id-length is not specified, a default length of 21 is used.
If --alphabet is not specified, the default ASCII alphabet is used.
If --count is not specified, one Nano ID is generated.`,
		RunE: runGenerate, // Use RunE to handle errors gracefully
	}

	// Define flags for the generate command
	cmd.Flags().IntVarP(&idLength, "id-length", "l", nanoid.DefaultLength, "Length of the Nano ID to generate")
	cmd.Flags().StringVarP(&alphabet, "alphabet", "a", nanoid.DefaultAlphabet, "Custom alphabet to use for Nano ID generation")
	cmd.Flags().IntVarP(&count, "count", "c", 1, "Number of Nano IDs to generate")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	return cmd
}

// runGenerate is the main execution function for the generate command
func runGenerate(cmd *cobra.Command, args []string) error {
	// Validate id-length
	if idLength <= 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStderr(), "--id-length must be a positive integer")
		return fmt.Errorf("--id-length must be a positive integer")
	}

	// Validate count
	if count <= 0 {
		_, _ = fmt.Fprintln(cmd.OutOrStderr(), "--count must be a positive integer")
		return fmt.Errorf("--count must be a positive integer")
	}

	// Configure the Nano ID generator using ConfigOptions
	var configOpts []nanoid.Option
	configOpts = append(configOpts, nanoid.WithLengthHint(uint16(idLength)))

	if alphabet != nanoid.DefaultAlphabet {
		configOpts = append(configOpts, nanoid.WithAlphabet(alphabet))
		if verbose {
			_, _ = fmt.Fprintln(cmd.OutOrStderr(), "Custom alphabet provided. Initializing custom generator.")
		}
	}

	// Initialize the Nano ID generator with the configured options
	generator, err := nanoid.NewGenerator(configOpts...)
	if err != nil {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "failed to initialize Nano ID generator: %v\n", err)
		return fmt.Errorf("failed to initialize Nano ID generator: %w", err)
	}

	// Use a buffered writer for efficient writing
	writer := bufio.NewWriter(cmd.OutOrStdout())
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Error flushing writer: %v\n", err)
		}
	}(writer)

	// Generate and write the specified number of Nano IDs
	for i := 0; i < count; i++ {
		id, err := generator.New(idLength)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "error generating Nano ID: %v\n", err)
			return fmt.Errorf("error generating Nano ID: %w", err)
		}

		_, err = writer.WriteString(id.String() + "\n")
		if err != nil {
			return fmt.Errorf("error writing Nano ID: %w", err)
		}

		if verbose {
			_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Generated ID %d: %s\n", i+1, id)
		}
	}

	return nil
}
