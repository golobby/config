package config_test

import (
	"testing"

	"github.com/golobby/config"
	"github.com/golobby/config/feeder"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Set(t *testing.T) {
	c, err := config.New()
	assert.NoError(t, err)

	c.Set("k", "v")
	v, err := c.Get("k")

	assert.NoError(t, err)
	assert.Equal(t, "v", v)
}

func TestConfig_GetAll(t *testing.T) {
	m := feeder.Map{
		"singer": "Pink Floyd",
		"albums": []struct {
			Name string
			Year int
		}{
			{Name: "Division Bell", Year: 1994},
			{Name: "The Wall", Year: 1979},
		},
	}

	c, err := config.New(m)
	assert.NoError(t, err)

	v := c.GetAll()
	assert.Equal(t, map[string]interface{}(m), v)
}

func TestConfig_Get(t *testing.T) {
	c, err := config.New(feeder.Map{
		"string": "String",
		"int":    13,
		"float":  3.14,
		"true":   true,
		"false":  false,
		"slice": map[string]interface{}{
			"item": "value",
		},
	})
	assert.NoError(t, err)

	v, err := c.Get("string")
	assert.NoError(t, err)
	assert.Equal(t, "String", v)

	v, err = c.Get("int")
	assert.NoError(t, err)
	assert.Equal(t, 13, v)

	v, err = c.Get("float")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v)

	v, err = c.Get("true")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.Get("false")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	v, err = c.Get("slice.item")
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	_, err = c.Get("slice.wrong")
	assert.Error(t, err)
	assert.Equal(t, "value not found for the key `wrong`", err.Error())

	_, err = c.Get("wrong")
	assert.Error(t, err)
	assert.Equal(t, "value not found for the key `wrong`", err.Error())

	_, err = c.Get("wrong.wrong")
	assert.Error(t, err)
	assert.Equal(t, "value not found for the key `wrong`", err.Error())
}

func TestConfig_GetBool(t *testing.T) {
	c, err := config.New(feeder.Map{
		"true":        true,
		"false":       false,
		"trueString":  "true",
		"falseString": "false",
		"trueInt":     1,
		"falseInt":    0,
		"string":      "String",
		"number":      13,
	})
	assert.NoError(t, err)

	v, err := c.GetBool("true")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetBool("false")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	v, err = c.GetBool("trueString")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetBool("falseString")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	v, err = c.GetBool("trueInt")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetBool("falseInt")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	_, err = c.GetBool("string")
	assert.Error(t, err)
	assert.Equal(t, "value `String` (`string`) is not `bool`", err.Error())

	_, err = c.GetBool("number")
	assert.Error(t, err)
	assert.Equal(t, "value `13` (`int`) is not `bool`", err.Error())
}

func TestConfig_GetFloat(t *testing.T) {
	c, err := config.New(feeder.Map{
		"float":       3.14,
		"int":         13,
		"floatString": "3.14",
		"intString":   "13",
		"string":      "String",
	})
	assert.NoError(t, err)

	v, err := c.GetFloat("float")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v)

	v, err = c.GetFloat("int")
	assert.NoError(t, err)
	assert.Equal(t, float64(13), v)

	v, err = c.GetFloat("floatString")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, v)

	v, err = c.GetFloat("intString")
	assert.NoError(t, err)
	assert.Equal(t, float64(13), v)

	_, err = c.GetFloat("string")
	assert.Error(t, err)
	assert.Equal(t, "value `String` (`string`) is not `float64`", err.Error())
}

func TestConfig_GetInt(t *testing.T) {
	c, err := config.New(feeder.Map{
		"int":         13,
		"float":       3.14,
		"intString":   "13",
		"floatString": "3.14",
		"string":      "String",
	})
	assert.NoError(t, err)

	v, err := c.GetInt("int")
	assert.NoError(t, err)
	assert.Equal(t, 13, v)

	v, err = c.GetInt("float")
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	v, err = c.GetInt("intString")
	assert.NoError(t, err)
	assert.Equal(t, 13, v)

	v, err = c.GetInt("floatString")
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	_, err = c.GetInt("string")
	assert.Error(t, err)
	assert.Equal(t, "value `String` (`string`) is not `int`", err.Error())
}

func TestConfig_GetString(t *testing.T) {
	c, err := config.New(feeder.Map{
		"int":    13,
		"float":  3.14,
		"false":  false,
		"true":   true,
		"string": "String",
	})
	assert.NoError(t, err)

	v, err := c.GetString("int")
	assert.Error(t, err)
	assert.Equal(t, "value `13` (`int`) is not `string`", err.Error())

	v, err = c.GetString("float")
	assert.Error(t, err)
	assert.Equal(t, "value `3.14` (`float64`) is not `string`", err.Error())

	v, err = c.GetString("false")
	assert.Error(t, err)
	assert.Equal(t, "value `false` (`bool`) is not `string`", err.Error())

	v, err = c.GetString("true")
	assert.Error(t, err)
	assert.Equal(t, "value `true` (`bool`) is not `string`", err.Error())

	v, err = c.GetString("string")
	assert.NoError(t, err)
	assert.Equal(t, "String", v)
}

func TestConfig_GetStrictBool(t *testing.T) {
	c, err := config.New(feeder.Map{
		"true":        true,
		"false":       false,
		"trueString":  "true",
		"falseString": "false",
		"trueInt":     1,
		"falseInt":    0,
		"string":      "String",
		"number":      13,
	})
	assert.NoError(t, err)

	v, err := c.GetStrictBool("true")
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	v, err = c.GetStrictBool("false")
	assert.NoError(t, err)
	assert.Equal(t, false, v)

	v, err = c.GetStrictBool("trueString")
	assert.Error(t, err)
	assert.Equal(t, "value `true` (`string`) is not `bool`", err.Error())

	v, err = c.GetStrictBool("falseString")
	assert.Error(t, err)
	assert.Equal(t, "value `false` (`string`) is not `bool`", err.Error())

	v, err = c.GetStrictBool("trueInt")
	assert.Error(t, err)
	assert.Equal(t, "value `1` (`int`) is not `bool`", err.Error())

	v, err = c.GetStrictBool("falseInt")
	assert.Error(t, err)
	assert.Equal(t, "value `0` (`int`) is not `bool`", err.Error())

	_, err = c.GetStrictBool("string")
	assert.Error(t, err)
	assert.Equal(t, "value `String` (`string`) is not `bool`", err.Error())

	_, err = c.GetStrictBool("number")
	assert.Error(t, err)
	assert.Equal(t, "value `13` (`int`) is not `bool`", err.Error())
}

func TestConfig_Feed_Multiple(t *testing.T) {
	c, err := config.New(
		feeder.Map{
			"url": "going to be overridden by the next feeders",
		},
		feeder.Map{
			"url": "going to be overridden by the next feeder",
		},
		feeder.Map{
			"url": "https://github.com/golobby/config",
		},
	)
	assert.NoError(t, err)

	v, err := c.Get("url")
	assert.Equal(t, "https://github.com/golobby/config", v)
	assert.NoError(t, err)
}
