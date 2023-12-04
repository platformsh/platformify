package utils

import (
	"bufio"
	"bytes"
	"cmp"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
	".git",
}

// FileExists checks if the file exists.
func FileExists(fileSystem fs.FS, searchPath, name string) bool {
	return FindFile(fileSystem, searchPath, name) != ""
}

// FindFile searches for the file inside the path recursively
// and returns the full path of the file if found.
// If multiple files exist, tries to return the one closest to root.
func FindFile(fileSystem fs.FS, searchPath, name string) string {
	files := FindAllFiles(fileSystem, searchPath, name)
	if len(files) == 0 {
		return ""
	}

	slices.SortFunc(files, func(a, b string) int {
		return cmp.Compare(strings.Count(a, string(os.PathSeparator)), strings.Count(b, string(os.PathSeparator)))
	})
	return files[0]
}

// FindAllFiles searches for the file inside the path recursively and returns all matches.
func FindAllFiles(fileSystem fs.FS, searchPath, name string) []string {
	found := make([]string, 0)
	if searchPath == "" {
		searchPath = "."
	}
	_ = fs.WalkDir(fileSystem, searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip vendor directories.
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

func GetMapValue(keyPath []string, data map[string]any) (value any, ok bool) {
	if len(keyPath) == 0 {
		return data, true
	}

	for _, key := range keyPath[:len(keyPath)-1] {
		if value, ok = data[key]; !ok {
			return nil, false
		}

		if data, ok = value.(map[string]any); !ok {
			return nil, false
		}
	}

	if value, ok = data[keyPath[len(keyPath)-1]]; !ok {
		return nil, false
	}

	return value, true
}

// GetJSONValue gets a value from a JSON file, by traversing the path given.
func GetJSONValue(fileSystem fs.FS, keyPath []string, filePath string, caseInsensitive bool) (value any, ok bool) {
	rawData, err := fs.ReadFile(fileSystem, filePath)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]any
	err = json.Unmarshal(rawData, &data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}

// ContainsStringInFile checks if the given file contains the given string.
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

// GetTOMLValue gets a value from a TOML file, by traversing the path given.
func GetTOMLValue(fileSystem fs.FS, keyPath []string, filePath string, caseInsensitive bool) (value any, ok bool) {
	rawData, err := fs.ReadFile(fileSystem, filePath)
	if err != nil {
		return nil, false
	}

	if caseInsensitive {
		rawData = bytes.ToLower(rawData)
		for i := range keyPath {
			keyPath[i] = strings.ToLower(keyPath[i])
		}
	}

	var data map[string]any
	err = toml.Unmarshal(rawData, &data)
	if err != nil {
		return nil, false
	}

	return GetMapValue(keyPath, data)
}
