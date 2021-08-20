package config_test

import (
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFeed(t *testing.T) {
	c := &struct{}{}
	err := config.New(feeder.Env{}).Feed(c)
	assert.NoError(t, err)
}

func TestFeed_With_Invalid_File_It_Should_Fail(t *testing.T) {
	c := &struct{}{}
	err := config.New(feeder.Json{}).Feed(c)
	assert.Error(t, err)
}

func TestFeed_WithMultiple_Feeders(t *testing.T) {
	_ = os.Setenv("PRODUCTION", "1")
	_ = os.Setenv("APP_PORT", "6969")

	c := &struct {
		App struct {
			Name string `dotenv:"APP_NAME" env:"APP_NAME"`
			Port int    `dotenv:"APP_PORT" env:"APP_PORT"`
		}
		Debug      bool    `dotenv:"DEBUG" env:"DEBUG"`
		Production bool    `dotenv:"PRODUCTION" env:"PRODUCTION"`
		Pi         float64 `dotenv:"PI" env:"PI"`
	}{}

	f1 := feeder.Json{Path: "assets/sample1.json"}
	f2 := feeder.DotEnv{Path: "assets/.env.sample2"}
	f3 := feeder.Env{}

	err := config.New(f1, f2, f3).Feed(c)
	assert.NoError(t, err)

	assert.Equal(t, "Blog", c.App.Name)
	assert.Equal(t, 6969, c.App.Port)
	assert.Equal(t, false, c.Debug)
	assert.Equal(t, true, c.Production)
	assert.Equal(t, 3.14, c.Pi)
}

func TestConfig_Refresh(t *testing.T) {
	_ = os.Setenv("NAME", "One")

	s := &struct{
		Name string `env:"NAME"`
	}{}

	c := config.New(feeder.Env{})
	err := c.Feed(s)
	assert.NoError(t, err)

	assert.Equal(t, "One", s.Name)

	_ = os.Setenv("NAME", "Two")

	err = c.Refresh()
	assert.NoError(t, err)

	assert.Equal(t, "Two", s.Name)
}
