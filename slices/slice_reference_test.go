package slices_test

import "reflect"

// This file contains reference slices functions implementations for unit-tests.

func equalTo(a1, a2 any) bool {
	s1, s2 := reflect.ValueOf(a1), reflect.ValueOf(a2)
	if s1.Len() != s2.Len() {
		return false
	}

	for i := 0; i < s1.Len(); i++ {
		if !s1.Index(i).Equal(s2.Index(i)) {
			return false
		}
	}

	return true
}

func equalFunc(a1, a2 any, eq func(any, any) bool) bool {
	s1, s2 := reflect.ValueOf(a1), reflect.ValueOf(a2)

	if s1.Len() != s2.Len() {
		return false
	}

	eqv := reflect.ValueOf(eq)

	for i := 0; i < s1.Len(); i++ {
		v1 := s1.Index(i)
		v2 := s2.Index(i)
		if !eqv.Call([]reflect.Value{v1, v2})[0].Bool() {
			return false
		}
	}

	return true
}

func index(a any, v any) int {
	s := reflect.ValueOf(a)
	vv := reflect.ValueOf(v)

	for i := 0; i < s.Len(); i++ {
		if vv.Equal(s.Index(i)) {
			return i
		}
	}

	return -1
}

func indexFunc(a any, f func(any) bool) int {
	s := reflect.ValueOf(a)
	fv := reflect.ValueOf(f)

	for i := 0; i < s.Len(); i++ {
		if fv.Call([]reflect.Value{s.Index(i)})[0].Bool() {
			return i
		}
	}

	return -1
}

func contains(a any, v any) bool {
	return index(a, v) >= 0
}

func containsFunc(s any, f func(any) bool) bool {
	return indexFunc(s, f) >= 0
}

func insert(a any, i int, v ...any) any {
	s := reflect.ValueOf(a)
	vv := reflect.MakeSlice(s.Type(), len(v), len(v))
	for i := range v {
		vv.Index(i).Set(reflect.ValueOf(v[i]))
	}

	tot := s.Len() + vv.Len()
	if tot <= s.Cap() {
		s2 := s.Slice(0, tot)
		reflect.Copy(s2.Slice(i+len(v), s2.Len()), s.Slice(i, s.Len()))
		reflect.Copy(s2.Slice(i, s2.Len()), vv)
		return s2.Interface()
	}

	s2 := reflect.MakeSlice(s.Type(), tot, tot)
	reflect.Copy(s2, s.Slice(0, i))
	reflect.Copy(s2.Slice(i, s2.Len()), vv)
	reflect.Copy(s2.Slice(i+len(v), s2.Len()), s.Slice(i, s.Len()))
	return s2.Interface()
}

func delete(a any, i, j int) any {
	s := reflect.ValueOf(a)
	_ = s.Slice(i, j)

	s2 := reflect.AppendSlice(s.Slice(0, i), s.Slice(j, s.Len()))
	reflect.Copy(s.Slice(s.Len()-(j-i), s.Len()), reflect.MakeSlice(s.Type(), (j-i), (j-i)))

	return s2.Interface()
}

func replace(a any, i, j int, v ...any) any {
	s := reflect.ValueOf(a)
	_ = s.Slice(i, j)

	vv := reflect.MakeSlice(s.Type(), len(v), len(v))
	for i := range v {
		vv.Index(i).Set(reflect.ValueOf(v[i]))
	}

	tot := s.Slice(0, i).Len() + len(v) + s.Slice(j, s.Len()).Len()
	if tot <= s.Cap() {
		s2 := s.Slice(0, tot)
		reflect.Copy(s2.Slice(i+len(v), s2.Len()), s.Slice(j, s.Len()))
		reflect.Copy(s2.Slice(i, s2.Len()), vv)
		return s2.Interface()
	}

	s2 := reflect.MakeSlice(s.Type(), tot, tot)
	reflect.Copy(s2, s.Slice(0, i))
	reflect.Copy(s2.Slice(i, s2.Len()), vv)
	reflect.Copy(s2.Slice(i+len(v), s2.Len()), s.Slice(j, s.Len()))
	return s2.Interface()
}

func clone(a any) any {
	s := reflect.ValueOf(a)
	if s.IsNil() {
		return s.Interface()
	}

	return reflect.AppendSlice(reflect.MakeSlice(s.Type(), 0, s.Len()), s).Interface()
}

func deepClone(a any) any {
	s := reflect.ValueOf(a)
	if s.IsNil() {
		return s.Interface()
	}

	typeOfE := s.Type().Elem()
	if method, ok := typeOfE.MethodByName("Clone"); ok {
		methodType := method.Type
		if methodType.NumIn() == 1 &&
			methodType.NumOut() == 1 && methodType.Out(0) == typeOfE {
			s2 := reflect.MakeSlice(s.Type(), s.Len(), s.Len())

			for i := 0; i < s.Len(); i++ {
				elem := method.Func.Call([]reflect.Value{s.Index(i)})[0]
				s2.Index(i).Set(elem)
			}

			return s2.Interface()
		}
	}

	return clone(a)
}

func compact(a any) any {
	s := reflect.ValueOf(a)
	if s.Len() < 2 {
		return s.Interface()
	}

	i := 1
	for j := 1; j < s.Len(); j++ {
		if !s.Index(j - 1).Equal(s.Index(j)) {
			if j != i {
				s.Index(i).Set(s.Index(j))
			}

			i++
		}
	}

	return s.Slice(0, i).Interface()
}

func compactFunc(a any, eq func(any, any) bool) any {
	s := reflect.ValueOf(a)
	if s.Len() < 2 {
		return s.Interface()
	}

	eqv := reflect.ValueOf(eq)

	i := 1
	for j := 1; j < s.Len(); j++ {
		if !eqv.Call([]reflect.Value{s.Index(j - 1), s.Index(j)})[0].Bool() {
			if j != i {
				s.Index(i).Set(s.Index(j))
			}

			i++
		}
	}

	return s.Slice(0, i).Interface()
}

func grow(a any, n int) any {
	if n < 0 {
		panic("cannot be negative")
	}

	s := reflect.ValueOf(a)
	if n -= s.Cap() - s.Len(); n > 0 {
		s = reflect.AppendSlice(s.Slice(0, s.Cap()), reflect.MakeSlice(s.Type(), n, n)).Slice(0, s.Len())
	}

	return s.Interface()
}

func clip(a any) any {
	s := reflect.ValueOf(a)
	return s.Slice3(0, s.Len(), s.Len()).Interface()
}
