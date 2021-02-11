package feeder_test

import (
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonDirectory_Feed(t *testing.T) {
	j := feeder.JsonDirectory{Path: "test/json"}

	m, err := j.Feed()

	assert.NoError(t, err)
	assert.Equal(t, m["app"].(map[string]interface{})["name"], "MyAppUsingConfig")
	assert.Equal(t, m["app"].(map[string]interface{})["version"], 3.14)
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
