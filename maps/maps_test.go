// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package maps_test

import (
	"math"
	"sort"
	"strconv"
	"testing"

	. "github.com/weiwenchen2022/utils/maps"
	"github.com/weiwenchen2022/utils/slices"
)

var (
	m1 = map[int]int{1: 2, 2: 4, 4: 8, 8: 16}
	m2 = map[int]string{1: "2", 2: "4", 4: "8", 8: "16"}
)

func TestMap_Keys(t *testing.T) {
	t.Parallel()

	want := []int{1, 2, 4, 8}

	got1 := NewMap(m1).Keys()
	sort.Ints(got1)
	if !slices.Equal(want, got1) {
		t.Errorf("Keys(%v) = %v, want %v", m1, got1, want)
	}

	got2 := NewMap(m2).Keys()
	sort.Ints(got2)
	if !slices.Equal(want, got2) {
		t.Errorf("Keys(%v) = %v, want %v", m2, got2, want)
	}
}

func TestMap_Values(t *testing.T) {
	t.Parallel()

	want1 := []int{2, 4, 8, 16}
	got1 := NewMap(m1).Values()
	sort.Ints(got1)
	if !slices.Equal(want1, got1) {
		t.Errorf("Values(%v) = %v, want %v", m1, got1, want1)
	}

	want2 := []string{"16", "2", "4", "8"}
	got2 := NewMap(m2).Values()
	sort.Strings(got2)
	if !slices.Equal(want2, got2) {
		t.Errorf("Values(%v) = %v, want %v", m2, got2, want2)
	}
}

func TestMap_EqualFunc(t *testing.T) {
	t.Parallel()

	m1 := NewMap(m1)
	if !m1.EqualFunc(m1, equal[int]) {
		t.Errorf("%v.EqualFunc(%[1]v, equal) = false, want true", m1)
	}
	if m1.EqualFunc(nil, equal[int]) {
		t.Errorf("%v.EqualFunc(nil, equal) = true, want false", m1)
	}
	if Map[int, int](nil).EqualFunc(m1, equal[int]) {
		t.Errorf("nil.EqualFunc(%v, equal) = true, want false", m1)
	}
	if !Map[int, int](nil).EqualFunc(map[int]int(nil), equal[int]) {
		t.Error("nil.EqualFunc(nil, equal) = false, want true")
	}
	if m := map[int]int{1: 2}; m1.EqualFunc(m, equal[int]) {
		t.Errorf("%v.EqualFunc(%v, equal) = true, want false", m1, m)
	}

	// Comparing NaN for equality is expected to fail.
	m := NewMap(map[int]float64{1: 0, 2: math.NaN()})
	if m.EqualFunc(m, equal[float64]) {
		t.Errorf("%v.EqualFunc(%[1]v, equal) = true, want false", m)
	}
	// But it should succeed using equalNaN.
	if !m.EqualFunc(m, equalNaN[float64]) {
		t.Errorf("%v.EqualFunc(%[1]v, equalNaN) = false, want true", m)
	}
}

func TestMap_Clear(t *testing.T) {
	t.Parallel()

	m := NewMap(map[int]int{1: 1, 2: 2, 3: 3})
	m.Clear()
	if got := len(m); got != 0 {
		t.Errorf("len(%v) = %d after Clear, want 0", m, got)
	}
	if !Equal(map[int]int(nil), m) {
		t.Errorf("Equal(nil, %v) = false, want true", m)
	}
}

func TestMap_Clone(t *testing.T) {
	t.Parallel()

	var m Map[int, int]
	mc := m.Clone()
	if mc != nil {
		t.Errorf("%v.Clone() = %v, want %[1]v", m, mc)
	}

	mc = NewMap(m1).Clone()
	if !Equal(m1, mc) {
		t.Errorf("%v.Clone() = %v, want %[1]v", m1, mc)
	}
	mc[16] = 32
	if Equal(m1, mc) {
		t.Errorf("Equal(%v, %v) = true, want false", m1, mc)
	}
}

func TestMap_Copy(t *testing.T) {
	t.Parallel()

	mc := NewMap(m1).Clone()
	mc.Copy(mc)
	if !Equal(m1, mc) {
		t.Errorf("%v.Copy(%[1]v) = %v, want %[1]v", m1, mc)
	}

	want := map[int]int{1: 2, 2: 4, 4: 8, 8: 16, 16: 32}
	mc.Copy(map[int]int{16: 32})
	if !Equal(want, mc) {
		t.Errorf("Copy result = %v, want %v", mc, want)
	}

	type M1 map[int]bool
	type M2 map[int]bool
	Copy(make(M1), make(M2))
}

