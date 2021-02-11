// Package Config is a lightweight yet powerful configuration management tool for Go projects.
// It takes advantage of env files and OS variables alongside Config files to be your only requirement.
package config

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// Feeder is an interface for config feeders that provide content of a config instance.
type Feeder interface {
	Feed() (map[string]interface{}, error)
}

// NotFoundError happens when it cannot find the requested key.
type NotFoundError struct {
	key string
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("value not found for the key `%s`", n.key)
}

// TypeError happens when it cannot cast a value to the requested type.
type TypeError struct {
	value  interface{}
	wanted string
}

func (t *TypeError) Error() string {
	return fmt.Sprintf("value `%v` (`%T`) is not `%s`", t.value, t.value, t.wanted)
}

// Config keeps all the Config instance data.
type Config struct {
	feeders []Feeder               // It keeps all the added feeders
	items   map[string]interface{} // It keeps all the key/value items
	sync    sync.RWMutex           // It's responsible for (un)locking the items
}

// StartListener makes the instance to listen to the SIGHUP and reload the feeders.
func (c *Config) StartListener() {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			_ = c.Reload()
		}
	}()
}

// Feed takes a feeder and feeds the instance with it.
// The built-in feeders are in the feeder sub-package.
func (c *Config) Feed(f Feeder) error {
	err := c.feedItems(f)
	if err != nil {
		return err
	}

	c.feeders = append(c.feeders, f)

	return nil
}

func (c *Config) feedItems(f Feeder) error {
	items, err := f.Feed()
	if err != nil {
		return err
	}

	for k, v := range items {
		c.Set(k, v)
	}

	return nil
}

// Reload reloads all the added feeders and applies new changes.
func (c *Config) Reload() error {
	for _, f := range c.feeders {
		if err := c.feedItems(f); err != nil {
			return err
		}
	}

	return nil
}

// Set stores the given key/value into the Config instance.
// It keeps all the changes in memory and won't change the Config files.
func (c *Config) Set(key string, value interface{}) {
	c.sync.Lock()
	defer c.sync.Unlock()

	if c.items == nil {
		c.items = map[string]interface{}{}
	}

	c.items[key] = value
}

// Get returns the value of the given key.
// The return type is "interface{}".
// It probably needs to be cast to the related data type.
// It returns an error if there is no value for the given key.
func (c *Config) Get(key string) (interface{}, error) {
	c.sync.RLock()
	defer c.sync.RUnlock()

	v, ok := c.items[key]

	if ok {
		return v, nil
	}

	if strings.Contains(key, ".") == false {
		return nil, &NotFoundError{key: key}
	}

	return lookup(c.items, key)
}

// GetAll returns all the configuration items (key/values).
func (c *Config) GetAll() map[string]interface{} {
	return c.items
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

	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		i, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, &TypeError{value: v, wanted: "float64"}
		}
		return i, nil
	}

	return 0, &TypeError{value: v, wanted: "float64"}
}

// GetBool returns the value of the given key.
// It also casts the value type to bool internally.
// It converts the "true", "false", 1 and 0 values to related boolean values.
// It returns an error if the related value is not a bool.
// It returns an error if there is no value for the given key.
func (c *Config) GetBool(key string) (bool, error) {
	v, err := c.Get(key)
	if err != nil {
		return false, err
	}

	if b, ok := v.(bool); ok {
		return b, nil
	} else if b, ok := v.(string); ok {
		if b == "true" {
			return true, nil
		} else if b == "false" {
			return false, nil
		}
	} else if b, ok := v.(int); ok {
		if b == 1 {
			return true, nil
		} else if b == 0 {
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

// lookup searches for the given key in deep and returns related value if exist.
func lookup(collection interface{}, key string) (interface{}, error) {
	keys := strings.Split(key, ".")

	if len(keys) == 1 {
		return find(collection, keys[0])
	}

	c, err := dig(collection, keys[0])
	if err != nil {
		return nil, err
	}

	return lookup(c, strings.Join(keys[1:], "."))
}

// find returns the value of given key in the given 1D collection
func find(collection interface{}, key string) (interface{}, error) {
	switch collection.(type) {
	case map[interface{}]interface{}:
		if v, ok := collection.(map[interface{}]interface{})[key]; ok {
			return v, nil
		}
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

	return nil, &NotFoundError{key: key}
}

// dig returns the sub-collection of the given collection by the given key.
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

	return nil, &NotFoundError{key: key}
}

// New makes a brand new instance of Config with the given feeders.
func New(feeders ...Feeder) (*Config, error) {
	c := &Config{}

	for _, fd := range feeders {
		if err := c.Feed(fd); err != nil {
			return nil, err
		}
	}

	c.StartListener()

	return c, nil
}
