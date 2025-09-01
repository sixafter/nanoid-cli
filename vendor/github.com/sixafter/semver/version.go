// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrEmptyVersionString indicates that the version string provided is empty.
	ErrEmptyVersionString = errors.New("version string is empty")

	// ErrMissingVersionElements indicates that one or more of the major, minor, or patch elements are missing in the version string.
	ErrMissingVersionElements = errors.New("missing major, minor, or patch elements")

	// ErrInvalidNumericIdentifier indicates that a numeric identifier (e.g., major, minor, or patch) is not a valid number.
	ErrInvalidNumericIdentifier = errors.New("invalid numeric identifier")

	// ErrLeadingZeroInNumericIdentifier indicates that a numeric identifier has a leading zero, which is not allowed.
	ErrLeadingZeroInNumericIdentifier = errors.New("leading zeros are not allowed in numeric identifiers")

	// ErrInvalidCharacterInIdentifier indicates that an identifier contains an invalid character.
	ErrInvalidCharacterInIdentifier = errors.New("invalid character in identifier")

	// ErrInvalidPrereleaseIdentifier indicates that a pre-release identifier contains invalid characters or is malformed.
	ErrInvalidPrereleaseIdentifier = errors.New("invalid pre-release identifier")

	// ErrEmptyPrereleaseIdentifier indicates that a pre-release identifier is empty, which is not allowed.
	ErrEmptyPrereleaseIdentifier = errors.New("empty pre-release identifier")

	// ErrEmptyBuildMetadata indicates that the build metadata portion of the version string is empty.
	ErrEmptyBuildMetadata = errors.New("build metadata is empty")

	// ErrInvalidBuildMetadataIdentifier indicates that the build metadata contains invalid characters or is malformed.
	ErrInvalidBuildMetadataIdentifier = errors.New("invalid build metadata identifier")

	// ErrUnexpectedCharacter indicates that an unexpected character was encountered in the version string.
	ErrUnexpectedCharacter = errors.New("unexpected character in version string")

	// ErrUnexpectedEndOfInput indicates that the version string ended unexpectedly during parsing.
	ErrUnexpectedEndOfInput = errors.New("unexpected end of input while parsing version string")

	// ErrUnsupportedType indicates that an unsupported type was provided for Version.
	ErrUnsupportedType = errors.New("unsupported type for Version")
)

// SupportedVersion is the latest fully supported Semantic Versioning specification version.
//
// Example:
//
//	fmt.Println(semver.SupportedVersion) // Output: 2.0.0
var SupportedVersion = Version{
	Major: 2,
	Minor: 0,
	Patch: 0,
}

// Version represents a Semantic Versioning 2.0.0 version.
//
// A Version includes major, minor, and patch numbers, as well as optional pre-release and build metadata.
//
// Example:
//
//	v := semver.Version{
//	  Major: 1,
//	  Minor: 2,
//	  Patch: 3,
//	  PreRelease: []PrereleaseVersion{{partString: "alpha"}, {partNumeric: 1, isNumeric: true}},
//	  BuildMetadata: []string{"build", "123"},
//	}
//	fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.123
type Version struct {
	BuildMetadata []string
	PreRelease    []PrereleaseVersion
	Major         uint64
	Minor         uint64
	Patch         uint64
}

var (
	// DefaultParser is a global, shared instance of a parser. It is safe for concurrent use.
	DefaultParser Parser
)

func init() {
	var err error
	DefaultParser, err = NewParser()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DefaultParser: %v", err))
	}
}

