package platformifier

import (
	"io"
	"os"
	"path"
)

//go:generate mockgen -destination=mock_fs_test.go -package=platformifier -source=fs.go
type FS interface {
	CreateFile(name string) (io.WriteCloser, error)
}

type OSFileSystem struct{}

func (fs *OSFileSystem) CreateFile(filePath string) (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Dir(filePath), os.ModeDir|os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(filePath)
}
