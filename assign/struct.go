// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package assign

import (
	"reflect"
)

// Assigns struct fields' value by its field's tag (such as the json tag).
// @param ptr The pointer of struct's instance to set
// @param data The data map that stores struct fields' tag/value pair
// @param tag Specify which struct field's tag name used to retrieve
// @return The count of fields that been assigned
func AssignStruct(ptr, data interface{}, tag string) int {

	objPtr, obj, ok := checkPtr(ptr)
	if !ok {
		return 0
	} else if obj.Kind() != reflect.Struct {
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
