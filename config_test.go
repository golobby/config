package config_test

import (
	json2 "encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
)

func Test_Config_Set_Get_With_A_Simple_Key_String_Value(t *testing.T) {
	c, err := config.New()
	assert.NoError(t, err)

	c.Set("k", "v")
	v, err := c.Get("k")

	assert.NoError(t, err)
	assert.Equal(t, "v", v)
}

func Test_Config_Feed_With_Map_Repo(t *testing.T) {
	m := feeder.Map{
		"name":     "Hey You",
		"band":     "Pink Floyd",
		"year":     1979,
		"duration": 4.6,
	}
	c, err := config.New(m)
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

	assert.Equal(t, map[string]interface{}(m), c.GetAll())
}

func Test_Config_GetBool(t *testing.T) {
	c, err := config.New(
		feeder.Map{
			"a": true,
			"b": "true",
			"c": false,
			"d": "false",
			"e": "error",
		},
	)
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
	c, err := config.New(
		feeder.Map{
			"a": true,
			"b": "true",
			"c": false,
			"d": "false",
			"e": "error",
		},
	)
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
	c, err := config.New(feeder.Map{
		"scores": map[string]interface{}{
			"A": 1,
			"B": 2,
			"C": 3,
		},
	})
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

	c, err := config.New(feeder.Map{
		"url": "going to be overrided by next feeder",
	}, &feeder.Env{Path: "feeder/test/.env"})
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

	c, err := config.New(feeder.Map{
		"url": "going to be overrided by next feeder",
	}, &feeder.Env{Path: "feeder/test/.env"})
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

	c, err := config.New(feeder.Map{
		"empty": "http://localhost",
	}, &feeder.Env{Keys: []string{"EMPTY"}})
	assert.NoError(t, err)

	v, err := c.Get("empty")
	assert.NoError(t, err)

	assert.Equal(t, "http://localhost", v)
}

func Test_Config_Feed_JSON(t *testing.T) {
	c, err := config.New(feeder.Json{Path: "feeder/test/config.json"})
	assert.NoError(t, err)

	v, err := c.Get("numbers.2")
	assert.NoError(t, err)
	assert.Equal(t, float64(3), v)

	v, err = c.Get("users.0.address.city")
	assert.NoError(t, err)
	assert.Equal(t, "Delfan", v)
}

func Test_Config_Feed_JSON_Directory(t *testing.T) {
	err := os.Setenv("APP_OS", "Linux")
	if err != nil {
		panic(err)
	}

	c, err := config.New(feeder.JsonDirectory{Path: "feeder/test/json"},
		&feeder.Env{Keys: []string{"APP_OS"}})
	assert.NoError(t, err)

	v, err := c.Get("app.name")
	assert.NoError(t, err)
	assert.Equal(t, "MyAppUsingConfig", v)

	v, err = c.Get("app.version")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v)

	v, err = c.Get("app.os")
	assert.NoError(t, err)
	assert.Equal(t, "Linux", v)
}

func Test_Config_Feed_Invalid_JSON(t *testing.T) {
	_, err := config.New(feeder.Json{Path: "feeder/test/invalid-json"})
	assert.Error(t, err)
}

func Test_Config_Env_With_Sample_Env_File(t *testing.T) {
	err := os.Setenv("URL", "")
	if err != nil {
		panic(err)
	}
	c, err := config.New(feeder.Map{
		"url": "",
	}, &feeder.Env{Path: "feeder/test/.env"})
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", v)
}

func Test_Config_Env_With_Empty_Env_It_Should_Use_OS_Vars(t *testing.T) {
	err := os.Setenv("NAME", "MyApp")
	if err != nil {
		panic(err)
	}

	c, err := config.New(feeder.Map{
		"name": "",
	}, &feeder.Env{Path: "feeder/test/.env"},
	)
	assert.NoError(t, err)

	v, err := c.Get("name")
	assert.NoError(t, err)
	assert.Equal(t, "MyApp", v)
}

func Test_Config_Env_With_Invalid_Env_It_Should_Raise_An_Error(t *testing.T) {
	_, err := config.New(feeder.Map{},
		&feeder.Env{Path: "env/test/.invalid.env"},
	)
	assert.Error(t, err)
}

func Test_Config_Reload_It_Should_Reload_The_Feeders(t *testing.T) {
	path := "feeder/test/runtime.json"

	json, err := json2.Marshal(map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, json, 0755)
	if err != nil {
		panic(err)
	}

	c, err := config.New(feeder.Json{Path: path})
	assert.NoError(t, err)

	v, err := c.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	json, err = json2.Marshal(map[string]interface{}{
		"key": "new-value",
	})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, json, 0755)
	if err != nil {
		panic(err)
	}

	err = c.Reload()
	assert.NoError(t, err)

	v, err = c.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, "new-value", v)
}

// https://github.com/golobby/config/issues/8
func Test_GetInt_From_JSON(t *testing.T) {
	c, err := config.New(feeder.Json{Path: "feeder/test/issue8.json"})
	assert.NoError(t, err)

	keys := []string{
		"int",
		"strInt",
	}

	for _, key := range keys {
		v, err := c.GetInt(key)
		if err != nil {
			t.Errorf(
				"\nkey: %v \nv: %v \nerr: %v",
				key, v, err.Error(),
			)
		}
	}
}
