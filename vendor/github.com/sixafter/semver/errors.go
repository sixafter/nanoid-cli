// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"errors"
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
