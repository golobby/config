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
// It holds the configuration feeders and structs.
type Config struct {
	Feeders []Feeder      // Feeders is the list of feeders that provides configuration data.
	Structs []interface{} // Structs is the list of structs that holds the configuration data.
}

// New creates a brand new instance of Config to use the package facilities.
func New() *Config {
	return &Config{}
}

// AddFeeder adds a feeder that provides configuration data.
func (c *Config) AddFeeder(f Feeder) *Config {
	c.Feeders = append(c.Feeders, f)
	return c
}

// AddStruct adds a struct that holds the configuration data.
func (c *Config) AddStruct(s interface{}) *Config {
	c.Structs = append(c.Structs, s)
	return c
}

// Feed binds configuration data from added feeders to the added structs.
func (c *Config) Feed() error {
	for _, s := range c.Structs {
		for _, f := range c.Feeders {
			if err := c.feedStruct(f, s); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetupListener adds an OS signal listener to the Config instance.
// The listener listens to the `SIGHUP` signal and refreshes the Config instance.
// It would call the provided fallback if the refresh process failed.
func (c *Config) SetupListener(fallback func(err error)) *Config {
	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGHUP)

	go func() {
		for {
			<-s
			if err := c.Feed(); err != nil {
				fallback(err)
			}
		}
	}()

	return c
}

// feedStruct feeds a struct using given feeder.
func (c *Config) feedStruct(f Feeder, s interface{}) error {
	if err := f.Feed(s); err != nil {
		return fmt.Errorf("config: failed to feed struct; err %v", err)
	}

	return nil
}
