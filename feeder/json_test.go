package feeder_test

import (
	"config/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Json_Extract_Not_Existing_File_It_Should_Return_Error(t *testing.T) {
	j := repository.Json{Path: "/404.json"}

	_, err := j.Extract()

	assert.Error(t, err)
}

func Test_Json_Extract_Invalid_JSON_It_Should_Return_Error(t *testing.T) {
	j := repository.Json{Path: "../test/invalid.json"}

	_, err := j.Extract()

	assert.Error(t, err)
}

func Test_Json_Extract_Sample1(t *testing.T) {
	j := repository.Json{Path: "../test/config.json"}

	m, err := j.Extract()

	assert.NoError(t, err)
	assert.Equal(t, "MyAppUsingConfig", m["name"])
	assert.Equal(t, 3.14, m["version"])
}