// NewParser creates a new Parser instance with the provided options.
// This function accepts a variadic number of Option parameters, allowing
// users to configure the behavior of the Parser as needed.
//
// Options can be used to customize various aspects of the Parser, such as
// specifying custom delimiters, enabling or disabling specific parsing features,
// or configuring error handling behavior.
//
// Example usage:
//
//	// Create a Parser with default settings
//	parser, err := NewParser()
//	if err != nil {
//	    log.Fatalf("Failed to create parser: %v", err)
//	}
//
//	// Create a Parser with custom options
//	parser, err := NewParser(WithDelimiter(','), WithStrictAdherence(true))
//	if err != nil {
//	    log.Fatalf("Failed to create custom parser: %v", err)
//	}
//
// Parameters:
// - options: A variadic list of Option functions used to configure the Parser.
//
// Returns:
// - Parser: An instance of the Parser configured with the specified options.
// - error: An error if there is an issue creating the Parser, otherwise nil.
func NewParser(options ...Option) (Parser, error) {
	// Initialize ConfigOptions with default values.
	// These defaults include the default alphabet, the default random reader,
	// and the default length hint for ID generation.
	configOpts := &ConfigOptions{
		Strict: true,
	}

	// Apply provided options to customize the configuration.
	// Each Option function modifies the ConfigOptions accordingly.
	for _, opt := range options {
		opt(configOpts)
	}

	config, err := buildRuntimeConfig(configOpts)
	if err != nil {
		return nil, err
	}

	return &parser{
		config: config,
	}, nil
}

// Parser defines an interface for parsing version strings into structured Version objects.
// Implementations of this interface are responsible for validating and converting a version
// string into a Version type that can be used programmatically.
//
// This interface can be useful when dealing with different version formats or when you need
// to standardize version parsing across multiple components of an application.
//
// Example usage:
//
//	var parser Parser = NewParser()
//	versionStr := "1.2.3-beta+build.123"
//	version, err := parser.Parse(versionStr)
//	if err != nil {
//	    log.Fatalf("Failed to parse version: %v", err)
//	}
//	fmt.Printf("Parsed version: %v\n", version)
//
// Methods:
//   - Parse(version string) (Version, error): Parses a version string and returns a Version object.
//     Returns an error if the version string is invalid or cannot be parsed.
type Parser interface {
	// Parse takes a version string as input and converts it into a structured Version object.
	// The input version string must follow a valid versioning format, and the implementation
	// of the method is responsible for handling the parsing logic.
	//
	// If the provided version string is invalid or cannot be parsed, an error will be returned.
	//
	// Parameters:
	// - version: A string representing the version to be parsed (e.g., "1.2.3", "1.0.0-alpha+build.123").
	//
	// Returns:
	// - Version: A Version object representing the parsed version information.
	// - error: An error if the version string is invalid or cannot be parsed.
	//
	// Example usage:
	//
	//    version, err := parser.Parse("1.2.3")
	//    if err != nil {
	//        log.Fatalf("Failed to parse version: %v", err)
	//    }
	//    fmt.Printf("Parsed version: %v\n", version)
	Parse(version string) (Version, error)
}

type parser struct {
	config *runtimeConfig
}

// New creates a new Version instance with the specified major, minor, patch components,
// optional prerelease identifiers, and optional build metadata.
//
// The prerelease identifiers determine the precedence of the version relative to other versions with the same
// major, minor, and patch numbers. Build metadata provides additional information but does not affect version precedence.
//
// Parameters:
//   - major: The major version number (uint64).
//   - minor: The minor version number (uint64).
//   - patch: The patch version number (uint64).
//   - preRelease: A slice of PrereleaseVersion structs representing the prerelease identifiers (optional).
//   - buildMetadata: A slice of strings representing the build metadata (optional).
//
// Returns:
//   - A Version struct populated with the provided values.
//
// Example:
//
//	package main
//
//	import (
//	    "fmt"
//	    "semver"
//	)
//
//	func main() {
//	    preRelease := NewPrereleaseVersion("alpha.1")
//	    buildMetadata := []string{"build", "2024"}
//
//	    v := semver.New(1, 2, 3, preRelease, buildMetadata)
//
//	    fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.2024
//	}
func New(major, minor, patch uint64, preRelease []PrereleaseVersion, buildMetadata []string) Version {
	return Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		PreRelease:    preRelease,
		BuildMetadata: buildMetadata,
	}
}

