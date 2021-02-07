package feeder

import "github.com/golobby/config/env"

type Env struct {
	Path string
}

func (e *Env) Feed() (map[string]interface{}, error) {
	values, err := env.Load(e.Path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	for k, v := range values {
		m[k] = v
	}
	return m, nil
}