func TestMap_DeleteFunc(t *testing.T) {
	t.Parallel()

	mc := NewMap(m1).Clone()

	mc.DeleteFunc(func(int, int) bool { return false })
	if !Equal(m1, mc) {
		t.Errorf("%v.DeleteFunc(false) = %v, want %[1]v", m1, mc)
	}

	mc.DeleteFunc(func(k, _ int) bool { return k > 3 })
	want := map[int]int{1: 2, 2: 4}
	if !Equal(want, mc) {
		t.Errorf("DeleteFunc result = %v, want %v", mc, want)
	}
}

func TestComparableMap_Equal(t *testing.T) {
	t.Parallel()

	if !NewComparableMap(m1).Equal(m1) {
		t.Errorf("%v.Equal(%[1]v) = false, want true", m1)
	}
	if NewComparableMap(m1).Equal(map[int]int(nil)) {
		t.Errorf("%v.Equal(nil) = true, want false", m1)
	}
	if ComparableMap[int, int](nil).Equal(m1) {
		t.Errorf("nil.Equal(%v) = true, want false", m1)
	}
	if !ComparableMap[int, int](nil).Equal(map[int]int(nil)) {
		t.Error("nil.Equal(nil) = false, want true")
	}
	if m := map[int]int{1: 2}; NewComparableMap(m1).Equal(m) {
		t.Errorf("%v.Equal(%v) = true, want false", m1, m)
	}

	// Comparing NaN for equality is expected to fail.
	m := map[int]float64{1: 0, 2: math.NaN()}
	if NewComparableMap(m).Equal(m) {
		t.Errorf("%v.Equal(%[1]v) = true, want false", m)
	}
}

func TestKeys(t *testing.T) {
	t.Parallel()

	want := []int{1, 2, 4, 8}

	got1 := Keys(m1)
	sort.Ints(got1)
	if !slices.Equal(want, got1) {
		t.Errorf("Keys(%v) = %v, want %v", m1, got1, want)
	}
	got1 = keys(m1).([]int)
	sort.Ints(got1)
	if !slices.Equal(want, got1) {
		t.Errorf("keys(%v) = %v, want %v", m1, got1, want)
	}

	got2 := Keys(m2)
	sort.Ints(got2)
	if !slices.Equal(want, got2) {
		t.Errorf("Keys(%v) = %v, want %v", m2, got2, want)
	}
	got2 = keys(m2).([]int)
	sort.Ints(got2)
	if !slices.Equal(want, got2) {
		t.Errorf("keys(%v) = %v, want %v", m2, got2, want)
	}
}

func TestValues(t *testing.T) {
	t.Parallel()

	want1 := []int{2, 4, 8, 16}
	got1 := Values(m1)
	sort.Ints(got1)
	if !slices.Equal(want1, got1) {
		t.Errorf("Values(%v) = %v, want %v", m1, got1, want1)
	}
	got1 = values(m1).([]int)
	sort.Ints(got1)
	if !slices.Equal(want1, got1) {
		t.Errorf("values(%v) = %v, want %v", m1, got1, want1)
	}

	want2 := []string{"16", "2", "4", "8"}
	got2 := Values(m2)
	sort.Strings(got2)
	if !slices.Equal(want2, got2) {
		t.Errorf("Values(%v) = %v, want %v", m2, got2, want2)
	}
	got2 = values(m2).([]string)
	sort.Strings(got2)
	if !slices.Equal(want2, got2) {
		t.Errorf("values(%v) = %v, want %v", m2, got2, want2)
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	if !Equal(m1, m1) {
		t.Errorf("Equal(%v, %[1]v) = false, want true", m1)
	}
	if Equal(m1, map[int]int(nil)) {
		t.Errorf("Equal(%v, nil) = true, want false", m1)
	}
	if Equal(map[int]int(nil), m1) {
		t.Errorf("Equal(nil, %v) = true, want false", m1)
	}
	if !Equal(map[int]int(nil), map[int]int(nil)) {
		t.Error("Equal(nil, nil) = false, want true")
	}
	if m := map[int]int{1: 2}; Equal(m1, m) {
		t.Errorf("Equal(%v, %v) = true, want false", m1, m)
	}

	if !equals(m1, m1) {
		t.Errorf("equals(%v, %[1]v) = false, want true", m1)
	}
	if equals(m1, map[int]int(nil)) {
		t.Errorf("equals(%v, nil) = true, want false", m1)
	}
	if equals(map[int]int(nil), m1) {
		t.Errorf("equals(nil, %v) = true, want false", m1)
	}
	if !equals(map[int]int(nil), map[int]int(nil)) {
		t.Error("equals(nil, nil) = false, want true")
	}
	if m := map[int]int{1: 2}; equals(m1, m) {
		t.Errorf("equals(%v, %v) = true, want false", m1, m)
	}

	// Comparing NaN for equality is expected to fail.
	m := map[int]float64{1: 0, 2: math.NaN()}
	if Equal(m, m) {
		t.Errorf("Equal(%v, %[1]v) = true, want false", m)
	}
	if equals(m, m) {
		t.Errorf("equals(%v, %[1]v) = true, want false", m)
	}
}

func BenchmarkEqual_Large(b *testing.B) {
	type Large [4 * 1024]byte

	xm := make(map[int]Large, 1024)
	for i := 0; i < 1024; i++ {
		xm[i] = Large{}
	}

	ym := make(map[int]Large, 1024)
	for i := 0; i < 1024; i++ {
		ym[i] = Large{}
	}

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = equals(xm, ym)
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Equal(xm, ym)
		}
	})
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

