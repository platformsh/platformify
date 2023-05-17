package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// GetJSONKey gets a value from a JSON file, by traversing the path given
func GetJSONKey(jsonPath []string, filePath string) (value interface{}, ok bool) {
	fin, err := os.Open(filePath)
	if err != nil {
		return nil, false
	}
	defer fin.Close()

	var data map[string]interface{}
	err = json.NewDecoder(fin).Decode(&data)
	if err != nil {
		return nil, false
	}

	if len(jsonPath) == 0 {
		return data, true
	}

	for _, key := range jsonPath[:len(jsonPath)-1] {
		if value, ok = data[key]; !ok {
			return nil, false
		}

		if data, ok = value.(map[string]interface{}); !ok {
			return nil, false
		}
	}

	if value, ok = data[jsonPath[len(jsonPath)-1]]; !ok {
		return nil, false
	}

	return value, true
}

// ContainsStringInFile checks if the given file contains the given string
func ContainsStringInFile(file io.Reader, target string) (bool, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), target) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}
