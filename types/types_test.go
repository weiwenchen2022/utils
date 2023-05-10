package types_test

import (
	"fmt"
	"testing"

	"github.com/weiwenchen2022/utils/types"
)

func TestZero(t *testing.T) {
	t.Parallel()

	if got, want := types.Zero[int](), 0; want != got {
		t.Errorf("Zero[int]() = %v, want %v", got, want)
	}
	if got, want := types.Zero[string](), ""; want != got {
		t.Errorf("Zero[string]() = %v, want %v", got, want)
	}
	if got, want := types.Zero[struct{}](), struct{}{}; want != got {
		t.Errorf("Zero[struct{}]() = %v, want %v", got, want)
	}
	if got, want := types.Zero[chan int](), (chan int)(nil); want != got {
		t.Errorf("Zero[chan int]() = %v, want %v", got, want)
	}
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	if got, want := types.IsZero(0), true; want != got {
		t.Errorf("IsZero(0) = %v, want %v", got, want)
	}
	if got, want := types.IsZero(42), false; want != got {
		t.Errorf("IsZero(42) = %v, want %v", got, want)
	}
	if got, want := types.IsZero(""), true; want != got {
		t.Errorf("IsZero(\"\") = %v, want %v", got, want)
	}
	if got, want := types.IsZero("foo"), false; want != got {
		t.Errorf("IsZero(\"foo\") = %v, want %v", got, want)
	}
	if got, want := types.IsZero(struct{ foo string }{}), true; want != got {
		t.Errorf("IsZero(struct{ foo string }{}) = %v, want %v", got, want)
	}
	if got, want := types.IsZero(struct{ foo string }{"bar"}), false; want != got {
		t.Errorf("IsZero(struct{ foo string }{\"bar\"}) = %v, want %v", got, want)
	}
}

func TestInterface(t *testing.T) {
	t.Parallel()

	if got, want := types.Interface(0), any(0); want != got {
		t.Errorf("Interface(0) = %v, want %v", got, want)
	}
	if got, want := types.Interface(42), any(42); want != got {
		t.Errorf("Interface(42) = %v, want %v", got, want)
	}
	if got, want := types.Interface(""), any(""); want != got {
		t.Errorf("Interface(\"\") = %v, want %v", got, want)
	}
	if got, want := types.Interface("foo"), any("foo"); want != got {
		t.Errorf("Interface(\"foo\") = %v, want %v", got, want)
	}

	var want any = struct{ foo string }{}
	if got := types.Interface(want); want != got {
		t.Errorf("Interface(%v) = %v, want %[1]v", want, got)
	}

	want = struct{ foo string }{"bar"}
	if got := types.Interface(want); want != got {
		t.Errorf("Interface(%v) = %v, want %[1]v", want, got)
	}
}

func TestConvert(t *testing.T) {
	t.Parallel()

	if got, want := types.Convert[int8](0), int8(0); want != got {
		t.Errorf("Convert[int8](0) = %v, want %v", got, want)
	}
	if got, want := types.Convert[int8](42), int8(42); want != got {
		t.Errorf("Convert[int8](42) = %v, want %v", got, want)
	}
	if got, want := types.Convert[[]byte](""), []byte(""); string(want) != string(got) {
		t.Errorf("Convert[[]byte](\"\") = %s, want %s", got, want)
	}
	if got, want := types.Convert[[]byte]("foo"), []byte("foo"); string(want) != string(got) {
		t.Errorf("Convert[[]byte](\"foo\") = %s, want %s", got, want)
	}
}

func TestCanConvert(t *testing.T) {
	t.Parallel()

	if got, want := types.CanConvert[int8](0), true; want != got {
		t.Errorf("CanConvert[int8](0) = %v, want %v", got, want)
	}
	if got, want := types.CanConvert[string](0), true; want != got {
		t.Errorf("CanConvert[string](0) = %v, want %v", got, want)
	}
	if got, want := types.CanConvert[int](""), false; want != got {
		t.Errorf("CanConvert[int](\"\") = %v, want %v", got, want)
	}
	if got, want := types.CanConvert[[]byte]("foo"), true; want != got {
		t.Errorf("Convert[[]byte](\"foo\") = %v, want %v", got, want)
	}
}

func TestToSliceOfAny(t *testing.T) {
	t.Parallel()

	if got, want := types.ToSliceOfAny([]int{0, 1, 2, 3}), []any{0, 1, 2, 3}; fmt.Sprint(want) != fmt.Sprint(got) {
		t.Errorf("ToSliceOfAny(%v) = %v, want %[1]v", want, got)
	}

	if got, want := types.ToSliceOfAny([]int(nil)), []any(nil); fmt.Sprintf("%#v", want) != fmt.Sprintf("%#v", got) {
		t.Errorf("ToSliceOfAny(%#v) = %#v, want %#[1]v", want, got)
	}
}

func TestFromSliceOfAny(t *testing.T) {
	t.Parallel()

	if got, want := types.FromSliceOfAny[int]([]any{0, 1, 2, 3}), []int{0, 1, 2, 3}; fmt.Sprint(want) != fmt.Sprint(got) {
		t.Errorf("FromSliceOfAny(%v) = %v, want %[1]v", want, got)
	}

	if got, want := types.FromSliceOfAny[int]([]any(nil)), []int(nil); fmt.Sprintf("%#v", want) != fmt.Sprintf("%#v", got) {
		t.Errorf("FromSliceOfAny(%#v) = %#v, want %#[1]v", want, got)
	}
}
