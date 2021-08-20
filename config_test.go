package config_test

import (
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFeed(t *testing.T) {
	c := &struct{}{}
	err := config.Feed(c, feeder.Env{})
	assert.NoError(t, err)
}

func TestFeed_With_Invalid_File_It_Should_Fail(t *testing.T) {
	c := &struct{}{}
	err := config.Feed(c, feeder.Json{})
	assert.Error(t, err)
}
