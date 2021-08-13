package feeder_test

import (
	"github.com/golobby/config/v3/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonDirectory_Feed(t *testing.T) {
	j := feeder.JsonDirectory{Path: "test/json"}

	m, err := j.Feed()

	assert.NoError(t, err)
	assert.Equal(t, "MyAppUsingConfig", m["app"].(map[string]interface{})["name"])
	assert.Equal(t, 3.14, m["app"].(map[string]interface{})["version"])
	assert.Equal(t, "mysql", m["db"].(map[string]interface{})["default"])
}

func TestJsonDirectory_Feed_With_Invalid_Dir_Path_It_Should_Fail(t *testing.T) {
	j := feeder.JsonDirectory{Path: "/404"}

	_, err := j.Feed()
	assert.Error(t, err)
}

func TestJsonDirectory_Feed_With_Invalid_Dir_File_It_Should_Fail(t *testing.T) {
	j := feeder.JsonDirectory{Path: "test/invalid-json"}

	_, err := j.Feed()
	assert.Error(t, err)
}
