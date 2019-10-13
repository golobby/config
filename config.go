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

	e := errors.New("find not found for the key " + key)

	if strings.Index(key, ".") == -1 {
		return nil, e
	}

	stack := []interface{}{c}

	keys := strings.Split(key, ".")

	for i, key := range keys {
		top := stack[len(stack)-1]

		if i == len(keys)-1 {
			return find(top, key)
		}

		if v, ok := top.(Config); ok {
			if v, ok := v[key]; ok {
				stack = append(stack, v)
				continue
			}
		} else if v, ok := top.(map[string]interface{}); ok {
			if v, ok := v[key]; ok {
				stack = append(stack, v)
				continue
			}
		} else if v, ok := top.([]interface{}); ok {
			i, err := strconv.Atoi(key)
			if err != nil {
				return nil, e
			} else if len(v) > i {
				stack = append(stack, v[i])
				continue
			}
		}

		return nil, e
	}

	return nil, e
}

func (c Config) GetString(key string) (string, error) {
	v, err := c.Get(key)
	if err != nil {
		return "", err
	}

	if v, ok := v.(string); ok {
		return v, nil
	}

	return "", errors.New("find for " + key + " is not string")
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

func find(collection interface{}, key string) (interface{}, error) {
	e := errors.New("find not found for the key " + key)

	switch collection.(type) {
	case map[string]interface{}:
		if v, ok := collection.(map[string]interface{})[key]; ok {
			return v, nil
		}
	case []interface{}:
		k, err := strconv.Atoi(key)
		if err != nil {
			return nil, e
		} else if len(collection.([]interface{})) > k {
			return collection.([]interface{})[k], nil
		}
	}

	return nil, e
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
