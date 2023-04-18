package utils

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/exp/slices"
)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
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

func NewFileCreator() *FileCreator {
	return &FileCreator{}
}

type FileCreator struct{}

func (f *FileCreator) Create(filePath string) (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filePath)
}
