package feeder_test

import (
	"github.com/golobby/config/v3/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYaml_Feed(t *testing.T) {
	y := feeder.Yaml{Path: "test/config.yaml"}

	m, err := y.Feed()
	assert.NoError(t, err)

	assert.Equal(t, "MyAppUsingGoLobbyConfig", m["name"])
	assert.Equal(t, 3.14, m["version"])
}

func TestYaml_Feed_With_Invalid_YAML_Path_It_Should_Fail(t *testing.T) {
	y := feeder.Yaml{Path: "/404.yaml"}

	_, err := y.Feed()
	assert.Error(t, err)
}

func TestYaml_Feed_With_Invalid_YAML_File_It_Should_Fail(t *testing.T) {
	y := feeder.Yaml{Path: "test/invalid.yaml"}

	_, err := y.Feed()
	assert.Error(t, err)
}
