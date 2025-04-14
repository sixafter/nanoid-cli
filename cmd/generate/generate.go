// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package generate

import (
	"bufio"
	"fmt"
	"math"
	"runtime"
	"time"

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
		return writeString(cmd, "--id-length must be a positive integer")
	}

	// Validate count
	if count <= 0 {
		return writeString(cmd, "--count must be a positive integer")
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
		return writeError(cmd, "failed to initialize Nano ID generator", err)
	}

	// Use a buffered writer for efficient writing
	writer := bufio.NewWriter(cmd.OutOrStdout())

	// Generate and write the specified number of Nano IDs
	start := time.Now()

	// Generate and write the specified number of Nano IDs
	for i := 0; i < count; i++ {
		var id nanoid.ID
		id, err = generator.New(idLength)
		if err != nil {
			return writeError(cmd, "error generating Nano ID", err)
		}

		_, err = writer.WriteString(id.String() + "\n")
		if err != nil {
			return writeError(cmd, "error generating Nano ID", err)
		}
	}

	duration := time.Since(start)

	err = writer.Flush()
	if err != nil {
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Error flushing writer: %v\n", err)
	}

	if verbose {
		// Gather memory stats
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Derived stats
		average := duration / time.Duration(count)
		throughput := float64(count) / duration.Seconds()
		estimatedBytes := count * (idLength + 1) // +1 for newline
		entropyPerChar := math.Log2(float64(len(alphabet)))
		estimatedEntropy := entropyPerChar * float64(idLength)

		// Print stats
		_, _ = fmt.Fprintln(cmd.OutOrStderr(), "")
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Start Time..............: %s\n", start.Format(time.RFC3339))
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Total IDs generated.....: %d\n", count)
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Total time taken........: %s\n", duration)
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Average time per ID.....: %s\n", average)
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Throughput..............: %.2f IDs/sec\n", throughput)
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Estimated output size...: %s\n", humanBytes(estimatedBytes))
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Estimated entropy per ID: %.2f bits\n", estimatedEntropy)
		_, _ = fmt.Fprintf(cmd.OutOrStderr(), "Memory used.............: %.2f MiB\n", float64(memStats.Alloc)/(1024*1024))
	}

	return nil
}

func writeError(cmd *cobra.Command, msg string, err error) error {
	// Flush stdout if necessary
	if w, ok := cmd.OutOrStdout().(*bufio.Writer); ok {
		_ = w.Flush()
	}

	_, _ = fmt.Fprintf(cmd.OutOrStderr(), "%s: %v", msg, err)
	return fmt.Errorf("%s: %w", msg, err)
}

func writeString(cmd *cobra.Command, msg string) error {
	// Flush stdout if necessary
	if w, ok := cmd.OutOrStdout().(*bufio.Writer); ok {
		_ = w.Flush()
	}

	_, _ = fmt.Fprintf(cmd.OutOrStderr(), "%s", msg)
	return fmt.Errorf("%s", msg)
}

func humanBytes(n int) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := unit, 0
	for n/unit >= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
