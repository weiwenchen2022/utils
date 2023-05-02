// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package slices defines various types and functions useful with slices of any type.
// Unless otherwise specified, these functions all apply to the elements of a slice at index 0 <= i < len(s).
package slices

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

// Slice attaches common methods to []E.
type Slice[E any] []E

// NewSlice creates and initializes a new Slice using s as its initial
// contents. The new Slice takes ownership of s, and the caller should not
// use s after this call.
func NewSlice[E any](s []E) *Slice[E] {
	ss := Slice[E](s)
	return &ss
}

// EqualFunc reports whether s equal to t using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func (s Slice[E]) EqualFunc(t []E, eq func(E, E) bool) bool {
	return EqualFunc(s, t, eq)
}

// CompareFunc compares the elements of s and t using a comparison function
// on each pair of elements. The elements are compared in increasing
// index order, and the comparisons stop after the first time cmp
// returns non-zero.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func (s Slice[E]) CompareFunc(t []E, cmp func(E, E) int) int {
	return CompareFunc(s, t, cmp)
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func (s Slice[E]) IndexFunc(f func(E) bool) int {
	return IndexFunc(s, f)
}

// ContainsFunc reports whether at least one
// element e of s satisfies f(e).
func (s Slice[E]) ContainsFunc(f func(E) bool) bool {
	return ContainsFunc(s, f)
}

// Insert inserts the values v... into s at index i,
// updating the slice s.
// In the updated slice s, s[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func (s *Slice[E]) Insert(i int, v ...E) {
	*s = Insert(*s, i, v...)
}

// Delete removes the elements s[i:j] from s, updating the slice s.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements is pointers or contain pointers, Delete zeroing those elements so that
// objects they reference can be garbage collected.
func (s *Slice[E]) Delete(i, j int) {
	*s = Delete(*s, i, j)
}

