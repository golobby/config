// Package config is a lightweight yet powerful config package for Go projects.
// It takes advantage of env files and OS variables alongside config files to be your ultimate requirement.
package config

import (
	"errors"
	"github.com/golobby/config/env"
	"os"
	"strconv"
	"strings"
)

// Feeder is an interface for feeders which can feed config instances (provider their contents).
type Feeder interface {
	Feed() (map[string]interface{}, error)
}

// Options is a struct that contains all the required data for instantiating a new config instance.
type Options struct {
	Feeder  Feeder // Feeder that is going to feed the config instance
	EnvFile string // EnvFile is the .env file that is going to be used in config file values
}

// config is the main struct that keeps all the config instance data.
type config struct {
	options  []Options              // options is the construction options
	envFiles map[string]string      // envFiles keeps all the given .env file paths
	items    map[string]interface{} // items keeps the config data
}

// addEnv will add given env items to the instance env items.
func (c config) addEnv(items map[string]string) {
	for k, v := range items {
		c.envFiles[k] = v
	}
}

// Feed will feed the config instance using the given feeder.
// It accepts all kinds of feeders that implement the Feeder interface.
// The built-in feeders are in the feeder subpackage.
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

// Env will return environment variable value for the given environment variable key.
func (c config) Env(key string) string {
	if v, ok := c.envFiles[key]; ok && v != "" {
		return v
	}

	return os.Getenv(key)
}

// Set will store the given key/value into the config instance.
// It keeps the key/values that have added on runtime in the memory.
// It won't change the config files.
func (c config) Set(key string, value interface{}) {
	c.items[key] = value
}

// Get will return the value of the given key.
// The return type is interface, so it probably should be cast to the related data type.
// It will return an error if there is no value for the given key.
func (c config) Get(key string) (interface{}, error) {
	if v, ok := c.items[key]; ok {
		return v, nil
	}

	if strings.Contains(key, ".") == false {
		return nil, errors.New("value not found for the key " + key)
	}

	return lookup(c.items, key)
}

// Get will return the value of the given key.
// It casts the value type to string.
// It will return an error if the related value is not string.
// It will return an error if there is no value for the given key.
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

// GetInt will return the value of the given key.
// It casts the value type to int.
// It will return an error if the related value is not bool.
// It will return an error if there is no value for the given key.
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

// GetFloat will return the value of the given key.
// It casts the value type to float64.
// It will return an error if the related value is not float.
// It will return an error if there is no value for the given key.
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

// GetBool will return the value of the given key.
// It casts the value type to bool.
// It considers the "true" and "false" string values like true and false boolean values respectively.
// It will return an error if the related value is not bool.
// It will return an error if there is no value for the given key.
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

// GetStrictBool will return the value of the given key.
// It casts the value type to bool.
// It checks the value type strictly so "true" and "false" string values won't considered boolean.
// It will return an error if the related value is not bool.
// It will return an error if there is no value for the given key.
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

// parse will replace the placeholders with env and OS values.
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

// lookup will search for the given key recursively.
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

// find will return the value of given key in the given collection
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

// dig will return sub-collection which the given partition of key points to.
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

// New will return a brand new instance of Config with given options.
func New(options ...Options) (*config, error) {
	c := &config{
		items:    map[string]interface{}{},
		envFiles: map[string]string{},
	}

	for _, o := range options {
		c.options = append(c.options, o)

		if o.EnvFile != "" {
			if items, err := env.Load(o.EnvFile); err != nil {
				return nil, err
			} else {
				c.addEnv(items)
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
