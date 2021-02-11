package feeder_test

import (
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYamlDirectory_Feed(t *testing.T) {
	j := feeder.YamlDirectory{Path: "test/yaml"}

	m, err := j.Feed()

	assert.NoError(t, err)
	assert.Equal(t, m["app"].(map[string]interface{})["name"], "MyAppUsingGoLobbyConfig")
	assert.Equal(t, m["app"].(map[string]interface{})["version"], 3.14)
}

func TestYamlDirectory_Feed_With_Invalid_Dir_Path_It_Should_Fail(t *testing.T) {
	j := feeder.YamlDirectory{Path: "/404"}

	_, err := j.Feed()
	assert.Error(t, err)
}

func TestYamlDirectory_Feed_With_Invalid_Dir_File_It_Should_Fail(t *testing.T) {
	j := feeder.YamlDirectory{Path: "test/invalid-yaml"}

	_, err := j.Feed()
	assert.Error(t, err)
}