// MustParse is a helper function that parses a version string and panics if invalid.
//
// Example:
//
//	v := semver.MustParse("1.2.3")
//	fmt.Println(v) // Output: 1.2.3
func MustParse(version string) Version {
	v, err := DefaultParser.Parse(version)
	if err != nil {
		panic(err)
	}
	return v
}

// Parse parses a version string into a Version struct.
//
// Returns an error if the version string is not a valid semantic version.
//
// Example:
//
//	v, err := semver.Parse("1.2.3-alpha.1+build.123")
//	if err != nil {
//	  fmt.Println("Error:", err)
//	} else {
//	  fmt.Println(v) // Output: 1.2.3-alpha.1+build.123
//	}
//
// Parse parses a version string into a Version struct.
// The version string must follow semantic versioning format, such as "1.0.0-alpha+001".
// It returns an error if the version string is invalid.
func Parse(version string) (Version, error) {
	return DefaultParser.Parse(version)
}

// Parse parses a version string into a Version struct.
//
// Returns an error if the version string is not a valid semantic version.
//
// Example:
//
//	v, err := semver.Parse("1.2.3-alpha.1+build.123")
//	if err != nil {
//	  fmt.Println("Error:", err)
//	} else {
//	  fmt.Println(v) // Output: 1.2.3-alpha.1+build.123
//	}
//
// Parse parses a version string into a Version struct.
// The version string must follow semantic versioning format, such as "1.0.0-alpha+001".
// It returns an error if the version string is invalid.
func (p *parser) Parse(version string) (Version, error) {
	if len(version) == 0 {
		return Version{}, ErrEmptyVersionString
	}

	var v Version
	var index int
	length := len(version)
	var err error

	// Parse Major
	v.Major, index, err = p.parseNumericIdentifier(version, index, length)
	if err != nil {
		return Version{}, err
	}

	// Expect a '.' after Major
	if index >= length || version[index] != '.' {
		return Version{}, ErrMissingVersionElements
	}
	index++ // Skip '.'

	// Parse Minor
	v.Minor, index, err = p.parseNumericIdentifier(version, index, length)
	if err != nil {
		return Version{}, err
	}

	// Expect a '.' after Minor
	if index >= length || version[index] != '.' {
		return Version{}, ErrMissingVersionElements
	}
	index++ // Skip '.'

	// Parse Patch
	v.Patch, index, err = p.parseNumericIdentifier(version, index, length)
	if err != nil {
		return Version{}, err
	}

	// Parse PreRelease and BuildMetadata if any
	if index < length {
		index, err = p.parsePreReleaseAndBuildMetadata(version, index, length, &v)
		if err != nil {
			return Version{}, err
		}
	}

	if index != length {
		return Version{}, ErrUnexpectedCharacter
	}

	return v, nil
}

// parseNumericIdentifier parses a numeric identifier from the version string.
// It returns the parsed value, the updated index, or an error if the parsing fails.
func (p *parser) parseNumericIdentifier(version string, index int, length int) (uint64, int, error) {
	if index >= length {
		return 0, index, ErrUnexpectedEndOfInput
	}

	start := index
	if version[index] == '0' {
		index++
		if index < length && version[index] >= '0' && version[index] <= '9' {
			return 0, index, ErrLeadingZeroInNumericIdentifier
		}
		return 0, index, nil
	}

	var n uint64
	for index < length && version[index] >= '0' && version[index] <= '9' {
		n = n*10 + uint64(version[index]-'0')
		index++
	}

	if start == index {
		return 0, index, ErrInvalidNumericIdentifier
	}

	return n, index, nil
}

