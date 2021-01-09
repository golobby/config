// Package Config is a lightweight yet powerful configuration management tool for Go projects.
// It takes advantage of env files and OS variables alongside Config files to be your only requirement.
package config

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

//NotFoundError happens when you try to access a key which is not defined in the configuration files.
type NotFoundError struct {
	key string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("value not found for the key %s", n.key)
}

//TypeError happens when you try to access a key using a helper function that casts value to a type which can't be done.
type TypeError struct {
	value  interface{}
	wanted string
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("value %s (%T) is not %s", t.value, t.value, t.wanted)
}

// Config keeps all the Config instance data.
type Config struct {
	ConfigBase
}

// StartListener makes the Config instance to listen to the SIGHUP signal and reload the feeders and environment files.
func (c *Config) StartListener() {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			_ = c.ReloadEnv()
			_ = c.Reload()
		}
	}()
}

// Get returns the value of the given key.
// The return type is "interface{}".
// It probably needs to be cast to the related data type.
// It returns an error if there is no value for the given key.
func (c *Config) Get(key string) (interface{}, error) {
	v, exists := c.ConfigBase.Get(key)
	if !exists {
		return nil, &NotFoundError{key: key}
	}

	return v, nil
}

// GetString returns the value of the given key.
// It also casts the value type to string internally.
// It returns an error if the related value is not a string.
// It returns an error if there is no value for the given key.
func (c *Config) GetString(key string) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if v, ok := v.(string); ok {
		return v, nil
	}

	return "", &TypeError{value: v, wanted: "string"}
}

// GetInt returns the value of the given key.
// It also casts the value type to int internally.
// It returns an error if the related value is not an int.
// It returns an error if there is no value for the given key.
func (c *Config) GetInt(key string) (int, error) {
	v, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	switch val := v.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		return strconv.Atoi(val)
	}

	return 0, &TypeError{value: v, wanted: "int"}
}

// GetFloat returns the value of the given key.
// It also casts the value type to float64 internally.
// It returns an error if the related value is not a float64.
// It returns an error if there is no value for the given key.
func (c *Config) GetFloat(key string) (float64, error) {
	v, err := c.Get(key)
	if err != nil {
		return 0, err
	}

	if v, ok := v.(float64); ok {
		return v, nil
	}

	return 0, &TypeError{value: v, wanted: "float"}
}

// GetBool returns the value of the given key.
// It also casts the value type to bool internally.
// It converts the "true" and "false" string values to related boolean values.
// It returns an error if the related value is not a bool.
// It returns an error if there is no value for the given key.
func (c *Config) GetBool(key string) (bool, error) {
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

	return false, &TypeError{value: v, wanted: "bool"}
}

// GetStrictBool returns the value of the given key.
// It also casts the value type to bool internally.
// It doesn't convert the "true" and "false" string values to related boolean values.
// It returns an error if the related value is not a bool.
// It returns an error if there is no value for the given key.
func (c *Config) GetStrictBool(key string) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if v, ok := v.(bool); ok {
		return v, nil
	}

	return false, &TypeError{value: v, wanted: "bool"}
}

// New returns a brand new instance of Config with the given options.
func New(ops ...Options) (*Config, error) {
	c := &Config{}

	for _, op := range ops {
		if op.Env != "" {
			err := c.FeedEnv(op.Env)
			if err != nil {
				return nil, err
			}
		}

		if op.Feeder != nil {
			if err := c.Feed(op.Feeder); err != nil {
				return nil, err
			}
		}
	}

	c.StartListener()

	return c, nil
}
