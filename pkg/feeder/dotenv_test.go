package feeder_test

import (
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDotEnv_Feed(t *testing.T) {
	type config struct {
		App struct {
			Name string `dotenv:"APP_NAME"`
			Port int    `dotenv:"APP_PORT"`
		}
		Debug      bool    `dotenv:"DEBUG"`
		Production bool    `dotenv:"PRODUCTION"`
		Pi         float64 `dotenv:"PI"`
	}

	c := config{}
	f := feeder.DotEnv{Path: "./../../assets/.env.sample1"}

	err := f.Feed(&c)
	assert.NoError(t, err)

	assert.Equal(t, "Shop", c.App.Name)
	assert.Equal(t, 8585, c.App.Port)
	assert.Equal(t, true, c.Debug)
	assert.Equal(t, false, c.Production)
	assert.Equal(t, 3.14, c.Pi)
}

func TestDotEnv_Feed_With_Invalid_File_It_Should_Fail(t *testing.T) {
	c := struct{}{}
	f := feeder.DotEnv{Path: "nowhere!"}

	err := f.Feed(&c)
	assert.Error(t, err)
}

func TestDotEnv_Feed_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	c := struct {
		App struct {
			Name float64 `dotenv:"APP_NAME"`
		}
	}{}
	f := feeder.DotEnv{Path: "./../../assets/.env.sample1"}

	err := f.Feed(&c)
	assert.Error(t, err)
}
