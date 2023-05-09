package slices_test

import (
	"math/rand"
	"reflect"
	"runtime"
	"sync"
)

// This file contains reference slices functions implementations for unit-tests.

func filter(a any, f any) any {
	s := reflect.ValueOf(a)
	if s.IsNil() {
		return reflect.Zero(s.Type()).Interface()
	}

	fv := reflect.ValueOf(f)
	r := reflect.MakeSlice(s.Type(), 0, s.Len())
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		if fv.Call([]reflect.Value{reflect.ValueOf(i), v})[0].Bool() {
			r = reflect.Append(r, v)
		}
	}
	return r.Interface()
}

func pfilter(a any, f any) any {
	s := reflect.ValueOf(a)
	if s.IsNil() {
		return reflect.Zero(s.Type()).Interface()
	}

	fv := reflect.ValueOf(f)
	ngoroutines := runtime.NumCPU()
	n := s.Len()
	step := n / ngoroutines
	if step == 0 {
		step = 1
	}

	c := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, s.Type().Elem()), ngoroutines)

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
		go func(s reflect.Value) {
			for i := 0; i < s.Len(); i++ {
				v := s.Index(i)
				if fv.Call([]reflect.Value{reflect.ValueOf(start + i), v})[0].Bool() {
					c.Send(v)
				}
			}
			wg.Done()
		}(s.Slice(start, end))
	}

	go func() {
		wg.Wait()
		c.Close()
	}()

	r := reflect.MakeSlice(s.Type(), 0, s.Len())
	for {
		v, ok := c.Recv()
		if !ok {
			break
		}

		r = reflect.Append(r, v)
	}
	return r.Interface()
}

func _map(a any, f any) any {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)
	if s.IsNil() {
		return reflect.Zero(reflect.SliceOf(fv.Type().Out(0))).Interface()
	}

	r := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		r.Index(i).Set(fv.Call([]reflect.Value{reflect.ValueOf(i), s.Index(i)})[0])
	}
	return r.Interface()
}

func pmap(a any, f any) any {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)
	if s.IsNil() {
		return reflect.Zero(reflect.SliceOf(fv.Type().Out(0))).Interface()
	}

	ngoroutines := runtime.NumCPU()
	n := s.Len()
	step := n / ngoroutines
	if step == 0 {
		step = 1
	}

	typeOfResult := reflect.StructOf([]reflect.StructField{
		{
			Name: "Index",
			Type: reflect.TypeOf(int(0)),
		},
		{
			Name: "Value",
			Type: fv.Type().Out(0),
		},
	})

	c := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, typeOfResult), ngoroutines)

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
		go func(s reflect.Value) {
			for i := 0; i < s.Len(); i++ {
				index := reflect.ValueOf(start + i)
				value := fv.Call([]reflect.Value{index, s.Index(i)})[0]

				v := reflect.New(typeOfResult).Elem()
				v.Field(0).SetInt(index.Int())
				v.Field(1).Set(value)
				c.Send(v)
			}
			wg.Done()
		}(s.Slice(start, end))
	}

	go func() {
		wg.Wait()
		c.Close()
	}()

	r := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), s.Len(), s.Len())
	for {
		v, ok := c.Recv()
		if !ok {
			break
		}

		r.Index(int(v.Field(0).Int())).Set(v.Field(1))
	}
	return r.Interface()
}

func reduce(a any, f any, init any) any {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)

	acc := reflect.ValueOf(init)
	for i := 0; i < s.Len(); i++ {
		acc = fv.Call([]reflect.Value{acc, reflect.ValueOf(i), s.Index(i)})[0]
	}
	return acc.Interface()
}

func forEach(a any, f any) {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)
	for i := 0; i < s.Len(); i++ {
		fv.Call([]reflect.Value{reflect.ValueOf(i), s.Index(i)})
	}
}

func pForEach(a any, f any) {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)
	ngoroutines := runtime.NumCPU()
	n := s.Len()
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
		go func(s reflect.Value) {
			for i := 0; i < s.Len(); i++ {
				fv.Call([]reflect.Value{reflect.ValueOf(start + i), s.Index(i)})
			}
			wg.Done()
		}(s.Slice(start, end))
	}

	wg.Wait()
}

func shuffle(a any) any {
	s := reflect.ValueOf(a)
	rand.Shuffle(s.Len(), reflect.Swapper(a))
	return s.Interface()
}

func reverse(a any) any {
	s := reflect.ValueOf(a)
	n := s.Len()
	h := n / 2
	swap := reflect.Swapper(a)
	for i := 0; i < h; i++ {
		j := n - i - 1
		swap(i, j)
	}
	return s.Interface()
}

func fill(a any, init any) any {
	s := reflect.ValueOf(a)
	v := reflect.ValueOf(init)
	for i := 0; i < s.Len(); i++ {
		s.Index(i).Set(v)
	}
	return s.Interface()
}

func fillFunc(a any, f any) any {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)
	for i := 0; i < s.Len(); i++ {
		v := fv.Call([]reflect.Value{reflect.ValueOf(i)})[0]
		s.Index(i).Set(v)
	}
	return s.Interface()
}

func repeat(init any, count int) any {
	if count == 0 {
		return reflect.Zero(reflect.SliceOf(reflect.TypeOf(init))).Interface()
	}

	initv := reflect.ValueOf(init)
	s := reflect.MakeSlice(reflect.SliceOf(initv.Type()), count, count)
	for i := 0; i < s.Len(); i++ {
		s.Index(i).Set(initv)
	}
	return s.Interface()
}

func repeatFunc(f any, count int) any {
	fv := reflect.ValueOf(f)
	if count == 0 {
		return reflect.Zero(reflect.SliceOf(fv.Type().Out(0))).Interface()
	}

	s := reflect.MakeSlice(reflect.SliceOf(fv.Type().Out(0)), count, count)
	for i := 0; i < s.Len(); i++ {
		v := fv.Call([]reflect.Value{reflect.ValueOf(i)})[0]
		s.Index(i).Set(v)
	}
	return s.Interface()
}

func count(a, v any) int {
	s := reflect.ValueOf(a)
	vv := reflect.ValueOf(v)

	count := 0
	for i := 0; i < s.Len(); i++ {
		if vv.Equal(s.Index(i)) {
			count++
		}
	}
	return count
}

func countFunc(a any, eq any) int {
	s := reflect.ValueOf(a)
	eqv := reflect.ValueOf(eq)

	count := 0
	for i := 0; i < s.Len(); i++ {
		if eqv.Call([]reflect.Value{s.Index(i)})[0].Bool() {
			count++
		}
	}
	return count
}
