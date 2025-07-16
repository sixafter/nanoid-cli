// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"fmt"
	"regexp"
	"strings"
)

// Operator represents a version comparison operator.
//
// Supported Operators:
//   - OpEq ("="): Equal
//   - OpGt (">"): Greater than
//   - OpGte (">="): Greater than or equal
//   - OpLt ("<"): Less than
//   - OpLte ("<="): Less than or equal
//   - OpNeq ("!="): Not equal
type Operator string

const (
	OpEq  Operator = "="
	OpGt  Operator = ">"
	OpGte Operator = ">="
	OpLt  Operator = "<"
	OpLte Operator = "<="
	OpNeq Operator = "!="
)

// Requirement represents a single version requirement.
// It consists of a comparison operator and a target version.
//
// Example:
//
//	req := Requirement{Op: OpGt, Ver: semver.MustParse("1.0.0")}
//	v := semver.MustParse("1.1.0")
//	fmt.Println(req.Contains(v)) // Output: true
type Requirement struct {
	Op  Operator
	Ver Version
}

// VersionRange represents a set of requirements separated by AND (space) and OR (||).
// A VersionRange can be used to check if a Version satisfies it.
//
// Example:
//
//	r, err := semver.ParseRange(">1.0.0 <2.0.0 || >=3.0.0 !=4.2.1")
//	if err != nil {
//	    // Handle error
//	}
//	v := semver.MustParse("1.2.3")
//	fmt.Println(r.Contains(v)) // Output: true
type VersionRange struct {
	Requirements [][]Requirement
}

// rangeRegex helps to parse individual range tokens.
var rangeRegex = regexp.MustCompile(`^(>=|<=|>|<|=|!=)?\s*([0-9A-Za-z.\-+]+)$`)

// ParseRange parses a range string into a VersionRange struct.
//
// Valid ranges are:
//   - "<1.0.0"
//   - "<=1.0.0"
//   - ">1.0.0"
//   - ">=1.0.0"
//   - "1.0.0", "=1.0.0", "==1.0.0"
//   - "!1.0.0", "!=1.0.0"
//
// Ranges can be combined with logical AND (space-separated) and logical OR (||):
//   - ">1.0.0 <2.0.0" matches between both versions.
//   - "<2.0.0 || >=3.0.0" matches either version ranges.
//
// Example:
//
//	r, err := semver.ParseRange(">1.0.0 <2.0.0")
//	if err != nil {
//	    fmt.Println("Error parsing range:", err)
//	    return
//	}
//	v := semver.MustParse("1.5.0")
//	fmt.Println(r.Contains(v)) // Output: true
func ParseRange(r string) (*VersionRange, error) {
	orParts := strings.Split(r, "||")
	var requirements [][]Requirement

	for _, part := range orParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		tokens := strings.Fields(part)
		var reqs []Requirement
		for _, token := range tokens {
			matches := rangeRegex.FindStringSubmatch(token)
			if matches == nil {
				return nil, fmt.Errorf("invalid range token: %s", token)
			}
			opStr := matches[1]
			verStr := matches[2]
			var op Operator
			if opStr == "" {
				op = OpEq
			} else {
				op = Operator(opStr)
			}
			ver, err := Parse(verStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in range: %s", verStr)
			}
			reqs = append(reqs, Requirement{
				Op:  op,
				Ver: ver,
			})
		}
		requirements = append(requirements, reqs)
	}

	return &VersionRange{
		Requirements: requirements,
	}, nil
}

// MustParseRange is like ParseRange but panics if the range cannot be parsed.
//
// This function is useful for scenarios where you are certain the input is valid
// and want to avoid handling errors explicitly. However, use it with caution
// in production code, as it will terminate the program if the range is invalid.
//
// Example:
//
//	r := semver.MustParseRange(">1.0.0 <2.0.0")
//	v := semver.MustParse("1.5.0")
//	fmt.Println(r.Contains(v)) // Output: true
func MustParseRange(s string) *VersionRange {
	r, err := ParseRange(s)
	if err != nil {
		panic(`semver: ParseRange(` + s + `): ` + err.Error())
	}
	return r
}

// Contains checks if a version satisfies the range.
//
// Example:
//
//	r, _ := semver.ParseRange(">1.0.0 <2.0.0")
//	v := semver.MustParse("1.5.0")
//	fmt.Println(r.Contains(v)) // Output: true
func (vr *VersionRange) Contains(v Version) bool {
	for _, andReqs := range vr.Requirements {
		matchesAll := true
		for _, req := range andReqs {
			if !req.Contains(v) {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			return true
		}
	}
	return false
}

// Contains checks if a version satisfies the requirement.
//
// Example:
//
//	req := Requirement{Op: OpGt, Ver: semver.MustParse("1.0.0")}
//	v := semver.MustParse("1.1.0")
//	fmt.Println(req.Contains(v)) // Output: true
func (r *Requirement) Contains(v Version) bool {
	switch r.Op {
	case OpEq:
		return v.Equal(r.Ver)
	case OpGt:
		return v.GreaterThan(r.Ver)
	case OpGte:
		return v.GreaterThanOrEqual(r.Ver)
	case OpLt:
		return v.LessThan(r.Ver)
	case OpLte:
		return v.LessThanOrEqual(r.Ver)
	case OpNeq:
		return !v.Equal(r.Ver)
	default:
		return false
	}
}

// OR combines the current VersionRange with another VersionRange using logical OR.
//
// Example:
//
//	r1, _ := semver.ParseRange(">1.0.0 <2.0.0")
//	r2, _ := semver.ParseRange(">=3.0.0 !=4.2.1")
//	combined := r1.OR(r2)
//	v := semver.MustParse("3.1.0")
//	fmt.Println(combined.Contains(v)) // Output: true
func (vr *VersionRange) OR(other *VersionRange) *VersionRange {
	combined := &VersionRange{
		Requirements: append(vr.Requirements, other.Requirements...),
	}
	return combined
}

// AND combines the current VersionRange with another VersionRange using logical AND.
//
// This function returns a new VersionRange that represents the intersection of the two ranges.
// It effectively creates a range that only matches versions satisfying both original ranges.
//
// Example:
//
//	r1, _ := semver.ParseRange(">1.0.0 <3.0.0")
//	r2, _ := semver.ParseRange("!=2.0.3-beta.2")
//	combined := r1.AND(r2)
//
//	v := semver.MustParse("2.1.0")
//	fmt.Println(combined.Contains(v)) // Output: true
//
//	v2 := semver.MustParse("2.0.3-beta.2")
//	fmt.Println(combined.Contains(v2)) // Output: false
func (vr *VersionRange) AND(other *VersionRange) *VersionRange {
	var combinedRequirements [][]Requirement
	for _, reqs1 := range vr.Requirements {
		for _, reqs2 := range other.Requirements {
			combinedReqs := append(reqs1, reqs2...)
			combinedRequirements = append(combinedRequirements, combinedReqs)
		}
	}
	return &VersionRange{
		Requirements: combinedRequirements,
	}
}
