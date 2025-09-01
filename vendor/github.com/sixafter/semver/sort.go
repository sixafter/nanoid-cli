// Copyright (c) 2024-2025 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package semver

import (
	"sort"
)

// Versions attaches the methods of sort.Interface to []*Version, allowing sorting in increasing order.
// This type implements the sort.Interface for slices of Version.
//
// Example:
//
//	versions := []*Version{
//	    MustParse("1.0.0"),
//	    MustParse("2.0.0"),
//	    MustParse("1.0.1"),
//	    MustParse("1.0.0-alpha"),
//	    MustParse("1.0.0-beta"),
//	}
//	sort.Sort(Versions(versions))
//	for _, v := range versions {
//	    fmt.Println(v)
//	}
//
// Output:
// 1.0.0-alpha
// 1.0.0-beta
// 1.0.0
// 1.0.1
// 2.0.0
type Versions []*Version

// Len returns the number of elements in the slice.
// It is a required method for implementing sort.Interface.
func (s Versions) Len() int {
	return len(s)
}

// Swap exchanges the elements at the specified indices.
// It is a required method for implementing sort.Interface.
func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the element at index i should sort before the element at index j.
// It uses the LessThan method of Version to determine order.
func (s Versions) Less(i, j int) bool {
	return s[i].LessThan(*s[j])
}

// Sort sorts a slice of Version instances in increasing order.
//
// Example:
//
//	versions := []*Version{
//	    MustParse("1.2.3"),
//	    MustParse("1.0.0"),
//	    MustParse("1.1.1"),
//	}
//	Sort(versions)
//	for _, v := range versions {
//	    fmt.Println(v)
//	}
//
// Output:
// 1.0.0
// 1.1.1
// 1.2.3
func Sort(versions []*Version) {
	sort.Sort(Versions(versions))
}

// Reverse sorts a slice of Version instances in decreasing order.
//
// Example:
//
//	versions := []*Version{
//	    MustParse("1.0.0"),
//	    MustParse("1.0.1"),
//	    MustParse("1.1.0"),
//	    MustParse("2.0.0"),
//	    MustParse("1.0.0-alpha"),
//	    MustParse("1.0.0-beta"),
//	}
//	Reverse(versions)
//	for _, v := range versions {
//	    fmt.Println(v)
//	}
//
// Output:
// 2.0.0
// 1.1.0
// 1.0.1
// 1.0.0
// 1.0.0-beta
// 1.0.0-alpha
func Reverse(versions []*Version) {
	sort.Sort(sort.Reverse(Versions(versions)))
}
