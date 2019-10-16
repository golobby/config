package feeder_test

import (
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Json_Feed_Not_Existing_File_It_Should_Return_Error(t *testing.T) {
	j := feeder.Json{Path: "/404.json"}

	_, err := j.Feed()

	assert.Error(t, err)
}

func Test_Json_Feed_Invalid_JSON_It_Should_Return_Error(t *testing.T) {
	j := feeder.Json{Path: "test/invalid.json"}

	_, err := j.Feed()

	assert.Error(t, err)
}

func Test_Json_Feed_Sample1(t *testing.T) {
	j := feeder.Json{Path: "test/config.json"}

	m, err := j.Feed()

	assert.NoError(t, err)
	assert.Equal(t, "MyAppUsingConfig", m["name"])
	assert.Equal(t, 3.14, m["version"])
}
