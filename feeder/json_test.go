package feeder_test

import (
	"github.com/golobby/config/v2/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson_Feed(t *testing.T) {
	j := feeder.Json{Path: "test/config.json"}

	m, err := j.Feed()
	assert.NoError(t, err)

	assert.Equal(t, "MyAppUsingConfig", m["name"])
	assert.Equal(t, 3.14, m["version"])
	assert.Equal(t, "Milad Rahimi", m["users"].([]interface{})[0].(map[string]interface{})["name"])
}

func TestJson_Feed_With_Invalid_JSON_Path_It_Should_Fail(t *testing.T) {
	j := feeder.Json{Path: "/404.json"}

	_, err := j.Feed()
	assert.Error(t, err)
}

func TestJson_Feed_With_Invalid_JSON_File_It_Should_Fail(t *testing.T) {
	j := feeder.Json{Path: "test/invalid.json"}

	_, err := j.Feed()
	assert.Error(t, err)
}
