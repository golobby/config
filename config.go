// Package config is a lightweight yet powerful configuration management library.
// It takes advantage of dot env files and OS variables alongside config files to be your only requirement.
package config

// Feeder is an interface for config feeders that provide content of a config instance.
type Feeder interface {
	Feed(structure interface{}) error
}

func Load(structure interface{}, feeders ...Feeder) error {
	for _, f := range feeders {
		if err := f.Feed(structure); err != nil {
			return err
		}
	}

	return nil
}