// equalStr compares ints and strings.
func equalIntStr(v1 int, v2 string) bool {
	return strconv.Itoa(v1) == v2
}

func TestEqualFunc(t *testing.T) {
	t.Parallel()

	if !EqualFunc(m1, m1, equal[int]) {
		t.Errorf("EqualFunc(%v, %[1]v, equal) = false, want true", m1)
	}
	if EqualFunc(m1, map[int]int(nil), equal[int]) {
		t.Errorf("EqualFunc(%v, nil, equal) = true, want false", m1)
	}
	if EqualFunc(map[int]int(nil), m1, equal[int]) {
		t.Errorf("EqualFunc(nil, %v, equal) = true, want false", m1)
	}
	if !EqualFunc(map[int]int(nil), map[int]int(nil), equal[int]) {
		t.Error("EqualFunc(nil, nil, equal) = false, want true")
	}
	if m := map[int]int{1: 2}; EqualFunc(m1, m, equal[int]) {
		t.Errorf("EqualFunc(%v, %v, equal) = true, want false", m1, m)
	}

	equalInt := func(v1, v2 any) bool {
		return equal(v1.(int), v2.(int))
	}
	if !equalFunc(m1, m1, equalInt) {
		t.Errorf("equalFunc(%v, %[1]v, equal) = false, want true", m1)
	}
	if equalFunc(m1, map[int]int(nil), equalInt) {
		t.Errorf("equalFunc(%v, nil, equal) = true, want false", m1)
	}
	if equalFunc(map[int]int(nil), m1, equalInt) {
		t.Errorf("equalFunc(nil, %v, equal) = true, want false", m1)
	}
	if !equalFunc(map[int]int(nil), map[int]int(nil), equalInt) {
		t.Error("equalFunc(nil, nil, equal) = false, want true")
	}
	if m := map[int]int{1: 2}; equalFunc(m1, m, equalInt) {
		t.Errorf("equalFunc(%v, %v, equal) = true, want false", m1, m)
	}

	// Comparing NaN for equality is expected to fail.
	m := map[int]float64{1: 0, 2: math.NaN()}
	if EqualFunc(m, m, equal[float64]) {
		t.Errorf("EqualFunc(%v, %[1]v, equal) = true, want false", m)
	}

	equalFloat := func(v1, v2 any) bool {
		return equal(v1.(float64), v2.(float64))
	}
	if equalFunc(m, m, equalFloat) {
		t.Errorf("equalFunc(%v, %[1]v, equal) = true, want false", m)
	}

	// But it should succeed using equalNaN.
	if !EqualFunc(m, m, equalNaN[float64]) {
		t.Errorf("EqualFunc(%v, %[1]v, equalNaN) = false, want true", m)
	}

	equalFloatNaN := func(v1, v2 any) bool {
		return equalNaN(v1.(float64), v2.(float64))
	}
	if !equalFunc(m, m, equalFloatNaN) {
		t.Errorf("equalFunc(%v, %[1]v, equalNaN) = false, want true", m)
	}

	if !EqualFunc(m1, m2, equalIntStr) {
		t.Errorf("EqualFunc(%v, %v, equalIntStr) = false, want true", m1, m2)
	}
	if !equalFunc(m1, m2, func(v1, v2 any) bool { return equalIntStr(v1.(int), v2.(string)) }) {
		t.Errorf("equalFunc(%v, %v, equalIntStr) = false, want true", m1, m2)
	}
}

