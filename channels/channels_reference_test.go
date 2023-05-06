package channels_test

import (
	"context"
	"reflect"
	"sync"
)

// This file contains reference channels functions implementations for unit-tests.

func sliceToChannel(a any) any {
	s := reflect.ValueOf(a)
	c := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, s.Type().Elem()), 0)

	go func() {
		for i := 0; i < s.Len(); i++ {
			c.Send(s.Index(i))
		}

		c.Close()
	}()

	return c.Interface()
}

func channelToSlice(a any) any {
	c := reflect.ValueOf(a)
	s := reflect.MakeSlice(reflect.SliceOf(c.Type().Elem()), 0, 0)

	for {
		v, ok := c.Recv()
		if !ok {
			break
		}

		s = reflect.Append(s, v)
	}

	return s.Interface()
}

func generator[E any](gen func(func(E))) any {
	genFn := reflect.ValueOf(gen)
	typeOfYield := genFn.Type().In(0)
	typeOfE := typeOfYield.In(0)
	c := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, typeOfE), 0)

	yieldImpl := func(in []reflect.Value) []reflect.Value {
		c.Send(in[0])
		return nil
	}

	makeYield := func(fptr any) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), yieldImpl)
		fn.Set(v)
	}

	go func() {
		var yield func(E)
		makeYield(&yield)

		genFn.Call([]reflect.Value{reflect.ValueOf(yield)})

		c.Close()
	}()

	return c.Interface()
}

func recv(ctx context.Context, a any, p any) (n int, closed bool, err error) {
	c := reflect.ValueOf(a)
	s := reflect.ValueOf(p)

	cases := []reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ctx.Done()),
		},
		{
			Dir:  reflect.SelectRecv,
			Chan: c,
		},
	}

	for i := 0; i < s.Len(); i++ {
		chosen, v, ok := reflect.Select(cases)
		switch chosen {
		case 0:
			return n, false, ctx.Err()
		case 1:
			if !ok {
				return n, true, nil
			}

			s.Index(i).Set(v)
			n++
		}
	}

	return n, false, nil
}

func fanIn(ctx context.Context, a any) any {
	cs := reflect.ValueOf(a)

	var wg sync.WaitGroup

	typeOfChan := reflect.ChanOf(reflect.BothDir, cs.Index(0).Type().Elem())
	out := reflect.MakeChan(typeOfChan, 0)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c or ctx.Done() is closed, then calls
	// wg.Done.
	output := func(c reflect.Value) {
		defer wg.Done()

		cases := []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ctx.Done()),
			},
			{
				Dir:  reflect.SelectSend,
				Chan: out,
			},
		}

		for {
			n, ok := c.Recv()
			if !ok {
				break
			}

			cases[1].Send = n
			chosen, _, _ := reflect.Select(cases)
			if chosen == 0 {
				return
			}
		}
	}

	wg.Add(cs.Len())
	for i := 0; i < cs.Len(); i++ {
		go output(cs.Index(i))
	}

	// Start a goroutine to close out once all the output goroutines are
	// done. This must start after the wg.Add call.
	go func() {
		wg.Wait()
		out.Close()
	}()

	return out.Interface()
}

func fanOut(ctx context.Context, n int, a any) any {
	input := reflect.ValueOf(a)
	typeOfChan := reflect.ChanOf(reflect.RecvDir, input.Type().Elem())
	cs := reflect.MakeSlice(reflect.SliceOf(typeOfChan), n, n)

	output := func(c reflect.Value) {
		defer c.Close()

		cases := []reflect.SelectCase{
			{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ctx.Done()),
			},
			{
				Dir:  reflect.SelectSend,
				Chan: c,
			},
		}

		for {
			v, ok := input.Recv()
			if !ok {
				break
			}

			cases[1].Send = v
			chosen, _, _ := reflect.Select(cases)
			if chosen == 0 {
				return
			}
		}
	}

	for i := 0; i < cs.Len(); i++ {
		typeOfChan := reflect.ChanOf(reflect.BothDir, input.Type().Elem())
		c := reflect.MakeChan(typeOfChan, 0)
		go output(c)

		cs.Index(i).Set(c)
	}

	return cs.Interface()
}
