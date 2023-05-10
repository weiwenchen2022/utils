// Package set implements a set.
package set

import (
	"fmt"
	"strings"
)

var emptyStruct = struct{}{}

// Set represents a set. To create a Set use New().
type Set[E comparable] map[E]struct{}

// New returns a initialized set that contains elements of xs.
func New[E comparable](xs ...E) Set[E] {
	s := make(Set[E], len(xs))

	for _, x := range xs {
		s[x] = emptyStruct
	}

	return s
}

// IsEmpty reports whether the set s is empty.
func (s Set[E]) IsEmpty() bool {
	return len(s) == 0
}

// Len returns the number of elements in the set s.
func (s Set[E]) Len() int {
	return len(s)
}

// Add adds x to the set s, and reports whether the set grew.
func (s Set[E]) Add(x E) bool {
	if _, ok := s[x]; ok {
		return false
	}

	s[x] = emptyStruct
	return true
}

// AddAll adds the elements of xs to the set s.
func (s Set[E]) AddAll(xs ...E) {
	for _, x := range xs {
		s[x] = emptyStruct
	}
}

// Remove removes x from the set s, and reports whether the set shrank.
func (s Set[E]) Remove(x E) bool {
	if _, ok := s[x]; !ok {
		return false
	}

	delete(s, x)
	return true
}

// RemoveAll removes the elements of xs from the set s.
func (s Set[E]) RemoveAll(xs ...E) {
	for _, x := range xs {
		delete(s, x)
	}
}

// Clear removes all elements from the set s.
func (s Set[E]) Clear() {
	for x := range s {
		delete(s, x)
	}
}

// Has reports whether x is an element of the set s.
func (s Set[E]) Has(x E) bool {
	_, ok := s[x]
	return ok
}

// Copy return a copy of the set s.
func (s Set[E]) Copy() Set[E] {
	copy := make(Set[E], len(s))

	for x := range s {
		copy[x] = emptyStruct
	}

	return copy
}

// Equals reports whether the sets s and t have the same elements.
func (s Set[E]) Equals(t Set[E]) bool {
	if len(s) != len(t) {
		return false
	}

	for x := range s {
		if _, ok := t[x]; !ok {
			return false
		}
	}

	return true
}

// String returns a human-readable description of the set s.
func (s Set[E]) String() string {
	var b strings.Builder

	b.WriteByte('{')
	for x := range s {
		if b.Len() > len("{") {
			b.WriteByte(' ')
		}

		fmt.Fprintf(&b, "%v", x)
	}
	b.WriteByte('}')

	return b.String()
}

// AppendTo returns the result of appending the elements of s to slice.
func (s Set[E]) AppendTo(slice []E) []E {
	tot := len(slice) + len(s)
	if tot <= cap(slice) {
		for x := range s {
			slice = append(slice, x)
		}
		return slice
	}

	slice = append(slice, make([]E, len(s))...)[:len(slice)]
	for x := range s {
		slice = append(slice, x)
	}
	return slice
}

// Elems returns the slice of the elements of s.
func (s Set[E]) Elems() []E {
	return s.AppendTo(nil)
}

// IntersectWith sets s to the intersection s ∩ t, and reports whether the set shrank.
func (s Set[E]) IntersectWith(t Set[E]) bool {
	var shrank bool

	for x := range s {
		if _, ok := t[x]; !ok {
			delete(s, x)
			shrank = true
		}
	}

	return shrank
}

// Intersects reports whether s ∩ t ≠ ∅.
func (s Set[E]) Intersects(t Set[E]) bool {
	for x := range s {
		if _, ok := t[x]; ok {
			return true
		}
	}

	return false
}

// UnionWith sets s to the union s ∪ t, and reports whether s grew.
func (s Set[E]) UnionWith(t Set[E]) bool {
	var grew bool

	for x := range t {
		if _, ok := s[x]; !ok {
			s[x] = emptyStruct
			grew = true
		}
	}

	return grew
}

// DifferenceWith sets s to the difference s ∖ t, and reports whether the set shrank.
func (s Set[E]) DifferenceWith(t Set[E]) bool {
	var shrank bool

	for x := range s {
		if _, ok := t[x]; ok {
			delete(s, x)
			shrank = true
		}
	}

	return shrank
}

// SubsetOf reports whether s ∖ t = ∅.
func (s Set[E]) SubsetOf(t Set[E]) bool {
	if len(s) > len(t) {
		return false
	}

	for x := range s {
		if _, ok := t[x]; !ok {
			return false
		}
	}

	return true
}

// SymmetricDifferenceWith sets s to the symmetric difference s ∆ t.
func (s Set[E]) SymmetricDifferenceWith(t Set[E]) {
	deleted := make([]E, 0, len(s))

	for x := range s {
		if _, ok := t[x]; ok {
			deleted = append(deleted, x)
		}
	}

	for x := range t {
		if _, ok := s[x]; !ok {
			s[x] = emptyStruct
		}
	}

	for _, x := range deleted {
		delete(s, x)
	}
}
