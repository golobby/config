package feeder_test

import (
	"config/repository"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_JsonDirectory_Extract_Not_Existing_Dir_It_Should_Return_Error(t *testing.T) {
	j := repository.JsonDirectory{Path: "/404"}

	_, err := j.Extract()

	assert.Error(t, err)
}

func Test_JsonDirectory_Extract_Invalid_JSON_Dir_It_Should_Return_Error(t *testing.T) {
	j := repository.JsonDirectory{Path: "../invalid-json"}

	_, err := j.Extract()

	assert.Error(t, err)
}

func Test_JsonDirectory_Extract_Sample1(t *testing.T) {
	j := repository.JsonDirectory{Path: "../test/json"}

	m, err := j.Extract()

	fmt.Println(m)

	assert.NoError(t, err)
	assert.Equal(t, m["app"].(map[string]interface{})["name"], "MyAppUsingConfig")
	assert.Equal(t, m["app"].(map[string]interface{})["version"], 3.14)
}
