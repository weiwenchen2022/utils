package slices

import (
	"math/rand"
	"runtime"
	"sync"

	"golang.org/x/exp/constraints"
)

// Filter returns a new slice of elements satisfies f(i, v).
func Filter[S ~[]E, E any](s S, f func(int, E) bool) S {
	if s == nil {
		return nil
	}

	r := make(S, 0, len(s))
	for i, v := range s {
		if f(i, v) {
			r = append(r, v)
		}
	}
	return r
}

// PFilter returns a new slice of elements satisfies f(i, v).
// f is call in a goroutine. Result may not keep the original order.
func PFilter[S ~[]E, E any](s S, f func(int, E) bool) S {
	if s == nil {
		return nil
	}

	ngoroutines := runtime.NumCPU()
	n := len(s)
	step := n / ngoroutines
	if step == 0 {
		step = 1
	}

	c := make(chan E, ngoroutines)

	var wg sync.WaitGroup
	for g := 0; g < ngoroutines; g++ {
		start := g * step
		if start >= n {
			break
		}

		end := start + step
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(s S) {
			for i, v := range s {
				if f(start+i, v) {
					c <- v
				}
			}
			wg.Done()
		}(s[start:end])
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	r := make(S, 0, len(s))
	for v := range c {
		r = append(r, v)
	}
	return r
}

// Map manipulates a slice and transforms it to a slice of another type.
func Map[S ~[]E1, E1, E2 any](s S, f func(int, E1) E2) []E2 {
	if s == nil {
		return nil
	}

	r := make([]E2, len(s))
	for i, v := range s {
		r[i] = f(i, v)
	}
	return r
}

// PMap manipulates a slice and transforms it to a slice of another type.
// f is call in a goroutine. Result keep the same order.
func PMap[S ~[]E1, E1, E2 any](s S, f func(int, E1) E2) []E2 {
	if s == nil {
		return nil
	}

	ngoroutines := runtime.NumCPU()
	n := len(s)
	step := n / ngoroutines
	if step == 0 {
		step = 1
	}

	type result struct {
		i int
		v E2
	}

	c := make(chan result, ngoroutines)

	var wg sync.WaitGroup
	for g := 0; g < ngoroutines; g++ {
		start := g * step
		if start >= n {
			break
		}

		end := start + step
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(s S) {
			for i, v := range s {
				c <- result{start + i, f(start+i, v)}
			}
			wg.Done()
		}(s[start:end])
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	r := make([]E2, len(s))
	for v := range c {
		r[v.i] = v.v
	}
	return r
}

// Reduce reduces the slice s to a signle value using a reduction function and a initial value.
func Reduce[S ~[]E1, E1, E2 any](s S, f func(E2, int, E1) E2, init E2) E2 {
	acc := init
	for i, v := range s {
		acc = f(acc, i, v)
	}
	return acc
}

// ForEach applies function f to each element of the slice s in order.
func ForEach[S ~[]E, E any](s S, f func(int, E)) {
	for i, v := range s {
		f(i, v)
	}
}

// PForEach applies function f to each element of the slice s in concurrency.
// f is call in a goroutine.
func PForEach[S ~[]E, E any](s S, f func(int, E)) {
	ngoroutines := runtime.NumCPU()
	n := len(s)
	step := n / ngoroutines
	if step == 0 {
		step = 1
	}

	var wg sync.WaitGroup
	for g := 0; g < ngoroutines; g++ {
		start := g * step
		if start >= n {
			break
		}

		end := start + step
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(s S) {
			for i, v := range s {
				f(start+i, v)
			}
			wg.Done()
		}(s[start:end])
	}

	wg.Wait()
}

// Shuffle returns a slice of shuffled elements of the slice s.
// Shuffle modifies the contents of the slice s; it does not create a new slice.
func Shuffle[S ~[]E, E any](s S) S {
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	return s
}

// Reverse reverses the slice s so that the first element becomes the last,
// the second element becomes the second to last, and so on.
// Reverse modifies the contents of the slice s; it does not create a new slice.
func Reverse[S ~[]E, E any](s S) S {
	n := len(s)
	h := n / 2
	for i := 0; i < h; i++ {
		j := n - i - 1
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Fill fills elements of the slice s with initial value.
// The elements are copied using assignment, so this is a shallow copy.
func Fill[S ~[]E, E any](s S, init E) S {
	for i := range s {
		s[i] = init
	}
	return s
}

// FillFunc is like Fill but uses the function f.
// The function f is invoked with index as argument.
func FillFunc[S ~[]E, E any](s S, f func(int) E) S {
	for i := range s {
		s[i] = f(i)
	}
	return s
}

// Repeat returns a slice with count copies of initial value.
// It return nil slice if count is zero.
func Repeat[E any](init E, count int) []E {
	if count == 0 {
		return nil
	}

	s := make([]E, count)
	for i := range s {
		s[i] = init
	}
	return s
}

// RepeatFunc is like Repeat but uses function f.
// The function f is invoked with index as argument.
func RepeatFunc[E any](f func(int) E, count int) []E {
	if count == 0 {
		return nil
	}

	s := make([]E, count)
	for i := range s {
		s[i] = f(i)
	}
	return s
}

// Count returns the number of elements in the slices s that equals to v.
func Count[S ~[]E, E comparable](s S, v E) int {
	count := 0
	for i := range s {
		if v == s[i] {
			count++
		}
	}
	return count
}

// CountFunc is like Count but uses a comparison function.
func CountFunc[S ~[]E, E any](s S, eq func(E) bool) int {
	count := 0
	for i := range s {
		if eq(s[i]) {
			count++
		}
	}
	return count
}

// Max returns the maximum element of the slice s, or panics if s is empty.
func Max[E constraints.Ordered](s ...E) E {
	max := s[0]
	for _, v := range s {
		if max < v {
			max = v
		}
	}
	return max
}

// Min returns the minimum element of the slices s, or panics if s is empty.
func Min[E constraints.Ordered](s ...E) E {
	min := s[0]
	for _, v := range s {
		if v < min {
			min = v
		}
	}
	return min
}

// SliceOf returns a slice which contains the element vs.
// If returns nil if len(vs) == 0.
func SliceOf[E any](vs ...E) []E {
	return append([]E(nil), vs...)
}

// Convenience wrappers for common cases.

// Filter returns the result of applying Filter to the receiver and f.
func (s Slice[E]) Filter(f func(int, E) bool) Slice[E] {
	return Filter(s, f)
}

// PFilter returns the result of applying PFilter to the receiver and f.
func (s Slice[E]) PFilter(f func(int, E) bool) Slice[E] {
	return PFilter(s, f)
}

// ForEach applies ForEach to the receiver and f.
func (s Slice[E]) ForEach(f func(int, E)) {
	ForEach(s, f)
}

// PForEach applies PForEach to the receiver and f.
func (s Slice[E]) PForEach(f func(int, E)) {
	PForEach(s, f)
}

// Shuffle returns the result of applying Shuffle to the receiver.
func (s Slice[E]) Shuffle() Slice[E] {
	return Shuffle(s)
}

// Reverse returns the result of applying Reverse to the receiver.
func (s Slice[E]) Reverse() Slice[E] {
	return Reverse(s)
}

// Fill returns the result applying Fill to the receiver and init.
func (s Slice[E]) Fill(init E) Slice[E] {
	return Fill(s, init)
}

// FillFunc returns the result applying FillFunc to the receiver and f.
func (s Slice[E]) FillFunc(f func(int) E) Slice[E] {
	return FillFunc(s, f)
}

// CountFunc returns the result applying CountFunc to the receiver and eq.
func (s Slice[E]) CountFunc(eq func(E) bool) int {
	return CountFunc(s, eq)
}

// Filter returns the result of applying Filter to the receiver and f.
func (s ComparableSlice[E]) Filter(f func(int, E) bool) ComparableSlice[E] {
	return Filter(s, f)
}

// PFilter returns the result of applying PFilter to the receiver and f.
func (s ComparableSlice[E]) PFilter(f func(int, E) bool) ComparableSlice[E] {
	return PFilter(s, f)
}

// ForEach applies ForEach to the receiver and f.
func (s ComparableSlice[E]) ForEach(f func(int, E)) {
	ForEach(s, f)
}

// PForEach applies PForEach to the receiver and f.
func (s ComparableSlice[E]) PForEach(f func(int, E)) {
	PForEach(s, f)
}

// Shuffle returns the result of applying Shuffle to the receiver.
func (s ComparableSlice[E]) Shuffle() ComparableSlice[E] {
	return Shuffle(s)
}

// Reverse returns the result of applying Reverse to the receiver.
func (s ComparableSlice[E]) Reverse() ComparableSlice[E] {
	return Reverse(s)
}

// Fill returns the result applying Fill to the receiver and init.
func (s ComparableSlice[E]) Fill(init E) ComparableSlice[E] {
	return Fill(s, init)
}

// FillFunc returns the result applying FillFunc to the receiver and f.
func (s ComparableSlice[E]) FillFunc(f func(int) E) ComparableSlice[E] {
	return FillFunc(s, f)
}

// Count returns the result applying Count to the receiver and v.
func (s ComparableSlice[E]) Count(v E) int {
	return Count(s, v)
}

// CountFunc returns the result applying CountFunc to the receiver and eq.
func (s ComparableSlice[E]) CountFunc(eq func(E) bool) int {
	return CountFunc(s, eq)
}
