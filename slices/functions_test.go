package slices_test

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/weiwenchen2022/utils/slices"
)

var filterTests = []struct {
	s    []int
	f    func(int, int) bool
	want []int
}{
	{nil, func(int, int) bool { return true }, nil},
	{[]int{0, 1, 2, 3, 4}, func(_ int, v int) bool { return v&0x1 == 0 }, []int{0, 2, 4}},
}

func TestFilter(t *testing.T) {
	t.Parallel()

	for _, tc := range filterTests {
		got := Filter(tc.s, tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("Filter(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = filter(tc.s, tc.f).([]int)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("filter(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = PFilter(tc.s, tc.f)
		sort.Ints(got)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("PFilter(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = pfilter(tc.s, tc.f).([]int)
		sort.Ints(got)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("pfilter(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		s    []int
		f    func(int, int) string
		want []string
	}{
		{
			s:    nil,
			f:    func(int, int) string { return "" },
			want: nil,
		},
		{
			s:    []int{1, 2, 3, 4},
			f:    func(_ int, v int) string { return strconv.Itoa(v) },
			want: []string{"1", "2", "3", "4"},
		},
	}

	for _, tc := range testcases {
		got := Map(tc.s, tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("Map(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = _map(tc.s, tc.f).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("_map(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = PMap(tc.s, tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("PMap(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		got = pmap(tc.s, tc.f).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("pmap(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}
	}
}

func sliceGenerator(size int) []int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	s := make([]int, size)
	for i := range s {
		s[i] = int(r.Int63()) % size
	}
	return s
}

func BenchmarkMap(b *testing.B) {
	s := sliceGenerator(1000_000)

	b.Run("for", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := make([]string, len(s))
			for i, v := range s {
				r[i] = strconv.Itoa(v)
			}
			_ = r
		}
	})

	b.Run("reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = _map(s, func(_ int, v int) string {
				return strconv.Itoa(v)
			}).([]string)
		}
	})

	b.Run("generic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Map(s, func(_ int, v int) string {
				return strconv.Itoa(v)
			})
		}
	})

	b.Run("reflect2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = pmap(s, func(_ int, v int) string {
				return strconv.Itoa(v)
			}).([]string)
		}
	})

	b.Run("generic2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = PMap(s, func(_ int, v int) string {
				return strconv.Itoa(v)
			})
		}
	})
}

