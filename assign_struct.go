package config

import (
	"reflect"
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
