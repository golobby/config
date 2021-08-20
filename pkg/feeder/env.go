package feeder

import (
	"fmt"
	"github.com/golobby/env/v2"
)

type Env struct{}

func (f Env) Feed(structure interface{}) error {
	if err := env.Feed(structure); err != nil {
		return fmt.Errorf("config: cannot feed struct; err: %v", err)
	}

	return nil
}
