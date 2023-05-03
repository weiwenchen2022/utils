package maps_test

import "reflect"

// This file contains reference maps functions implementations for unit-tests.

func keys(a any) any {
	m := reflect.ValueOf(a)
	typeOfSliceOfK := reflect.SliceOf(m.Type().Key())
	if m.Len() == 0 {
		return reflect.Zero(typeOfSliceOfK).Interface()
	}

	r := reflect.MakeSlice(typeOfSliceOfK, 0, m.Len())
	for iter := m.MapRange(); iter.Next(); {
		r = reflect.Append(r, iter.Key())
	}
	return r.Interface()
}

func values(a any) any {
	m := reflect.ValueOf(a)
	typeOfSliceOfE := reflect.SliceOf(m.Type().Elem())
	if m.Len() == 0 {
		return reflect.Zero(typeOfSliceOfE).Interface()
	}

	r := reflect.MakeSlice(typeOfSliceOfE, 0, m.Len())
	for iter := m.MapRange(); iter.Next(); {
		r = reflect.Append(r, iter.Value())
	}
	return r.Interface()
}

func equals(a1, a2 any) bool {
	m1 := reflect.ValueOf(a1)
	m2 := reflect.ValueOf(a2)

	if m1.Len() != m2.Len() {
		return false
	}

	for iter := m1.MapRange(); iter.Next(); {
		k := iter.Key()
		if v2 := m2.MapIndex(k); !v2.IsValid() {
			return false
		} else if v1 := iter.Value(); !v1.Equal(v2) {
			return false
		}
	}

	return true
}

func equalFunc(a1, a2 any, eq func(any, any) bool) bool {
	m1 := reflect.ValueOf(a1)
	m2 := reflect.ValueOf(a2)
	eqfn := reflect.ValueOf(eq)

	if m1.Len() != m2.Len() {
		return false
	}

	for iter := m1.MapRange(); iter.Next(); {
		k := iter.Key()
		if v2 := m2.MapIndex(k); !v2.IsValid() {
			return false
		} else if v1 := iter.Value(); !eqfn.Call([]reflect.Value{v1, v2})[0].Bool() {
			return false
		}
	}

	return true
}

func clear(a any) {
	m := reflect.ValueOf(a)
	for iter := m.MapRange(); iter.Next(); {
		m.SetMapIndex(iter.Key(), reflect.Value{})
	}
}

func clone(a any) any {
	m := reflect.ValueOf(a)
	if m.IsNil() {
		return reflect.Zero(m.Type()).Interface()
	}

	copy := reflect.MakeMapWithSize(m.Type(), m.Len())
	for iter := m.MapRange(); iter.Next(); {
		copy.SetMapIndex(iter.Key(), iter.Value())
	}
	return copy.Interface()
}

func copy(dst, src any) {
	dstv := reflect.ValueOf(dst)
	srcv := reflect.ValueOf(src)

	for iter := srcv.MapRange(); iter.Next(); {
		dstv.SetMapIndex(iter.Key(), iter.Value())
	}
}

func deleteFunc(a any, del func(any, any) bool) {
	delfn := reflect.ValueOf(del)
	m := reflect.ValueOf(a)

	for iter := m.MapRange(); iter.Next(); {
		k := iter.Key()
		if delfn.Call([]reflect.Value{k, iter.Value()})[0].Bool() {
			m.SetMapIndex(k, reflect.Value{})
		}
	}
}
