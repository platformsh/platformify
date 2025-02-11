package discovery

import (
	"io/fs"

	"github.com/platformsh/platformify/platformifier"
)

const (
	settingsPyFile   = "settings.py"
	managePyFile     = "manage.py"
	composerJSONFile = "composer.json"
	packageJSONFile  = "package.json"
	symfonyLockFile  = "symfony.lock"
)

type Discoverer struct {
	fileSystem fs.FS
	memory     map[string]any
}

func New(fileSystem fs.FS) *Discoverer {
	return &Discoverer{fileSystem: fileSystem, memory: make(map[string]any)}
}

func (d *Discoverer) DeployCommand() ([]string, error) {
	return nil, nil
}
func (d *Discoverer) Locations() (map[string]map[string]any, error) {
	return nil, nil
}
func (d *Discoverer) Dependencies() (map[string]map[string]string, error) {
	// TODO: Add back dependencies
	// answers.Dependencies["nodejs"]["n"] = "*"
	// answers.Dependencies["nodejs"]["npx"] = "*"
	return nil, nil
}
func (d *Discoverer) Mounts() (map[string]map[string]string, error) {
	return nil, nil
}
func (d *Discoverer) Services() ([]platformifier.Service, error) {
	return nil, nil
}
