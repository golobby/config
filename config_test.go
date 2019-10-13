package config_test

import (
	"config"
	"config/repository"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Config_With_Simple_String(t *testing.T) {
	c := config.Config{}
	c["k"] = "v"
	assert.Equal(t, "v", c["k"])
}

func Test_Config_Feed_With_Map_Repo(t *testing.T) {
	m := repository.Map{
		"name":     "Hey You",
		"band":     "Pink Floyd",
		"year":     1979,
		"duration": 4.6,
	}

	c, err := config.New(m)
	assert.NoError(t, err)

	assert.Equal(t, "Hey You", c["name"])

	band, err := c.Get("band")
	assert.NoError(t, err)
	assert.Equal(t, "Pink Floyd", band)

	year, err := c.Get("year")
	assert.NoError(t, err)
	assert.Equal(t, 1979, year)

	duration, err := c.Get("duration")
	assert.NoError(t, err)
	assert.Equal(t, 4.6, duration)

	wrong, err := c.Get("wrong.nested")
	assert.Error(t, err)
	assert.Equal(t, nil, wrong)
}

func Test_Config_Feed_With_Map_Repo_Includes_Slice(t *testing.T) {
	m := repository.Map{
		"scores": map[string]int{
			"A": 1,
			"B": 2,
			"C": 3,
		},
	}

	c, err := config.New(m)

	assert.NoError(t, err)
	assert.Equal(t, 1, c["scores"].(map[string]int)["A"])
	assert.Equal(t, 2, c["scores"].(map[string]int)["B"])
	assert.Equal(t, 3, c["scores"].(map[string]int)["C"])
}

func Test_Config_Feed_With_Map_Repo_It_Should_Parse_Nested_Maps(t *testing.T) {
	m := repository.Map{
		"scores": map[string]interface{}{
			"a": 1,
			"b": 2,
		},
	}

	c, err := config.New(m)
	assert.NoError(t, err)

	a, err := c.Get("scores.a")
	assert.NoError(t, err)
	assert.Equal(t, 1, a)

	b, err := c.Get("scores.b")
	assert.NoError(t, err)
	assert.Equal(t, 2, b)

	assert.Equal(t, 1, c["scores"].(map[string]interface{})["a"])
	assert.Equal(t, 2, c["scores"].(map[string]interface{})["b"])
}

func Test_Config_Feed_It_Should_Get_Env_From_OS(t *testing.T) {
	_ = os.Setenv("URL", "https://miladrahimi.com")

	m := repository.Map{
		"url": "${ URL }",
	}

	c, err := config.New(m)

	assert.NoError(t, err)
	assert.Equal(t, os.Getenv("URL"), c["url"])
}

func Test_Config_Feed_It_Should_Get_Env_Default_When_Not_In_OS(t *testing.T) {
	_ = os.Setenv("EMPTY", "")

	m := repository.Map{
		"url": "${ EMPTY | http://localhost }",
	}

	c, err := config.New(m)

	assert.NoError(t, err)
	assert.Equal(t, "http://localhost", c["url"])
}

func Test_Config_Feed_JSON(t *testing.T) {
	j := repository.Json{Path: "test/config.json"}

	c, err := config.New(j)
	assert.NoError(t, err)

	v, err := c.Get("numbers.2")
	assert.NoError(t, err)
	assert.Equal(t, float64(3), v)

	v, err = c.Get("users.0.address.city")
	assert.NoError(t, err)
	assert.Equal(t, "Delfan", v)
}
