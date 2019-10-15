package config

import (
	"errors"
	"github.com/golobby/config/env"
	"os"
	"strconv"
	"strings"
)

type Feeder interface {
	Feed() (map[string]interface{}, error)
}

type Options struct {
	Feeder  Feeder
	EnvFile string
}

type config struct {
	options []Options
	envs    map[string]string
	items   map[string]interface{}
}

func (c config) feedEnv(items map[string]string) {
	for k, v := range items {
		c.envs[k] = v
	}
}

func (c config) Feed(r Feeder) error {
	if items, err := r.Feed(); err != nil {
		return err
	} else {
		for k, v := range items {
			c.items[k] = c.parse(v)
		}

		return nil
	}
}

func (c config) Env(key string) string {
	if v, ok := c.envs[key]; ok {
		return v
	}

	return os.Getenv(key)
}

func (c config) Set(key string, value interface{}) {
	c.items[key] = value
}

func (c config) Get(key string) (interface{}, error) {
	if v, ok := c.items[key]; ok {
		return v, nil
	}

	if strings.Contains(key, ".") == false {
		return nil, errors.New("value not found for the key " + key)
	}

	return lookup(c.items, key)
}

func (c config) GetString(key string) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if v, ok := v.(string); ok {
		return v, nil
	}

	return "", errors.New("value for " + key + " is not string")
}

func (c config) GetInt(key string) (int, error) {
	v, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	if v, ok := v.(int); ok {
		return v, nil
	}

	return 0, errors.New("value for " + key + " is not int")
}

func (c config) GetFloat(key string) (float64, error) {
	v, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	if v, ok := v.(float64); ok {
		return v, nil
	}

	return 0, errors.New("value for " + key + " is not float")
}

func (c config) GetBool(key string) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if v, ok := v.(bool); ok {
		return v, nil
	}

	if v, ok := v.(string); ok {
		if v == "true" {
			return true, nil
		} else if v == "false" {
			return false, nil
		}
	}

	return false, errors.New("value for " + key + " is not bool")
}

func (c config) GetStrictBool(key string) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if v, ok := v.(bool); ok {
		return v, nil
	}

	return false, errors.New("value for " + key + " is not bool")
}

func (c config) parse(value interface{}) interface{} {
	if stmt, ok := value.(string); ok {
		if stmt[0:2] == "${" && stmt[len(stmt)-1:] == "}" {
			pipe := strings.Index(stmt, "|")

			if pipe == -1 {
				key := strings.TrimSpace(stmt[2 : len(stmt)-1])
				return c.Env(key)
			}

			key := strings.TrimSpace(stmt[2:pipe])
			if v := c.Env(key); v != "" {
				return v
			}

			return strings.TrimSpace(stmt[pipe+1 : len(stmt)-1])
		}
	}

	return value
}

func lookup(collection interface{}, key string) (interface{}, error) {
	keys := strings.Split(key, ".")

	if len(keys) == 1 {
		return find(collection, keys[0])
	} else {
		c, err := dig(collection, keys[0])
		if err != nil {
			return nil, err
		}

		return lookup(c, strings.Join(keys[1:], "."))
	}
}

func find(collection interface{}, key string) (interface{}, error) {
	switch collection.(type) {
	case map[string]interface{}:
		if v, ok := collection.(map[string]interface{})[key]; ok {
			return v, nil
		}
	case []interface{}:
		k, err := strconv.Atoi(key)
		if err == nil && len(collection.([]interface{})) > k {
			return collection.([]interface{})[k], nil
		}
	}

	return nil, errors.New("value not found for the key " + key)
}

func dig(collection interface{}, key string) (interface{}, error) {
	if v, ok := collection.(map[string]interface{}); ok {
		if v, ok := v[key]; ok {
			return v, nil
		}
	} else if v, ok := collection.([]interface{}); ok {
		i, err := strconv.Atoi(key)
		if err == nil && len(v) > i {
			return v[i], nil
		}
	}

	return nil, errors.New("value not found for the key " + key)
}

func New(options ...Options) (*config, error) {
	c := &config{
		items: map[string]interface{}{},
	}

	for _, o := range options {
		c.options = append(c.options, o)

		if o.EnvFile != "" {
			if items, err := env.Load(o.EnvFile); err != nil {
				return nil, err
			} else {
				c.feedEnv(items)
			}
		}

		if o.Feeder != nil {
			if err := c.Feed(o.Feeder); err != nil {
				return nil, err
			}
		}
	}

	return c, nil
}
