package config

import (
	"encoding/json"
	"reflect"
	"testing"
	"strings"

	"github.com/stretchr/testify/assert"
)

type Data = map[string]interface{}

type User struct {
	Name string `json:"name"`
	Year int `json:"year"`
}

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

func prepareData() (Data, error) {
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

func Test_AssignStructFields_Plain(t *testing.T) {
	assert := assert.New(t)

	data := Data{"name": "Milad Rahimi", "year": 1993}

	var ptr *User

	assert.Equal(0, assignStruct(ptr, data, "json"))

	assert.Equal(2, assignStruct(&ptr, data, "json"))
	assert.Equal(data["name"], ptr.Name)
	assert.Equal(data["year"], ptr.Year)

	ptr = &User{}
	assert.Equal(2, assignStruct(ptr, data, "json"))
	assert.Equal(data["name"], ptr.Name)
	assert.Equal(data["year"], ptr.Year)
}

func Test_AssignStructFields_Nested_0(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareData()
	assert.NoError(err)

	ptr := &UserWithAddr_0{}
	assert.Equal(3, assignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
}

func Test_AssignStructFields_Nested_1(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareData()
	assert.NoError(err)

	ptr := &UserWithAddr_1{}
	assert.Equal(3, assignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
}

func Test_AssignStructFields_Nested_2(t *testing.T) {
	assert := assert.New(t)

	user, err := prepareData()
	assert.NoError(err)

	ptr := &UserWithAddr_2{}
	assert.Equal(3, assignStruct(ptr, user, "json"))
	assert.Equal(user["name"], ptr.Name)
	assert.Equal(user["year"], float64(ptr.Year))

	addr := user["address"].(Data)
	assert.Equal(addr["city"], ptr.Addr.City)
	assert.Equal(addr["country"], ptr.Addr.Country)
	assert.Equal(addr["state"], ptr.Addr.State)
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
	objPtr, _, ok := checkPtr(ptr)
	if !ok {
		return 0
	}

	if objPtr.Kind() == reflect.Ptr {
		return 2
	} else {
		return 1
	}
}
