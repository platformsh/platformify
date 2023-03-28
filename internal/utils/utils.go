package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// FileExists checks if the file exists
func FileExists(searchPath, name string) bool {
	return FindFile(searchPath, name) != ""
}

// FindFile inside the path and return the full path of the file if found
func FindFile(searchPath, name string) string {
	var found string
	//nolint: errcheck
	filepath.WalkDir(searchPath, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
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
