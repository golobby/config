// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package assign

import (
	"reflect"
)

// Assigns slice elements.
// @param ptr The pointer of slice instance to appent elements
// @param data The data that stores elements' value
// @param tag If element's type is struct, using the tag name to retrieve struct fields
// @return The count of elements that been assigned
func AssignSlice(ptr, data interface{}, tag string) int {

	objPtr, obj, ok := checkPtr(ptr)
	if !ok {
		return 0
	} else if obj.Kind() != reflect.Slice {
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
