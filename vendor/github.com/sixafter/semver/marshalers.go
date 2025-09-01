// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"database/sql/driver"
	"encoding/json"
)

// MarshalText implements encoding.TextMarshaler.
// It returns the string representation of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha+build.456")
//	text, err := v.MarshalText()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(text)) // Output: 1.2.3-alpha+build.456
func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It parses the given text into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalText([]byte("1.2.3-alpha+build.456"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-alpha+build.456
func (v *Version) UnmarshalText(text []byte) error {
	parsed, err := Parse(string(text))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
// It returns the binary encoding of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha")
//	binaryData, err := v.MarshalBinary()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%s\n", binaryData) // Output: 1.2.3-alpha
func (v Version) MarshalBinary() ([]byte, error) {
	return v.MarshalText()
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
// It decodes the given binary data into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalBinary([]byte("1.2.3+build.123"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3+build.123
func (v *Version) UnmarshalBinary(data []byte) error {
	return v.UnmarshalText(data)
}

// MarshalJSON implements json.Marshaler.
// It returns the JSON encoding of the Version.
//
// Example:
//
//	v := semver.MustParse("1.2.3-beta")
//	jsonData, err := v.MarshalJSON()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(jsonData)) // Output: "1.2.3-beta"
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler.
// It decodes JSON data into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.UnmarshalJSON([]byte("\"1.2.3-beta+build\""))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-beta+build
func (v *Version) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}
	return v.UnmarshalText([]byte(text))
}

// Value implements database/sql/driver.Valuer.
// It returns the string representation of the Version as a database value.
//
// Example:
//
//	v := semver.MustParse("1.2.3-alpha")
//	dbValue, err := v.Value()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(dbValue) // Output: 1.2.3-alpha
func (v Version) Value() (driver.Value, error) {
	return v.String(), nil
}

// Scan implements database/sql.Scanner.
// It scans a database value into a Version.
//
// Example:
//
//	var v semver.Version
//	err := v.Scan("1.2.3-alpha+build")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(v) // Output: 1.2.3-alpha+build
func (v *Version) Scan(value interface{}) error {
	switch t := value.(type) {
	case string:
		return v.UnmarshalText([]byte(t))
	case []byte:
		return v.UnmarshalText(t)
	default:
		return ErrUnsupportedType
	}
}
