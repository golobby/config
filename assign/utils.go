// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package assign

import (
	"reflect"
	"strings"
)

func CheckPtr(ptr interface{}) (refPtr, refVal reflect.Value, ok bool) {
	return checkPtr(ptr)
}

func checkPtr(ptr interface{}) (refPtr, refVal reflect.Value, ok bool) {
	v := reflect.ValueOf(ptr)
	if !v.IsValid() {
		return v, v, false
	}

	ok = true

	// Find the final value that pointed to, or nil ptr
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	refPtr = v

	tv := v
	// For nil ptr
	if tv.Kind() == reflect.Ptr && tv.CanSet() {
		// Init the nil ptr
		tv.Set(reflect.New(tv.Type().Elem()))
		tv = tv.Elem()
	}

	refVal = tv

	return
}

// --- assign method ---

func assign(obj reflect.Value, data interface{}, tag string) int {

	src := reflect.ValueOf(data)
	if !src.IsValid() {
		return 0
	}

	if obj.Kind() == reflect.Ptr {
		// Init the nil ptr
		if obj.IsNil() && obj.CanSet() {
			obj.Set(reflect.New(obj.Type().Elem()))
		}
		obj = obj.Elem()
	}

	if !obj.CanSet() {
		return 0
	}

	switch obj.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		return 0

	case reflect.Slice:
		return assign_slice(obj, src, tag)

	case reflect.Map:
		// TODO
		return 0
	case reflect.Array:
		// TODO
		return 0
	}

	if obj.Kind() != reflect.Struct {
		dataType := src.Type()
		objType := obj.Type()

		if objType == dataType {
			obj.Set(src)
		} else if dataType.ConvertibleTo(objType) {
			obj.Set(src.Convert(objType))
		} else {
			return 0
		}
		return 1
	}

	// case reflect.Struct:

	return assign_struct(obj, src, data, tag)
}

// --- assign_struct ---

func assign_struct(obj, src reflect.Value, data interface{}, tag string) int {

	if src.Kind() != reflect.Map || src.Type().Key().Kind() != reflect.String {
		return 0
	}

	count := 0

	for i, fields := 0, obj.NumField(); i < fields; i++ {
		fv := obj.Field(i)

		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		ft := obj.Type().Field(i)

		name := strings.Split(ft.Tag.Get(tag), ",")[0]
		if len(name) == 0 {
			name = strings.ToLower(ft.Name)
		} else if name == "-" {
			continue
		}

		fsv := src.MapIndex(reflect.ValueOf(name))
		if fsv.IsValid() {
			if pv, need := needNewPtr(fv); need {
				if assign(pv, fsv.Interface(), tag) > 0 {
					if fv.Kind() == reflect.Ptr {
						fv.Set(pv)
					}
					count++
				}
			} else {
				if assign(fv, fsv.Interface(), tag) > 0 {
					count++
				}
			}
		} else if ft.Anonymous {
			if pv, need := needNewPtr(fv); need {
				if ret := assign(pv, data, tag); ret > 0 {
					if fv.Kind() == reflect.Ptr {
						fv.Set(pv)
					}
					count += ret
				}
			} else {
				if ret := assign(fv, data, tag); ret > 0 {
					count += ret
				}
			}
		}
	}

	return count
}

func needNewPtr(fv reflect.Value) (reflect.Value, bool) {
	switch fv.Kind() {
	case reflect.Slice:
		return fv.Addr(), true

	case reflect.Ptr:
		if fv.IsNil() {
			return reflect.New(fv.Type().Elem()), true
		}
	}
	return fv, false
}

// --- assign_slice ---

func assign_slice(obj, src reflect.Value, tag string) int {

	if src.Kind() != reflect.Slice {
		return 0
	}

	srcLen := src.Len()
	if srcLen == 0 {
		return 0
	}

	dst := make([]reflect.Value, 0, srcLen)

	objType := obj.Type().Elem()

	count := 0

	for i := 0; i < srcLen; i++ {
		sv := src.Index(i)
		ov := reflect.New(objType)

		if assign(ov, sv.Interface(), tag) > 0 {
			dst = append(dst, ov.Elem())
			count++
		}
	}

	obj.Set(reflect.Append(obj, dst...))

	return count
}
