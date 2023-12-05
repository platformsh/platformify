package discovery

import (
	"io/fs"
)

const (
	settingsPyFile   = "settings.py"
	managePyFile     = "manage.py"
	composerJSONFile = "composer.json"
	packageJSONFile  = "package.json"
	symfonyLockFile  = "symfony.lock"
)

// Discoverer detects project characteristics from the filesystem.
type Discoverer struct {
	fileSystem fs.FS
	memory     map[string]any
}

// New creates a Discoverer for the given filesystem.
func New(fileSystem fs.FS) *Discoverer {
	return &Discoverer{fileSystem: fileSystem, memory: make(map[string]any)}
}
