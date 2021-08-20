package feeder_test

import (
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestEnv_Feed(t *testing.T) {
	_ = os.Setenv("APP_NAME", "Shop")
	_ = os.Setenv("APP_PORT", "8585")
	_ = os.Setenv("DEBUG", "true")
	_ = os.Setenv("PRODUCTION", "false")
	_ = os.Setenv("PI", "3.14")

	type config struct {
		App struct {
			Name string `env:"APP_NAME"`
			Port int    `env:"APP_PORT"`
		}
		Debug      bool    `env:"DEBUG"`
		Production bool    `env:"PRODUCTION"`
		Pi         float64 `env:"PI"`
	}

	c := config{}
	f := feeder.Env{}

	err := f.Feed(&c)
	assert.NoError(t, err)

	assert.Equal(t, "Shop", c.App.Name)
	assert.Equal(t, 8585, c.App.Port)
	assert.Equal(t, true, c.Debug)
	assert.Equal(t, false, c.Production)
	assert.Equal(t, 3.14, c.Pi)
}

func TestEnv_Feed_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	_ = os.Setenv("APP_NAME", "Shop")

	c := struct {
		App struct {
			Name float64 `env:"APP_NAME"`
		}
	}{}
	f := feeder.Env{}

	err := f.Feed(&c)
	assert.Error(t, err)
}
