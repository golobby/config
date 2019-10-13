package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Repository interface {
	Extract() (map[string]interface{}, error)
}

type Config map[string]interface{}

func (c Config) Feed(r Repository) error {
	if items, err := r.Extract(); err != nil {
		return err
	} else {
		for k, v := range items {
			c[k] = parse(v)
		}

		return nil
	}
}

func (c Config) Get(key string) (interface{}, error) {
	if v, ok := c[key]; ok {
		return v, nil
	}

	if strings.Contains(key, ".") == false {
		return nil, errors.New("value not found for the key " + key)
	}

	return lookup(c, key)
}

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

func parse(value interface{}) interface{} {
	if stmt, ok := value.(string); ok {
		if stmt[0:2] == "${" && stmt[len(stmt)-1:] == "}" {
			pipe := strings.Index(stmt, "|")

			if pipe == -1 {
				name := strings.TrimSpace(stmt[2 : len(stmt)-1])
				return os.Getenv(name)
			}

			name := strings.TrimSpace(stmt[2:pipe])
			if v := os.Getenv(name); v != "" {
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

func New(rs ...Repository) (Config, error) {
	c := Config{}

	for _, r := range rs {
		if err := c.Feed(r); err != nil {
			return nil, err
		}
	}

	return c, nil
}
