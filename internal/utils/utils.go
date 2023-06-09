package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"
)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
	".git",
}

// FileExists checks if the file exists
func FileExists(searchPath, name string) bool {
	return FindFile(searchPath, name) != ""
}

// FindFile searches for the file inside the path recursively
// and returns the full path of the file if found
func FindFile(searchPath, name string) string {
	var found string
	//nolint:errcheck
	filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			if slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() == name {
			found = p
			return errors.New("found")
		}

		return nil
	})

	return found
}

// FindAllFiles searches for the file inside the path recursively and returns all matches
func FindAllFiles(searchPath, name string) []string {
	found := make([]string, 0)
	_ = filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories
			if slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if d.Name() == name {
			found = append(found, p)
		}

		return nil
	})

	return found
}

func GetMapValue(keyPath []string, data map[string]interface{}) (value interface{}, ok bool) {
	if len(keyPath) == 0 {
		return data, true
	}

	for _, key := range keyPath[:len(keyPath)-1] {
		if value, ok = data[key]; !ok {
			return nil, false
		}

		if data, ok = value.(map[string]interface{}); !ok {
			return nil, false
		}
	}

	if value, ok = data[keyPath[len(keyPath)-1]]; !ok {
		return nil, false
	}

	return value, true
}

// GetJSONValue gets a value from a JSON file, by traversing the path given
func GetJSONValue(keyPath []string, filePath string, caseInsensitive bool) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	rawData, err := io.ReadAll(fin)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]interface{}
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}

// ContainsStringInFile checks if the given file contains the given string
func ContainsStringInFile(file io.Reader, target string, caseInsensitive bool) (bool, error) {
	if caseInsensitive {
		target = strings.ToLower(target)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if caseInsensitive {
			if strings.Contains(strings.ToLower(scanner.Text()), target) {
				return true, nil
			}
		} else {
			if strings.Contains(scanner.Text(), target) {
				return true, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

// GetTOMLValue gets a value from a TOML file, by traversing the path given
func GetTOMLValue(keyPath []string, filePath string, caseInsensitive bool) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	rawData, err := io.ReadAll(fin)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]interface{}
	err = toml.Unmarshal(rawData, data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}
