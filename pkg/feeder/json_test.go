package feeder_test

import (
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson_Feed(t *testing.T) {
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
	f := feeder.Json{Path: "./../../assets/sample1.json"}

	err := f.Feed(&c)
	assert.NoError(t, err)

	assert.Equal(t, "Shop", c.App.Name)
	assert.Equal(t, 8585, c.App.Port)
	assert.Equal(t, true, c.Debug)
	assert.Equal(t, false, c.Production)
	assert.Equal(t, 3.14, c.Pi)
}

func TestJson_Feed_With_Invalid_File_It_Should_Fail(t *testing.T) {
	c := struct{}{}
	f := feeder.Json{Path: "nowhere!"}

	err := f.Feed(&c)
	assert.Error(t, err)
}

func TestJson_Feed_With_Invalid_Struct_It_Should_Fail(t *testing.T) {
	c := struct {
		App struct{
			Name float64
		}
	}{}
	f := feeder.Json{Path: "./../../assets/sample1.json"}

	err := f.Feed(&c)
	assert.Error(t, err)
}
