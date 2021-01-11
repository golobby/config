// Copyright 2021 Zhaoping Yu.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config_test

import (
	"encoding/json"
	"testing"
	"strings"

	"github.com/golobby/config"

	"github.com/stretchr/testify/assert"
)

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	State   string `json:"state"`
}

type UserWithAddr_0 struct {
	Name string   `json:"name"`
	Year int      `json:"year"`
	Addr Address  `json:"address"`
}

type UserWithAddr_1 struct {
	User
	Addr Address  `json:"address"`
}

type UserWithAddr_2 struct {
	User
	Addr *Address `json:"address"`
}

func prepareStructData() (Data, error) {
	jsonStr := `
    {
      "name": "Milad Rahimi",
      "year": 1993,
      "address": {
        "country": "Iran",
        "state": "Lorestan",
        "city": "Delfan"
      }
    }
	`

	dec := json.NewDecoder(strings.NewReader(jsonStr))

	data := Data{}

	return data, dec.Decode(&data)
}

func Test_AssignStruct_Plain(t *testing.T) {
	assert := assert.New(t)

	data := Data{"name": "Milad Rahimi", "year": 1993}

	var ptr *User

	assert.Equal(0, config.AssignStruct(ptr, data, "json"))

	assert.Equal(2, config.AssignStruct(&ptr, data, "json"))
	assert.Equal(data["name"], ptr.Name)
	assert.Equal(data["year"], ptr.Year)

	ptr = &User{}
	assert.Equal(2, config.AssignStruct(ptr, data, "json"))
	assert.Equal(data["name"], ptr.Name)
	assert.Equal(data["year"], ptr.Year)
}

func Test_AssignStruct_Nested_0(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareStructData()
	assert.NoError(err)

	ptr := &UserWithAddr_0{}
	assert.Equal(3, config.AssignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
}

func Test_AssignStruct_Nested_1(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareStructData()
	assert.NoError(err)

	ptr := &UserWithAddr_1{}
	assert.Equal(3, config.AssignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
}

func Test_AssignStruct_Nested_2(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareStructData()
	assert.NoError(err)

	ptr := &UserWithAddr_2{}
	assert.Equal(3, config.AssignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
}
