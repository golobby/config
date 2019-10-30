// Package env is a simple package to read environment variable files.
// It parses env files and extracts their key/values as a string map.
package env

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

// Load reads the given env file and extracts the variables as a string map
func Load(filename string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path := wd + string(os.PathSeparator) + filename
	file, err := os.Open(path)
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
func read(file *os.File) (map[string]string, error) {
	items := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if key, value, err := parse(scanner.Text()); err != nil {
			return nil, err
		} else if key != "" {
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