// Replace replaces the elements s[i:j] by the given v, and updates the
// slice s. Replace panics if s[i:j] is not a valid slice of s.
func (s *Slice[E]) Replace(i, j int, v ...E) {
	*s = Replace(*s, i, j, v...)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func (s Slice[E]) Clone() *Slice[E] {
	copy := Clone(s)
	return &copy
}

// DeepClone is like Clone but if elements has an Clone method of the form "Clone() E",
// it use the result of e.Clone() as assignment right operand.
func (s Slice[E]) DeepClone() *Slice[E] {
	copy := DeepClone(s)
	return &copy
}

// CompactFunc is like Compact but uses a comparison function.
func (s *Slice[E]) CompactFunc(eq func(E, E) bool) {
	*s = CompactFunc(*s, eq)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func (s *Slice[E]) Grow(n int) {
	*s = Grow(*s, n)
}

// Clip removes unused capacity from the slice, updating s tp s[:len(s):len(s)].
func (s *Slice[E]) Clip() {
	*s = Clip(*s)
}

// IsNil reports whether s is nil.
func (s Slice[E]) IsNil() bool {
	return s == nil
}

// Len returns the number of elements in s.
func (s Slice[E]) Len() int {
	return len(s)
}

// Cap returns the capacity of s.
func (s Slice[E]) Cap() int {
	return cap(s)
}

// Slice returns s[i:j]. It panics if s[i:j] is not a valid slice of s.
func (s Slice[E]) Slice(i, j int) Slice[E] {
	return s[i:j]
}

// Slice3 is the 3-index form of the slice operation: it returns s[i:j:k].
// It panics if s[i:j:k] is not a valid slice of s.
func (s Slice[E]) Slice3(i, j, k int) Slice[E] {
	return s[i:j:k]
}

// Append appends the values x to s and updates the slice s.
func (s *Slice[E]) Append(x ...E) {
	*s = append(*s, x...)
}

// AppendSlice appends a slice t to s and updates the slice s.
func (s *Slice[E]) AppendSlice(t []E) {
	*s = append(*s, t...)
}

// ComparableSlice is like Slice but element type requires comparable.
type ComparableSlice[E comparable] []E

// NewComparableSlice creates and initializes a new ComparableSlice using s as its initial
// contents. The new ComparableSlice takes ownership of s, and the caller should not
// use s after this call.
func NewComparableSlice[E comparable](s []E) *ComparableSlice[E] {
	ss := ComparableSlice[E](s)
	return &ss
}

// Equal reports whether s equal to t: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func (s ComparableSlice[E]) Equal(t []E) bool {
	return Equal(s, t)
}

// EqualFunc reports whether s equal to t using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
func (s ComparableSlice[E]) EqualFunc(t []E, eq func(E, E) bool) bool {
	return EqualFunc(s, t, eq)
}

// CompareFunc is like Compare but uses a comparison function
// on each pair of elements. The elements are compared in increasing
// index order, and the comparisons stop after the first time cmp
// returns non-zero.
// The result is the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func (s ComparableSlice[E]) CompareFunc(t []E, cmp func(E, E) int) int {
	return CompareFunc(s, t, cmp)
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func (s ComparableSlice[E]) Index(v E) int {
	return Index(s, v)
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func (s ComparableSlice[E]) IndexFunc(f func(E) bool) int {
	return IndexFunc(s, f)
}

// Contains reports whether v is present in s.
func (s ComparableSlice[E]) Contains(v E) bool {
	return Contains(s, v)
}

// ContainsFunc reports whether at least one
// element e of s satisfies f(e).
func (s ComparableSlice[E]) ContainsFunc(f func(E) bool) bool {
	return ContainsFunc(s, f)
}

// Insert inserts the values v... into s at index i,
// updating the slice s.
// In the updated slice s, s[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func (s *ComparableSlice[E]) Insert(i int, v ...E) {
	*s = Insert(*s, i, v...)
}

// Delete removes the elements s[i:j] from s, updating the slice s.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements is pointers or contain pointers, Delete zeroing those elements so that
// objects they reference can be garbage collected.
func (s *ComparableSlice[E]) Delete(i, j int) {
	*s = Delete(*s, i, j)
}

// Replace replaces the elements s[i:j] by the given v, and updates the
// slice s. Replace panics if s[i:j] is not a valid slice of s.
func (s *ComparableSlice[E]) Replace(i, j int, v ...E) {
	*s = Replace(*s, i, j, v...)
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func (s ComparableSlice[E]) Clone() *ComparableSlice[E] {
	copy := Clone(s)
	return &copy
}

// DeepClone is like Clone but if elements has an Clone method of the form "Clone() E",
// it use the result of e.Clone() as assignment right operand.
func (s ComparableSlice[E]) DeepClone() *ComparableSlice[E] {
	copy := DeepClone(s)
	return &copy
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
// When Compact discards m elements in total, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements is pointers or contain pointers, Compact
// zeroing those elements so that objects they reference can be garbage collected.
func (s *ComparableSlice[E]) Compact() {
	*s = Compact(*s)
}

// CompactFunc is like Compact but uses a comparison function.
func (s *ComparableSlice[E]) CompactFunc(eq func(E, E) bool) {
	*s = CompactFunc(*s, eq)
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func (s *ComparableSlice[E]) Grow(n int) {
	*s = Grow(*s, n)
}

// Clip removes unused capacity from the slice, updating s to s[:len(s):len(s)].
func (s *ComparableSlice[E]) Clip() {
	*s = Clip(*s)
}

// IsNil reports whether s is nil.
func (s ComparableSlice[E]) IsNil() bool {
	return s == nil
}

// Len returns the number of elements in s.
func (s ComparableSlice[E]) Len() int {
	return len(s)
}

// Cap returns the capacity of s.
func (s ComparableSlice[E]) Cap() int {
	return cap(s)
}

// Slice returns s[i:j]. It panics if s[i:j] is not a valid slice of s.
func (s ComparableSlice[E]) Slice(i, j int) ComparableSlice[E] {
	return s[i:j]
}

// Slice3 is the 3-index form of the slice operation: it returns s[i:j:k].
// It panics if s[i:j:k] is not a valid slice of s.
func (s ComparableSlice[E]) Slice3(i, j, k int) ComparableSlice[E] {
	return s[i:j:k]
}

// Append appends the values x to s and updates the slice s.
func (s *ComparableSlice[E]) Append(x ...E) {
	*s = append(*s, x...)
}

// AppendSlice appends a slice t to s and updates the slice s.
func (s *ComparableSlice[E]) AppendSlice(t []E) {
	*s = append(*s, t...)
}

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in increasing index order, and the
// comparison stops at the first unequal pair.
// nil slices and empty non-nil slices are considered equal.
// Floating point NaNs are not considered equal.
func Equal[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}

	return true
}

// EqualFunc is like Equal but using a comparison function on each pair of elements.
func EqualFunc[E1, E2 any](s1 []E1, s2 []E2, eq func(E1, E2) bool) bool {
	if len(s1) != len(s2) {
		return false
	}

	for i := range s1 {
		if !eq(s1[i], s2[i]) {
			return false
		}
	}
	// for i, v1 := range s1 {
	// 	v2 := s2[i]
	// 	if !eq(v1, v2) {
	// 		return false
	// 	}
	// }

	return true
}

// Compare compares the elements of s1 and s2.
// The elements are compared sequentially, starting at index 0,
// until one element is not equal to the other.
// The result of comparing the first non-matching elements is returned.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one.
// The result is 0 if s1 == s2, -1 if s1 < s2, and +1 if s1 > s2.
// Comparisons involving floating point NaNs are ignored.
func Compare[E constraints.Ordered](s1, s2 []E) int {
	l2 := len(s2)

	for i, v1 := range s1 {
		if i >= l2 {
			return +1
		}

		v2 := s2[i]
		switch {
		case v1 < v2:
			return -1
		case v2 < v1:
			return +1
		}
	}

	if len(s1) < l2 {
		return -1
	}

	return 0
}

// CompareFunc is like Compare but uses a comparison function on each pair of elements.
// The elements are compared in increasing index order, and the comparisons stop after the first time cmp returns non-zero.
// The result is the first non-zero result of cmp;
// if cmp always returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2), and +1 if len(s1) > len(s2).
func CompareFunc[E1, E2 any](s1 []E1, s2 []E2, cmp func(E1, E2) int) int {
	l2 := len(s2)

	for i, v1 := range s1 {
		if i >= l2 {
			return +1
		}

		v2 := s2[i]
		if r := cmp(v1, v2); r != 0 {
			return r
		}
	}

	if len(s1) < l2 {
		return -1
	}

	return 0
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}

	return -1
}

// IndexFunc returns the first index i satisfying f(s[i]),
// or -1 if none do.
func IndexFunc[E any](s []E, f func(E) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}

	return -1
}

// Contains reports whether v is present in s.
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

// ContainsFunc reports whether at least one
// element e of s satisfies f(e).
func ContainsFunc[E any](s []E, f func(E) bool) bool {
	return IndexFunc(s, f) >= 0
}

// Insert inserts the values v... into s at index i,
// returning the modified slice.
// In the returned slice r, r[i] == v[0].
// Insert panics if i is out of range.
// This function is O(len(s) + len(v)).
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	tot := len(s) + len(v)
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[i:])
		copy(s2[i:], v)
		return s2
	}

	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[i:])
	return s2
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-j), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
// Delete might not modify the elements s[len(s)-(j-i):len(s)]. If those
// elements is pointers or contain pointers, Delete zeroing those elements so that
// objects they reference can be garbage collected.
func Delete[S ~[]E, E any](s S, i, j int) S {
	_ = s[i:j] // bounds check

	s2 := append(s[:i], s[j:]...)

	if containsPointer(*new(E)) {
		_ = append([]E(s[len(s)-(j-i):]), make([]E, j-i)...)
	}

	return s2
}

