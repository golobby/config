package config

import (
	"strconv"
	"strings"

	"github.com/golobby/config/assign"
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

// The setter interface is used to solve polymorphism of Set() method using for Reload() method.
type setter interface {
	Set(key string, value interface{})
}

// ConfigBase keeps all the Config instance data.
// The ConfigBase is NOT goroutine safe.
type ConfigBase struct {
	EnvConfig
	feeders []Feeder               // It keeps all the added feeders
	items   map[string]interface{} // It keeps all the key/value items (excluding environment ones).
}

// Get the Config's env instance.
func (c *ConfigBase) Env() *EnvConfig {
	return &c.EnvConfig
}

// Feed takes a feeder and feeds the Config instance with it.
// The built-in feeders are in the feeder subpackage.
func (c *ConfigBase) Feed(f Feeder) error {
	err := c.doFeed(f, c)
	if err != nil {
		return err
	}

	c.feeders = append(c.feeders, f)

	return nil
}

func (c *ConfigBase) doFeed(f Feeder, s setter) error {
	items, err := f.Feed()
	if err != nil {
		return err
	}

	for k, v := range items {
		s.Set(k, c.parse(v))
	}

	return nil
}

// Reload reloads all the added feeders and applies new changes.
func (c *ConfigBase) Reload() error {
	return c.doReload(c)
}

func (c *ConfigBase) doReload(s setter) error {
	for _, f := range c.feeders {
		if err := c.doFeed(f, s); err != nil {
			return err
		}
	}

	return nil
}

// Set stores the given key/value into the Config instance.
// It keeps all the changes in memory and won't change the Config files.
func (c *ConfigBase) Set(key string, value interface{}) {
	if c.items == nil {
		c.items = map[string]interface{}{}
	}

	c.items[key] = value
}

// Get returns the value of the given key.
// The return type is "interface{}".
// It probably needs to be cast to the related data type.
// It returns false if there is no value for the given key.
func (c *ConfigBase) Get(key string) (interface{}, bool) {
	if v, ok := c.items[key]; ok {
		return v, true
	}

	if strings.IndexByte(key, '.') < 0 {
		return nil, false
	}

	return lookup(c.items, key)
}

// GetAll returns all the configuration items (key/values).
func (c *ConfigBase) GetAll() map[string]interface{} {
	return c.items
}

// Assigns struct fields' value by its field's tag (such as the json tag).
// @param ptr The pointer of struct's instance to set
// @param key Specify where to get the struct's value
// @param tag Specify which struct field's tag name used to retrieve
// @return The count of fields that been assigned, -1 if struct's value not found by the key
func (c *ConfigBase) AssignStruct(ptr interface{}, key, tag string) int {
	if data, found := c.Get(key); found {
		return assign.AssignStruct(ptr, data, tag)
	}

	return -1
}

// Assigns slice elements.
// @param ptr The pointer of slice instance to appent elements
// @param key Specify where to get the slice elements's value
// @param tag If element's type is struct, using the tag name to retrieve struct fields
// @return The count of elements that been assigned, -1 if slice's value not found by the key
func (c *ConfigBase) AssignSlice(ptr interface{}, key, tag string) int {
	if data, found := c.Get(key); found {
		return assign.AssignSlice(ptr, data, tag)
	}
	return -1
}

// parse replaces the placeholders with environment and OS variables.
func (c *ConfigBase) parse(value interface{}) interface{} {
	if stmt, ok := value.(string); ok {
		if sLen := len(stmt); sLen > 3 && stmt[0] == '$' && stmt[sLen-1] == '}' && stmt[0:2] == "${" {
			pipe := strings.IndexByte(stmt, '|')

			if pipe == -1 {
				key := strings.TrimSpace(stmt[2 : sLen-1])
				return c.GetEnv(key)
			}

			key := strings.TrimSpace(stmt[2 : pipe])
			if v := c.GetEnv(key); v != "" {
				return v
			}

			return strings.TrimSpace(stmt[pipe+1 : sLen-1])
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

		return collection
	}

	return value
}

// lookup searches for the given key in deep and returns related value.
func lookup(collection interface{}, key string) (interface{}, bool) {
	rest := key
	key, rest = segmentKey(rest)

	if rest == "" {
		return find(collection, key)
	}

	c, ok := dig(collection, key)
	if !ok {
		return nil, false
	}

	return lookup(c, rest)
}

// segment the key by dot('.'), returns the first segment before dot and the rest.
func segmentKey(rest string) (string, string) {
	key := ""

	for key == "" {
		key = rest
		i := strings.IndexByte(key, '.')
		if i < 0 {
			return key, ""
		}

		rest, key = key[i+1:], key[:i]
	}

	return key, rest
}

// find returns the value of given key in the given 1D collection.
func find(collection interface{}, key string) (interface{}, bool) {
	switch collection.(type) {
	case map[string]interface{}:
		if v, ok := collection.(map[string]interface{})[key]; ok {
			return v, true
		}
	case []interface{}:
		k, err := strconv.Atoi(key)
		if err == nil && k >= 0 && len(collection.([]interface{})) > k {
			return collection.([]interface{})[k], true
		}
	case map[interface{}]interface{}:
		if v, ok := collection.(map[interface{}]interface{})[key]; ok {
			return v, true
		}
	}

	return nil, false
}

// dig returns the sub-collection of the given collection by the given key.
func dig(collection interface{}, key string) (interface{}, bool) {
	if v, ok := collection.(map[string]interface{}); ok {
		if v, ok := v[key]; ok {
			return v, true
		}
	} else if v, ok := collection.([]interface{}); ok {
		i, err := strconv.Atoi(key)
		if err == nil && i >= 0 && len(v) > i {
			return v[i], true
		}
	}

	return nil, false
}

// Initialize the instance of ConfigBase with the given options.
func (c *ConfigBase) init(ops ...Options) error {
	for _, op := range ops {
		if op.Env != "" {
			err := c.FeedEnv(op.Env)
			if err != nil {
				return err
			}
		}

		if op.Feeder != nil {
			if err := c.Feed(op.Feeder); err != nil {
				return err
			}
		}
	}
	return nil
}

// New returns a brand new instance of ConfigBase with the given options.
func NewBase(ops ...Options) (*ConfigBase, error) {
	c := &ConfigBase{}

	err := c.init(ops...)

	return c, err
}
