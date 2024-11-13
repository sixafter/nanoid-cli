// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sixafter/nanoid"
	"github.com/spf13/cobra"
)

// Variables to hold flag values
var (
	idLength int
	alphabet string
	count    int
	output   string
	verbose  bool
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate one or more Nano IDs",
	Long: `Generate one or more Nano IDs with customizable length and alphabet.

If --id-length is not specified, a default length of 21 is used.
If --alphabet is not specified, the default ASCII alphabet is used.
If --count is not specified, one Nano ID is generated.`,
	RunE: runGenerate, // Use RunE to handle errors gracefully
}

func init() {
	// Add generateCmd as a subcommand to rootCmd
	RootCmd.AddCommand(generateCmd)

	// Define flags for the generate command
	generateCmd.Flags().IntVarP(&idLength, "id-length", "l", nanoid.DefaultLength, "Length of the Nano ID to generate")
	generateCmd.Flags().StringVarP(&alphabet, "alphabet", "a", nanoid.DefaultAlphabet, "Custom alphabet to use for Nano ID generation")
	generateCmd.Flags().IntVarP(&count, "count", "c", 1, "Number of Nano IDs to generate")
	generateCmd.Flags().StringVarP(&output, "output", "o", "", "Output file to write the generated Nano IDs")
	generateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
}

// runGenerate is the main execution function for the generate command
func runGenerate(cmd *cobra.Command, args []string) error {
	// Validate id-length
	if idLength <= 0 {
		return fmt.Errorf("--id-length must be a positive integer")
	}

	// Validate count
	if count <= 0 {
		return fmt.Errorf("--count must be a positive integer")
	}

	// Configure the Nano ID generator using ConfigOptions
	var configOpts []nanoid.Option
	configOpts = append(configOpts, nanoid.WithLengthHint(uint16(idLength)))

	if alphabet != nanoid.DefaultAlphabet {
		configOpts = append(configOpts, nanoid.WithAlphabet(alphabet))
		if verbose {
			fmt.Println("Custom alphabet provided. Initializing custom generator.")
		}
	}

	// Initialize the Nano ID generator with the configured options
	generator, err := nanoid.NewGenerator(configOpts...)
	if err != nil {
		return fmt.Errorf("failed to initialize Nano ID generator: %w", err)
	}

	// Determine the output destination
	var outputDest *os.File
	if output != "" {
		outputDest, err = os.Create(output)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func(outputDest *os.File) {
			_ = outputDest.Close()
		}(outputDest)
		if verbose {
			fmt.Printf("Output will be written to file: %s\n", output)
		}
	} else {
		outputDest = os.Stdout
		if verbose {
			fmt.Println("Output will be printed to stdout.")
		}
	}

	// Use a buffered writer for efficient writing
	writer := bufio.NewWriter(outputDest)
	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error flushing writer: %v\n", err)
		}
	}(writer)

	// Generate and write the specified number of Nano IDs
	for i := 0; i < count; i++ {
		id, err := generator.New(idLength)
		if err != nil {
			return fmt.Errorf("error generating Nano ID: %w", err)
		}

		_, err = writer.WriteString(id + "\n")
		if err != nil {
			return fmt.Errorf("error writing Nano ID: %w", err)
		}

		if verbose {
			fmt.Printf("Generated ID %d: %s\n", i+1, id)
		}
	}

	return nil
}