// reports whether a is a pointer or contains pointers.
func containsPointer(a any) bool {
	t := reflect.TypeOf(a)
	if t.Kind() == reflect.Pointer {
		return true
	}

	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := reflect.Zero(t.Field(i).Type).Interface()
			if containsPointer(f) {
				return true
			}
		}
	}

	return false
}

// Replace replaces the elements s[i:j] by the given v, and returns the
// modified slice. Replace panics if s[i:j] is not a valid slice of s.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	_ = s[i:j] // verify that i:j is a valid subslice

	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot <= cap(s) {
		s2 := s[:tot]
		copy(s2[i+len(v):], s[j:])
		copy(s2[i:], v)
		return s2
	}

	s2 := make(S, tot)
	copy(s2, s[:i])
	copy(s2[i:], v)
	copy(s2[i+len(v):], s[j:])
	return s2
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}

	return append(S{}, s...)
}

// DeepClone is like Clone but if elements has an Clone method of the form "Clone() E",
// it use the result of e.Clone() as assignment right operand.
func DeepClone[S ~[]E, E any](s S) S {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}

	type cloneable interface {
		Clone() E
	}

	var v any = *new(E)
	if _, ok := v.(cloneable); ok {
		s2 := make(S, len(s))

		var i int
		for i, v = range s {
			s2[i] = v.(cloneable).Clone()
		}

		return s2
	}

	return Clone(s)
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
// When Compact discards m elements in total, it might not modify the elements
// s[len(s)-m:len(s)]. If those elements is pointers or contain pointers, Compact
// zeroing those elements so that objects they reference can be garbage collected.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) < 2 {
		return s
	}

	i := 1
	for j := 1; j < len(s); j++ {
		if s[j-1] != s[j] {
			if j != i {
				s[i] = s[j]
			}

			i++
		}
	}

	if containsPointer(*new(E)) {
		_ = append(s[i:], make([]E, len(s)-i)...)
	}

	return s[:i]
}

// CompactFunc is like Compact but uses a comparison function.
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	if len(s) < 2 {
		return s
	}

	i := 1
	for j := 1; j < len(s); j++ {
		if !eq(s[j-1], s[j]) {
			if j != i {
				s[i] = s[j]
			}

			i++
		}
	}

	if containsPointer(*new(E)) {
		_ = append(s[i:], make([]E, len(s)-i)...)
	}

	return s[:i]
}

// Grow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
func Grow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}

	if n -= cap(s) - len(s); n > 0 {
		s = append(s[:cap(s)], make([]E, n)...)[:len(s)]
	}

	return s
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}
