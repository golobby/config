package config

import (
	"reflect"
	"strings"
)

// Assigns struct fields' value by its field's tag (such as the json tag).
// @param ptr The pointer of struct's instance to set
// @param key Specify where to get the struct's value
// @param tag Specify which struct field's tag name used to retrieve
// @return The count of fields that been assigned, -1 if struct's value not found by the key
func (c *ConfigBase) AssignStruct(ptr interface{}, key, tag string) int {
	if data, found := c.Get(key); found {
		return assignStruct(ptr, data, tag)
	}
	return -1
}

// Assigns struct fields' value by its field's tag (such as the json tag).
// @param ptr The pointer of struct's instance to set
// @param key Specify where to get the struct's value
// @param tag Specify which struct field's tag name used to retrieve
// @return The count of fields that been assigned, -1 if struct's value not found by the key
func (c *Config) AssignStruct(ptr interface{}, key, tag string) int {
	c.sync.RLock()
	defer c.sync.RUnlock()

	return c.ConfigBase.AssignStruct(ptr, key, tag)
}

// Assigns struct fields' value by its field's tag (such as the json tag).
// @param ptr The pointer of struct's instance to set
// @param data The data map that stores struct fields' tag/value pair
// @param tag Specify which struct field's tag name used to retrieve
// @return The count of fields that been assigned
func AssignStruct(ptr interface{}, data map[string]interface{}, tag string) int {
	return assignStruct(ptr, data, tag)
}

func assignStruct(ptr, data interface{}, tag string) int {
	objPtr, obj, ok := checkPtr(ptr)
	if !ok {
		return 0
	}

	count := assign(obj, data, tag)

	if count > 0 {
		if objPtr.Kind() == reflect.Ptr {
			objPtr.Set(obj.Addr())
		} else {
			objPtr.Set(obj)
		}
	}

	return count
}

func checkPtr(ptr interface{}) (refPtr, refVal reflect.Value, ok bool) {
	v := reflect.ValueOf(ptr)
	if !v.IsValid() {
		return v, v, false
	}

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

	ok = (tv.Kind() == reflect.Struct)

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