func TestReduce(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		s    []int
		f    func(int, int, int) int
		init int
		want int
	}{
		{[]int(nil), func(acc int, _ int, v int) int { return acc + v }, 0, 0},
		{[]int{1, 2, 3, 4}, func(acc int, _ int, v int) int { return acc + v }, 0, 10},
	}

	for _, tc := range testcases {
		got := Reduce(tc.s, tc.f, tc.init)
		if tc.want != got {
			t.Errorf("Reduce(%#v) = %v, want %v", tc.s, got, tc.want)
		}

		got = reduce(tc.s, tc.f, tc.init).(int)
		if tc.want != got {
			t.Errorf("reduce(%#v) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

var forEachTests = []struct {
	s []int
}{
	{nil},
	{[]int{0, 1, 2, 3, 4}},
}

func TestForEach(t *testing.T) {
	t.Parallel()

	for _, tc := range forEachTests {
		var got []int
		f := func(_ int, v int) {
			got = append(got, v)
		}

		ForEach(tc.s, f)
		if fmt.Sprintf("%#v", tc.s) != fmt.Sprintf("%#v", got) {
			t.Errorf("ForEach(%#v) = %#v, want %#[1]v", tc.s, got)
		}

		got = nil
		forEach(tc.s, f)
		if fmt.Sprintf("%#v", tc.s) != fmt.Sprintf("%#v", got) {
			t.Errorf("forEach(%#v) = %#v, want %#[1]v", tc.s, got)
		}
	}
}

var pForEachTests = []struct {
	s    []int
	want int
}{
	{nil, 0},
	{[]int{0, 1, 2, 3, 4}, 5},
}

func TestPForEach(t *testing.T) {
	t.Parallel()

	for _, tc := range pForEachTests {
		var n atomic.Int64
		f := func(int, int) {
			n.Add(1)
		}

		PForEach(tc.s, f)
		if got := int(n.Load()); tc.want != got {
			t.Errorf("PForEach(%#v) = %v, want %v", tc.s, got, tc.want)
		}

		n.Store(0)
		pForEach(tc.s, f)
		if got := int(n.Load()); tc.want != got {
			t.Errorf("pForEach(%#v) = %v, want %v", tc.s, got, tc.want)
		}
	}
}

var shuffleTests = []struct {
	s []int
}{
	{nil},
	{[]int{0, 1, 2, 3, 4}},
}

func TestShuffle(t *testing.T) {
	t.Parallel()

	for _, tc := range shuffleTests {
		s := Clone(tc.s)
		got := Shuffle(s)
		if (tc.s == nil) != (got == nil) &&
			fmt.Sprintf("%#v", tc.s) == fmt.Sprintf("%#v", got) {
			t.Errorf("Shuffle(%#v) = %#v", tc.s, got)
		}

		s = Clone(tc.s)
		got = shuffle(s).([]int)
		if (tc.s == nil) != (got == nil) &&
			fmt.Sprintf("%#v", tc.s) == fmt.Sprintf("%#v", got) {
			t.Errorf("shuffle(%#v) = %#v", tc.s, got)
		}
	}
}

var reverseTests = []struct {
	s []int
}{
	{nil},
	{[]int{0, 1, 2, 3, 4}},
}

func TestReverse(t *testing.T) {
	t.Parallel()

	for _, tc := range reverseTests {
		s := Clone(tc.s)
		got := Reverse(s)
		for i, n := 0, len(tc.s); i < n; i++ {
			if got[n-i-1] != tc.s[i] {
				t.Errorf("Reverse(%#v) = %#v", tc.s, got)
			}
		}

		s = Clone(tc.s)
		got = reverse(s).([]int)
		for i, n := 0, len(tc.s); i < n; i++ {
			if got[n-i-1] != tc.s[i] {
				t.Errorf("reverse(%#v) = %#v", tc.s, got)
			}
		}
	}
}

var fillTests = []struct {
	s    []string
	init string
	want []string
}{
	{nil, "", nil},
	{make([]string, 2), "foo", []string{"foo", "foo"}},
}

func TestFill(t *testing.T) {
	t.Parallel()

	for _, tc := range fillTests {
		s := Clone(tc.s)
		got := Fill(s, tc.init)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("Fill(%#v, %v) = %#v, want %#v", tc.s, tc.init, got, tc.want)
		}

		s = Clone(tc.s)
		got = fill(s, tc.init).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("fill(%#v, %v) = %#v, want %#v", tc.s, tc.init, got, tc.want)
		}
	}
}

var fillFuncTests = []struct {
	s    []string
	f    func(int) string
	want []string
}{
	{nil, func(int) string { return "" }, nil},
	{make([]string, 2), func(i int) string { return strconv.Itoa(i) }, []string{"0", "1"}},
}

func TestFillFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range fillFuncTests {
		s := Clone(tc.s)
		got := FillFunc(s, tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("FillFunc(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}

		s = Clone(tc.s)
		got = fillFunc(s, tc.f).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("fillFunc(%#v) = %#v, want %#v", tc.s, got, tc.want)
		}
	}
}

func TestRepeat(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		init  string
		count int
		want  []string
	}{
		{"", 0, nil},
		{"foo", 2, []string{"foo", "foo"}},
	}

	for _, tc := range testcases {
		got := Repeat(tc.init, tc.count)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("Repeat(%v, %v) = %#v, want %#v", tc.init, tc.count, got, tc.want)
		}

		got = repeat(tc.init, tc.count).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("repeat(%v, %v) = %#v, want %#v", tc.init, tc.count, got, tc.want)
		}
	}
}

func TestRepeatFunc(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		f     func(int) string
		count int
		want  []string
	}{
		{func(i int) string { return strconv.Itoa(i) }, 0, nil},
		{func(i int) string { return strconv.Itoa(i) }, 5, []string{"0", "1", "2", "3", "4"}},
	}

	for _, tc := range testcases {
		got := RepeatFunc(tc.f, tc.count)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("RepeatFunc() = %#v, want %#v", got, tc.want)
		}

		got = repeatFunc(tc.f, tc.count).([]string)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", got) {
			t.Errorf("repeatFunc() = %#v, want %#v", got, tc.want)
		}
	}
}

var countTests = []struct {
	s    []int
	v    int
	want int
}{
	{[]int{1, 2, 1}, 1, 2},
	{[]int{1, 2, 1}, 3, 0},
	{[]int(nil), 1, 0},
}

func TestCount(t *testing.T) {
	t.Parallel()

	for _, tc := range countTests {
		if got := Count(tc.s, tc.v); tc.want != got {
			t.Errorf("Count(%#v, %v) = %v, want %v", tc.s, tc.v, got, tc.want)
		}
		if got := count(tc.s, tc.v); tc.want != got {
			t.Errorf("count(%#v, %v) = %v, want %v", tc.s, tc.v, got, tc.want)
		}
	}
}

