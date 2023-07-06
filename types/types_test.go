package types_test

import (
	"fmt"
	"testing"

	. "github.com/weiwenchen2022/utils/types"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type T struct{ A int }

	if got, want := New(23), 23; want != *got {
		t.Errorf("New(%#v) = %#v, want %#[1]v", want, *got)
	}
	if got, want := New("foo"), "foo"; want != *got {
		t.Errorf("New(%#v) = %#v, want %#[1]v", want, *got)
	}
	if got, want := New(T{23}), (T{23}); want != *got {
		t.Errorf("New(%#v) = %#v, want %#[1]v", want, *got)
	}
}

func TestZero(t *testing.T) {
	t.Parallel()

	type T struct{ A int }

	if got, want := Zero[int](), 0; want != got {
		t.Errorf("Zero[int]() = %v, want %v", got, want)
	}
	if got, want := Zero[string](), ""; want != got {
		t.Errorf("Zero[string]() = %v, want %v", got, want)
	}
	if got, want := Zero[T](), (T{}); want != got {
		t.Errorf("Zero[T]() = %v, want %v", got, want)
	}
	if got, want := Zero[chan int](), (chan int)(nil); want != got {
		t.Errorf("Zero[chan int]() = %v, want %v", got, want)
	}
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	type T struct{ A int }
	tests := []struct {
		v    any
		want bool
	}{
		{0, true},
		{23, false},
		{"", true},
		{"foo", false},
		{T{}, true},
		{T{23}, false},
	}

	for _, tt := range tests {
		if got := IsZero(tt.v); tt.want != got {
			t.Errorf("IsZero(%#v) = %t, want %t", tt.v, got, tt.want)
		}
	}
}

func TestInterface(t *testing.T) {
	t.Parallel()

	type T struct{ A int }
	tests := []struct {
		v any
	}{
		{0},
		{23},
		{""},
		{"foo"},
		{T{23}},
	}

	for _, tt := range tests {
		if got := Interface(tt.v); tt.v != got {
			t.Errorf("Interface(%#v) = %#v, want %#[1]v", tt.v, got)
		}
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	if got, want := Convert[int8](0), int8(0); want != got {
		t.Errorf("Convert[int8](0) = %v, want %v", got, want)
	}
	if got, want := Convert[int8](42), int8(42); want != got {
		t.Errorf("Convert[int8](42) = %v, want %v", got, want)
	}
	if got, want := Convert[[]byte](""), []byte(""); string(want) != string(got) {
		t.Errorf(`Convert[[]byte]("") = %s, want %s`, got, want)
	}
	if got, want := Convert[[]byte]("foo"), []byte("foo"); string(want) != string(got) {
		t.Errorf(`Convert[[]byte]("foo") = %s, want %s`, got, want)
	}
}

func TestCanConvert(t *testing.T) {
	t.Parallel()

	if got, want := CanConvert[int8](0), true; want != got {
		t.Errorf("CanConvert[int8](0) = %v, want %v", got, want)
	}
	if got, want := CanConvert[string](0), true; want != got {
		t.Errorf("CanConvert[string](0) = %v, want %v", got, want)
	}
	if got, want := CanConvert[int](""), false; want != got {
		t.Errorf(`CanConvert[int]("") = %v, want %v`, got, want)
	}
	if got, want := CanConvert[[]byte]("foo"), true; want != got {
		t.Errorf(`Convert[[]byte]("foo") = %v, want %v`, got, want)
	}
}

func TestToSliceOfAny(t *testing.T) {
	t.Parallel()

	if got, want := ToSliceOfAny([]int{0, 1, 2, 3}), []any{0, 1, 2, 3}; fmt.Sprint(want) != fmt.Sprint(got) {
		t.Errorf("ToSliceOfAny(%v) = %v, want %[1]v", want, got)
	}

	if got, want := ToSliceOfAny([]int(nil)), []any(nil); fmt.Sprintf("%#v", want) != fmt.Sprintf("%#v", got) {
		t.Errorf("ToSliceOfAny(%#v) = %#v, want %#[1]v", want, got)
	}
}

func TestFromSliceOfAny(t *testing.T) {
	t.Parallel()

	if got, want := FromSliceOfAny[int]([]any{0, 1, 2, 3}), []int{0, 1, 2, 3}; fmt.Sprint(want) != fmt.Sprint(got) {
		t.Errorf("FromSliceOfAny(%v) = %v, want %[1]v", want, got)
	}

	if got, want := FromSliceOfAny[int]([]any(nil)), []int(nil); fmt.Sprintf("%#v", want) != fmt.Sprintf("%#v", got) {
		t.Errorf("FromSliceOfAny(%#v) = %#v, want %#[1]v", want, got)
	}
}
