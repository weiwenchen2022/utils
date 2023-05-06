package channels_test

import (
	"context"
	"sort"
	"testing"
	"time"

	. "github.com/weiwenchen2022/utils/channels"
	. "github.com/weiwenchen2022/utils/slices"
)

func TestSliceToChannel(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		want := []int{1, 2, 3}

		c := SliceToChannel([]int{1, 2, 3})
		var got []int
		for v := range c {
			got = append(got, v)
		}

		if !Equal(want, got) {
			t.Errorf("SliceToChannel(%v) = %v, want %[1]v", want, got)
		}

		c = sliceToChannel([]int{1, 2, 3}).(chan int)
		got = nil
		for v := range c {
			got = append(got, v)
		}

		if !Equal(want, got) {
			t.Errorf("sliceToChannel(%v) = %v, want %[1]v", want, got)
		}
	})
}

func TestChannelToSlice(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		want := []int{1, 2, 3}

		c := SliceToChannel(want)
		s := ChannelToSlice(c)
		if !Equal(want, s) {
			t.Errorf("ChannelToSlice(%v) = %v, want %[1]v", want, s)
		}

		c = SliceToChannel(want)
		s = channelToSlice(c).([]int)
		if !Equal(want, s) {
			t.Errorf("channelToSlice(%v) = %v, want %[1]v", want, s)
		}
	})
}

func TestGenerator(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		want := []int{0, 1, 2, 3, 4}

		gen := func(yield func(int)) {
			for i := range want {
				yield(i)
			}
		}

		var got []int
		for v := range Generator(gen) {
			got = append(got, v)
		}
		if !Equal(want, got) {
			t.Errorf("Generator() = %v, want %v", got, want)
		}

		got = nil
		for v := range generator(gen).(chan int) {
			got = append(got, v)
		}
		if !Equal(want, got) {
			t.Errorf("generator() = %v, want %v", got, want)
		}
	})
}

func TestRecv(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		var c <-chan int
		buf := make([]int, 2)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, _, err := Recv(ctx, c, buf)
		cancel()
		if err != context.DeadlineExceeded {
			t.Error("expect Recv timeout")
		}

		ctx, cancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		_, _, err = recv(ctx, c, buf)
		cancel()
		if err != context.DeadlineExceeded {
			t.Error("expect Recv timeout")
		}

		c = SliceToChannel([]int{1, 2, 3})
		for i := 0; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			n, closed, err := Recv(ctx, c, buf)
			cancel()

			var want []int
			if i == 0 {
				want = []int{1, 2}
			} else if i == 1 {
				want = []int{3}
			} else {
				t.Error("too many element to recv")
				break
			}

			if !Equal(want, buf[:n]) {
				t.Errorf("Recv() = %v, want %v", buf[:n], want)
			}

			if err != nil {
				t.Errorf("Recv() error: %v", err)
			}

			if closed {
				break
			}
		}

		c = SliceToChannel([]int{1, 2, 3})
		for i := 0; ; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			n, closed, err := recv(ctx, c, buf)
			cancel()

			var want []int
			if i == 0 {
				want = []int{1, 2}
			} else if i == 1 {
				want = []int{3}
			} else {
				t.Error("too many element to recv")
				break
			}

			if !Equal(want, buf[:n]) {
				t.Errorf("Recv() = %v, want %v", buf[:n], want)
			}

			if err != nil {
				t.Errorf("Recv() error: %v", err)
			}

			if closed {
				break
			}
		}
	})
}

func TestFanIn(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		want := []int{0, 1, 2, 3, 4, 5}

		cs := make([]<-chan int, 3)
		for i := range cs {
			cs[i] = SliceToChannel([]int{i * 2, i*2 + 1})
		}

		var got []int
		c := FanIn(context.Background(), cs...)
		for v := range c {
			got = append(got, v)
		}
		sort.Ints(got)

		if !Equal(want, got) {
			t.Errorf("FanIn() = %v, want %v", got, want)
		}

		for i := range cs {
			cs[i] = SliceToChannel([]int{i * 2, i*2 + 1})
		}

		got = nil
		c = fanIn(context.Background(), cs).(chan int)
		for v := range c {
			got = append(got, v)
		}
		sort.Ints(got)

		if !Equal(want, got) {
			t.Errorf("fanIn() = %v, want %v", got, want)
		}
	})
}

func TestFanOut(t *testing.T) {
	t.Parallel()

	testWithTimeout(t, 5*time.Second, func(t *testing.T) {
		want := []int{0, 1, 2, 3, 4, 5}
		c := SliceToChannel(want)
		cs := FanOut(context.Background(), 3, c)
		if got := len(cs); got != 3 {
			t.Errorf("len(cs) = %d, want %d", got, 3)
		}

		var got []int
		for _, c := range cs {
			for v := range c {
				got = append(got, v)
			}
		}
		sort.Ints(got)

		if !Equal(want, got) {
			t.Errorf("FanOut() = %v, want %v", got, want)
		}

		c = SliceToChannel(want)
		cs = fanOut(context.Background(), 3, c).([]<-chan int)
		if got := len(cs); got != 3 {
			t.Errorf("len(cs) = %d, want %d", got, 3)
		}

		got = nil
		for _, c := range cs {
			for v := range c {
				got = append(got, v)
			}
		}
		sort.Ints(got)

		if !Equal(want, got) {
			t.Errorf("fanOut() = %v, want %v", got, want)
		}
	})
}

func testWithTimeout(t *testing.T, d time.Duration, f func(*testing.T)) {
	t.Helper()

	done := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	go func() {
		f(t)
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		t.Fatalf("test timed out after %v", d)
	}
}
