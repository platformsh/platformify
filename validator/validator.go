package validator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ValidateFile checks the file exists and is valid yaml, then returns the unmarshalled data.
func ValidateFile(path string) (map[string]interface{}, error) {
	// Does the file exist?
	rawData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read %v file: %v", path, err)
	}

	// Does the yamlData unmarshal as yaml?
	var yamlData map[string]interface{}
	err = yaml.Unmarshal(rawData, &yamlData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML yamlData for %v: %v", path, err)
	}

	// Return the unmarshalled yaml yamlData.
	return yamlData, nil
}

// ValidateConfig uses ValidateFile and to check config for a given directory is valid Platform.sh config.
func ValidateConfig(path string) error {
	var errs error
	files := [3]string{".platform.app.yaml", ".platform/routes.yaml", ".platform/services.yaml"}
	for _, file := range files {
		_, fileErr := ValidateFile(filepath.Join(path, file))
		if fileErr != nil {
			errs = errors.Join(errs, fmt.Errorf("problem with %v: %v", file, fileErr))
		}
	}

	return errs
}
