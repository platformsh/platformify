package discovery

import (
	"path"
	"slices"

	"github.com/platformsh/platformify/internal/utils"
)

// Returns the application root, either from memory or by discovering it on the spot
func (d *Discoverer) ApplicationRoot() (string, error) {
	if applicationRoot, ok := d.memory["application_root"]; ok {
		return applicationRoot.(string), nil
	}

	appRoot, err := d.discoverApplicationRoot()
	if err != nil {
		return "", err
	}

	d.memory["application_root"] = appRoot
	return appRoot, nil
}

func (d *Discoverer) discoverApplicationRoot() (string, error) {
	depManagers, err := d.DependencyManagers()
	if err != nil {
		return "", err
	}

	for _, dependencyManager := range dependencyManagersMap {
		if !slices.Contains(depManagers, dependencyManager.name) {
			continue
		}

		lockPath := utils.FindFile(d.fileSystem, "", dependencyManager.lockFile)
		if lockPath == "" {
			continue
		}

		return path.Dir(lockPath), nil
	}

	return "", nil
}
