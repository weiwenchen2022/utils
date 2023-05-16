// Package channels defines various types and functions useful with channels of any type.
package channels

import (
	"context"
	"sync"
	"time"
)

// ErrTimeout is the error returned by SendTimed.
var ErrTimeout error = timeoutError{}

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

// Channel attaches the common methods to chan E.
type Channel[E any] chan E

// NewChannel converts c to Channel type.
func NewChannel[C ~chan E, E any](c C) Channel[E] {
	return Channel[E](c)
}

// Recv is a convenience method: c.Recv(ctx, s) returns Recv(ctx, c, s).
func (c Channel[E]) Recv(ctx context.Context, s []E) (n int, closed bool, err error) {
	return Recv(ctx, c, s)
}

// SendTimed is a convenience method: c.SendTimed(x, s, d) returns SendTimed(c, x, d).
func (c Channel[E]) SendTimed(x E, d time.Duration) error {
	return SendTimed(c, x, d)
}

// RecvOnlyChannel attaches the common methods to <-chan E.
type RecvOnlyChannel[E any] <-chan E

// NewRecvOnlyChannel converts c to RecvOnlyChannel type.
func NewRecvOnlyChannel[C ~<-chan E, E any](c C) RecvOnlyChannel[E] {
	return RecvOnlyChannel[E](c)
}

// Recv is a convenience method: c.Recv(ctx, s) returns Recv(ctx, c, s).
func (c RecvOnlyChannel[E]) Recv(ctx context.Context, s []E) (n int, closed bool, err error) {
	return Recv(ctx, c, s)
}

// SliceToChannel returns a receive only channel of elements of s.
// Channel is closed after s has been exhausted.
func SliceToChannel[S ~[]E, E any](s S) <-chan E {
	c := make(chan E)

	go func() {
		for _, v := range s {
			c <- v
		}

		close(c)
	}()

	return c
}

// ChannelToSlice returns a slice consists of elements from channel c.
// Returned after c closed.
func ChannelToSlice[E any](c <-chan E) []E {
	var r []E
	for v := range c {
		r = append(r, v)
	}
	return r
}

// Generator implements the generator design pattern.
// The returned channel will closed after generator returned.
func Generator[E any](generator func(yield func(E))) <-chan E {
	c := make(chan E)

	go func() {
		generator(func(e E) {
			c <- e
		})

		close(c)
	}()

	return c
}

// Recv reads up to len(s) elements into s. It returns the number of elements recv (0
// <= n <= len(s)), and boolean value closed indicate the channel c is closed.
// When Recv encounters a closed condition after successfully
// reading n > 0 elements, it returns the number of elements receive.
// Callers should always process the n > 0 elements returned before considering
// the closed condition. Doing so correctly handles closed condition that happen after
// receiving some elements and also both of the allowed closed behaviors.
func Recv[E any](ctx context.Context, c <-chan E, s []E) (n int, closed bool, err error) {
	for i := range s {
		select {
		case <-ctx.Done():
			return n, false, ctx.Err()
		case v, ok := <-c:
			if !ok {
				return n, true, nil
			}

			s[i] = v
			n++
		}
	}

	return n, false, nil
}

// SendTimed sends x on the channel c, with timeout.
func SendTimed[E any](c chan<- E, x E, d time.Duration) error {
	t := time.NewTimer(d)

	select {
	case <-t.C:
		return ErrTimeout
	case c <- x:
	}

	if !t.Stop() {
		<-t.C
	}

	return nil
}

// FanIn returns a channel that's elements from provided input channels.
// They're receive concurrently. Once all inputs channels have closed, then returned channel will close.
func FanIn[C ~<-chan E, E any](ctx context.Context, cs ...C) <-chan E {
	var wg sync.WaitGroup
	out := make(chan E)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or ctx.Done() is closed, then calls
	// wg.Done.
	output := func(c <-chan E) {
		defer wg.Done()

		for n := range c {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// FanOut returns n channels that's elements from provided input channel.
// Once input channel have closed, then the returned channels will close.
func FanOut[E any](ctx context.Context, n int, input <-chan E) []<-chan E {
	cs := make([]<-chan E, n)

	output := func(c chan<- E) {
		defer close(c)

		for v := range input {
			select {
			case <-ctx.Done():
				return
			case c <- v:
			}
		}
	}

	for i := range cs {
		c := make(chan E)
		go output(c)
		cs[i] = c
	}

	return cs
}
