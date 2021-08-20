// Package config is a lightweight yet powerful configuration management library.
// It takes advantage of dot env files and OS variables alongside config files to be your only requirement.
package config

import (
	"os"
	"os/signal"
	"syscall"
)

// Feeder is an interface for config Feeders that provide content of a config instance.
type Feeder interface {
	Feed(structure interface{}) error
}

type Config struct {
	Feeders    []Feeder
	Structures []interface{}
	Fallback   func(err error)
}

func New(feeders ...Feeder) *Config {
	return &Config{Feeders: feeders}
}

func (c *Config) Feed(structure interface{}) error {
	c.Structures = append(c.Structures, structure)
	return c.feedStructure(structure)
}

func (c *Config) Refresh() error {
	for _, s := range c.Structures {
		if err := c.feedStructure(s); err != nil {
			return err
		}
	}

	return nil
}

// WithListener makes the instance to listen to the SIGHUP and reload the Feeders.
func (c *Config) WithListener(fallback func(err error)) *Config {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			if err := c.Refresh(); err != nil {
				fallback(err)
			}
		}
	}()

	return c
}

func (c *Config) feedStructure(structure interface{}) error {
	for _, f := range c.Feeders {
		if err := f.Feed(structure); err != nil {
			return err
		}
	}

	return nil
}