// parsePreReleaseAndBuildMetadata parses the pre-release and build metadata components from the version string.
// It updates the Version struct with the parsed values and returns the updated index or an error.
func (p *parser) parsePreReleaseAndBuildMetadata(version string, index int, length int, v *Version) (int, error) {
	var err error

	// Parse PreRelease if present
	if index < length && version[index] == '-' {
		index++ // Skip '-'
		start := index
		for index < length && version[index] != '+' {
			if version[index] > 127 {
				return index, ErrInvalidCharacterInIdentifier
			}
			index++
		}
		prerelease := version[start:index]
		v.PreRelease, err = p.parsePrerelease(prerelease)
		if err != nil {
			return index, err
		}
	}

	// Parse BuildMetadata if present
	if index < length && version[index] == '+' {
		index++ // Skip '+'
		start := index
		build := version[start:]
		v.BuildMetadata, err = p.parseBuildMetadata(build)
		if err != nil {
			return index, err
		}
		index = length // End of string
	}

	return index, nil
}

// parsePrerelease parses the given string into a slice of PrereleaseVersion components.
// The string is expected to contain prerelease identifiers separated by dots.
//
// Prerelease identifiers must conform to the following rules:
//   - Identifiers must not be empty.
//   - Identifiers must only contain alphanumeric characters or hyphens.
//   - Numeric identifiers must not have leading zeros.
//
// Returns an error if the input string is empty, contains invalid characters, or contains empty identifiers.
//
// Example:
//
//	s := "alpha.1.0-beta"
//	prerelease, err := parsePrerelease(s)
//	if err != nil {
//	    // handle error
//	}
func (p *parser) parsePrerelease(s string) ([]PrereleaseVersion, error) {
	if len(s) == 0 {
		return nil, ErrEmptyPrereleaseIdentifier
	}

	var prerelease []PrereleaseVersion
	length := len(s)
	start := 0

	for i := 0; i <= length; i++ {
		if i == length || s[i] == '.' {
			if start == i {
				return nil, ErrEmptyPrereleaseIdentifier
			}
			part := s[start:i]

			if !p.isValidPrereleaseIdentifier(part) {
				return nil, ErrInvalidPrereleaseIdentifier
			}
			component, err := NewPrereleaseVersion(part)

			if err != nil {
				return nil, err
			}

			prerelease = append(prerelease, component)
			start = i + 1
		} else if s[i] > 127 || !p.isAllowedInIdentifier(s[i]) {
			return nil, ErrInvalidCharacterInIdentifier
		}
	}
	return prerelease, nil
}

// parseBuildMetadata parses the given string into a slice of build metadata components.
// The string is expected to contain build identifiers separated by dots.
//
// Build identifiers must conform to the following rules:
//   - Identifiers must not be empty.
//   - Identifiers must only contain alphanumeric characters or hyphens.
//
// Returns an error if the input string is empty, contains invalid characters, or contains empty identifiers.
//
// Example:
//
//	s := "001.alpha"
//	buildMetadata, err := parseBuildMetadata(s)
//	if err != nil {
//	    // handle error
//	}
func (p *parser) parseBuildMetadata(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, ErrEmptyBuildMetadata
	}

	var buildMetadata []string
	length := len(s)
	start := 0

	for i := 0; i <= length; i++ {
		if i == length || s[i] == '.' {
			if start == i {
				return nil, ErrEmptyBuildMetadata
			}
			part := s[start:i]
			if !p.isValidBuildIdentifier(part) {
				return nil, ErrInvalidBuildMetadataIdentifier
			}
			buildMetadata = append(buildMetadata, part)
			start = i + 1
		} else if s[i] > 127 || !p.isAllowedInIdentifier(s[i]) {
			return nil, ErrInvalidCharacterInIdentifier
		}
	}
	return buildMetadata, nil
}

// isAllowedInIdentifier checks if a character is allowed in a semantic version identifier.
// Allowed characters are:
//   - Digits ('0'-'9')
//   - Uppercase letters ('A'-'Z')
//   - Lowercase letters ('a'-'z')
//   - Hyphen ('-')
func (p *parser) isAllowedInIdentifier(ch byte) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= 'a' && ch <= 'z') ||
		ch == '-'
}

// isValidPrereleaseIdentifier checks if a prerelease identifier is valid.
// The identifier must not be empty and must contain only allowed characters.
// Numeric identifiers must not have leading zeros.
func (p *parser) isValidPrereleaseIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if !p.isAllowedInIdentifier(ch) {
			return false
		}
	}
	if p.config.StrictAdherence() && isNumeric(s) && s[0] == '0' && len(s) > 1 {
		return false // Leading zeros are not allowed in numeric identifiers
	}

	return true
}

