package config

import (
	"os"
	"sync"

	feed "github.com/golobby/config/env"
)

// EnvConfig keeps all the Config environments' data.
type EnvConfig struct {
	paths []string          // It keeps all the added environment files' paths
	items map[string]string // It keeps all the given environment key/value items.
	sync  sync.RWMutex      // It's responsible for (un)locking the items
}

// FeedEnv reads the given environment file path, extract key/value items, and add them to the Config instance.
func (env *EnvConfig) FeedEnv(path string) error {
	err := env.feedEnvItems(path)
	if err != nil {
		return err
	}

	env.paths = append(env.paths, path)

	return nil
}

func (env *EnvConfig) feedEnvItems(path string) error {
	items, err := feed.Load(path)
	if err != nil {
		return err
	}

	for k, v := range items {
		env.SetEnv(k, v)
	}

	return nil
}

// ReloadEnv reloads all the added environment files and applies new changes.
func (env *EnvConfig) ReloadEnv() error {
	for _, p := range env.paths {
		if err := env.feedEnvItems(p); err != nil {
			return err
		}
	}

	return nil
}

// GetEnv returns the environment variable value for the given environment variable key.
func (env *EnvConfig) GetEnv(key string) string {
	env.sync.RLock()
	defer env.sync.RUnlock()

	if env.items != nil {
		v, ok := env.items[key]

		if ok && v != "" {
			return v
		}
	}

	return os.Getenv(key)
}

// GetAllEnvs returns all the environment variables (key/values)
func (env *EnvConfig) GetAllEnvs() map[string]string {
	return env.items
}

// SetEnv sets the given value for the given env key
func (env *EnvConfig) SetEnv(key, value string) {
	env.sync.Lock()
	defer env.sync.Unlock()

	if env.items == nil {
		env.items = map[string]string{}
	}

	env.items[key] = value
}
