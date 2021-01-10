package config

import (
	"reflect"
	"strings"
)

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
	case reflect.Map:
		// TODO
		return 0
	case reflect.Array, reflect.Slice:
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
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				pv := reflect.New(fv.Type().Elem())
				if assign(pv, fsv.Interface(), tag) > 0 {
					fv.Set(pv)
					count++
				}
			} else {
				if assign(fv, fsv.Interface(), tag) > 0 {
					count++
				}
			}
		} else if ft.Anonymous {
			if fv.Kind() == reflect.Ptr && fv.IsNil() {
				pv := reflect.New(fv.Type().Elem())
				if ret := assign(pv, data, tag); ret > 0 {
					fv.Set(pv)
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