// isValidBuildIdentifier checks if a build metadata identifier is valid.
// The identifier must not be empty and must contain only allowed characters.
func (p *parser) isValidBuildIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if !p.isAllowedInIdentifier(ch) {
			return false
		}
	}
	return true
}

// isNumeric checks if a string consists only of numeric characters ('0'-'9').
func isNumeric(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}

// String returns the string representation of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha.1+build.123")
//	fmt.Println(v.String()) // Output: 1.2.3-alpha.1+build.123
func (v Version) String() string {
	var sb strings.Builder
	sb.Grow(16) // Preallocate a reasonable capacity

	sb.WriteString(strconv.FormatUint(v.Major, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Minor, 10))
	sb.WriteByte('.')
	sb.WriteString(strconv.FormatUint(v.Patch, 10))

	if len(v.PreRelease) > 0 {
		sb.WriteByte('-')
		for i, pr := range v.PreRelease {
			if i > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(pr.String())
		}
	}

	if len(v.BuildMetadata) > 0 {
		sb.WriteByte('+')
		for i, bm := range v.BuildMetadata {
			if i > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(bm)
		}
	}

	return sb.String()
}

// Compare compares two Version instances.
// Returns -1 if v < other, 0 if v == other, +1 if v > other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.4")
//	fmt.Println(v1.Compare(v2)) // Output: -1
func (v Version) Compare(other Version) int {
	// Compare Major
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}

	// Compare Minor
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}

	// Compare Patch
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	// Handle pre-release comparison
	if len(v.PreRelease) == 0 && len(other.PreRelease) == 0 {
		return 0
	} else if len(v.PreRelease) == 0 && len(other.PreRelease) > 0 {
		return 1
	} else if len(v.PreRelease) > 0 && len(other.PreRelease) == 0 {
		return -1
	}

	// Compare pre-release identifiers one by one
	length := len(v.PreRelease)
	if len(other.PreRelease) < length {
		length = len(other.PreRelease)
	}

	for i := 0; i < length; i++ {
		c := v.PreRelease[i].Compare(other.PreRelease[i])
		if c != 0 {
			return c
		}
	}

	// If all compared identifiers are equal, the one with more identifiers has higher precedence
	if len(v.PreRelease) < len(other.PreRelease) {
		return -1
	}
	if len(v.PreRelease) > len(other.PreRelease) {
		return 1
	}

	return 0
}

// Equal checks if two versions are equal.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.Equal(v2)) // Output: true
func (v Version) Equal(other Version) bool {
	return v.Compare(other) == 0
}

// LessThan checks if v is less than other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.4")
//	fmt.Println(v1.LessThan(v2)) // Output: true
func (v Version) LessThan(other Version) bool {
	return v.Compare(other) == -1
}

// LessThanOrEqual checks if v is less than or equal to other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.LessThanOrEqual(v2)) // Output: true
func (v Version) LessThanOrEqual(other Version) bool {
	return v.Compare(other) <= 0
}

// GreaterThan checks if v is greater than other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.4")
//	v2 := semver.MustParse("1.2.3")
//	fmt.Println(v1.GreaterThan(v2)) // Output: true
func (v Version) GreaterThan(other Version) bool {
	return v.Compare(other) == 1
}

// GreaterThanOrEqual checks if v is greater than or equal to other.
//
// Example:
//
//	v1 := semver.MustParse("1.2.3")
//	v2 := semver.MustParse("1.2.2")
//	fmt.Println(v1.GreaterThanOrEqual(v2)) // Output: true
func (v Version) GreaterThanOrEqual(other Version) bool {
	return v.Compare(other) >= 0
}

// Config holds the runtime configuration for the parser.
//
// It is immutable after initialization.
func (p *parser) Config() Config {
	return p.config
}
