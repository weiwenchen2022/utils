// Package types defines various functions useful with any type.
package types

import "reflect"

// Zero returns a zero value for the specified type.
func Zero[T any]() (v T) {
	return v
}

// IsZero reports whether v is the zero value for its type.
func IsZero[T any](v T) bool {
	return reflect.ValueOf(v).IsZero()
}

// Interface returns v as an any. It is equivalent to:
// var a any = v
func Interface[T any](v T) (a any) {
	return v
}

// Convert returns the value v converted to type T2.
// If the usual Go conversion rules do not allow conversion of the value v to type T2,
// or if converting v to type T2 panics,
// or if the result value was obtained by accessing unexported struct fields, Convert panics.
func Convert[T2, T1 any](v T1) T2 {
	typeOfT2 := reflect.TypeOf(*new(T2))
	return reflect.ValueOf(v).Convert(typeOfT2).Interface().(T2)
}

// CanConvert reports whether the value v can be converted to type T2.
// If CanConvert[T2](v) returns true then Convert[T2](v) will not panic.
func CanConvert[T2, T1 any](v T1) bool {
	v2 := reflect.ValueOf(*new(T2))
	return reflect.ValueOf(v).CanConvert(v2.Type()) && v2.CanInterface()
}

// ToSliceOfAny returns a new slice of any with elements of the slice s.
func ToSliceOfAny[S ~[]E, E any](s S) []any {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}

	r := make([]any, len(s))
	for i, v := range s {
		r[i] = v
	}
	return r
}

// FromSliceOfAny returns a slice of E with elements of the slice s.
// If type of the underlying value is not E, FromSliceOfAny panics.
func FromSliceOfAny[E any](s []any) []E {
	// Preserve nil in case it matters.
	if s == nil {
		return nil
	}

	r := make([]E, len(s))
	for i, v := range s {
		r[i] = v.(E)
	}
	return r
}
