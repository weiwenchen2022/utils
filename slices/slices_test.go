// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"math"
	"strings"
	"testing"

	. "github.com/weiwenchen2022/utils/slices"

	"golang.org/x/exp/constraints"
)

var raceEnabled bool

func TestSlice_EqualFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range equalIntTests {
		s1 := NewSlice(tc.s1)
		if got := s1.EqualFunc(tc.s2, equal[int]); tc.want != got {
			t.Errorf("%v.EqualFunc(%v, equal[int]) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}
	}

	for _, tc := range equalFloatTests {
		s1 := NewSlice(tc.s1)
		if got := s1.EqualFunc(tc.s2, equal[float64]); tc.wantEqual != got {
			t.Errorf("%v.EqualFunc(%v, equal[float64]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}
		if got := s1.EqualFunc(tc.s2, equalNaN[float64]); tc.wantEqualNaN != got {
			t.Errorf("%v.EqualFunc(%v, equalNaN[float64]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqualNaN)
		}
	}

	s1 := NewSlice([]int{1, 2, 3})
	s2 := []int{2, 3, 4}
	if s1.EqualFunc(*s1, offByOne[int]) {
		t.Errorf("%v.EqualFunc(%v, offByOne) = true, want false", *s1, *s1)
	}
	if !s1.EqualFunc(s2, offByOne[int]) {
		t.Errorf("%v.EqualFunc(%v, offByOne) = false, want true", *s1, s2)
	}

	s3 := NewSlice([]string{"a", "b", "c"})
	s4 := []string{"A", "B", "C"}
	if !s3.EqualFunc(s4, strings.EqualFold) {
		t.Errorf("%v.EqualFunc(%v, strings.EqualFold) = false, want true", *s3, s4)
	}
}

func TestSlice_CompareFunc(t *testing.T) {
	t.Parallel()

	intWant := func(want bool) string {
		if want {
			return "0"
		}

		return "!= 0"
	}

	for _, test := range equalIntTests {
		s1 := NewSlice(test.s1)
		if got := s1.CompareFunc(test.s2, equalToCmp(equal[int])); test.want != (got == 0) {
			t.Errorf("%v.CompareFunc(%v, equalToCmp(equal[int])) = %d, want %s", test.s1, test.s2, got, intWant(test.want))
		}
	}

	for _, test := range equalFloatTests {
		s1 := NewSlice(test.s1)
		if got := s1.CompareFunc(test.s2, equalToCmp(equal[float64])); test.wantEqual != (got == 0) {
			t.Errorf("%v.CompareFunc(%v, equalToCmp(equal[float64])) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqual))
		}
	}

	for _, test := range compareIntTests {
		s1 := NewSlice(test.s1)
		if got := s1.CompareFunc(test.s2, cmp[int]); test.want != got {
			t.Errorf("%v.CompareFunc(%v, cmp[int]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range compareFloatTests {
		s1 := NewSlice(test.s1)
		if got := s1.CompareFunc(test.s2, cmp[float64]); test.want != got {
			t.Errorf("%v.CompareFunc(%v, cmp[float64]) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	s1 := NewSlice([]int{1, 2, 3})
	s2 := []int{2, 3, 4}
	if got := s1.CompareFunc(s2, equalToCmp(offByOne[int])); got != 0 {
		t.Errorf("%v.CompareFunc(%v, offByOne) = %d, want 0", *s1, s2, got)
	}

	s3 := NewSlice([]string{"a", "b", "c"})
	s4 := []string{"A", "B", "C"}
	if got := s3.CompareFunc(s4, strings.Compare); got != 1 {
		t.Errorf("%v.CompareFunc(%v, strings.Compare) = %d, want 1", *s3, s4, got)
	}

	compareLower := func(v1, v2 string) int {
		return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
	}
	if got := s3.CompareFunc(s4, compareLower); got != 0 {
		t.Errorf("%v.CompareFunc(%v, compareLower) = %d, want 0", *s3, s4, got)
	}
}

func TestSlice_IndexFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		s := NewSlice(tc.s)
		if got := s.IndexFunc(equalToIndex(equal[int], tc.v)); tc.want != got {
			t.Errorf("%v.IndexFunc(equalToIndex(equal[int], %v)) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
	}

	s1 := NewSlice([]string{"hi", "HI"})
	if got := s1.IndexFunc(equalToIndex(equal[string], "HI")); got != 1 {
		t.Errorf("%v.IndexFunc(equalToIndex(equal[string], %q)) = %d, want %d", *s1, "HI", got, 1)
	}
	if got := s1.IndexFunc(equalToIndex(strings.EqualFold, "HI")); got != 0 {
		t.Errorf("%v.IndexFunc(equalToIndex(strings.EqualFold, %q)) = %d, want %d", *s1, "HI", got, 0)
	}
}

func TestSlice_ContainsFunc(t *testing.T) {
	t.Parallel()

	for _, test := range indexTests {
		s := NewSlice(test.s)
		if got := s.ContainsFunc(equalToIndex(equal[int], test.v)); (test.want != -1) != got {
			t.Errorf("%v.ContainsFunc(equalToIndex(equal[int], %v)) = %t, want %t", test.s, test.v, got, test.want != -1)
		}
	}

	s1 := NewSlice([]string{"hi", "HI"})
	if got := s1.ContainsFunc(equalToIndex(equal[string], "HI")); !got {
		t.Errorf("%v.ContainsFunc(equalToContains(equal[string], %q)) = %t, want %t", *s1, "HI", got, true)
	}
	if got := s1.ContainsFunc(equalToIndex(equal[string], "hI")); got {
		t.Errorf("%v.ContainsFunc(equalToContains(equal[string], %q)) = %t, want %t", *s1, "hI", got, false)
	}
	if got := s1.ContainsFunc(equalToIndex(strings.EqualFold, "hI")); !got {
		t.Errorf("%v.ContainsFunc(equalToContains(strings.EqualFold, %q)) = %t, want %t", *s1, "hI", got, true)
	}
}

func TestSlice_Insert(t *testing.T) {
	t.Parallel()

	s := NewSlice([]int{1, 2, 3})
	want := []int{1, 2, 3}
	s.Insert(0)
	if !Equal(want, *s) {
		t.Errorf("%[2]v.Insert(0) = %[1]v, want %[2]v", s, want)
	}

	for _, tc := range insertTests {
		copy := NewSlice(Clone(tc.s))
		copy.Insert(tc.i, tc.add...)
		if !Equal(tc.want, *copy) {
			t.Errorf("%v.Insert(%d, %v...) = %v, want %v", tc.s, tc.i, tc.add, *copy, tc.want)
		}
	}
}

func TestSlice_Delete(t *testing.T) {
	t.Parallel()

	for _, tc := range deleteTests {
		copy := NewSlice(Clone(tc.s))
		copy.Delete(tc.i, tc.j)
		if !Equal(tc.want, *copy) {
			t.Errorf("%v.Delete(%d, %d) = %v, want %v", tc.s, tc.i, tc.j, *copy, tc.want)
		}
	}

	for _, tc := range deletePanicsTests {
		if !panics(func() { NewSlice(tc.s).Delete(tc.i, tc.j) }) {
			t.Errorf("Delete %s: got no panic, want panic", tc.name)
		}
	}
}

func TestSlic_Replace(t *testing.T) {
	t.Parallel()

	for _, tc := range replaceTests {
		ss, s := Clone(tc.s), NewSlice(Clone(tc.s))
		want := naiveReplace(ss, tc.i, tc.j, tc.v...)
		s.Replace(tc.i, tc.j, tc.v...)
		if !Equal(want, *s) {
			t.Errorf("%v.Replace(%v, %v, %v...) = %v, want %v", tc.s, tc.i, tc.j, tc.v, *s, want)
		}
	}

	for _, tc := range replacePanicsTests {
		ss := NewSlice(Clone(tc.s))
		if !panics(func() { ss.Replace(tc.i, tc.j, tc.v...) }) {
			t.Errorf("Replace %s: should have panicked", tc.name)
		}
	}
}

func TestSlice_Clone(t *testing.T) {
	t.Parallel()

	s1 := NewSlice([]int{1, 2, 3})
	s2 := s1.Clone()
	if !Equal(*s1, *s2) {
		t.Errorf("%v.Clone() = %v, want %[1]v", *s1, *s2)
	}

	(*s1)[0] = 4
	want := []int{1, 2, 3}
	if !Equal(*s2, want) {
		t.Errorf("%v.Clone() changed unexpectedly to %v", want, *s2)
	}

	if got := NewSlice([]int(nil)).Clone(); !got.IsNil() {
		t.Errorf("nil.Clone() = %#v, want nil", *got)
	}

	if got := (*s1)[:0].Clone(); got.IsNil() || got.Len() != 0 {
		t.Errorf("%v.Clone() = %#v, want %#[1]v", (*s1)[:0], *got)
	}
}

func TestSlice_DeepClone(t *testing.T) {
	t.Parallel()

	s1 := NewSlice([]*foo{{1}, {2}, {3}})
	s2 := s1.DeepClone()
	if !EqualFunc(*s1, *s2, equalStructPointer[foo]) {
		t.Errorf("%v.DeepClone() = %v, want %[1]v", s1, s2)
	}

	(*s1)[0] = &foo{4}
	want := []*foo{{1}, {2}, {3}}
	if !EqualFunc(want, *s2, equalStructPointer[foo]) {
		t.Errorf("%v.DeepClone() changed unexpectedly to %v", want, *s2)
	}

	if got := NewSlice([]*foo(nil)).DeepClone(); !got.IsNil() {
		t.Errorf("nil.DeepClone() = %#v, want nil", *got)
	}

	if got := (*s1)[:0].DeepClone(); got.IsNil() || got.Len() != 0 {
		t.Errorf("%v.DeepClone() = %#v, want %#[1]v", (*s1)[:0], *got)
	}
}

func TestSlice_CompactFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range compactTests {
		copy := NewSlice(Clone(tc.s))
		copy.CompactFunc(equal[int])
		if !Equal(tc.want, *copy) {
			t.Errorf("%v.CompactFunc(equal[int]) = %v, want %v", tc.s, *copy, tc.want)
		}
	}

	s1 := NewSlice([]string{"a", "a", "A", "B", "b"})
	copy := s1.Clone()
	want := []string{"a", "B"}
	copy.CompactFunc(strings.EqualFold)
	if !Equal(want, *copy) {
		t.Errorf("%v.CompactFunc(strings.EqualFold) = %v, want %v", *s1, *copy, want)
	}
}

func TestSlice_Grow(t *testing.T) {
	t.Parallel()

	s1 := NewSlice([]int{1, 2, 3})
	copy := s1.Clone()
	copy.Grow(1000)
	s2 := *copy
	if !Equal(*s1, s2) {
		t.Errorf("%v.Grow() = %v, want %[1]v", *s1, s2)
	}
	if s2.Cap() < 1000+s1.Len() {
		t.Errorf("after %v.Grow() cap = %d, want >= %d", *s1, s2.Cap(), 1000+s1.Len())
	}

	// Test mutation of elements between length and capacity.
	copy = s1.Clone()
	s3 := (*copy)[:1]
	s3.Grow(2)
	s3 = s3[:3]
	if !Equal(*s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}

	s3 = (*copy)[:1]
	s3.Grow(1000)
	s3 = s3[:3]
	if !Equal(*s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}

	// Test number of allocations.
	if n := testing.AllocsPerRun(100, func() {
		saved := s2
		s2.Grow(s2.Cap() - s2.Len())
		s2 = saved
	}); n != 0 {
		t.Errorf("Grow should not allocate when given sufficient capacity; allocated %v times", n)
	}

	if n := testing.AllocsPerRun(100, func() {
		saved := s2
		s2.Grow(s2.Cap() - s2.Len() + 1)
		s2 = saved
	}); n != 1 {
		errorf := t.Errorf
		if raceEnabled {
			errorf = t.Logf // this allocates multiple times in race detector mode
		}
		errorf("Grow should allocate once when given insufficient capacity; allocated %v times", n)
	}

	// Test for negative growth sizes.
	if !panics(func() { s1.Grow(-1) }) {
		t.Errorf("Grow(-1) did not panic; expected a panic")
	}
}

func TestSlice_Clip(t *testing.T) {
	t.Parallel()

	s1 := NewSlice([]int{1, 2, 3, 4, 5, 6}).Slice(0, 3)
	orig := s1.Clone()
	if s1.Len() != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, s1.Len())
	}

	if s1.Cap() < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", *orig, s1.Cap())
	}

	s2 := s1
	s2.Clip()
	if !Equal(s1, s2) {
		t.Errorf("%v.Clip() = %v, want %[1]v", s1, s2)
	}
	if s2.Cap() != 3 {
		t.Errorf("cap(%v.Clip()) = %d, want 3", *orig, s2.Cap())
	}
}

func TestComparableSlice_Equal(t *testing.T) {
	t.Parallel()

	for _, tc := range equalIntTests {
		s1 := NewComparableSlice(tc.s1)
		if got := s1.Equal(tc.s2); tc.want != got {
			t.Errorf("%v.Equal(%v) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}
	}

	for _, tc := range equalFloatTests {
		s1 := NewComparableSlice(tc.s1)
		if got := s1.Equal(tc.s2); tc.wantEqual != got {
			t.Errorf("%v.Equal(%v) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}
	}
}

func TestComparableSlice_Index(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		s := NewComparableSlice(tc.s)
		if got := s.Index(tc.v); tc.want != got {
			t.Errorf("%v.Index(%v) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
	}
}

func TestComparableSlice_Contains(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		s := NewComparableSlice(tc.s)
		if got := s.Contains(tc.v); (tc.want != -1) != got {
			t.Errorf("%v.Contains(%v) = %t, want %t", tc.s, tc.v, got, tc.want != -1)
		}
	}
}

func TestComparableSlice_Compact(t *testing.T) {
	t.Parallel()

	for _, tc := range compactTests {
		copy := NewComparableSlice(Clone(tc.s))
		copy.Compact()
		if !Equal(tc.want, *copy) {
			t.Errorf("%v.Compact() = %v, want %v", tc.s, *copy, tc.want)
		}
	}
}

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

func TestEqual(t *testing.T) {
	t.Parallel()

	for _, tc := range equalIntTests {
		if got := Equal(tc.s1, tc.s2); tc.want != got {
			t.Errorf("Equal(%v, %v) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}

		if got := equalTo(tc.s1, tc.s2); tc.want != got {
			t.Errorf("equalTo(%v, %v) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}
	}

	for _, tc := range equalFloatTests {
		if got := Equal(tc.s1, tc.s2); tc.wantEqual != got {
			t.Errorf("Equal(%v, %v) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}

		if got := equalTo(tc.s1, tc.s2); tc.wantEqual != got {
			t.Errorf("equalTo(%v, %v) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}
	}
}

// equal is simply ==.
func equal[T comparable](v1, v2 T) bool {
	return v1 == v2
}

// equalNaN is like == except that all NaNs are equal.
func equalNaN[T comparable](v1, v2 T) bool {
	isNaN := func(f T) bool { return f != f }
	return v1 == v2 || (isNaN(v1) && isNaN(v2))
}

// offByOne returns true if integers v1 and v2 differ by 1.
func offByOne[E constraints.Integer](v1, v2 E) bool {
	return v1 == v2+1 || v1 == v2-1
}

func TestEqualFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range equalIntTests {
		if got := EqualFunc(tc.s1, tc.s2, equal[int]); tc.want != got {
			t.Errorf("EqualFunc(%v, %v, equal[int]) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}

		if got := equalFunc(tc.s1, tc.s2, equal[any]); tc.want != got {
			t.Errorf("equalFunc(%v, %v, equal[any]) = %t, want %t", tc.s1, tc.s2, got, tc.want)
		}
	}

	for _, tc := range equalFloatTests {
		if got := EqualFunc(tc.s1, tc.s2, equal[float64]); tc.wantEqual != got {
			t.Errorf("EqualFunc(%v, %v, equal[float64]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}
		if got := equalFunc(tc.s1, tc.s2, equal[any]); tc.wantEqual != got {
			t.Errorf("equalFunc(%v, %v, equal[any]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqual)
		}

		if got := EqualFunc(tc.s1, tc.s2, equalNaN[float64]); tc.wantEqualNaN != got {
			t.Errorf("EqualFunc(%v, %v, equalNaN[float64]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqualNaN)
		}
		if got := equalFunc(tc.s1, tc.s2, equalNaN[any]); tc.wantEqualNaN != got {
			t.Errorf("equalFunc(%v, %v, equalNaN[float64]) = %t, want %t", tc.s1, tc.s2, got, tc.wantEqualNaN)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if EqualFunc(s1, s1, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}
	if equalFunc(s1, s1, func(a1, a2 any) bool { return offByOne(a1.(int), a2.(int)) }) {
		t.Errorf("equalFunc(%v, %v, offByOne) = true, want false", s1, s1)
	}

	if !EqualFunc(s1, s2, offByOne[int]) {
		t.Errorf("EqualFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}
	if !equalFunc(s1, s2, func(a1, a2 any) bool { return offByOne(a1.(int), a2.(int)) }) {
		t.Errorf("equalFunc(%v, %v, offByOne) = false, want true", s1, s2)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if !EqualFunc(s3, s4, strings.EqualFold) {
		t.Errorf("EqualFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}
	if !equalFunc(s3, s4, func(a1, a2 any) bool { return strings.EqualFold(a1.(string), a2.(string)) }) {
		t.Errorf("equalFunc(%v, %v, strings.EqualFold) = false, want true", s3, s4)
	}

	cmpIntString := func(v1 int, v2 string) bool {
		return string(rune(v1)-1+'a') == v2
	}
	if !EqualFunc(s1, s3, cmpIntString) {
		t.Errorf("EqualFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
	if !equalFunc(s1, s3, func(a1, a2 any) bool { return cmpIntString(a1.(int), a2.(string)) }) {
		t.Errorf("equalFunc(%v, %v, cmpIntString) = false, want true", s1, s3)
	}
}

func BenchmarkEqualFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	xs := make([]Large, 1024)
	ys := make([]Large, 1024)

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = equalFunc(xs, ys, func(a1, a2 any) bool {
				return equal(a1.(Large), a2.(Large))
			})
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = EqualFunc(xs, ys, equal[Large])
		}
	})
}

var compareIntTests = []struct {
	s1, s2 []int
	want   int
}{
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		-1,
	},
	{
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3},
		+1,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 4, 3},
		-1,
	},
	{
		[]int{1, 4, 3},
		[]int{1, 2, 3},
		+1,
	},
}

var compareFloatTests = []struct {
	s1, s2 []float64
	want   int
}{
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, math.NaN(), 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, 4},
		-1,
	},
	{
		[]float64{1, math.NaN(), 3},
		[]float64{1, 2, math.NaN()},
		0,
	},
	{
		[]float64{1, math.NaN(), 3, 4},
		[]float64{1, 2, math.NaN()},
		+1,
	},
}

func TestCompare(t *testing.T) {
	t.Parallel()

	intWant := func(want bool) string {
		if want {
			return "0"
		}

		return "!= 0"
	}

	for _, tc := range equalIntTests {
		if got := Compare(tc.s1, tc.s2); tc.want != (got == 0) {
			t.Errorf("Compare(%v, %v) = %d, want %s", tc.s1, tc.s2, got, intWant(tc.want))
		}
	}

	for _, test := range equalFloatTests {
		if got := Compare(test.s1, test.s2); test.wantEqualNaN != (got == 0) {
			t.Errorf("Compare(%v, %v) = %d, want %s", test.s1, test.s2, got, intWant(test.wantEqualNaN))
		}
	}

	for _, test := range compareIntTests {
		if got := Compare(test.s1, test.s2); test.want != got {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}

	for _, test := range compareFloatTests {
		if got := Compare(test.s1, test.s2); test.want != got {
			t.Errorf("Compare(%v, %v) = %d, want %d", test.s1, test.s2, got, test.want)
		}
	}
}

func equalToCmp[T comparable](eq func(T, T) bool) func(T, T) int {
	return func(v1, v2 T) int {
		if eq(v1, v2) {
			return 0
		}

		return 1
	}
}

func cmp[T constraints.Ordered](v1, v2 T) int {
	if v1 < v2 {
		return -1
	} else if v1 > v2 {
		return 1
	} else {
		return 0
	}
}

func TestCompareFunc(t *testing.T) {
	t.Parallel()

	intWant := func(want bool) string {
		if want {
			return "0"
		}

		return "!= 0"
	}

	for _, tc := range equalIntTests {
		if got := CompareFunc(tc.s1, tc.s2, equalToCmp(equal[int])); tc.want != (got == 0) {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[int])) = %d, want %s", tc.s1, tc.s2, got, intWant(tc.want))
		}
	}

	for _, tc := range equalFloatTests {
		if got := CompareFunc(tc.s1, tc.s2, equalToCmp(equal[float64])); tc.wantEqual != (got == 0) {
			t.Errorf("CompareFunc(%v, %v, equalToCmp(equal[float64])) = %d, want %s", tc.s1, tc.s2, got, intWant(tc.wantEqual))
		}
	}

	for _, tc := range compareIntTests {
		if got := CompareFunc(tc.s1, tc.s2, cmp[int]); tc.want != got {
			t.Errorf("CompareFunc(%v, %v, cmp[int]) = %d, want %d", tc.s1, tc.s2, got, tc.want)
		}
	}

	for _, tc := range compareFloatTests {
		if got := CompareFunc(tc.s1, tc.s2, cmp[float64]); tc.want != got {
			t.Errorf("CompareFunc(%v, %v, cmp[float64]) = %d, want %d", tc.s1, tc.s2, got, tc.want)
		}
	}

	s1 := []int{1, 2, 3}
	s2 := []int{2, 3, 4}
	if got := CompareFunc(s1, s2, equalToCmp(offByOne[int])); got != 0 {
		t.Errorf("CompareFunc(%v, %v, offByOne) = %d, want 0", s1, s2, got)
	}

	s3 := []string{"a", "b", "c"}
	s4 := []string{"A", "B", "C"}
	if got := CompareFunc(s3, s4, strings.Compare); got != 1 {
		t.Errorf("CompareFunc(%v, %v, strings.Compare) = %d, want 1", s3, s4, got)
	}

	compareLower := func(v1, v2 string) int {
		return strings.Compare(strings.ToLower(v1), strings.ToLower(v2))
	}
	if got := CompareFunc(s3, s4, compareLower); got != 0 {
		t.Errorf("CompareFunc(%v, %v, compareLower) = %d, want 0", s3, s4, got)
	}

	cmpIntString := func(v1 int, v2 string) int {
		return strings.Compare(string(rune(v1)-1+'a'), v2)
	}
	if got := CompareFunc(s1, s3, cmpIntString); got != 0 {
		t.Errorf("CompareFunc(%v, %v, cmpIntString) = %d, want 0", s1, s3, got)
	}
}

var indexTests = []struct {
	s    []int
	v    int
	want int
}{
	{
		nil,
		0,
		-1,
	},
	{
		[]int{},
		0,
		-1,
	},
	{
		[]int{1, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 2, 3},
		2,
		1,
	},
	{
		[]int{1, 2, 3, 2},
		2,
		1,
	},
}

func TestIndex(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		if got := Index(tc.s, tc.v); tc.want != got {
			t.Errorf("Index(%v, %v) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
		if got := index(tc.s, tc.v); tc.want != got {
			t.Errorf("index(%v, %v) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
	}
}

func BenchmarkIndex_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = index(ss, Large{1})
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Index(ss, Large{1})
		}
	})
}

func equalToIndex[T any](f func(T, T) bool, v1 T) func(T) bool {
	return func(v2 T) bool {
		return f(v1, v2)
	}
}

func TestIndexFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		if got := IndexFunc(tc.s, equalToIndex(equal[int], tc.v)); tc.want != got {
			t.Errorf("IndexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", tc.s, tc.v, got, tc.want)
		}

		eq := equalToIndex(equal[int], tc.v)
		if got := indexFunc(tc.s, func(a any) bool {
			return eq(a.(int))
		}); tc.want != got {
			t.Errorf("indexFunc(%v, equalToIndex(equal[int], %v)) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
	}

	s1 := []string{"hi", "HI"}
	if got := IndexFunc(s1, equalToIndex(equal[string], "HI")); 1 != got {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, "HI", got, 1)
	}
	eq := equalToIndex(equal[string], "HI")
	if got := indexFunc(s1, func(a any) bool { return eq(a.(string)) }); 1 != got {
		t.Errorf("IndexFunc(%v, equalToIndex(equal[string], %q)) = %d, want %d", s1, "HI", got, 1)
	}

	if got := IndexFunc(s1, equalToIndex(strings.EqualFold, "HI")); 0 != got {
		t.Errorf("IndexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, "HI", got, 0)
	}

	eq = equalToIndex(strings.EqualFold, "HI")
	if got := indexFunc(s1, func(a any) bool { return eq(a.(string)) }); 0 != got {
		t.Errorf("indexFunc(%v, equalToIndex(strings.EqualFold, %q)) = %d, want %d", s1, "HI", got, 0)
	}
}

func BenchmarkIndexFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)

	b.Run("reflect", func(b *testing.B) {
		eq := equalToIndex(equal[Large], Large{1})
		for i := 0; i < b.N; i++ {
			_ = indexFunc(ss, func(a any) bool { return eq(a.(Large)) })
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = IndexFunc(ss, equalToIndex(equal[Large], Large{1}))
		}
	})
}

func TestContains(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		if got := Contains(tc.s, tc.v); (tc.want != -1) != got {
			t.Errorf("Contains(%v, %v) = %t, want %t", tc.s, tc.v, got, tc.want != -1)
		}

		if got := contains(tc.s, tc.v); (tc.want != -1) != got {
			t.Errorf("contains(%v, %v) = %t, want %t", tc.s, tc.v, got, tc.want != -1)
		}
	}
}

func TestContainsFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range indexTests {
		if got := ContainsFunc(tc.s, equalToIndex(equal[int], tc.v)); (tc.want != -1) != got {
			t.Errorf("ContainsFunc(%v, equalToIndex(equal[int], %v)) = %t, want %t", tc.s, tc.v, got, tc.want != -1)
		}

		if got := containsFunc(tc.s, func(a any) bool { return equalToIndex(equal[int], tc.v)(a.(int)) }); (tc.want != -1) != got {
			t.Errorf("containsFunc(%v, equalToIndex(equal[int], %v)) = %t, want %t", tc.s, tc.v, got, tc.want != -1)
		}
	}

	s1 := []string{"hi", "HI"}
	if got := ContainsFunc(s1, equalToIndex(equal[string], "HI")); !got {
		t.Errorf("ContainsFunc(%v, equalToContains(equal[string], %q)) = %t, want %t", s1, "HI", got, true)
	}
	if got := containsFunc(s1, func(a any) bool { return equalToIndex(equal[string], "HI")(a.(string)) }); !got {
		t.Errorf("containsFunc(%v, equalToContains(equal[string], %q)) = %t, want %t", s1, "HI", got, true)
	}

	if got := ContainsFunc(s1, equalToIndex(equal[string], "hI")); got {
		t.Errorf("ContainsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, false)
	}
	if got := containsFunc(s1, func(a any) bool { return equalToIndex(equal[string], "hI")(a.(string)) }); got {
		t.Errorf("containsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, false)
	}

	if got := ContainsFunc(s1, equalToIndex(strings.EqualFold, "hI")); !got {
		t.Errorf("ContainsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, true)
	}
	if got := containsFunc(s1, func(a any) bool { return equalToIndex(strings.EqualFold, "hI")(a.(string)) }); !got {
		t.Errorf("containsFunc(%v, equalToContains(strings.EqualFold, %q)) = %t, want %t", s1, "hI", got, true)
	}
}

var insertTests = []struct {
	s    []int
	i    int
	add  []int
	want []int
}{
	{
		[]int{1, 2, 3},
		0, []int{4},
		[]int{4, 1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		1, []int{4},
		[]int{1, 4, 2, 3},
	},
	{
		[]int{1, 2, 3},
		3, []int{4},
		[]int{1, 2, 3, 4},
	},
	{
		[]int{1, 2, 3},
		2, []int{4, 5},
		[]int{1, 2, 4, 5, 3},
	},
}

func TestInsert(t *testing.T) {
	t.Parallel()

	s := []int{1, 2, 3}
	if got := Insert(s, 0); !Equal(s, got) {
		t.Errorf("Insert(%v, 0) = %v, want %[1]v", s, got)
	}
	if got := insert(s, 0).([]int); !Equal(s, got) {
		t.Errorf("insert(%v, 0) = %v, want %[1]v", s, got)
	}

	for _, tc := range insertTests {
		copy := Clone(tc.s)
		if got := Insert(copy, tc.i, tc.add...); !Equal(tc.want, got) {
			t.Errorf("Insert(%v, %d, %v...) = %v, want %v", tc.s, tc.i, tc.add, got, tc.want)
		}

		copy = Clone(tc.s)
		add := make([]any, len(tc.add))
		for i := range add {
			add[i] = tc.add[i]
		}
		if got := insert(copy, tc.i, add...).([]int); !Equal(tc.want, got) {
			t.Errorf("insert(%v, %d, %v...) = %v, want %v", tc.s, tc.i, tc.add, got, tc.want)
		}
	}
}

var deleteTests = []struct {
	s    []int
	i, j int
	want []int
}{
	{
		[]int{1, 2, 3},
		0, 0,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0, 1,
		[]int{2, 3},
	},
	{
		[]int{1, 2, 3},
		3, 3,
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		0, 2,
		[]int{3},
	},
	{
		[]int{1, 2, 3},
		0, 3,
		[]int{},
	},
}

func TestDelete(t *testing.T) {
	t.Parallel()

	for _, tc := range deleteTests {
		copy := Clone(tc.s)
		if got := Delete(copy, tc.i, tc.j); !Equal(tc.want, got) {
			t.Errorf("Delete(%v, %d, %d) = %v, want %v", tc.s, tc.i, tc.j, got, tc.want)
		}

		copy = Clone(tc.s)
		if got := delete(copy, tc.i, tc.j).([]int); !Equal(tc.want, got) {
			t.Errorf("delete(%v, %d, %d) = %v, want %v", tc.s, tc.i, tc.j, got, tc.want)
		}
	}
}

func panics(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()

	f()
	return false
}

var deletePanicsTests = []struct {
	name string
	s    []int
	i, j int
}{
	{"with negative first index", []int{42}, -2, 1},
	{"with negative second index", []int{42}, 1, -1},
	{"with out-of-bounds first index", []int{42}, 2, 3},
	{"with out-of-bounds second index", []int{42}, 0, 2},
	{"with invalid i>j", []int{42}, 1, 0},
}

func TestDeletePanics(t *testing.T) {
	t.Parallel()

	for _, tc := range deletePanicsTests {
		if !panics(func() { Delete(tc.s, tc.i, tc.j) }) {
			t.Errorf("Delete %s: got no panic, want panic", tc.name)
		}

		if !panics(func() { delete(tc.s, tc.i, tc.j) }) {
			t.Errorf("delete %s: got no panic, want panic", tc.name)
		}
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	s1 := []int{1, 2, 3}
	s2 := Clone(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %[1]v", s1, s2)
	}

	s2 = clone(s1).([]int)
	if !Equal(s1, s2) {
		t.Errorf("clone(%v) = %v, want %[1]v", s1, s2)
	}

	s1[0] = 4
	want := []int{1, 2, 3}
	if !Equal(want, s2) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}

	if got := Clone([]int(nil)); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := clone([]int(nil)).([]int); got != nil {
		t.Errorf("clone(nil) = %#v, want nil", got)
	}

	if got := Clone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
	if got := clone(s1[:0]).([]int); got == nil || len(got) != 0 {
		t.Errorf("clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
}

type foo struct {
	i int
}

func (f *foo) Clone() *foo {
	return &foo{f.i}
}

func equalStructPointer[T comparable](p1, p2 *T) bool {
	if p1 == p2 {
		return false
	}

	return *p1 == *p2
}

func TestDeepClone(t *testing.T) {
	t.Parallel()

	s1 := []*foo{{1}, {2}, {3}}
	s2 := DeepClone(s1)
	if !EqualFunc(s1, s2, equalStructPointer[foo]) {
		t.Errorf("DeepClone(%v) = %v, want %[1]v", s1, s2)
	}

	s2 = deepClone(s1).([]*foo)
	if !EqualFunc(s1, s2, equalStructPointer[foo]) {
		t.Errorf("deepClone(%v) = %v, want %[1]v", s1, s2)
	}

	s1[0] = &foo{4}
	want := []*foo{{1}, {2}, {3}}
	if !EqualFunc(want, s2, equalStructPointer[foo]) {
		t.Errorf("DeepClone(%v) changed unexpectedly to %v", want, s2)
	}

	if got := DeepClone([]*foo(nil)); got != nil {
		t.Errorf("DeepClone(nil) = %#v, want nil", got)
	}
	if got := deepClone([]*foo(nil)).([]*foo); got != nil {
		t.Errorf("deepClone(nil) = %#v, want nil", got)
	}

	if got := DeepClone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("DeepClone(%v) = %#v, want %#[1]v", s1[:0], got)
	}
	if got := deepClone(s1[:0]).([]*foo); got == nil || len(got) != 0 {
		t.Errorf("DeepClone(%v) = %#v, want %#[1]v", s1[:0], got)
	}
}

var compactTests = []struct {
	name string
	s    []int
	want []int
}{
	{
		"nil",
		nil,
		nil,
	},
	{
		"one",
		[]int{1},
		[]int{1},
	},
	{
		"sorted",
		[]int{1, 2, 3},
		[]int{1, 2, 3},
	},
	{
		"1 item",
		[]int{1, 1, 2},
		[]int{1, 2},
	},
	{
		"unsorted",
		[]int{1, 2, 1},
		[]int{1, 2, 1},
	},
	{
		"many",
		[]int{1, 2, 2, 3, 3, 4},
		[]int{1, 2, 3, 4},
	},
}

func TestCompact(t *testing.T) {
	t.Parallel()

	for _, tc := range compactTests {
		copy := Clone(tc.s)
		if got := Compact(copy); !Equal(tc.want, got) {
			t.Errorf("Compact(%v) = %v, want %v", tc.s, got, tc.want)
		}

		copy = Clone(tc.s)
		if got := compact(copy).([]int); !Equal(tc.want, got) {
			t.Errorf("compact(%v) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func BenchmarkCompact(b *testing.B) {
	for _, tc := range compactTests {
		b.Run(tc.name, func(b *testing.B) {
			ss := make([]int, 0, 64)

			for k := 0; k < b.N; k++ {
				ss = append(ss[:0], tc.s...)
				_ = Compact(ss)
			}
		})
	}
}

func BenchmarkCompact_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = compact(ss)
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Compact(ss)
		}
	})
}

func TestCompactFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range compactTests {
		copy := Clone(tc.s)
		if got := CompactFunc(copy, equal[int]); !Equal(tc.want, got) {
			t.Errorf("CompactFunc(%v, equal[int]) = %v, want %v", tc.s, got, tc.want)
		}

		copy = Clone(tc.s)
		eq := equal[int]
		if got := compactFunc(copy, func(a1, a2 any) bool {
			return eq(a1.(int), a2.(int))
		}); !Equal(tc.want, got.([]int)) {
			t.Errorf("compactFunc(%v, equal[int]) = %v, want %v", tc.s, got, tc.want)
		}
	}

	s1 := []string{"a", "a", "A", "B", "b"}
	copy := Clone(s1)
	want := []string{"a", "B"}
	if got := CompactFunc(copy, strings.EqualFold); !Equal(want, got) {
		t.Errorf("CompactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}

	copy = Clone(s1)
	if got := compactFunc(copy, func(a1, a2 any) bool {
		return strings.EqualFold(a1.(string), a2.(string))
	}); !Equal(want, got.([]string)) {
		t.Errorf("compactFunc(%v, strings.EqualFold) = %v, want %v", s1, got, want)
	}
}

func BenchmarkCompactFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	ss := make([]Large, 1024)

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = compactFunc(ss, func(a1, a2 any) bool { return equal(a1.(Large), a2.(Large)) })
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = CompactFunc(ss, equal[Large])
		}
	})
}

func TestGrow(t *testing.T) {
	t.Parallel()

	s1 := []int{1, 2, 3}

	copy := Clone(s1)
	s2 := Grow(copy, 1000)
	if !Equal(s1, s2) {
		t.Errorf("Grow(%v) = %v, want %[1]v", s1, s2)
	}
	if cap(s2) < 1000+len(s1) {
		t.Errorf("after Grow(%v) cap = %d, want >= %d", s1, cap(s2), 1000+len(s1))
	}

	copy = Clone(s1)
	s2 = grow(copy, 1000).([]int)
	if !Equal(s1, s2) {
		t.Errorf("grow(%v) = %v, want %[1]v", s1, s2)
	}
	if cap(s2) < 1000+len(s1) {
		t.Errorf("after grow(%v) cap = %d, want >= %d", s1, cap(s2), 1000+len(s1))
	}

	// Test mutation of elements between length and capacity.
	copy = Clone(s1)
	s3 := Grow(copy[:1], 2)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}
	s3 = Grow(copy[:1], 1000)[:3]
	if !Equal(s1, s3) {
		t.Errorf("Grow should not mutate elements between length and capacity")
	}

	copy = Clone(s1)
	s3 = grow(copy[:1], 2).([]int)[:3]
	if !Equal(s1, s3) {
		t.Errorf("grow should not mutate elements between length and capacity")
	}
	s3 = grow(copy[:1], 1000).([]int)[:3]
	if !Equal(s1, s3) {
		t.Errorf("grow should not mutate elements between length and capacity")
	}

	// Test number of allocations.
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)) }); n != 0 {
		t.Errorf("Grow should not allocate when given sufficient capacity; allocated %v times", n)
	}
	if n := testing.AllocsPerRun(100, func() { Grow(s2, cap(s2)-len(s2)+1) }); n != 1 {
		errorf := t.Errorf
		if raceEnabled {
			errorf = t.Logf // this allocates multiple times in race detector mode
		}

		errorf("Grow should allocate once when given insufficient capacity; allocated %v times", n)
	}

	// Test for negative growth sizes.
	if !panics(func() { Grow(s1, -1) }) {
		t.Errorf("Grow(-1) did not panic; expected a panic")
	}
	if !panics(func() { grow(&s1, -1) }) {
		t.Errorf("grow(-1) did not panic; expected a panic")
	}
}

func TestClip(t *testing.T) {
	t.Parallel()

	s1 := []int{1, 2, 3, 4, 5, 6}[:3]
	orig := Clone(s1)

	if len(s1) != 3 {
		t.Errorf("len(%v) = %d, want 3", s1, len(s1))
	}
	if cap(s1) < 6 {
		t.Errorf("cap(%v[:3]) = %d, want >= 6", orig, cap(s1))
	}

	s2 := Clip(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clip(%v) = %v, want %[1]v", s1, s2)
	}
	if cap(s2) != 3 {
		t.Errorf("cap(Clip(%v)) = %d, want 3", orig, cap(s2))
	}

	s2 = clip(s1).([]int)
	if !Equal(s1, s2) {
		t.Errorf("clip(%v) = %v, want %[1]v", s1, s2)
	}
	if cap(s2) != 3 {
		t.Errorf("cap(clip(%v)) = %d, want 3", orig, cap(s2))
	}
}

// naiveReplace is a baseline implementation to the Replace function.
func naiveReplace[S ~[]E, E any](s S, i, j int, v ...E) S {
	s = Delete(s, i, j)
	s = Insert(s, i, v...)
	return s
}

var replaceTests = []struct {
	s    []int
	i, j int
	v    []int
}{
	{}, // all zero value
	{
		[]int{1, 2, 3, 4},
		1, 2,
		[]int{5},
	},
	{
		[]int{1, 2, 3, 4},
		1, 2,
		[]int{5, 6, 7, 8},
	},
	{
		append(make([]int, 0, 20), []int{0, 1, 2}...),
		0, 1,
		[]int{3, 4, 5, 6, 7},
	},
}

func TestReplace(t *testing.T) {
	t.Parallel()

	for _, tc := range replaceTests {
		ss, s := Clone(tc.s), Clone(tc.s)
		want := naiveReplace(ss, tc.i, tc.j, tc.v...)
		got := Replace(s, tc.i, tc.j, tc.v...)
		if !Equal(want, got) {
			t.Errorf("Replace(%v, %v, %v, %v...) = %v, want %v", tc.s, tc.i, tc.j, tc.v, got, want)
		}

		s = Clone(tc.s)
		v := make([]any, len(tc.v))
		for i := range v {
			v[i] = tc.v[i]
		}
		got = replace(s, tc.i, tc.j, v...).([]int)
		if !Equal(want, got) {
			t.Errorf("replace(%v, %v, %v, %v...) = %v, want %v", tc.s, tc.i, tc.j, tc.v, got, want)
		}
	}
}

var replacePanicsTests = []struct {
	name string
	s    []int
	i, j int
	v    []int
}{
	{
		"indexes out of order",
		[]int{1, 2},
		2, 1,
		[]int{3},
	},
	{
		"large index",
		[]int{1, 2},
		1, 10,
		[]int{3},
	},
	{
		"negative index",
		[]int{1, 2},
		-1, 2,
		[]int{3},
	},
}

func TestReplacePanics(t *testing.T) {
	t.Parallel()

	for _, tc := range replacePanicsTests {
		ss := Clone(tc.s)
		if !panics(func() { Replace(ss, tc.i, tc.j, tc.v...) }) {
			t.Errorf("Replace %s: should have panicked", tc.name)
		}

		ss = Clone(tc.s)
		v := make([]any, len(tc.v))
		for i := range v {
			v[i] = tc.v[i]
		}
		if !panics(func() { replace(ss, tc.i, tc.j, v...) }) {
			t.Errorf("replace %s: should have panicked", tc.name)
		}
	}
}

func BenchmarkReplace(b *testing.B) {
	testcases := []struct {
		name string
		s    func() []int
		i, j int
		v    func() []int
		vv   func() []any
	}{
		{
			"fast",
			func() []int {
				return make([]int, 100)
			},
			10, 40,
			func() []int {
				return make([]int, 20)
			},
			func() []any {
				s := make([]any, 20)
				for i := range s {
					s[i] = 0
				}

				return s
			},
		},
		{
			"slow",
			func() []int {
				return make([]int, 100)
			},
			0, 2,
			func() []int {
				return make([]int, 20)
			},
			func() []any {
				s := make([]any, 20)
				for i := range s {
					s[i] = 0
				}

				return s
			},
		},
	}

	for _, tc := range testcases {
		b.Run("naive-"+tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = naiveReplace(tc.s(), tc.i, tc.j, tc.v()...)
			}
		})

		b.Run("reflect-"+tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = replace(tc.s(), tc.i, tc.j, tc.vv()...)
			}
		})

		b.Run("optimized-"+tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Replace(tc.s(), tc.i, tc.j, tc.v()...)
			}
		})
	}
}
