// Package Config is a lightweight yet powerful configuration management tool for Go projects.
// It takes advantage of env files and OS variables alongside Config files to be your only requirement.
package config

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/golobby/config/env"
)

// Feeder is an interface for config feeders that provide content of a config instance.
type Feeder interface {
	Feed() (map[string]interface{}, error)
}

// Options will contain all the required data for instantiating a new Config instance.
type Options struct {
	Feeder Feeder // Feeder is the feeder that is going to feed the Config instance.
	Env    string // Env is the file path that locates the environment file.
}

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
	env struct {
		paths []string          // It keeps all the added environment files' paths
		items map[string]string // It keeps all the given environment key/value items.
		sync  sync.RWMutex      // It's responsible for (un)locking the items
	}
	feeders []Feeder               // It keeps all the added feeders
	items   map[string]interface{} // It keeps all the key/value items (excluding environment ones).
	sync    sync.RWMutex           // It's responsible for (un)locking the items
}

// FeedEnv reads the given environment file path, extract key/value items, and add them to the Config instance.
func (c *Config) FeedEnv(path string) error {
	items, err := env.Load(path)
	if err != nil {
		return err
	}

	for k, v := range items {
		c.SetEnv(k, v)
	}

	c.env.paths = append(c.env.paths, path)

	return nil
}

// ReloadEnv reloads all the added environment files and applies new changes.
func (c *Config) ReloadEnv() error {
	for _, p := range c.env.paths {
		if err := c.FeedEnv(p); err != nil {
			return err
		}
	}

	return nil
}

// GetEnv returns the environment variable value for the given environment variable key.
func (c *Config) GetEnv(key string) string {
	c.env.sync.RLock()
	defer c.env.sync.RUnlock()

	v, ok := c.env.items[key]

	if ok && v != "" {
		return v
	}

	return os.Getenv(key)
}

// GetAllEnvs returns all the environment variables (key/values)
func (c *Config) GetAllEnvs() map[string]string {
	return c.env.items
}

// SetEnv sets the given value for the given env key
func (c *Config) SetEnv(key, value string) {
	c.env.sync.Lock()
	defer c.env.sync.Unlock()

	if c.env.items == nil {
		c.env.items = map[string]string{}
	}

	c.env.items[key] = value
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

// Feed takes a feeder and feeds the Config instance with it.
// The built-in feeders are in the feeder subpackage.
func (c *Config) Feed(f Feeder) error {
	items, err := f.Feed()
	if err != nil {
		return err
	}

	for k, v := range items {
		c.Set(k, c.parse(v))
	}

	c.feeders = append(c.feeders, f)

	return nil
}

// Reload reloads all the added feeders and applies new changes.
func (c *Config) Reload() error {
	for _, f := range c.feeders {
		if err := c.Feed(f); err != nil {
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

	v, err := lookup(c.items, key)

	return v, err
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

	if v, ok := v.(int); ok {
		return v, nil
	}

	if v, ok := v.(float64); ok {
		return int(v), nil
	}

	if v, ok := v.(string); ok {
		return strconv.Atoi(v)
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

// parse replaces the placeholders with environment and OS variables.
func (c *Config) parse(value interface{}) interface{} {
	if stmt, ok := value.(string); ok {
		if len(stmt) > 3 && stmt[0:2] == "${" && stmt[len(stmt)-1:] == "}" {
			pipe := strings.Index(stmt, "|")

			if pipe == -1 {
				key := strings.TrimSpace(stmt[2 : len(stmt)-1])
				return c.GetEnv(key)
			}

			key := strings.TrimSpace(stmt[2:pipe])
			if v := c.GetEnv(key); v != "" {
				return v
			}

			return strings.TrimSpace(stmt[pipe+1 : len(stmt)-1])
		}
	} else if collection, ok := value.(map[string]interface{}); ok {
		for k, v := range collection {
			collection[k] = c.parse(v)
		}

		return collection
	} else if collection, ok := value.(map[interface{}]interface{}); ok {
		for k, v := range collection {
			collection[k] = c.parse(v)
		}
	}

	return value
}

// lookup searches for the given key in deep and returns related value.
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

	return nil, errors.New("value not found for the key " + key)
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

	return nil, errors.New("value not found for the key " + key)
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
