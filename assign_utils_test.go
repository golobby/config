// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data = map[string]interface{}

type User struct {
	Name string `json:"name"`
	Year int `json:"year"`
}

func Test_CheckPtr(t *testing.T) {
	assert := assert.New(t)

	var ptr *User
	assert.Equal(0, testCheckPtr(ptr))
	assert.Equal(2, testCheckPtr(&ptr))

	ptr = &User{}
	assert.Equal(1, testCheckPtr(ptr))
}

func testCheckPtr(ptr interface{}) int {
	objPtr, obj, ok := checkPtr(ptr)
	if !ok {
		return 0
	} else if obj.Kind() != reflect.Struct {
		return 0
	}

	if objPtr.Kind() == reflect.Ptr {
		return 2
	} else {
		return 1
	}
}
