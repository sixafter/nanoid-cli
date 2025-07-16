// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

// ConfigOptions holds the configurable options for the Parser.
// It is used with the Function Options pattern.
type ConfigOptions struct {
	Strict bool
}

// Config holds the runtime configuration for the parser.
//
// It is immutable after initialization.
type Config interface {
	// StrictAdherence returns a boolean value indicating whether strict adherence is enabled.
	// This method is used to determine if the configuration should follow strict rules,
	// such as requiring full compliance with the Semantic Versioning specification.
	// When enabled, parsing or processing might be more stringent, rejecting inputs that
	// do not fully comply with the expected standards.
	//
	// Returns:
	// - bool: true if strict adherence is enabled, false otherwise.
	//
	// Example usage:
	//
	//    var config Config = NewConfig()
	//    if config.StrictAdherence() {
	//        fmt.Println("Strict adherence is enabled.")
	//    } else {
	//        fmt.Println("Strict adherence is disabled.")
	//    }
	StrictAdherence() bool
}

// Configuration defines the interface for retrieving parser configuration.
type Configuration interface {
	// Config returns the runtime configuration of the parser.
	Config() Config
}

type runtimeConfig struct {
	strict bool
}

// Option defines a function type for configuring the Parser.
// It allows for flexible and extensible configuration by applying
// various settings to the ConfigOptions during Parser initialization.
type Option func(*ConfigOptions)

// WithStrictAdherence sets the strict adherence value for the configuration.
// This option can be used to enable or disable strict mode, which affects the way
// certain rules are enforced during parsing or processing.
//
// Setting strict adherence to true can be used to enforce more rigid compliance
// with versioning rules or configuration standards. When set to false, the parser
// may allow some flexibility in handling certain inputs.
//
// Parameters:
// - value: A boolean indicating whether strict adherence should be enabled (true) or disabled (false).
//
// Returns:
// - Option: A functional option that can be passed to a configuration function to modify behavior.
//
// Example usage:
//
//	parser, err := NewParser(WithStrictAdherence(true))
//	if err != nil {
//	    log.Fatalf("Failed to create parser: %v", err)
//	}
//
//	// Use the parser with strict adherence enabled
//	version, err := parser.Parse("1.0.0")
//	if err != nil {
//	    log.Fatalf("Failed to parse version: %v", err)
//	}
//	fmt.Printf("Parsed version: %v\n", version)
func WithStrictAdherence(value bool) Option {
	return func(o *ConfigOptions) {
		o.Strict = value
	}
}

// StrictAdherence returns a boolean value indicating whether strict adherence is enabled.
// This method is used to determine if the configuration should follow strict rules,
// such as requiring full compliance with the Semantic Versioning specification.
// When enabled, parsing or processing might be more stringent, rejecting inputs that
// do not fully comply with the expected standards.
//
// Returns:
// - bool: true if strict adherence is enabled, false otherwise.
//
// Example usage:
//
//	var config Config = NewConfig()
//	if config.StrictAdherence() {
//	    fmt.Println("Strict adherence is enabled.")
//	} else {
//	    fmt.Println("Strict adherence is disabled.")
//	}
func (c *runtimeConfig) StrictAdherence() bool {
	return c.strict
}

func buildRuntimeConfig(opts *ConfigOptions) (*runtimeConfig, error) {
	return &runtimeConfig{
		strict: opts.Strict,
	}, nil
}
