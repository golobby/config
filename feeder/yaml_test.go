package feeder_test

import (
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Yaml_Feed_Not_Existing_File_It_Should_Return_Error(t *testing.T) {
	y := feeder.Yaml{Path: "/404.yaml"}

	_, err := y.Feed()

	assert.Error(t, err)
}

func Test_Yaml_Feed_Invalid_JSON_It_Should_Return_Error(t *testing.T) {
	y := feeder.Yaml{Path: "test/invalid.yaml"}

	_, err := y.Feed()

	assert.Error(t, err)
}

func Test_Yaml_Feed_Sample1(t *testing.T) {
	y := feeder.Yaml{Path: "test/config.yaml"}

	m, err := y.Feed()

	assert.NoError(t, err)
	assert.Equal(t, "MyAppUsingGoLobbyConfig", m["name"])
	assert.Equal(t, 3.14, m["version"])
}