var countFuncTests = []struct {
	s    []int
	eq   func(int) bool
	want int
}{
	{[]int(nil), func(v int) bool { return v < 2 }, 0},
	{[]int{1, 2, 1}, func(v int) bool { return v < 2 }, 2},
	{[]int{1, 2, 1}, func(v int) bool { return v > 2 }, 0},
}

func TestCountFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range countFuncTests {
		if got := CountFunc(tc.s, tc.eq); tc.want != got {
			t.Errorf("%#v.CountFunc() = %v, want %v", tc.s, got, tc.want)
		}
		if got := countFunc(tc.s, tc.eq); tc.want != got {
			t.Errorf("%#v.countFunc() = %v, want %v", tc.s, got, tc.want)
		}
	}
}

// Tests for convenience wrappers.

func TestSlice_Filter(t *testing.T) {
	t.Parallel()

	for _, tc := range filterTests {
		got := NewSlice(tc.s).Filter(tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", []int(got)) {
			t.Errorf("%#v.Filter() = %#v, want %#v", tc.s, got, tc.want)
		}

		got = NewSlice(tc.s).PFilter(tc.f)
		sort.Ints(got)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", []int(got)) {
			t.Errorf("%#v.PFilter() = %#v, want %#v", tc.s, got, tc.want)
		}
	}
}

func TestSlice_ForEach(t *testing.T) {
	t.Parallel()

	for _, tc := range forEachTests {
		var got []int
		f := func(_ int, v int) {
			got = append(got, v)
		}

		NewSlice(tc.s).ForEach(f)
		if fmt.Sprintf("%#v", tc.s) != fmt.Sprintf("%#v", got) {
			t.Errorf("%#v.ForEach() = %#v, want %#[1]v", tc.s, got)
		}
	}
}

func TestSlice_PForEach(t *testing.T) {
	t.Parallel()

	for _, tc := range pForEachTests {
		var n atomic.Int64
		f := func(int, int) {
			n.Add(1)
		}

		NewSlice(tc.s).PForEach(f)
		if got := int(n.Load()); tc.want != got {
			t.Errorf("%#v.PForEach() = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestSlice_Shuffle(t *testing.T) {
	t.Parallel()

	for _, tc := range shuffleTests {
		s := NewSlice(Clone(tc.s))
		got := s.Shuffle()
		if (tc.s == nil) != (got == nil) &&
			fmt.Sprintf("%#v", tc.s) == fmt.Sprintf("%#v", ([]int)(got)) {
			t.Errorf("%#v.Shuffle() = %#v", tc.s, ([]int)(got))
		}
	}
}

func TestSlice_Reverse(t *testing.T) {
	t.Parallel()

	for _, tc := range reverseTests {
		s := NewSlice(Clone(tc.s))
		got := s.Reverse()
		for i, n := 0, len(tc.s); i < n; i++ {
			if got[n-i-1] != tc.s[i] {
				t.Errorf("%#v.Reverse() = %#v", tc.s, got)
			}
		}
	}
}

func TestSlice_Fill(t *testing.T) {
	t.Parallel()

	for _, tc := range fillTests {
		s := NewSlice(Clone(tc.s))
		got := s.Fill(tc.init)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", ([]string)(got)) {
			t.Errorf("%#v.Fill(%v) = %#v, want %#v", tc.s, tc.init, ([]string)(got), tc.want)
		}
	}
}

func TestSlice_FillFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range fillFuncTests {
		s := NewSlice(Clone(tc.s))
		got := s.FillFunc(tc.f)
		if fmt.Sprintf("%#v", tc.want) != fmt.Sprintf("%#v", ([]string)(got)) {
			t.Errorf("%#v.FillFunc() = %#v, want %#v", tc.s, got, tc.want)
		}
	}
}

func TestSlice_CountFunc(t *testing.T) {
	t.Parallel()

	for _, tc := range countFuncTests {
		if got := NewSlice(tc.s).CountFunc(tc.eq); tc.want != got {
			t.Errorf("%#v.CountFunc() = %v, want %v", tc.s, got, tc.want)
		}
	}
}

func TestComparableSlice_Count(t *testing.T) {
	t.Parallel()

	for _, tc := range countTests {
		if got := NewComparableSlice(tc.s).Count(tc.v); tc.want != got {
			t.Errorf("%#v.Count(%v) = %v, want %v", tc.s, tc.v, got, tc.want)
		}
	}
}
