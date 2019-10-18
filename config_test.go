package config_test

import (
	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Config_Set_Get_With_A_Simple_Key_String_Value(t *testing.T) {
	c, err := config.New(config.Options{})
	assert.NoError(t, err)

	c.Set("k", "v")
	v, err := c.Get("k")

	assert.NoError(t, err)
	assert.Equal(t, "v", v)
}

func Test_Config_Feed_With_Map_Repo(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"name":     "Hey You",
			"band":     "Pink Floyd",
			"year":     1979,
			"duration": 4.6,
		},
	})
	assert.NoError(t, err)

	v, err := c.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "Hey You", v)

	v, err = c.GetString("name")
	assert.NoError(t, err)
	assert.Equal(t, "Hey You", v)

	_, err = c.GetInt("name")
	assert.Error(t, err)

	band, err := c.Get("band")
	assert.NoError(t, err)
	assert.Equal(t, "Pink Floyd", band)

	_, err = c.GetFloat("band")
	assert.Error(t, err)

	year, err := c.Get("year")
	assert.NoError(t, err)
	assert.Equal(t, 1979, year)

	_, err = c.GetString("year")
	assert.Error(t, err)

	year, err = c.GetInt("year")
	assert.NoError(t, err)
	assert.Equal(t, 1979, year)

	duration, err := c.Get("duration")
	assert.NoError(t, err)
	assert.Equal(t, 4.6, duration)

	_, err = c.GetBool("duration")
	assert.Error(t, err)

	duration, err = c.GetFloat("duration")
	assert.NoError(t, err)
	assert.Equal(t, 4.6, duration)

	_, err = c.GetStrictBool("duration")
	assert.Error(t, err)

	_, err = c.Get("wrong")
	assert.Error(t, err)

	_, err = c.GetString("wrong")
	assert.Error(t, err)

	_, err = c.GetInt("wrong")
	assert.Error(t, err)

	_, err = c.GetFloat("wrong")
	assert.Error(t, err)

	_, err = c.GetBool("wrong")
	assert.Error(t, err)

	_, err = c.GetStrictBool("wrong")
	assert.Error(t, err)

	_, err = c.Get("wrong.nested")
	assert.Error(t, err)
}

func Test_Config_GetBool(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"a": true,
			"b": "true",
			"c": false,
			"d": "false",
			"e": "error",
		},
	})
	assert.NoError(t, err)

	v, err := c.GetBool("a")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetBool("b")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetBool("c")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	v, err = c.GetBool("d")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	_, err = c.GetBool("e")
	assert.Error(t, err)
}

func Test_Config_GetStrictBool(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"a": true,
			"b": "true",
			"c": false,
			"d": "false",
			"e": "error",
		},
	})
	assert.NoError(t, err)

	v, err := c.GetStrictBool("a")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	_, err = c.GetStrictBool("b")
	assert.Error(t, err)

	v, err = c.GetStrictBool("c")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	_, err = c.GetStrictBool("d")
	assert.Error(t, err)

	_, err = c.GetStrictBool("e")
	assert.Error(t, err)
}

func Test_Config_Feed_With_Map_Repo_Includes_A_Slice(t *testing.T) {
	c, err := config.New(config.Options{Feeder: feeder.Map{
		"scores": map[string]interface{}{
			"A": 1,
			"B": 2,
			"C": 3,
		},
	}})
	assert.NoError(t, err)

	v, err := c.Get("scores.A")
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = c.Get("scores.B")
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	_, err = c.Get("scores.Wrong")
	assert.Error(t, err)
}

func Test_Config_Feed_It_Should_Get_Env_From_OS(t *testing.T) {
	err := os.Setenv("URL", "https://miladrahimi.com")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{Feeder: feeder.Map{
		"url": "${ URL }",
	}})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)

	assert.Equal(t, os.Getenv("URL"), v)
}

func Test_Config_Feed_It_Should_Get_Env_From_OS_With_Default_Value(t *testing.T) {
	err := os.Setenv("URL", "https://miladrahimi.com")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{Feeder: feeder.Map{
		"url": "${ URL | DEFAULT_BUT_NOT_USED }",
	}})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)

	assert.Equal(t, os.Getenv("URL"), v)
}

func Test_Config_Feed_It_Should_Get_Env_Default_When_Not_In_OS(t *testing.T) {
	err := os.Setenv("EMPTY", "")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"url": "${ EMPTY | http://localhost }",
		},
	})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost", v)
}

func Test_Config_Feed_JSON(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Json{Path: "feeder/test/Config.json"},
	})
	assert.NoError(t, err)

	v, err := c.Get("numbers.2")
	assert.NoError(t, err)
	assert.Equal(t, float64(3), v)

	v, err = c.Get("users.0.address.city")
	assert.NoError(t, err)
	assert.Equal(t, "Delfan", v)
}

func Test_Config_Feed_Invalid_JSON(t *testing.T) {
	_, err := config.New(config.Options{
		Feeder: feeder.Json{Path: "feeder/test/invalid-json"},
	})
	assert.Error(t, err)
}

func Test_Config_Env_With_Sample_Env_File(t *testing.T) {
	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"url": "${ APP_URL }",
		},
		EnvFile: "env/test/.env",
	})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", v)
}

func Test_Config_Env_With_Empty_Env_It_Should_Use_OS_Vars(t *testing.T) {
	err := os.Setenv("APP_NAME", "MyApp")
	if err != nil {
		panic(err)
	}

	c, err := config.New(config.Options{
		Feeder: feeder.Map{
			"name": "${ APP_NAME }",
		},
		EnvFile: "env/test/.env",
	})
	assert.NoError(t, err)

	v, err := c.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "MyApp", v)
}

func Test_Config_Env_With_Invalid_Env_It_Should_Raise_An_Error(t *testing.T) {
	_, err := config.New(config.Options{
		Feeder:  feeder.Map{},
		EnvFile: "env/test/.invalid.env",
	})
	assert.Error(t, err)
}
