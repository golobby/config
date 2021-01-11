// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package assign_test

import (
	"encoding/json"
	"testing"
	"strings"

	"github.com/golobby/config/assign"

	"github.com/stretchr/testify/assert"
)

type UserS struct {
	Users []User `json:"users"`
}

func prepareSliceSimpleData() (Data, error) {
	jsonStr := `
    { "names": [
      "Milad Rahimi",
      "Amirreza Askarpour"
    ]}
	`

	dec := json.NewDecoder(strings.NewReader(jsonStr))

	data := Data{}
	return data, dec.Decode(&data)
}

func prepareSliceStructData() (Data, error) {
	jsonStr := `
    { "users": [
      {"name": "Milad Rahimi", "year": 1993},
      {"name": "Amirreza Askarpour", "year": 1998}
    ]}
	`

	dec := json.NewDecoder(strings.NewReader(jsonStr))

	data := Data{}
	return data, dec.Decode(&data)
}

func Test_AssignSlice_Simple(t *testing.T) {
	assert := assert.New(t)

	data, err := prepareSliceSimpleData()
	assert.NoError(err)

	names := data["names"].([]interface{})
	namesLen := len(names)

	var ptr []string

	assert.Equal(namesLen, assign.AssignSlice(&ptr, data["names"], "json"))
	assert.Equal(namesLen, len(ptr))

	for i := 0; i < namesLen; i++ {
		assert.Equal(ptr[i], names[i])
	}
}

func Test_AssignSlice_Struct(t *testing.T) {
	assert := assert.New(t)

	data, err := prepareSliceStructData()
	assert.NoError(err)

	users := data["users"].([]interface{})
	usersLen := len(users)

	var ptr []User

	assert.Equal(usersLen, assign.AssignSlice(&ptr, data["users"], "json"))
	assert.Equal(usersLen, len(ptr))

	for i := 0; i < usersLen; i++ {
		user := users[i].(map[string]interface{})

		assert.Equal(user["name"], ptr[i].Name)
		assert.Equal(user["year"], float64(ptr[i].Year))
	}
}

func Test_AssignStruct_Slice(t *testing.T) {
	assert := assert.New(t)

	data, err := prepareSliceStructData()
	assert.NoError(err)

	ptr := &UserS{}

	assert.Equal(1, assign.AssignStruct(ptr, data, "json"))

	users := data["users"].([]interface{})
	usersLen := len(users)

	assert.Equal(usersLen, len(ptr.Users))

	for i := 0; i < usersLen; i++ {
		user := users[i].(map[string]interface{})

		assert.Equal(user["name"], ptr.Users[i].Name)
		assert.Equal(user["year"], float64(ptr.Users[i].Year))
	}
}
