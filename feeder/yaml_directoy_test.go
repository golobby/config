package feeder_test

import (
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_YamlDirectory_Feed_Not_Existing_Dir_It_Should_Return_Error(t *testing.T) {
	j := feeder.YamlDirectory{Path: "/404"}

	_, err := j.Feed()

	assert.Error(t, err)
}

func Test_YamlDirectory_Feed_Invalid_JSON_Dir_It_Should_Return_Error(t *testing.T) {
	j := feeder.YamlDirectory{Path: "test/invalid-yaml"}

	_, err := j.Feed()

	assert.Error(t, err)
}

func Test_YamlDirectory_Feed_Sample1(t *testing.T) {
	j := feeder.YamlDirectory{Path: "test/yaml"}

	m, err := j.Feed()

	assert.NoError(t, err)
	assert.Equal(t, m["app"].(map[string]interface{})["name"], "MyAppUsingGoLobbyConfig")
	assert.Equal(t, m["app"].(map[string]interface{})["version"], 3.14)
}
