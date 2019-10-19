// Package Config is a lightweight yet powerful Config package for Go projects.
// It takes advantage of env files and OS variables alongside Config files to be your ultimate requirement.
package config

import (
	"errors"
	"github.com/golobby/config/env"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

// Feeder is an interface for paths which can feed Config instances (provider their contents).
type Feeder interface {
	Feed() (map[string]interface{}, error)
}

// Options is a struct that contains all the required data for instantiating a new Config instance.
type Options struct {
	Feeder   Feeder // Feeder that is going to feed the Config instance
	Env      string // GetEnv is the .env file that is going to be used in Config file values
	Listener bool   // If true it listens to OS signal to reload Config end env files
}

// Config is the main struct that keeps all the Config instance data.
type Config struct {
	env struct {
		paths []string          // paths keeps all the added env file paths
		items map[string]string // Items keeps all the given .env key/value Items
		sync  sync.RWMutex      // sync is responsible for lock/unlock the env
	}
	feeders []Feeder               // paths keeps all the added paths
	Items   map[string]interface{} // Items keeps the Config data
	sync    sync.RWMutex           // sync is responsible for lock/unlock the config
}

// FeedEnv will add key/value Items from given env file to the config instance
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

// ReloadEnv will reload all the added env files and apply new changes
func (c *Config) ReloadEnv() error {
	for _, p := range c.env.paths {
		if err := c.FeedEnv(p); err != nil {
			return err
		}
	}

	return nil
}

// GetEnv will return environment variable value for the given environment variable key.
func (c Config) GetEnv(key string) string {
	c.env.sync.RLock()
	defer c.env.sync.RUnlock()

	v, ok := c.env.items[key]

	if ok && v != "" {
		return v
	}

	return os.Getenv(key)
}

// SetEnv will set value for the given env key
func (c *Config) SetEnv(key, value string) {
	c.env.sync.Lock()
	defer c.env.sync.Unlock()

	if c.env.items == nil {
		c.env.items = map[string]string{}
	}

	c.env.items[key] = value
}

// StartListener will make Config to listen to SIGINFO signal and reload the feeders and env files
func (c *Config) StartListener() {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINFO)

	go func() {
		for {
			<-s
			_ = c.ReloadEnv()
			_ = c.Reload()
		}
	}()
}

// Feed will feed the Config instance using the given feeder.
// It accepts all kinds of paths that implement the Feeder interface.
// The built-in paths are in the feeder subpackage.
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

// Reload will reload all the added feeders and applies new changes
func (c *Config) Reload() error {
	for _, f := range c.feeders {
		if err := c.Feed(f); err != nil {
			return err
		}
	}

	return nil
}

// Set will store the given key/value into the Config instance.
// It keeps the key/values that have added on runtime in the memory.
// It won't change the Config files.
func (c *Config) Set(key string, value interface{}) {
	c.sync.Lock()
	defer c.sync.Unlock()

	if c.Items == nil {
		c.Items = map[string]interface{}{}
	}

	c.Items[key] = value
}

// Get will return the value of the given key.
// The return type is interface, so it probably should be cast to the related data type.
// It will return an error if there is no value for the given key.
func (c Config) Get(key string) (interface{}, error) {
	c.sync.RLock()
	defer c.sync.RUnlock()

	v, ok := c.Items[key]

	if ok {
		return v, nil
	}

	if strings.Contains(key, ".") == false {
		return nil, errors.New("value not found for the key " + key)
	}

	v, err := lookup(c.Items, key)

	return v, err
}

// Get will return the value of the given key.
// It casts the value type to string.
// It will return an error if the related value is not string.
// It will return an error if there is no value for the given key.
func (c Config) GetString(key string) (string, error) {
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
func (c Config) GetInt(key string) (int, error) {
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
func (c Config) GetFloat(key string) (float64, error) {
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
func (c Config) GetBool(key string) (bool, error) {
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
func (c Config) GetStrictBool(key string) (bool, error) {
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
func (c Config) parse(value interface{}) interface{} {
	if stmt, ok := value.(string); ok {
		if stmt[0:2] == "${" && stmt[len(stmt)-1:] == "}" {
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
	}

	return value
}

// lookup will search for the given key recursively.
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

// New will return a brand new instance of Config with given option.
func New(ops Options) (*Config, error) {
	c := &Config{}

	if ops.Env != "" {
		err := c.FeedEnv(ops.Env)
		if err != nil {
			return nil, err
		}
	}

	if ops.Feeder != nil {
		if err := c.Feed(ops.Feeder); err != nil {
			return nil, err
		}
	}

	if ops.Listener {
		c.StartListener()
	}

	return c, nil
}
