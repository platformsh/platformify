package platformifier

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

var skipDirs = []string{
	"vendor",
	"node_modules",
	".next",
	".git",
}

//go:generate mockgen -destination=fs_mock_test.go -package=platformifier -source=fs.go
type FS interface {
	Create(name string) (io.WriteCloser, error)
	Find(root, name string, firstMatch bool) []string
	Open(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error)
}

func NewOSFileSystem(root string) *OSFileSystem {
	return &OSFileSystem{
		root: root,
	}
}

type OSFileSystem struct {
	root string
}

func (f *OSFileSystem) Open(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	return os.OpenFile(f.fullPath(name), flag, perm)
}

func (f *OSFileSystem) Create(name string) (io.WriteCloser, error) {
	filePath := f.fullPath(name)
	if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filePath)
}

// Find searches for the file inside the path recursively and returns all matches
func (f *OSFileSystem) Find(root, name string, firstMatch bool) []string {
	if root == "" {
		root = "."
	}
	root = strings.TrimPrefix(root, "/")
	found := make([]string, 0)
	_ = fs.WalkDir(f.readonly(), root, func(p string, d os.DirEntry, err error) error {
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
			if firstMatch {
				return errors.New("found")
			}
		}

		return nil
	})

	return found
}

func (f *OSFileSystem) readonly() fs.FS {
	return os.DirFS(f.root)
}

func (f *OSFileSystem) fullPath(name string) string {
	return filepath.Join(f.root, name)
}
