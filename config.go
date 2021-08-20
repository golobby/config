// Package config is a lightweight yet powerful configuration management library.
// It takes advantage of dot env files and OS variables alongside config files to be your ultimate requirement.
package config

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Feeder is an interface for configuration Feeders that provide configuration data.
type Feeder interface {
	// Feed gets a struct and feeds it using configuration data.
	Feed(structure interface{}) error
}

// Config is the configuration manager.
// To use the package facilities, there should be at least one instance of it.
// It holds the configuration feeders and structures that it is going to feed them.
type Config struct {
	Feeders    []Feeder      // Feeders is the list of configuration feeders that provides configuration data.
	Structures []interface{} // Structures is the list of structures that are going to be fed.
}

// New creates a brand new instance of Config to use the package facilities.
// It gets feeders that are going to feed the configuration structures.
func New(feeders ...Feeder) *Config {
	return &Config{Feeders: feeders}
}

// Feed gets a structure and feeds it using the provided feeders.
func (c *Config) Feed(structure interface{}) error {
	c.Structures = append(c.Structures, structure)
	return c.feedStructure(structure)
}

// Refresh refreshes registered structures using the provided feeders.
func (c *Config) Refresh() error {
	for _, s := range c.Structures {
		if err := c.feedStructure(s); err != nil {
			return err
		}
	}

	return nil
}

// WithListener adds an OS signal listener to the Config instance.
// The listener listens to the SIGHUP signal and refreshes the Config instance.
// It would call the provided fallback if the refresh process failed.
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

// feedStructure gets a structure and feeds it using all the provided feeders.
func (c *Config) feedStructure(structure interface{}) error {
	for _, f := range c.Feeders {
		if err := f.Feed(structure); err != nil {
			return fmt.Errorf("config: faild to feed struct; err %v", err)
		}
	}

	return nil
}
