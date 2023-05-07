package set_test

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/weiwenchen2022/utils/set"
)

func TestBasics(t *testing.T) {
	t.Parallel()

	s := set.New[int]()
	if len := s.Len(); len != 0 {
		t.Errorf("{}.Len(): got %d, want 0", len)
	}
	if s := s.String(); s != "{}" {
		t.Errorf("{}.String(): got %q, want \"{}\"", s)
	}
	if s.Has(3) {
		t.Errorf("Has(3): got true, want false")
	}

	if !s.Add(3) {
		t.Errorf("Add(3): got false, want true")
	}

	if !s.Add(435) {
		t.Errorf("Add(435): got false, want true")
	}

	elems := s.AppendTo(nil)
	sort.Ints(elems)
	if s := fmt.Sprint(elems); s != "[3 435]" {
		t.Errorf("{3 435}.AppendTo: got %q, want \"[3 435]\"", s)
	}
	if len := s.Len(); len != 2 {
		t.Errorf("Len: got %d, want 2", len)
	}

	if !s.Remove(435) {
		t.Errorf("Remove(435): got false, want true")
	}

	elems = s.AppendTo(nil)
	sort.Ints(elems)
	if s := fmt.Sprint(elems); s != "[3]" {
		t.Errorf("{3}.AppendTo: got %q, want \"[3]\"", s)
	}
}

// Add, Len, IsEmpty, Has, Clear, AppendTo.
func TestBasicsMore(t *testing.T) {
	t.Parallel()

	s := set.New[int]()
	s.Add(456)
	s.Add(123)
	s.Add(789)
	if s.Len() != 3 {
		t.Errorf("%s.Len: got %d, want 3", s, s.Len())
	}
	if s.IsEmpty() {
		t.Errorf("%s.IsEmpty: got true", s)
	}
	if !s.Has(123) {
		t.Errorf("%s.Has(123): got false", s)
	}
	if s.Has(1234) {
		t.Errorf("%s.Has(1234): got true", s)
	}
	got := s.AppendTo([]int{-1})
	sort.Ints(got)
	if want := []int{-1, 123, 456, 789}; fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("%s.AppendTo: got %v, want %v", s, got, want)
	}

	s.Clear()

	if s.Len() != 0 {
		t.Errorf("Clear: got %d, want 0", s.Len())
	}
	if !s.IsEmpty() {
		t.Errorf("IsEmpty: got false")
	}
	if s.Has(123) {
		t.Errorf("%s.Has: got false", s)
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()

	s1 := set.New[int]()
	s1.Add(456)
	s1.Add(123)
	s1.Add(789)

	if !s1.Equals(s1) {
		t.Errorf("%s.Equals(%s): got false", s1, s1)
	}

	s2 := set.New[int]()
	s2.Add(789)
	s2.Add(456)
	s2.Add(123)

	if !s1.Equals(s2) {
		t.Errorf("%s.Equals(%s): got false", s1, s2)
	}

	s2.Add(1)
	if s1.Equals(s2) {
		t.Errorf("%s.Equals(%s): got true", s1, s2)
	}

	empty := set.New[int]()
	if s2.Equals(empty) {
		t.Errorf("%s.Equals(%s): got true", s2, empty)
	}

	s2.Remove(123)
	if s1.Equals(s2) {
		t.Errorf("%s.Equals(%s): got true", s1, s2)
	}
}

func TestIntersectWith(t *testing.T) {
	t.Parallel()

	check := func(s1, s2 set.Set[int]) {
		s1str := s1.String()
		oldLen := s1.Len()
		if got, want := s1.IntersectWith(s2), s1.Len() < oldLen; want != got {
			t.Errorf("%s.IntersectWith(%s) = %t, want %t", s1str, s2, got, want)
		}
	}

	check(set.New(1, 2), set.New(1, 2))
	check(set.New(1, 2, 3), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 3))
	check(set.New(1, 2), set.New[int]())

	check(set.New(1, 1000000), set.New(1, 1000000))
	check(set.New(1, 2, 1000000), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 1000000))
	check(set.New(1, 1000000), set.New[int]())
}

func TestIntersectWith2(t *testing.T) {
	t.Parallel()

	s1, s2 := set.New[int](), set.New[int]()
	s1.Add(1)
	s1.Add(1000)
	s1.Add(8000)
	s2.Add(1)
	s2.Add(2000)
	s2.Add(4000)

	if got, want := s1.IntersectWith(s2), true; want != got {
		t.Errorf("IntersectWith: got %t, want %t", got, want)
	}
	if got, want := s1.String(), "{1}"; want != got {
		t.Errorf("IntersectWith: got %s, want %s", got, want)
	}
}

