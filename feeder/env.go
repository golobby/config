package feeder

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Env struct {
	Path string
}

func (e *Env) Feed() (map[string]interface{}, error) {
	values, err := Load(e.Path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	for k, v := range values {
		m[k] = v
	}
	return m, nil
}

// Load reads the given env file and extracts the variables as a string map
func Load(filename string) (map[string]string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	variables, err := read(file)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// read opens the given env file and extracts the variables as a string map
func read(file io.Reader) (map[string]string, error) {
	items := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value, err := parse(scanner.Text())
		if err != nil {
			return nil, err
		}

		if key != "" {
			items[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// parse extracts the key/value from the given line
func parse(line string) (string, string, error) {
	ln := strings.TrimSpace(line)

	if len(ln) == 0 {
		return "", "", nil
	}

	if ln[0] == '#' {
		return "", "", nil
	}

	s := strings.Index(ln, "=")
	if s == -1 {
		return "", "", errors.New("Invalid line: " + ln)
	}

	k := strings.TrimSpace(ln[:s])
	v := strings.TrimSpace(ln[s+1:])

	return k, v, nil
}
