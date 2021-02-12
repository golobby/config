package feeder

import "os"

// OS is a feeder that feeds using a OS variables.
type OS struct {
	Variables []string // The variable names that should be imported
	Strict    bool     // Set true to import empty variables
}

// Feed returns all the content.
func (s *OS) Feed() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if s.Variables != nil {
		for _, k := range s.Variables {
			v := os.Getenv(k)
			if s.Strict || v != "" {
				m[standardize(k)] = v
			}
		}
	}

	return m, nil
}