// randomSet returns a set of random size and elements.
func randomSet(r *rand.Rand, maxSize int) set.Set[int] {
	s := set.New[int]()

	size := int(r.Int63()) % maxSize
	for i := 0; i < size; i++ {
		n := int(r.Int63()) % 10000
		s.Add(n)
	}

	return s
}

func TestIntersects(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := uint(0); i < 12; i++ {
		s1, s2 := randomSet(r, 1<<i), randomSet(r, 1<<i)

		// test the slow way
		s3 := s1.Copy()
		s3.IntersectWith(s2)

		if got, want := s1.Intersects(s2), !s3.IsEmpty(); want != got {
			t.Errorf("%s.Intersects(%s): got %v, want %v (%s)", s1, s2, got, want, s3)
		}

		// make it false
		a := s1.AppendTo(nil)
		for _, x := range a {
			s2.Remove(x)
		}

		if got, want := s1.Intersects(s2), false; want != got {
			t.Errorf("Intersects: got %v, want %v", got, want)
		}

		// make it true
		if s1.IsEmpty() {
			continue
		}

		i := r.Intn(len(a))
		s2.Add(a[i])

		if got, want := s1.Intersects(s2), true; want != got {
			t.Errorf("Intersects: got %v, want %v", got, want)
		}
	}
}

func TestUnionWith(t *testing.T) {
	t.Parallel()

	check := func(s1, s2 set.Set[int]) {
		s1str := s1.String()
		oldLen := s1.Len()
		if got, want := s1.UnionWith(s2), s1.Len() > oldLen; want != got {
			t.Errorf("%s.UnionWith(%s) = %t, want %t", s1str, s2, got, want)
		}
	}

	check(set.New(1, 2), set.New(1, 2))
	check(set.New(1, 2, 3), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 3))
	check(set.New(1, 2), set.New[int]())

	check(set.New(1, 1000000), set.New(1, 1000000))
	check(set.New(1, 2, 1000000), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 1000000))
	check(set.New(1, 1000000), set.New[int]())
}

func TestDifferenceWith(t *testing.T) {
	t.Parallel()

	check := func(s1, s2 set.Set[int]) {
		s1str := s1.String()
		oldLen := s1.Len()
		if got, want := s1.DifferenceWith(s2), s1.Len() < oldLen; want != got {
			t.Errorf("%s.DifferenceWith(%s) = %t, want %t", s1str, s2, got, want)
		}
	}

	check(set.New(1, 2), set.New(1, 2))
	check(set.New(1, 2, 3), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 3))
	check(set.New(1, 2), set.New[int]())

	check(set.New(1, 1000000), set.New(1, 1000000))
	check(set.New(1, 2, 1000000), set.New(1, 2))
	check(set.New(1, 2), set.New(1, 2, 1000000))
	check(set.New(1, 1000000), set.New[int]())
}

// -- Benchmarks -------------------------------------------------------
func BenchmarkAdd(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Add(int(r.Int63()) % 10000)
	}
}

func BenchmarkRemove(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if s.Remove(int(r.Int63()) % 10000) {
			s.Add(int(r.Int63()) % 10000)
		}
	}
}

func BenchmarkHas(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Has(int(r.Int63()) % 10000)
	}
}

func BenchmarkCopy(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Copy()
	}
}

func BenchmarkEquals(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	for i := 0; i < b.N; i++ {
		_ = s.Equals(s)
	}
}

func BenchmarkString(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.String()
	}
}

func BenchmarkAppendTo(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := set.New[int]()
	for i := 0; i < 1000; i++ {
		s.Add(int(r.Int63()) % 10000)
	}

	elems := make([]int, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		elems = s.AppendTo(elems[:0])
	}
}

func BenchmarkIntersectWith(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		_ = s1.IntersectWith(s2)
	}
}

func BenchmarkIntersects(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		_ = s1.Intersects(s2)
	}
}

func BenchmarkUnionWith(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		_ = s1.UnionWith(s2)
	}
}

func BenchmarkDifferenceWith(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		_ = s1.DifferenceWith(s2)
	}
}

func BenchmarkSubsetOf(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		_ = s1.SubsetOf(s2)
	}
}

func BenchmarkSymmetricDifferenceWith(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		s1, s2 := set.New[int](), set.New[int]()
		for i := 0; i < 1000; i++ {
			x := int(r.Int63()) % 100000
			if i%2 == 0 {
				s1.Add(x)
			} else {
				s2.Add(x)
			}
		}

		s1.SymmetricDifferenceWith(s2)
	}
}
