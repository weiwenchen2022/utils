// Package channels defines various types and functions useful with channels of any type.
package channels

import (
	"context"
	"sync"
)

// Channel attaches the common methods to chan E.
type Channel[E any] chan E

// NewChannel converts c to Channel type.
func NewChannel[C ~chan E, E any](c C) Channel[E] {
	return Channel[E](c)
}

// Recv is a convenience method: c.Recv(ctx, s) returns Recv(ctx, c, s).
func (c Channel[E]) Recv(ctx context.Context, s []E) (n int, closed bool, err error) {
	return Recv(ctx, (<-chan E)(c), s)
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
func ChannelToSlice[C ~<-chan E, E any](c C) []E {
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
func Recv[C ~<-chan E, E any](ctx context.Context, c C, s []E) (n int, closed bool, err error) {
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
func FanOut[C ~<-chan E, E any](ctx context.Context, n int, input C) []<-chan E {
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
