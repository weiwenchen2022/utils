// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package maps defines various types and functions useful with maps of any type.
package maps

// Map attaches the common methods to map[K]V
type Map[K comparable, V any] map[K]V

// NewMap converts m to Map type.
func NewMap[M ~map[K]V, K comparable, V any](m M) Map[K, V] {
	return Map[K, V](m)
}

// Keys is a convenience method: m.Keys() returns Keys(m).
func (m Map[K, V]) Keys() []K {
	return Keys(m)
}

// Values is a convenience method: m.Values() returns Values(m).
func (m Map[K, V]) Values() []V {
	return Values(m)
}

// EqualFunc is a convenience method: m.EqualFunc(m2, eq) returns EqualFunc(m, m2, eq).
func (m Map[K, V]) EqualFunc(m2 map[K]V, eq func(V, V) bool) bool {
	return EqualFunc(m, m2, eq)
}

// Clear is a convenience method: m.Clear() calls Clear(m).
func (m Map[K, V]) Clear() {
	Clear(m)
}

// Clone is a convenience method: m.Clone() returns Clone(m).
func (m Map[K, V]) Clone() Map[K, V] {
	return Clone(m)
}

// Copy is a convenience method: m.Copy(srt) calls Copy(m, src).
func (m Map[K, V]) Copy(src map[K]V) {
	Copy(m, src)
}

// DeleteFunc is a convenience method: m.DeleteFunc(del) calls DeleteFunc(m, del).
func (m Map[K, V]) DeleteFunc(del func(K, V) bool) {
	DeleteFunc(m, del)
}

// ComparableMap is like Map but values requires comparable
type ComparableMap[K, V comparable] map[K]V

// NewComparableMap converts m to ComparableMap type.
func NewComparableMap[M ~map[K]V, K, V comparable](m M) ComparableMap[K, V] {
	return ComparableMap[K, V](m)
}

// Keys is a convenience method: m.Keys() returns Keys(m).
func (m ComparableMap[K, V]) Keys() []K {
	return Keys(m)
}

// Values is a convenience method: m.Values() returns Values(m).
func (m ComparableMap[K, V]) Values() []V {
	return Values(m)
}

// Equal is a convenience method: m.Equal(m2) returns Equal(m, m2).
func (m ComparableMap[K, V]) Equal(m2 map[K]V) bool {
	return Equal(m, m2)
}

// EqualFunc is a convenience method: m.EqualFunc(m2, eq) returns EqualFunc(m, m2, eq).
func (m ComparableMap[K, V]) EqualFunc(m2 map[K]V, eq func(V, V) bool) bool {
	return EqualFunc(m, m2, eq)
}

// Clear is a convenience method: m.Clear() calls Clear(m).
func (m ComparableMap[K, V]) Clear() {
	Clear(m)
}

// Clone is a convenience method: m.Clone() returns Clone(m).
func (m ComparableMap[K, V]) Clone() ComparableMap[K, V] {
	return Clone(m)
}

// Copy is a convenience method: m.Copy(srt) calls Copy(m, src).
func (m ComparableMap[K, V]) Copy(src map[K]V) {
	Copy(m, src)
}

// DeleteFunc is a convenience method: m.DeleteFunc(del) calls DeleteFunc(m, del).
func (m ComparableMap[K, V]) DeleteFunc(del func(K, V) bool) {
	DeleteFunc(m, del)
}

// Keys returns the keys of the map m.
// The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// Values returns the values of the map m.
// The values will be in an indeterminate order.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))

	for _, v := range m {
		values = append(values, v)
	}

	return values
}

// Equal reports whether two maps contain the same key/value pairs.
// Values are compared using ==.
func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}

	return true
}

// EqualFunc is like Equal, but compares values using eq.
// Keys are still compared with ==.
func EqualFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool {
	if len(m1) != len(m2) {
		return false
	}

	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !eq(v1, v2) {
			return false
		}
	}

	return true
}

// Clear removes all entries from m, leaving it empty.
func Clear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}

// Clone returns a copy of m. This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	// Preserve nil in case it matters.
	if m == nil {
		return nil
	}

	copy := make(M, len(m))
	for k, v := range m {
		copy[k] = v
	}
	return copy
}

// Copy copies all key/value pairs in src adding them to dst.
// When a key in src is already present in dst,
// the value in dst will be overwritten by the value associated
// with the key in src.
func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

// DeleteFunc deletes any key/value pairs from m for which del returns true.
func DeleteFunc[M ~map[K]V, K comparable, V any](m M, del func(K, V) bool) {
	for k, v := range m {
		if del(k, v) {
			delete(m, k)
		}
	}
}