func BenchmarkEqualFunc_Large(b *testing.B) {
	type Large [4 * 1024]byte

	xm := make(map[int]Large, 1024)
	for i := 0; i < 1024; i++ {
		xm[i] = Large{}
	}
	ym := make(map[int]Large, 1024)
	for i := 0; i < 1024; i++ {
		ym[i] = Large{}
	}

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = equalFunc(xm, ym, func(v1, v2 any) bool { return equal(v1.(Large), v2.(Large)) })
		}
	})
	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = EqualFunc(xm, ym, equal[Large])
		}
	})
}

func TestClear(t *testing.T) {
	t.Parallel()

	m := map[int]int{1: 1, 2: 2, 3: 3}
	mc := Clone(m)
	Clear(m)
	if got := len(m); got != 0 {
		t.Errorf("len(%v) = %d after Clear, want 0", m, got)
	}
	if !Equal(map[int]int(nil), m) {
		t.Errorf("Equal(nil, %v) = false, want true", m)
	}

	m = Clone(mc)
	clear(m)
	if got := len(m); got != 0 {
		t.Errorf("len(%v) = %d after clear, want 0", m, got)
	}
	if !Equal(map[int]int(nil), m) {
		t.Errorf("Equal(nil, %v) = false, want true", m)
	}
}

func TestClone(t *testing.T) {
	t.Parallel()

	var m map[int]int
	mc := Clone(m)
	if mc != nil {
		t.Errorf("Clone(%v) = %v, want %[1]v", m, mc)
	}
	mc = clone(m).(map[int]int)
	if mc != nil {
		t.Errorf("clone(%v) = %v, want %[1]v", m, mc)
	}

	mc = Clone(m1)
	if !Equal(m1, mc) {
		t.Errorf("Clone(%v) = %v, want %[1]v", m1, mc)
	}
	mc[16] = 32
	if Equal(m1, mc) {
		t.Errorf("Equal(%v, %v) = true, want false", m1, mc)
	}

	mc = clone(m1).(map[int]int)
	if !Equal(m1, mc) {
		t.Errorf("clone(%v) = %v, want %[1]v", m1, mc)
	}
	mc[16] = 32
	if Equal(m1, mc) {
		t.Errorf("Equal(%v, %v) = true, want false", m1, mc)
	}
}

func TestCopy(t *testing.T) {
	t.Parallel()

	mc := Clone(m1)
	Copy(mc, mc)
	if !Equal(m1, mc) {
		t.Errorf("Copy(%v, %[1]v) = %v, want %[1]v", m1, mc)
	}
	mc = Clone(m1)
	Copy(mc, mc)
	if !Equal(m1, mc) {
		t.Errorf("copy(%v, %[1]v) = %v, want %[1]v", m1, mc)
	}

	want := map[int]int{1: 2, 2: 4, 4: 8, 8: 16, 16: 32}
	Copy(mc, map[int]int{16: 32})
	if !Equal(want, mc) {
		t.Errorf("Copy result = %v, want %v", mc, want)
	}
	mc = Clone(m1)
	copy(mc, map[int]int{16: 32})
	if !Equal(want, mc) {
		t.Errorf("copy result = %v, want %v", mc, want)
	}

	type M1 map[int]bool
	type M2 map[int]bool
	Copy(make(M1), make(M2))
	copy(make(M1), make(M2))
}

func TestDeleteFunc(t *testing.T) {
	t.Parallel()

	mc := Clone(m1)

	DeleteFunc(mc, func(int, int) bool { return false })
	if !Equal(m1, mc) {
		t.Errorf("DeleteFunc(%v, false) = %v, want %[1]v", m1, mc)
	}
	deleteFunc(mc, func(any, any) bool { return false })
	if !Equal(m1, mc) {
		t.Errorf("deleteFunc(%v, false) = %v, want %[1]v", m1, mc)
	}

	DeleteFunc(mc, func(k, _ int) bool { return k > 3 })
	want := map[int]int{1: 2, 2: 4}
	if !Equal(want, mc) {
		t.Errorf("DeleteFunc result = %v, want %v", mc, want)
	}
	mc = Clone(m1)
	deleteFunc(mc, func(k, _ any) bool { return k.(int) > 3 })
	if !Equal(want, mc) {
		t.Errorf("deleteFunc result = %v, want %v", mc, want)
	}
}
