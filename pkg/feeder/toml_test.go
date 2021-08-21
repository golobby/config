package feeder_test

import (
	"testing"

	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
)

func TestToml_Feed(t *testing.T) {
	type config struct {
		App struct {
			Name string
			Port int
		}
		Debug      bool
		Production bool
		Pi         float64
	}

	c := config{}
	f := feeder.Toml{Path: "./../../assets/sample1.toml"}

	err := f.Feed(&c)
	assert.NoError(t, err)

	assert.Equal(t, "Shop", c.App.Name)
	assert.Equal(t, 8585, c.App.Port)
	assert.Equal(t, true, c.Debug)
	assert.Equal(t, false, c.Production)
	assert.Equal(t, 3.14, c.Pi)
}

func TestToml_Feed_With_Invalid_File_It_Should_Fail(t *testing.T) {
	c := struct{}{}
	f := feeder.Toml{Path: "nowhere!"}

	err := f.Feed(&c)
	assert.Error(t, err)
}

func TestToml_Feed_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	c := struct {
		App struct {
			Name float64
		}
	}{}
	f := feeder.Toml{Path: "./../../assets/sample1.toml"}

	err := f.Feed(&c)
	assert.Error(t, err)
}
